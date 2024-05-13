package main

import (
	"context"
	"database/sql"
	"net"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/steve-mir/diivix_backend/db/sqlc"
	"github.com/steve-mir/diivix_backend/internal/app/authentication/gapi/interceptor"
	"github.com/steve-mir/diivix_backend/internal/app/posts/gapi"
	"github.com/steve-mir/diivix_backend/internal/app/posts/pb"
	"github.com/steve-mir/diivix_backend/utils"
	"github.com/steve-mir/diivix_backend/worker"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func accessibleRoles() map[string][]int8 {
	return map[string][]int8{
		"/pb.SocialMedia/PostStream":             {1, 2, 9}, // -> admin and users
		"/pb.SocialMedia/CreatePost":             {3, 4, 5, 6, 7, 8, 9},
		"/pb.SocialMedia/LiveComments":           {3, 4, 5, 6, 7, 8, 9},
		"/pb.SocialMedia/ViewPost":               {2, 1},
		"/pb.SocialMedia/ViewUserProfile":        {3, 4, 5, 6, 7, 8, 9},
		"/pb.SocialMedia/GenericSearch":          {3, 4, 5, 6, 7, 8, 9},
		"/pb.SocialMedia/GetNotifications":       {9, 8, 7, 6, 5, 4, 3},
		"/pb.SocialMedia/DeleteNotifications":    {9, 8, 7, 6, 5, 4, 3},
		"/pb.SocialMedia/MarkNotificationAsRead": {9, 8, 7, 6, 5, 4, 3},
		"/pb.SocialMedia/GetPost":                {9, 8, 7, 6, 5, 4, 3},

		"/grpc.reflection.v1alpha.ServerReflection/ServerReflectionInfo": {2, 1},

		// Only users
		"/pb.SocialMedia/UnFollowUser":         {9, 8, 7, 6, 5, 4, 3},
		"/pb.SocialMedia/FollowUser":           {9, 8, 7, 6, 5, 4, 3},
		"/pb.SocialMedia/GetFollowers":         {1},
		"/pb.SocialMedia/GetFollowing":         {1},
		"/pb.SocialMedia/GetPostsByUserId":     {1},
		"/pb.SocialMedia/SuggestUsersToFollow": {1, 9},
		"/pb.SocialMedia/StreamLikes":          {1},
		"/pb.SocialMedia/LikePost":             {3, 4, 5, 6, 7, 8, 9},
		"/pb.SocialMedia/FetchPostComment":     {3, 4, 5, 6, 7, 8, 9},
		"/pb.SocialMedia/SendPostComment":      {3, 4, 5, 6, 7, 8, 9},
		"/pb.SocialMedia/CreateRepost":         {1},
		"/pb.SocialMedia/GetRepost":            {1},
		"/pb.SocialMedia/GetRepostsByUser":     {1},
		"/pb.SocialMedia/DeleteRepost":         {1},
		"/pb.SocialMedia/BookmarkPost":         {1},
		"/pb.SocialMedia/GetBookmarkedPosts":   {1},
		"/pb.SocialMedia/DeleteBookmark":       {1},
		"/pb.SocialMedia/BlockUser":            {1},
		"/pb.SocialMedia/GetBlockedUsers":      {1},
		"/pb.SocialMedia/UnblockUser":          {1},

		// Only admins
		"/pb.SocialMedia/BanUser": {2}, // -> admin

	}
}

func main() {

	// Use Viper for configuration management
	config, err := utils.LoadConfig(".")
	if err != nil {
		log.Fatal().Msg("cannot load config " + err.Error())
	}

	if config.Environment == "development" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	// Run db migrations
	// runDbMigration(config.MigrationUrl, config.DBSource)

	// Create the routes
	db, err := sqlc.CreateDbPool(config)
	if err != nil {
		log.Fatal().Msg("cannot create db pool")
		return
	}

	// Connect to redis
	redisOpt := asynq.RedisClientOpt{
		Addr:     config.RedisAddress,
		Username: config.RedisUsername,
		Password: config.RedisPwd,
	}

	taskDistributor := worker.NewRedisTaskDistributor(redisOpt)

	runSocialMediaServer(db, config, taskDistributor)

}

func setupNewServer(config utils.Config) *grpc.Server {
	// Define Interceptors
	intercpts := interceptor.NewAuthInterceptor(config, accessibleRoles())
	// Create a limiter that allows 10 requests per second with a burst size of 2
	ratelimitInterceptor := interceptor.NewRateLimitInterceptor(10, 2)

	// Chain the unary and stream interceptors
	chainedUnary := grpc_middleware.ChainUnaryServer(
		interceptor.GrpcLogger,
		ratelimitInterceptor.Unary(),
		intercpts.Unary(),
	)

	chainedStream := grpc_middleware.ChainStreamServer(
		interceptor.GrpcStreamLogger,
		ratelimitInterceptor.Stream(),
		intercpts.Stream(),
	)

	return grpc.NewServer(
		grpc.UnaryInterceptor(chainedUnary),
		grpc.StreamInterceptor(chainedStream),
	)
}

func runSocialMediaServer(db *sql.DB, config utils.Config, taskDistributor worker.TaskDistributor) {
	ctx, cancel := context.WithCancel(context.Background())

	// Handle graceful shutdown on receiving termination signals
	go func() {
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
		<-signals
		log.Info().Msg("Received termination signal. Shutting down gracefully...")
		cancel()
	}()

	// Auth server
	server, err := gapi.NewSocialMediaServer(db, config, taskDistributor)
	if err != nil {
		log.Fatal().Msg("cannot create a server:")
	}

	grpcServer := setupNewServer(config)

	pb.RegisterSocialMediaServer(grpcServer, server)

	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GRPCSocialMediaServerAddress)
	if err != nil {
		log.Fatal().Msg("cannot create listener:")
	}

	log.Info().Msgf("start grpc server at %s", listener.Addr().String())

	// Start the gRPC server in a goroutine
	go func() {
		err := grpcServer.Serve(listener)
		if err != nil {
			log.Fatal().Msg("cannot start grpc server")
		}
	}()

	// Wait for the context to be canceled (either by the termination signal or an error)
	<-ctx.Done()

	// Stop the gRPC server
	grpcServer.GracefulStop()

	// Log a message indicating a graceful shutdown
	log.Info().Msg("gRPC server stopped gracefully")
}

// func runTaskProcessor(redisOpt asynq.RedisClientOpt, store *sqlc.Store, db *sql.DB) {
// 	mailer := worker.NewRedisTaskProcessor(redisOpt, store, db)
// 	log.Info().Msg("Starting task processor")
// 	err := mailer.Start()
// 	if err != nil {
// 		log.Fatal().Msg("cannot start task processor")
// 	}

// }

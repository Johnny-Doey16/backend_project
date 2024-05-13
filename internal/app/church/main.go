package main

import (
	"context"
	"database/sql"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/steve-mir/diivix_backend/db/sqlc"
	"github.com/steve-mir/diivix_backend/internal/app/authentication/gapi/interceptor"
	"github.com/steve-mir/diivix_backend/internal/app/church/gapi"
	"github.com/steve-mir/diivix_backend/internal/app/church/pb"
	"github.com/steve-mir/diivix_backend/utils"
	"github.com/steve-mir/diivix_backend/worker"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func accessibleRoles() map[string][]int8 {
	return map[string][]int8{
		"/pb.ChurchService/GetDenominationList": {2, 9}, // -> admin and users

		// Only users
		"/pb.SocialMedia/UnFollowUser":                   {1},
		"/pb.ChurchService/ChangeDenominationMembership": {3, 4, 5, 6, 7, 8, 9},
		"/pb.ChurchService/ChangeChurchMembership":       {3, 4, 5, 6, 7, 8, 9},
		"/pb.ChurchService/SearchChurches":               {3, 4, 5, 6, 7, 8, 9},
		"/pb.ChurchService/SearchNearbyChurches":         {3, 4, 5, 6, 7, 8, 9},
		"/pb.ChurchService/GetChurchProfile":             {3, 4, 5, 6, 7, 8, 9},
		"/pb.ChurchService/GetUserChurch":                {3, 4, 5, 6, 7, 8, 9},
		"/pb.ChurchService/GetChurchMembers":             {3, 4, 5, 6, 7, 8, 9},
		"/pb.ChurchService/CreateProjectDonate":          {3, 4, 5, 6, 7, 8, 9},
		"/pb.ChurchService/GetProjectDetails":            {3, 4, 5, 6, 7, 8, 9},
		"/pb.ChurchService/GetProjectContributors":       {3, 4, 5, 6, 7, 8, 9},
		"/pb.ChurchService/GetChurchAnnouncements":       {3, 4, 5, 6, 7, 8, 9},
		"/pb.ChurchService/GetAnnouncementsForUser":      {3, 4, 5, 6, 7, 8, 9},

		// Only church admin
		"/pb.ChurchService/CreateProject":              {3, 4, 5, 6, 7, 8, 9},
		"/pb.ChurchService/CreateAnnouncement":         {3, 4, 5, 6, 7, 8, 9},
		"/pb.ChurchService/UpdateProject":              {3, 4, 5, 6, 7, 8, 9},
		"/pb.ChurchService/MarkProjectCompleted":       {3, 4, 5, 6, 7, 8, 9},
		"/pb.ChurchService/GetChurchProjects":          {3, 4, 5, 6, 7, 8, 9},
		"/pb.ChurchService/GetDonationAnalytics":       {3, 4, 5, 6, 7, 8, 9},
		"/pb.ChurchService/UpdateChurchAccountDetails": {3, 4, 5, 6, 7, 8, 9},
		"/pb.ChurchService/CreateProgram":              {3, 4, 5, 6, 7, 8, 9},
		"/pb.ChurchService/GetChurchPrograms":          {3, 4, 5, 6, 7, 8, 9},
		"/pb.ChurchService/EditChurchProgram":          {3, 4, 5, 6, 7, 8, 9},
		"/pb.ChurchService/DeleteChurchProgram":        {3, 4, 5, 6, 7, 8, 9},

		// Only admins
		"/pb.SocialMedia/BanUser":           {2}, // -> admin
		"/pb.ChurchService/AddDenomination": {2},
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
	runDbMigration(config.MigrationUrl, config.DBSource)

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

	td := worker.NewRedisTaskDistributor(redisOpt)

	runChurchServer(db, config, td)

}

func setupNewServer(config utils.Config) *grpc.Server {
	// Define Interceptors
	intercpts := interceptor.NewAuthInterceptor(config, accessibleRoles())

	// Chain the unary and stream interceptors
	chainedUnary := grpc_middleware.ChainUnaryServer(
		interceptor.GrpcLogger,
		intercpts.Unary(),
	)

	chainedStream := grpc_middleware.ChainStreamServer(
		interceptor.GrpcStreamLogger,
		intercpts.Stream(),
	)

	return grpc.NewServer(
		grpc.UnaryInterceptor(chainedUnary),
		grpc.StreamInterceptor(chainedStream),
	)
}

func runChurchServer(db *sql.DB, config utils.Config, td worker.TaskDistributor) {
	ctx, cancel := context.WithCancel(context.Background())

	// Handle graceful shutdown on receiving termination signals
	go func() {
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
		<-signals
		log.Info().Msg("Received termination signal. Shutting down gracefully...")
		cancel()
	}()

	// Church server
	server, err := gapi.NewChurchServer(db, config, td)
	if err != nil {
		log.Fatal().Msg("cannot create a server:")
	}

	grpcServer := setupNewServer(config)

	pb.RegisterChurchServiceServer(grpcServer, server)

	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GRPCChurchServerAddress)
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
	log.Info().Msg("church gRPC server stopped gracefully")
}

func runDbMigration(migrationUrl string, dbSource string) {
	migration, err := migrate.New(migrationUrl, dbSource)
	if err != nil {
		log.Fatal().Msg("cannot create new migration instance:") //, err)
	}

	if err = migration.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal().Msg("failed to run migrate up:") //, err)
	}
	log.Info().Msg("db migrated successfully")

}

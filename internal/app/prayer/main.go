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
	"github.com/steve-mir/diivix_backend/internal/app/prayer/gapi"
	"github.com/steve-mir/diivix_backend/internal/app/prayer/pb"
	"github.com/steve-mir/diivix_backend/utils"
	"github.com/steve-mir/diivix_backend/worker"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func accessibleRoles() map[string][]int8 {
	return map[string][]int8{
		// "/pb.ChurchService/GetDenominationList": {2, 1}, // -> admin and users

		// Only users
		"/pb.PrayerService/CreatePrayerRoom":   {9, 8, 7, 6, 5, 4, 3},
		"/pb.PrayerService/UpdatePrayerRoom":   {9, 8, 7, 6, 5, 4, 3},
		"/pb.PrayerService/DeletePrayerRoom":   {9, 8, 7, 6, 5, 4, 3},
		"/pb.PrayerService/InviteParticipant":  {9, 8, 7, 6, 5, 4, 3},
		"/pb.PrayerService/AcceptInvitation":   {9, 8, 7, 6, 5, 4, 3},
		"/pb.PrayerService/DeclineInvitation":  {9, 8, 7, 6, 5, 4, 3},
		"/pb.PrayerService/GetUserPrayerRooms": {9, 8, 7, 6, 5, 4, 3},
		"/pb.PrayerService/JoinPrayerRoom":     {9, 8, 7, 6, 5, 4, 3},

		// Only church admin

		// Only admins
		// "/pb.SocialMedia/BanUser":           {2}, // -> admin
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

	runPrayerServer(db, config, td)
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

func runPrayerServer(db *sql.DB, config utils.Config, td worker.TaskDistributor) {
	ctx, cancel := context.WithCancel(context.Background())

	// Handle graceful shutdown on receiving termination signals
	go func() {
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
		<-signals
		log.Info().Msg("Received termination signal. Shutting down gracefully...")
		cancel()
	}()

	// Prayer server
	server, err := gapi.NewPrayerServer(db, config, td)
	if err != nil {
		log.Fatal().Msg("cannot create a server:")
	}

	grpcServer := setupNewServer(config)

	pb.RegisterPrayerServiceServer(grpcServer, server)

	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GRPCPrayerSeverAddress)
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
	log.Info().Msg("prayer gRPC server stopped gracefully")
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

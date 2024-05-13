package main

import (
	"context"
	"database/sql"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/steve-mir/diivix_backend/constants"
	"github.com/steve-mir/diivix_backend/db/sqlc"
	"github.com/steve-mir/diivix_backend/internal/app/authentication/api"
	"github.com/steve-mir/diivix_backend/internal/app/authentication/gapi"
	"github.com/steve-mir/diivix_backend/internal/app/authentication/gapi/interceptor"
	"github.com/steve-mir/diivix_backend/internal/app/authentication/pb"
	"github.com/steve-mir/diivix_backend/utils"

	"github.com/steve-mir/diivix_backend/worker"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
)

// user - 1, admin - 2,
func accessibleRoles() map[string][]int8 {
	return map[string][]int8{
		// For All except guests
		"/pb.UserAuth/RegisterMFA":     {9, 8, 7, 6, 5, 4, 3, 2, 1},
		"/pb.UserAuth/VerifyMFAWorks":  {9, 8, 7, 6, 5, 4, 3, 2, 1},
		"/pb.UserAuth/VerifyMFA":       {9, 8, 7, 6, 5, 4, 3, 2, 1},
		"/pb.UserAuth/ByPassMFA":       {9, 8, 7, 6, 5, 4, 3, 2, 1},
		"/pb.UserAuth/UserSuggestions": {9, 8, 7, 6, 5, 4, 3, 2, 1},

		// For users
		"/pb.UserProfile/UploadImage":                                    {1},                         // -> only users
		"/pb.UserProfile/UpdateProfile":                                  {2, constants.RegularUsers}, // -> admin and users
		"/pb.UserAuth/ChangePassword":                                    {2, constants.RegularUsers}, // -> admin and users
		"/pb.UserAuth/RequestPasswordReset":                              {2, constants.RegularUsers}, // -> admin and users
		"/pb.UserAuth/RequestPassword":                                   {2, 1},                      // -> admin and users
		"/pb.UserAuth/ResetPassword":                                     {2, 1},
		"/pb.UserAuth/VerifyEmail":                                       {2, 1},
		"/grpc.reflection.v1alpha.ServerReflection/ServerReflectionInfo": {2, 1},
		"/pb.UserAuth/UpdateUserProfile":                                 {1, 2, 9},
		"/pb.UserAuth/DeleteAccount":                                     {9, 8, 7, 6, 5, 4, 3, 2, 1},
		"/pb.UserAuth/IncreaseTotalCoin":                                 {9, 8, 7, 6, 5, 4, 3, 2, 1},

		// Only users
		"/pb.UserAuth/InitiateChangeEmail": {constants.RegularUsers}, // -> users
		"/pb.UserAuth/ConfirmChangeEmail":  {constants.RegularUsers}, // -> users
		"/pb.UserAuth/InitiateChangePhone": {constants.RegularUsers}, // -> users
		"/pb.UserAuth/ConfirmChangePhone":  {constants.RegularUsers}, // -> users
		"/pb.UserAuth/ChangeUsername":      {constants.RegularUsers}, // -> users

		// Only admins

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

	taskDistributor := worker.NewRedisTaskDistributor(redisOpt)

	go runTaskProcessor(redisOpt, sqlc.NewStore(db), db) // TODO: Change to use db instead of store
	//**************** GRPC Server **********************/
	go runGrpcGatewayServer(db, config, taskDistributor)
	go ssoServer()
	runGrpcServer(db, config, taskDistributor)

}

func newGRPCServer(config utils.Config) *grpc.Server {
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

func runGrpcServer(db *sql.DB, config utils.Config, taskDistributor worker.TaskDistributor) {
	// Create a context that listens for termination signals
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
	server, err := gapi.NewServer(db, config, taskDistributor)
	if err != nil {
		log.Fatal().Msg("cannot create a server:")
	}

	// UserProfile
	profileServer, err := gapi.NewUserProfileServer(db, config)
	if err != nil {
		log.Fatal().Msg("cannot create a user profile server:")
	}

	// Social media server
	// socialMediaServer, err := gapi.NewSocialMediaServer(db, config)
	// if err != nil {
	// 	log.Fatal().Msg("cannot create a server:")
	// }

	grpcServer := newGRPCServer(config)

	pb.RegisterUserAuthServer(grpcServer, server)
	pb.RegisterUserProfileServer(grpcServer, profileServer)
	// pb.RegisterSocialMediaServiceServer(grpcServer, socialMediaServer)

	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GRPCServerAddress)
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

func runTaskProcessor(redisOpt asynq.RedisClientOpt, store *sqlc.Store, db *sql.DB) {
	mailer := worker.NewRedisTaskProcessor(redisOpt, store, db)
	log.Info().Msg("Starting task processor")
	err := mailer.Start()
	if err != nil {
		log.Fatal().Msg("cannot start task processor")
	}

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

func runGrpcGatewayServer(db *sql.DB, config utils.Config, taskDistributor worker.TaskDistributor) {
	server, err := gapi.NewServer(db, config, taskDistributor)
	if err != nil {
		log.Fatal().Msg("cannot create a server:") //, err)
	}

	jsonOptions := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	})

	grpcMux := runtime.NewServeMux(jsonOptions)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = pb.RegisterUserAuthHandlerServer(ctx, grpcMux, server)
	if err != nil {
		log.Fatal().Msg("cannot register handler server")
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	listener, err := net.Listen("tcp", config.HTTPServerAddress)
	if err != nil {
		log.Fatal().Msg("cannot create listener:") //, err)
	}

	log.Info().Msgf("start HTTP gateway server at %s", listener.Addr().String())
	handler := interceptor.HttpLogger(mux)
	err = http.Serve(listener, handler)
	if err != nil {
		log.Fatal().Msg("cannot start HTTP Gateway server")
	}
}

/*func createGinServer(db *sql.DB, config utils.Config, l *zap.Logger) *http.Server {
	port := config.HTTPServerAddress
	router := gin.New()

	setupRouter(db, config, router, l)

	return &http.Server{
		Addr:         port, //":" + port,
		Handler:      router,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}
}

func setupRouter(db *sql.DB, config utils.Config, route *gin.Engine, l *zap.Logger) {
	// Create db store and pass as injector
	store := sqlc.NewStore(db)
	// Create cors
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"https://localhost:3000"}
	route.Use(cors.New(corsConfig))

	// Use structured logger middleware
	route.Use(gin.Logger())
	routers.Auth(config, db, store, l, route)
	security.Security(config, db, store, l, route)
	profiles.Profile(config, store, l, route)
}*/

func ssoServer() {
	// fs := http.FileServer(http.Dir("public"))
	http.HandleFunc("/", handleMain)
	http.HandleFunc("/signin", api.Signin)
	http.HandleFunc("/callback", api.Callback)
	http.ListenAndServe(":7000", nil)
}

func handleMain(w http.ResponseWriter, r *http.Request) {
	var html = `<html><body><a href="/signin">Google Login</a></body></html>`
	fmt.Fprint(w, html)
}

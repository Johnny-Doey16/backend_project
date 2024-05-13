package gapi

import (
	"database/sql"

	"github.com/steve-mir/diivix_backend/cache"
	"github.com/steve-mir/diivix_backend/db/sqlc"
	"github.com/steve-mir/diivix_backend/internal/app/authentication/pb"
	"github.com/steve-mir/diivix_backend/internal/app/authentication/services"
	"github.com/steve-mir/diivix_backend/utils"
	"github.com/steve-mir/diivix_backend/worker"
)

type Server struct {
	pb.UnimplementedUserAuthServer
	config          utils.Config
	store           *sqlc.Store
	db              *sql.DB
	taskDistributor worker.TaskDistributor
	imageStore      services.ImageStore
	redisCache      cache.Cache
}

type ProfileServer struct {
	pb.UnimplementedUserProfileServer
	config     utils.Config
	store      *sqlc.Store
	db         *sql.DB
	imageStore services.ImageStore
}

// GRPC server
func NewServer(db *sql.DB, config utils.Config, taskDistributor worker.TaskDistributor) (*Server, error) {
	// Create db store and pass as injector
	return &Server{
		config:          config,
		db:              db,
		store:           sqlc.NewStore(db),
		taskDistributor: taskDistributor,
		imageStore:      services.NewDiskImageStore("img"),
		redisCache:      *cache.NewCache(config.RedisAddress, config.RedisUsername, config.RedisPwd, 0),
	}, nil
}

// Profile Server
func NewUserProfileServer(db *sql.DB, config utils.Config) (*ProfileServer, error) {
	// Create db store and pass as injector
	return &ProfileServer{
		config:     config,
		db:         db,
		store:      sqlc.NewStore(db),
		imageStore: services.NewDiskImageStore("img"),
	}, nil
}

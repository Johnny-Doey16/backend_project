package gapi

import (
	"database/sql"

	"github.com/steve-mir/diivix_backend/db/sqlc"
	"github.com/steve-mir/diivix_backend/internal/app/church/pb"
	"github.com/steve-mir/diivix_backend/utils"
	"github.com/steve-mir/diivix_backend/worker"
)

type ChurchServer struct {
	pb.UnimplementedChurchServiceServer
	config          utils.Config
	store           *sqlc.Store
	db              *sql.DB
	taskDistributor worker.TaskDistributor
}

func NewChurchServer(db *sql.DB, config utils.Config, taskDistributor worker.TaskDistributor) (*ChurchServer, error) {
	s := &ChurchServer{
		config:          config,
		db:              db,
		store:           sqlc.NewStore(db),
		taskDistributor: taskDistributor,
	}

	return s, nil
}

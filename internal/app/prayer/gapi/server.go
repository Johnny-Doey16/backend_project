package gapi

import (
	"database/sql"

	"github.com/steve-mir/diivix_backend/db/sqlc"
	"github.com/steve-mir/diivix_backend/internal/app/prayer/pb"
	"github.com/steve-mir/diivix_backend/utils"
	"github.com/steve-mir/diivix_backend/worker"
)

type PrayerServer struct {
	pb.UnimplementedPrayerServiceServer
	config          utils.Config
	store           *sqlc.Store
	db              *sql.DB
	taskDistributor worker.TaskDistributor
}

func NewPrayerServer(db *sql.DB, config utils.Config, taskDistributor worker.TaskDistributor) (*PrayerServer, error) {
	s := &PrayerServer{
		config:          config,
		db:              db,
		store:           sqlc.NewStore(db),
		taskDistributor: taskDistributor,
	}

	return s, nil
}

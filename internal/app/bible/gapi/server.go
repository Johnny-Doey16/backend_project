package gapi

import (
	"database/sql"

	"github.com/steve-mir/diivix_backend/db/sqlc"
	"github.com/steve-mir/diivix_backend/internal/app/bible/pb"
	"github.com/steve-mir/diivix_backend/utils"
)

type BibleServer struct {
	pb.UnimplementedBibleServiceServer
	config utils.Config
	store  *sqlc.Store
	db     *sql.DB
}

func NewBibleServer(db *sql.DB, config utils.Config) (*BibleServer, error) {
	s := &BibleServer{
		config: config,
		db:     db,
		store:  sqlc.NewStore(db),
	}

	return s, nil
}

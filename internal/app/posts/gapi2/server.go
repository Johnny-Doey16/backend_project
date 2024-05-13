package gapi2

/**
import (
	"database/sql"
	"sync"

	"github.com/steve-mir/diivix_backend/cache"
	"github.com/steve-mir/diivix_backend/db/sqlc"
	"github.com/steve-mir/diivix_backend/internal/app/authentication/services"
	"github.com/steve-mir/diivix_backend/internal/app/posts/pb"
	"github.com/steve-mir/diivix_backend/utils"
)

type SocialMediaServer struct {
	pb.UnimplementedSocialMediaServer
	config      utils.Config
	store       *sqlc.Store
	db          *sql.DB
	comments    []pb.Comment                      // ! start of removed
	mu          sync.Mutex                        // For thread-safety
	postStreams []pb.SocialMedia_PostStreamServer // * end
	// postCache   cache.PostCache // Add this line

	// posts         *sync.Map         // Use a concurrent map for posts
	subscribers *sync.Map // Map to track subscribers
	// signalSubscribe chan chan Signal  // Channel for subscribing clients to signals
	subscribe   chan chan pb.Post // Channel to subscribe new clients
	unsubscribe chan chan pb.Post // Channel to unsubscribe clients
	shutdown    chan struct{}     // Channel to shut down the server
	imageStore  services.ImageStore
	redisCache  cache.Cache
}

func NewSocialMediaServer(db *sql.DB, config utils.Config) (*SocialMediaServer, error) {

	// Create db store and pass as injector
	s := &SocialMediaServer{
		config:      config,
		db:          db,
		store:       sqlc.NewStore(db),
		comments:    []pb.Comment{},
		postStreams: []pb.SocialMedia_PostStreamServer{},
		mu:          sync.Mutex{},
		// postCache:   cache.NewRedisCache(config.RedisAddress, 0, 5*time.Second), // Add this line
		// posts:         &sync.Map{},
		subscribers: &sync.Map{},
		subscribe:   make(chan chan pb.Post),
		// signalSubscribe: make(chan chan Signal),
		unsubscribe: make(chan chan pb.Post),
		shutdown:    make(chan struct{}),
		imageStore:  services.NewDiskImageStore("img"),
		redisCache:  *cache.NewCache("localhost:6379", "", 0),
	}

	go s.run()
	return s, nil
}

func (s *SocialMediaServer) run() {
	for {
		select {
		case sub := <-s.subscribe:
			// Subscribe new client
			s.subscribers.Store(sub, true)
		case unsub := <-s.unsubscribe:
			// Unsubscribe client
			s.subscribers.Delete(unsub)
		case <-s.shutdown:
			return
		}
	}
}
**/

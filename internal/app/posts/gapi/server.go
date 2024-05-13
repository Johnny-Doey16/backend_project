package gapi

import (
	"context"
	"database/sql"
	"sync"

	"github.com/steve-mir/diivix_backend/cache"
	"github.com/steve-mir/diivix_backend/db/sqlc"
	"github.com/steve-mir/diivix_backend/internal/app/authentication/services"
	"github.com/steve-mir/diivix_backend/internal/app/posts/pb"
	"github.com/steve-mir/diivix_backend/utils"
	"github.com/steve-mir/diivix_backend/worker"
)

// TODO: reported_posts, moderation_queue. Read more on how reporting posts should be like.
const (
	PostCacheKey      = "posts"           // key for storing post to cache
	PostChannelKey    = "posts_channel"   // key for broadcasting to posts channel subscribers
	LikeChannelKey    = "likes_channel"   // key for broadcasting to like channel subscribers
	CommentChannelKey = "comment_channel" // key for broadcasting to comments channel subscribers
	RepostChannelKey  = "repost_channel"  // key for broadcasting to repost channel subscribers
	ViewChannelKey    = "views_channel"   // key for broadcasting to views channel subscribers
	PostCommentsKey   = "post_comments"   // key for storing a post's comment metrics (total comments)
	PostLikesKey      = "post_likes"      // key for storing a post's likes metrics (total likes)
	RepostsKey        = "repost"          // key for storing a post's reposts metrics (total reposts)
)

type SocialMediaServer struct {
	pb.UnimplementedSocialMediaServer
	config          utils.Config
	store           *sqlc.Store
	db              *sql.DB
	mu              sync.Mutex    // For thread-safety
	shutdown        chan struct{} // Channel to shut down the server
	imageStore      services.ImageStore
	redisCache      cache.Cache
	taskDistributor worker.TaskDistributor
}

func NewSocialMediaServer(db *sql.DB, config utils.Config, taskDistributor worker.TaskDistributor) (*SocialMediaServer, error) {
	s := &SocialMediaServer{
		config:          config,
		db:              db,
		store:           sqlc.NewStore(db),
		mu:              sync.Mutex{},
		shutdown:        make(chan struct{}),
		taskDistributor: taskDistributor,
		imageStore:      services.NewDiskImageStore("img"),
		redisCache:      *cache.NewCache(config.RedisAddress, config.RedisUsername, config.RedisPwd, 0),
	}

	go s.run()
	return s, nil
}

func (s *SocialMediaServer) run() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start a goroutine that listens for subscription messages
	go s.redisCache.SubscribeAndProcess(ctx, PostChannelKey)

	go s.redisCache.SubscribeAndProcess(ctx, CommentChannelKey)

	go s.redisCache.SubscribeAndProcess(ctx, RepostChannelKey)

	// Handle server shutdown logic here if necessary.
	<-s.shutdown
	cancel() // Cancel the context to shut down all goroutines
	// Perform any cleanup if needed.
}

func (s *SocialMediaServer) Shutdown() {
	close(s.shutdown) // Signal all goroutines to stop
}

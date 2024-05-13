package interceptor

import (
	"context"

	"golang.org/x/time/rate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type RateLimitInterceptor struct {
	limiter *rate.Limiter
}

// NewRateLimitInterceptor creates a new RateLimitInterceptor with a given rate and burst settings
func NewRateLimitInterceptor(r rate.Limit, b int) *RateLimitInterceptor {
	return &RateLimitInterceptor{
		limiter: rate.NewLimiter(r, b),
	}
}

// Unary returns a server interceptor function to enforce rate limiting
func (i *RateLimitInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		_ *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		if !i.limiter.Allow() {
			return nil, status.Errorf(codes.ResourceExhausted, "rate limit exceeded")
		}

		return handler(ctx, req)
	}
}

// Stream returns a server interceptor function to enforce rate limiting
func (i *RateLimitInterceptor) Stream() grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		ss grpc.ServerStream,
		_ *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		if !i.limiter.Allow() {
			return status.Errorf(codes.ResourceExhausted, "rate limit exceeded")
		}

		return handler(srv, ss)
	}
}

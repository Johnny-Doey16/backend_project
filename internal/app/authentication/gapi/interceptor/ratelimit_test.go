package interceptor

import (
	"context"
	"sync"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// dummyUnaryHandler is a placeholder for the real handler function. It does nothing in this test.
func dummyUnaryHandler(ctx context.Context, req interface{}) (interface{}, error) {
	return req, nil
}

// TestRateLimitInterceptor_Unary tests the Unary() interceptor method for rate limiting.
func TestRateLimitInterceptor_Unary(t *testing.T) {
	// Create an interceptor that allows 2 requests per second with a burst size of 1.
	limiter := NewRateLimitInterceptor(2, 1)

	// Define a dummy request and info.
	req := struct{}{}
	info := &grpc.UnaryServerInfo{}

	// Call the interceptor manually and expect it to succeed.
	_, err := limiter.Unary()(context.Background(), req, info, dummyUnaryHandler)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Call the interceptor again quickly and expect it to fail due to the rate limit.
	_, err = limiter.Unary()(context.Background(), req, info, dummyUnaryHandler)
	if err == nil {
		t.Errorf("Expected rate limiting error, got none")
	}

	// Check if the error is the correct type and code.
	st, ok := status.FromError(err)
	if !ok || st.Code() != codes.ResourceExhausted {
		t.Errorf("Expected ResourceExhausted error, got %v", err)
	}

	// Wait for longer than the rate limit refresh interval and try again, expecting success.
	time.Sleep(time.Second / 2)
	_, err = limiter.Unary()(context.Background(), req, info, dummyUnaryHandler)
	if err != nil {
		t.Errorf("Unexpected error after waiting: %v", err)
	}
}

// TestRateLimitInterceptor_ExceedLimit tests that the rate limit interceptor blocks calls that exceed limits
func TestRateLimitInterceptor_ExceedLimit(t *testing.T) {
	// Create an interceptor that allows 10 requests per second with a burst size of 2.
	limiter := NewRateLimitInterceptor(10, 2)

	// Define a dummy request and info.
	req := struct{}{}
	info := &grpc.UnaryServerInfo{}

	// Wait group to wait for all goroutines to finish.
	var wg sync.WaitGroup

	// Counter for the number of calls that get blocked.
	var blockedCalls int32

	// Number of calls to simulate.
	totalCalls := 50

	// Simulate concurrent calls to the interceptor.
	for i := 0; i < totalCalls; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := limiter.Unary()(context.Background(), req, info, dummyUnaryHandler)
			if err != nil {
				// Check if the error is due to rate limiting.
				st, ok := status.FromError(err)
				if ok && st.Code() == codes.ResourceExhausted {
					// Increment the counter for blocked calls.
					blockedCalls++
				}
			}
		}()
	}

	// Wait for all calls to finish.
	wg.Wait()

	// We allow 10 requests per second with a burst size of 2, so we expect at least
	// totalCalls - (rate + burst size) calls to be blocked.
	expectedBlocked := totalCalls - (10 + 2)
	if int(blockedCalls) < expectedBlocked {
		t.Errorf("Rate limit did not block enough calls, only blocked %d out of %d, expected at least %d", blockedCalls, totalCalls, expectedBlocked)
	}

	// To observe the rate limiting in action, we could log the result or make an assertion here.
	t.Logf("Total blocked calls due to rate limiting: %d", blockedCalls)
}

// In a real test, you'd also want to test the Stream() interceptor method similarly.

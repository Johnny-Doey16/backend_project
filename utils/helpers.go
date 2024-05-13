package utils

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"math/big"
	"net"
	"runtime"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/sqlc-dev/pqtype"
)

// getIpAddr returns a pqtype.Inet representation of the client's IP address.
//
// It takes a string parameter, clientIP, which represents the client's IP address.
// It returns a pqtype.Inet value.
func GetIpAddr(clientIP string) pqtype.Inet {
	if clientIP == "::1" {
		clientIP = "127.0.0.1"
	}

	ip := net.ParseIP(clientIP)

	if ip == nil {
		// Handle the case where ctx.ClientIP() doesn't return a valid IP address
		return pqtype.Inet{}
	}

	inet := pqtype.Inet{
		IPNet: net.IPNet{
			IP:   ip,
			Mask: net.CIDRMask(32, 32), // If you're dealing with IPv4 addresses
		},
		Valid: true,
	}
	return inet
}

func GetKeyForToken(config Config, isRefresh bool) string {
	var key string
	if isRefresh {
		key = config.RefreshTokenSymmetricKey
	} else {
		key = config.AccessTokenSymmetricKey
	}

	return key
}

// GenerateUniqueToken generates a unique verification token.
func GenerateUniqueToken(len int) (string, error) {
	// Generate a cryptographically secure random value
	randomBytes := make([]byte, len)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}

	// Create a unique token by combining user ID, timestamp, and random value
	// timestamp := time.Now().Unix()
	// token := fmt.Sprintf("%s-%d-%s", userID, timestamp, formatConsistentToken(timestamp, base64.URLEncoding.EncodeToString(randomBytes)))
	token := base64.URLEncoding.EncodeToString(randomBytes)

	return token, nil
}

func GenerateSecureRandomNumber(max int64) (int64, error) {
	nBig, err := rand.Int(rand.Reader, big.NewInt(max))
	if err != nil {
		return 0, err
	}
	return nBig.Int64(), nil
}

func StartMemCal() time.Time {
	return time.Now()
}

func EndMemCal(startTime time.Time) {
	// Measure execution time.
	executionTime := time.Since(startTime).Seconds()

	// Get memory statistics.
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	peakMemoryBytes := memStats.Sys

	// Convert peak memory usage to GB.
	peakMemoryGB := float64(peakMemoryBytes) / (1024 * 1024 * 1024)

	// Calculate GB-seconds.
	gbSeconds := peakMemoryGB * executionTime

	fmt.Printf("Peak Memory Usage: %.6f GB\n", peakMemoryGB)
	fmt.Printf("Execution Time: %.6f seconds\n", executionTime)
	fmt.Printf("GB-Seconds: %.6f\n", gbSeconds)

}

func StartMemCalOld() (time.Time, uint64, runtime.MemStats) {
	startTime := time.Now()
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	startMemory := memStats.Alloc

	return startTime, startMemory, memStats
}

func EndMemCalOld(startTime time.Time, startMemory uint64, memStats *runtime.MemStats) {
	endTime := time.Now()
	runtime.ReadMemStats(memStats)
	endMemory := memStats.Alloc

	executionTime := endTime.Sub(startTime).Seconds()
	memoryUsed := (endMemory - startMemory) / (1024 * 1024 * 1024)

	gbSeconds := memoryUsed * uint64(executionTime)
	log.Info().Msgf("GB-seconds: %v", gbSeconds)

}

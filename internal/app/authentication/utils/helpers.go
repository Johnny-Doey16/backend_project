package utils

import (
	"crypto/rand"
	"encoding/base64"
	"math/big"
	"net"

	"github.com/sqlc-dev/pqtype"
	"github.com/steve-mir/diivix_backend/utils"
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

func GetKeyForToken(config utils.Config, isRefresh bool) string {
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

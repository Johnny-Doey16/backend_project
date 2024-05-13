package utils

import (
	"context"
	"net"

	"github.com/sqlc-dev/pqtype"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

const (
	grpcGatewayUserAgentHeader = "grpcgateway-user-agent"
	userAgentHeader            = "user-agent"
	xForwardedForHeader        = "x-forwarded-for"
)

type Metadata struct {
	UserAgent string
	ClientIP  string
}

func ExtractMetadata(ctx context.Context) *Metadata {
	metaData := &Metadata{}

	if md, ok := metadata.FromIncomingContext(ctx); ok {

		// for grpc gateway
		if userAgents := md.Get(grpcGatewayUserAgentHeader); len(userAgents) > 0 {
			metaData.UserAgent = userAgents[0]
		}

		// for grpc
		if userAgents := md.Get(userAgentHeader); len(userAgents) > 0 {
			metaData.UserAgent = userAgents[0]
		}

		if clientIps := md.Get(xForwardedForHeader); len(clientIps) > 0 {
			metaData.ClientIP = clientIps[0]
		}
	}

	// for grpc
	if peer, ok := peer.FromContext(ctx); ok {
		host, _, _ := net.SplitHostPort(peer.Addr.String())
		if host == "::1" {
			host = "127.0.0.1"
		}
		metaData.ClientIP = host
	}

	return metaData
}

func GetMetadata(ctx context.Context) (string, pqtype.Inet) {
	agent := ExtractMetadata(ctx).UserAgent
	ip := ExtractMetadata(ctx).ClientIP
	clientIP := GetIpAddr(ip)

	return agent, clientIP
}

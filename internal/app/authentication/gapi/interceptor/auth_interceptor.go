package interceptor

import (
	"context"
	"strings"

	"github.com/steve-mir/diivix_backend/internal/app/authentication/token"
	"github.com/steve-mir/diivix_backend/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const (
	authorizationHeader = "authorization"
	authorizationBearer = "bearer"
)

type AuthContextKey struct {
	Name string
}

// Create a new struct that embeds grpc.ServerStream
type wrappedStream struct {
	grpc.ServerStream
	ctx context.Context
}

// Override the Context method to return the new context
func (w *wrappedStream) Context() context.Context {
	return w.ctx
}

type AuthInterceptor struct {
	tokenMaker      token.Maker
	accessibleRoles map[string][]int8
}

func NewAuthInterceptor(config utils.Config, accessibleRoles map[string][]int8) *AuthInterceptor {
	tokenMaker, err := token.NewPasetoMaker(config.AccessTokenSymmetricKey)
	if err != nil {
		return nil
	}

	return &AuthInterceptor{
		tokenMaker:      tokenMaker,
		accessibleRoles: accessibleRoles}
}

func (interceptor *AuthInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		// Skip authentication for specific methods (e.g., login and register)
		if shouldSkipAuthentication(info.FullMethod) {
			return handler(ctx, req)
		}

		ctx, err := interceptor.authorize(ctx, info.FullMethod)
		if err != nil {
			return nil, err
		}

		return handler(ctx, req)
	}
}

// Stream returns a server interceptor function to authenticate and authorize stream RPC
func (interceptor *AuthInterceptor) Stream() grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		stream grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		// Skip authentication for specific methods (e.g., login and register)
		if shouldSkipAuthentication(info.FullMethod) {
			return handler(srv, stream)
		}

		ctx, err := interceptor.authorize(stream.Context(), info.FullMethod)
		if err != nil {
			return err
		}

		// Create a new wrappedStream with the new context
		wrapped := &wrappedStream{ServerStream: stream, ctx: ctx}

		return handler(srv, wrapped)
	}
}

func (interceptor *AuthInterceptor) authorize(ctx context.Context, method string) (context.Context, error) {
	accessibleRoles, ok := interceptor.accessibleRoles[method]
	if !ok {
		// everyone can access
		return ctx, status.Error(codes.PermissionDenied, "sorry you do not have enough permissions to access: "+method)
	}

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ctx, status.Error(codes.InvalidArgument, "missing metadata")
	}

	accessToken, err := interceptor.extractToken(md)
	if err != nil {
		return ctx, status.Error(codes.Unknown, err.Error())
	}

	payload, err := interceptor.tokenMaker.VerifyToken(accessToken)
	if err != nil {
		return ctx, status.Error(codes.PermissionDenied, "invalid token: "+err.Error())
	}

	if !payload.IsEmailVerified {
		return ctx, status.Error(codes.PermissionDenied, "please verify your email")
	}

	if !payload.MfaPassed && (method != "/pb.UserAuth/RegisterMFA" && method != "/pb.UserAuth/VerifyMFA" && method != "/pb.UserAuth/ByPassMFA") {
		return ctx, status.Error(codes.PermissionDenied, "Multi factor Authentication failed")
	}

	// Storing the uid in the context within the interceptor
	ctx2 := context.WithValue(ctx, "payloadKey", payload)

	for _, role := range accessibleRoles {
		// log.Println("Role", role)
		if role == payload.Role {
			return ctx2, nil
		}
	}

	return ctx, status.Error(codes.PermissionDenied, "no permission to access this RPC")
}

func (interceptor *AuthInterceptor) extractToken(md metadata.MD) (string, error) {

	values := md.Get(authorizationHeader)
	if len(values) == 0 {
		return "", status.Error(codes.InvalidArgument, "missing authorization header")
	}

	authHeader := values[0]
	fields := strings.Fields(authHeader)
	if len(fields) < 2 {
		return "", status.Error(codes.InvalidArgument, "invalid authorization header format")
	}

	authorizationType := strings.ToLower(fields[0])
	if authorizationType != authorizationBearer {
		return "", status.Error(codes.PermissionDenied, "unsupported authorization type: "+authorizationType)
	}
	return fields[1], nil
}

func shouldSkipAuthentication(method string) bool {
	// Add logic to determine whether authentication should be skipped for the given method
	return method == "/pb.UserAuth/LoginUser" || method == "/pb.UserAuth/CreateUser" || method == "/pb.UserAuth/RotateToken" ||
		method == "/pb.UserAuth/RequestPasswordReset" || method == "/pb.UserAuth/RequestAccountRecovery" ||
		method == "/pb.UserAuth/ResetPassword" || method == "/pb.UserAuth/CompleteAccountRecovery" ||

		method == "/pb.UserAuth/RegisterWithSocial" || method == "/pb.UserAuth/LoginWithSocial" ||
		method == "/pb.UserAuth/VerifyEmail" || method == "/pb.ChurchService/CreateChurch" ||

		method == "/grpc.reflection.v1alpha.ServerReflection/ServerReflectionInfo" ||

		// TODO: Remove when do testing Social media
		method == "/pb.SocialMedia/LiveComments" ||

		method == "/pb.SocialMedia/GetFollowers" ||
		method == "/pb.SocialMedia/GetFollowing" || method == "/pb.SocialMedia/GetPostsByUserId" ||
		method == "/pb.SocialMedia/ViewPost" ||

		method == "/pb.SocialMedia/CreateRepost" || method == "/pb.SocialMedia/GetRepost" || method == "/pb.SocialMedia/GetRepostsByUser" ||
		method == "/pb.SocialMedia/DeleteRepost" || method == "/pb.SocialMedia/BookmarkPost" || method == "/pb.SocialMedia/GetBookmarkedPosts" || method == "/pb.SocialMedia/DeleteBookmark" ||

		method == "/pb.SocialMedia/BlockUser" ||
		// method == "/pb.SocialMedia/PostStream" ||
		method == "/pb.SocialMedia/GetBlockedUsers" || method == "/pb.SocialMedia/UnblockUser" ||
		method == "/pb.ChurchService/CreateProject" || method == "/pb.ChurchService/UpdateProject" || method == "/pb.ChurchService/MarkProjectCompleted" ||
		method == "/pb.ChurchService/GetChurchProjects" || method == "/pb.ChurchService/CreateProjectDonate" || method == "/pb.ChurchService/GetProjectDetails" ||
		method == "/pb.ChurchService/GetDonationAnalytics" || method == "/pb.ChurchService/GetProjectContributors" || method == "/pb.ChurchService/CreateAnnouncement" ||

		method == "/pb.ChurchService/GetChurchAnnouncements" || method == "/pb.ChurchService/GetAnnouncementsForUser" ||
		method == "/pb.UserAuth/UserSuggestions" || method == "/pb.BibleService/DownloadBible" ||

		method == "/pb.ChurchService/CreateProgram" || method == "/pb.ChurchService/GetChurchPrograms" || method == "/pb.ChurchService/EditChurchProgram" || method == "/pb.ChurchService/DeleteChurchProgram"

	// method == "/pb.SocialMedia/GetNotifications"

	//  || method == "/pb.SocialMedia/GenericSearch" || method == "/pb.SocialMedia/GenericSearch2"

	// /grpc.reflection.v1alpha.ServerReflection/ServerReflectionInfo

}

package gapi

import (
	"context"
	"fmt"
	"reflect"

	"github.com/steve-mir/diivix_backend/db/sqlc"
	"github.com/steve-mir/diivix_backend/internal/app/authentication/token"
	"github.com/steve-mir/diivix_backend/internal/app/posts/pb"
	"github.com/steve-mir/diivix_backend/internal/app/posts/services"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// ! Get user's notifications RPC
func (s *SocialMediaServer) GetNotifications(req *pb.GetNotificationRequest, stream pb.SocialMedia_GetNotificationsServer) error {
	claims, ok := stream.Context().Value("payloadKey").(*token.Payload)
	if !ok {
		return status.Error(codes.Unauthenticated, "unable to retrieve UID from context")
	}
	// userId, err := services.StrToUUID("558c60aa-977f-4d38-885b-e813656371ac")
	// if err != nil {
	// 	return status.Errorf(codes.Internal, "parsing uid: %s", err)
	// }

	notifications, err := s.store.GetUserNotifications(stream.Context(), sqlc.GetUserNotificationsParams{
		UserID: claims.UserId,
		Type:   req.GetType(),
	})
	if err != nil {
		return status.Errorf(codes.Internal, "cannot get notification data: %s", err)
	}

	fmt.Println("Notification Content type: ", reflect.TypeOf(notifications[0].PostAnnouncementContent))
	for _, notification := range notifications {
		var content *string

		s := string(notification.InvitationStatus.ParticipantStatus)
		stream.Send(
			&pb.Notification{
				Type:                    notification.Type,
				NotificationId:          notification.NotificationID,
				Read:                    notification.Read,
				UserType:                notification.UserType.String,
				Id:                      notification.ID.UUID.String(),
				AuthorUsername:          notification.AuthorUsername.String,
				AuthorProfileImage:      notification.AuthorProfileImage.String,
				PostAnnouncementId:      notification.PostAnnouncementID,
				Time:                    timestamppb.New(notification.Time.Time),
				StartTime:               timestamppb.New(notification.StartTime.Time),
				EndTime:                 timestamppb.New(notification.EndTime.Time),
				MeetingTitle:            services.InterfaceToStr(notification.PostAnnouncementContent, content),
				MeetingId:               &notification.RoomID.String,
				MeetingInviteStat:       &s,
				PostAnnouncementContent: *services.InterfaceToStr(notification.PostAnnouncementContent, content),
			},
		)
	}
	return nil
}

// ! Mark notification as read RPC
func (s *SocialMediaServer) MarkNotificationAsRead(ctx context.Context, req *pb.NotificationId) (*pb.MarkNotificationAsReadResponse, error) {
	err := s.store.MarkNotificationAsRead(ctx, req.GetNotificationId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error updating notification status: %s", err)
	}

	return &pb.MarkNotificationAsReadResponse{
		Msg:  "Notification read",
		Read: true,
	}, nil
}

// ! Delete notification RPC
func (s *SocialMediaServer) DeleteNotifications(ctx context.Context, req *pb.NotificationIds) (*pb.DeleteNotificationResponse, error) {
	err := s.store.DeleteNotifications(ctx, req.GetNotificationIds())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error deleting notifications: %s", err)
	}

	return &pb.DeleteNotificationResponse{
		Msg:     "Notifications deleted",
		Deleted: true,
	}, nil
}

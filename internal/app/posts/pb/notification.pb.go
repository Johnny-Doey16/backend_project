// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v3.19.6
// source: notification.proto

package pb

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type GetNotificationRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Type string `protobuf:"bytes,1,opt,name=type,proto3" json:"type,omitempty"`
}

func (x *GetNotificationRequest) Reset() {
	*x = GetNotificationRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_notification_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetNotificationRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetNotificationRequest) ProtoMessage() {}

func (x *GetNotificationRequest) ProtoReflect() protoreflect.Message {
	mi := &file_notification_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetNotificationRequest.ProtoReflect.Descriptor instead.
func (*GetNotificationRequest) Descriptor() ([]byte, []int) {
	return file_notification_proto_rawDescGZIP(), []int{0}
}

func (x *GetNotificationRequest) GetType() string {
	if x != nil {
		return x.Type
	}
	return ""
}

type Notification struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	NotificationId          int32                  `protobuf:"varint,1,opt,name=notification_id,json=notificationId,proto3" json:"notification_id,omitempty"`
	Type                    string                 `protobuf:"bytes,2,opt,name=type,proto3" json:"type,omitempty"`
	Read                    bool                   `protobuf:"varint,3,opt,name=read,proto3" json:"read,omitempty"`
	Time                    *timestamppb.Timestamp `protobuf:"bytes,4,opt,name=time,proto3" json:"time,omitempty"`
	Id                      string                 `protobuf:"bytes,5,opt,name=id,proto3" json:"id,omitempty"`
	AuthorUsername          string                 `protobuf:"bytes,6,opt,name=author_username,json=authorUsername,proto3" json:"author_username,omitempty"`
	AuthorProfileImage      string                 `protobuf:"bytes,7,opt,name=author_profile_image,json=authorProfileImage,proto3" json:"author_profile_image,omitempty"`
	PostAnnouncementContent string                 `protobuf:"bytes,8,opt,name=post_announcement_content,json=postAnnouncementContent,proto3" json:"post_announcement_content,omitempty"`
	PostAnnouncementId      string                 `protobuf:"bytes,9,opt,name=post_announcement_id,json=postAnnouncementId,proto3" json:"post_announcement_id,omitempty"`
	UserType                string                 `protobuf:"bytes,10,opt,name=user_type,json=userType,proto3" json:"user_type,omitempty"`
	MeetingInviteStat       *string                `protobuf:"bytes,11,opt,name=meeting_invite_stat,json=meetingInviteStat,proto3,oneof" json:"meeting_invite_stat,omitempty"`
	MeetingTitle            *string                `protobuf:"bytes,12,opt,name=meeting_title,json=meetingTitle,proto3,oneof" json:"meeting_title,omitempty"`
	MeetingId               *string                `protobuf:"bytes,13,opt,name=meeting_id,json=meetingId,proto3,oneof" json:"meeting_id,omitempty"`
	StartTime               *timestamppb.Timestamp `protobuf:"bytes,14,opt,name=start_time,json=startTime,proto3,oneof" json:"start_time,omitempty"`
	EndTime                 *timestamppb.Timestamp `protobuf:"bytes,15,opt,name=end_time,json=endTime,proto3,oneof" json:"end_time,omitempty"`
}

func (x *Notification) Reset() {
	*x = Notification{}
	if protoimpl.UnsafeEnabled {
		mi := &file_notification_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Notification) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Notification) ProtoMessage() {}

func (x *Notification) ProtoReflect() protoreflect.Message {
	mi := &file_notification_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Notification.ProtoReflect.Descriptor instead.
func (*Notification) Descriptor() ([]byte, []int) {
	return file_notification_proto_rawDescGZIP(), []int{1}
}

func (x *Notification) GetNotificationId() int32 {
	if x != nil {
		return x.NotificationId
	}
	return 0
}

func (x *Notification) GetType() string {
	if x != nil {
		return x.Type
	}
	return ""
}

func (x *Notification) GetRead() bool {
	if x != nil {
		return x.Read
	}
	return false
}

func (x *Notification) GetTime() *timestamppb.Timestamp {
	if x != nil {
		return x.Time
	}
	return nil
}

func (x *Notification) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Notification) GetAuthorUsername() string {
	if x != nil {
		return x.AuthorUsername
	}
	return ""
}

func (x *Notification) GetAuthorProfileImage() string {
	if x != nil {
		return x.AuthorProfileImage
	}
	return ""
}

func (x *Notification) GetPostAnnouncementContent() string {
	if x != nil {
		return x.PostAnnouncementContent
	}
	return ""
}

func (x *Notification) GetPostAnnouncementId() string {
	if x != nil {
		return x.PostAnnouncementId
	}
	return ""
}

func (x *Notification) GetUserType() string {
	if x != nil {
		return x.UserType
	}
	return ""
}

func (x *Notification) GetMeetingInviteStat() string {
	if x != nil && x.MeetingInviteStat != nil {
		return *x.MeetingInviteStat
	}
	return ""
}

func (x *Notification) GetMeetingTitle() string {
	if x != nil && x.MeetingTitle != nil {
		return *x.MeetingTitle
	}
	return ""
}

func (x *Notification) GetMeetingId() string {
	if x != nil && x.MeetingId != nil {
		return *x.MeetingId
	}
	return ""
}

func (x *Notification) GetStartTime() *timestamppb.Timestamp {
	if x != nil {
		return x.StartTime
	}
	return nil
}

func (x *Notification) GetEndTime() *timestamppb.Timestamp {
	if x != nil {
		return x.EndTime
	}
	return nil
}

type NotificationId struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	NotificationId int32 `protobuf:"varint,1,opt,name=notification_id,json=notificationId,proto3" json:"notification_id,omitempty"`
}

func (x *NotificationId) Reset() {
	*x = NotificationId{}
	if protoimpl.UnsafeEnabled {
		mi := &file_notification_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NotificationId) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NotificationId) ProtoMessage() {}

func (x *NotificationId) ProtoReflect() protoreflect.Message {
	mi := &file_notification_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NotificationId.ProtoReflect.Descriptor instead.
func (*NotificationId) Descriptor() ([]byte, []int) {
	return file_notification_proto_rawDescGZIP(), []int{2}
}

func (x *NotificationId) GetNotificationId() int32 {
	if x != nil {
		return x.NotificationId
	}
	return 0
}

type NotificationIds struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	NotificationIds []int32 `protobuf:"varint,1,rep,packed,name=notification_ids,json=notificationIds,proto3" json:"notification_ids,omitempty"`
}

func (x *NotificationIds) Reset() {
	*x = NotificationIds{}
	if protoimpl.UnsafeEnabled {
		mi := &file_notification_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NotificationIds) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NotificationIds) ProtoMessage() {}

func (x *NotificationIds) ProtoReflect() protoreflect.Message {
	mi := &file_notification_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NotificationIds.ProtoReflect.Descriptor instead.
func (*NotificationIds) Descriptor() ([]byte, []int) {
	return file_notification_proto_rawDescGZIP(), []int{3}
}

func (x *NotificationIds) GetNotificationIds() []int32 {
	if x != nil {
		return x.NotificationIds
	}
	return nil
}

type MarkNotificationAsReadResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Msg  string `protobuf:"bytes,1,opt,name=msg,proto3" json:"msg,omitempty"`
	Read bool   `protobuf:"varint,2,opt,name=read,proto3" json:"read,omitempty"`
}

func (x *MarkNotificationAsReadResponse) Reset() {
	*x = MarkNotificationAsReadResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_notification_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MarkNotificationAsReadResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MarkNotificationAsReadResponse) ProtoMessage() {}

func (x *MarkNotificationAsReadResponse) ProtoReflect() protoreflect.Message {
	mi := &file_notification_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MarkNotificationAsReadResponse.ProtoReflect.Descriptor instead.
func (*MarkNotificationAsReadResponse) Descriptor() ([]byte, []int) {
	return file_notification_proto_rawDescGZIP(), []int{4}
}

func (x *MarkNotificationAsReadResponse) GetMsg() string {
	if x != nil {
		return x.Msg
	}
	return ""
}

func (x *MarkNotificationAsReadResponse) GetRead() bool {
	if x != nil {
		return x.Read
	}
	return false
}

type DeleteNotificationResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Msg     string `protobuf:"bytes,1,opt,name=msg,proto3" json:"msg,omitempty"`
	Deleted bool   `protobuf:"varint,2,opt,name=deleted,proto3" json:"deleted,omitempty"`
}

func (x *DeleteNotificationResponse) Reset() {
	*x = DeleteNotificationResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_notification_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeleteNotificationResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteNotificationResponse) ProtoMessage() {}

func (x *DeleteNotificationResponse) ProtoReflect() protoreflect.Message {
	mi := &file_notification_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteNotificationResponse.ProtoReflect.Descriptor instead.
func (*DeleteNotificationResponse) Descriptor() ([]byte, []int) {
	return file_notification_proto_rawDescGZIP(), []int{5}
}

func (x *DeleteNotificationResponse) GetMsg() string {
	if x != nil {
		return x.Msg
	}
	return ""
}

func (x *DeleteNotificationResponse) GetDeleted() bool {
	if x != nil {
		return x.Deleted
	}
	return false
}

var File_notification_proto protoreflect.FileDescriptor

var file_notification_proto_rawDesc = []byte{
	0x0a, 0x12, 0x6e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x12, 0x02, 0x70, 0x62, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74,
	0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x2c, 0x0a, 0x16, 0x47, 0x65, 0x74,
	0x4e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x22, 0xd9, 0x05, 0x0a, 0x0c, 0x4e, 0x6f, 0x74, 0x69,
	0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x27, 0x0a, 0x0f, 0x6e, 0x6f, 0x74, 0x69,
	0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x05, 0x52, 0x0e, 0x6e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x49,
	0x64, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x72, 0x65, 0x61, 0x64, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x08, 0x52, 0x04, 0x72, 0x65, 0x61, 0x64, 0x12, 0x2e, 0x0a, 0x04, 0x74, 0x69, 0x6d,
	0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74,
	0x61, 0x6d, 0x70, 0x52, 0x04, 0x74, 0x69, 0x6d, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18,
	0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x27, 0x0a, 0x0f, 0x61, 0x75, 0x74,
	0x68, 0x6f, 0x72, 0x5f, 0x75, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x06, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0e, 0x61, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x55, 0x73, 0x65, 0x72, 0x6e, 0x61,
	0x6d, 0x65, 0x12, 0x30, 0x0a, 0x14, 0x61, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x5f, 0x70, 0x72, 0x6f,
	0x66, 0x69, 0x6c, 0x65, 0x5f, 0x69, 0x6d, 0x61, 0x67, 0x65, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x12, 0x61, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x50, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x49,
	0x6d, 0x61, 0x67, 0x65, 0x12, 0x3a, 0x0a, 0x19, 0x70, 0x6f, 0x73, 0x74, 0x5f, 0x61, 0x6e, 0x6e,
	0x6f, 0x75, 0x6e, 0x63, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x5f, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e,
	0x74, 0x18, 0x08, 0x20, 0x01, 0x28, 0x09, 0x52, 0x17, 0x70, 0x6f, 0x73, 0x74, 0x41, 0x6e, 0x6e,
	0x6f, 0x75, 0x6e, 0x63, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74,
	0x12, 0x30, 0x0a, 0x14, 0x70, 0x6f, 0x73, 0x74, 0x5f, 0x61, 0x6e, 0x6e, 0x6f, 0x75, 0x6e, 0x63,
	0x65, 0x6d, 0x65, 0x6e, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x09, 0x20, 0x01, 0x28, 0x09, 0x52, 0x12,
	0x70, 0x6f, 0x73, 0x74, 0x41, 0x6e, 0x6e, 0x6f, 0x75, 0x6e, 0x63, 0x65, 0x6d, 0x65, 0x6e, 0x74,
	0x49, 0x64, 0x12, 0x1b, 0x0a, 0x09, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18,
	0x0a, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x75, 0x73, 0x65, 0x72, 0x54, 0x79, 0x70, 0x65, 0x12,
	0x33, 0x0a, 0x13, 0x6d, 0x65, 0x65, 0x74, 0x69, 0x6e, 0x67, 0x5f, 0x69, 0x6e, 0x76, 0x69, 0x74,
	0x65, 0x5f, 0x73, 0x74, 0x61, 0x74, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x09, 0x48, 0x00, 0x52, 0x11,
	0x6d, 0x65, 0x65, 0x74, 0x69, 0x6e, 0x67, 0x49, 0x6e, 0x76, 0x69, 0x74, 0x65, 0x53, 0x74, 0x61,
	0x74, 0x88, 0x01, 0x01, 0x12, 0x28, 0x0a, 0x0d, 0x6d, 0x65, 0x65, 0x74, 0x69, 0x6e, 0x67, 0x5f,
	0x74, 0x69, 0x74, 0x6c, 0x65, 0x18, 0x0c, 0x20, 0x01, 0x28, 0x09, 0x48, 0x01, 0x52, 0x0c, 0x6d,
	0x65, 0x65, 0x74, 0x69, 0x6e, 0x67, 0x54, 0x69, 0x74, 0x6c, 0x65, 0x88, 0x01, 0x01, 0x12, 0x22,
	0x0a, 0x0a, 0x6d, 0x65, 0x65, 0x74, 0x69, 0x6e, 0x67, 0x5f, 0x69, 0x64, 0x18, 0x0d, 0x20, 0x01,
	0x28, 0x09, 0x48, 0x02, 0x52, 0x09, 0x6d, 0x65, 0x65, 0x74, 0x69, 0x6e, 0x67, 0x49, 0x64, 0x88,
	0x01, 0x01, 0x12, 0x3e, 0x0a, 0x0a, 0x73, 0x74, 0x61, 0x72, 0x74, 0x5f, 0x74, 0x69, 0x6d, 0x65,
	0x18, 0x0e, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61,
	0x6d, 0x70, 0x48, 0x03, 0x52, 0x09, 0x73, 0x74, 0x61, 0x72, 0x74, 0x54, 0x69, 0x6d, 0x65, 0x88,
	0x01, 0x01, 0x12, 0x3a, 0x0a, 0x08, 0x65, 0x6e, 0x64, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x0f,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70,
	0x48, 0x04, 0x52, 0x07, 0x65, 0x6e, 0x64, 0x54, 0x69, 0x6d, 0x65, 0x88, 0x01, 0x01, 0x42, 0x16,
	0x0a, 0x14, 0x5f, 0x6d, 0x65, 0x65, 0x74, 0x69, 0x6e, 0x67, 0x5f, 0x69, 0x6e, 0x76, 0x69, 0x74,
	0x65, 0x5f, 0x73, 0x74, 0x61, 0x74, 0x42, 0x10, 0x0a, 0x0e, 0x5f, 0x6d, 0x65, 0x65, 0x74, 0x69,
	0x6e, 0x67, 0x5f, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x42, 0x0d, 0x0a, 0x0b, 0x5f, 0x6d, 0x65, 0x65,
	0x74, 0x69, 0x6e, 0x67, 0x5f, 0x69, 0x64, 0x42, 0x0d, 0x0a, 0x0b, 0x5f, 0x73, 0x74, 0x61, 0x72,
	0x74, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x42, 0x0b, 0x0a, 0x09, 0x5f, 0x65, 0x6e, 0x64, 0x5f, 0x74,
	0x69, 0x6d, 0x65, 0x22, 0x39, 0x0a, 0x0e, 0x4e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x49, 0x64, 0x12, 0x27, 0x0a, 0x0f, 0x6e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0e,
	0x6e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x64, 0x22, 0x3c,
	0x0a, 0x0f, 0x4e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x64,
	0x73, 0x12, 0x29, 0x0a, 0x10, 0x6e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x5f, 0x69, 0x64, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x05, 0x52, 0x0f, 0x6e, 0x6f, 0x74,
	0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x64, 0x73, 0x22, 0x46, 0x0a, 0x1e,
	0x4d, 0x61, 0x72, 0x6b, 0x4e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x41, 0x73, 0x52, 0x65, 0x61, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x10,
	0x0a, 0x03, 0x6d, 0x73, 0x67, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6d, 0x73, 0x67,
	0x12, 0x12, 0x0a, 0x04, 0x72, 0x65, 0x61, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x08, 0x52, 0x04,
	0x72, 0x65, 0x61, 0x64, 0x22, 0x48, 0x0a, 0x1a, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x4e, 0x6f,
	0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x6d, 0x73, 0x67, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x03, 0x6d, 0x73, 0x67, 0x12, 0x18, 0x0a, 0x07, 0x64, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x64, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x64, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x64, 0x42, 0x28,
	0x5a, 0x26, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x73, 0x74, 0x65,
	0x76, 0x65, 0x2d, 0x6d, 0x69, 0x72, 0x2f, 0x64, 0x69, 0x69, 0x76, 0x69, 0x78, 0x5f, 0x62, 0x61,
	0x63, 0x6b, 0x65, 0x6e, 0x64, 0x2f, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_notification_proto_rawDescOnce sync.Once
	file_notification_proto_rawDescData = file_notification_proto_rawDesc
)

func file_notification_proto_rawDescGZIP() []byte {
	file_notification_proto_rawDescOnce.Do(func() {
		file_notification_proto_rawDescData = protoimpl.X.CompressGZIP(file_notification_proto_rawDescData)
	})
	return file_notification_proto_rawDescData
}

var file_notification_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_notification_proto_goTypes = []interface{}{
	(*GetNotificationRequest)(nil),         // 0: pb.GetNotificationRequest
	(*Notification)(nil),                   // 1: pb.Notification
	(*NotificationId)(nil),                 // 2: pb.NotificationId
	(*NotificationIds)(nil),                // 3: pb.NotificationIds
	(*MarkNotificationAsReadResponse)(nil), // 4: pb.MarkNotificationAsReadResponse
	(*DeleteNotificationResponse)(nil),     // 5: pb.DeleteNotificationResponse
	(*timestamppb.Timestamp)(nil),          // 6: google.protobuf.Timestamp
}
var file_notification_proto_depIdxs = []int32{
	6, // 0: pb.Notification.time:type_name -> google.protobuf.Timestamp
	6, // 1: pb.Notification.start_time:type_name -> google.protobuf.Timestamp
	6, // 2: pb.Notification.end_time:type_name -> google.protobuf.Timestamp
	3, // [3:3] is the sub-list for method output_type
	3, // [3:3] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_notification_proto_init() }
func file_notification_proto_init() {
	if File_notification_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_notification_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetNotificationRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_notification_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Notification); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_notification_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NotificationId); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_notification_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NotificationIds); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_notification_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MarkNotificationAsReadResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_notification_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeleteNotificationResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	file_notification_proto_msgTypes[1].OneofWrappers = []interface{}{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_notification_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_notification_proto_goTypes,
		DependencyIndexes: file_notification_proto_depIdxs,
		MessageInfos:      file_notification_proto_msgTypes,
	}.Build()
	File_notification_proto = out.File
	file_notification_proto_rawDesc = nil
	file_notification_proto_goTypes = nil
	file_notification_proto_depIdxs = nil
}
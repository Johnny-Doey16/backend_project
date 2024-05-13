// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v3.19.6
// source: rpc_news.proto

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

type CreateAnnouncementRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Title   string `protobuf:"bytes,1,opt,name=title,proto3" json:"title,omitempty"`
	Content string `protobuf:"bytes,2,opt,name=content,proto3" json:"content,omitempty"`
}

func (x *CreateAnnouncementRequest) Reset() {
	*x = CreateAnnouncementRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rpc_news_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateAnnouncementRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateAnnouncementRequest) ProtoMessage() {}

func (x *CreateAnnouncementRequest) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_news_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateAnnouncementRequest.ProtoReflect.Descriptor instead.
func (*CreateAnnouncementRequest) Descriptor() ([]byte, []int) {
	return file_rpc_news_proto_rawDescGZIP(), []int{0}
}

func (x *CreateAnnouncementRequest) GetTitle() string {
	if x != nil {
		return x.Title
	}
	return ""
}

func (x *CreateAnnouncementRequest) GetContent() string {
	if x != nil {
		return x.Content
	}
	return ""
}

type CreateAnnouncementResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	NewsId string `protobuf:"bytes,1,opt,name=news_id,json=newsId,proto3" json:"news_id,omitempty"`
	Msg    string `protobuf:"bytes,2,opt,name=msg,proto3" json:"msg,omitempty"`
}

func (x *CreateAnnouncementResponse) Reset() {
	*x = CreateAnnouncementResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rpc_news_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateAnnouncementResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateAnnouncementResponse) ProtoMessage() {}

func (x *CreateAnnouncementResponse) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_news_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateAnnouncementResponse.ProtoReflect.Descriptor instead.
func (*CreateAnnouncementResponse) Descriptor() ([]byte, []int) {
	return file_rpc_news_proto_rawDescGZIP(), []int{1}
}

func (x *CreateAnnouncementResponse) GetNewsId() string {
	if x != nil {
		return x.NewsId
	}
	return ""
}

func (x *CreateAnnouncementResponse) GetMsg() string {
	if x != nil {
		return x.Msg
	}
	return ""
}

type GetChurchAnnouncementsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId     string `protobuf:"bytes,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	PageNumber int32  `protobuf:"varint,2,opt,name=page_number,json=pageNumber,proto3" json:"page_number,omitempty"`
	PageSize   int32  `protobuf:"varint,3,opt,name=page_size,json=pageSize,proto3" json:"page_size,omitempty"`
}

func (x *GetChurchAnnouncementsRequest) Reset() {
	*x = GetChurchAnnouncementsRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rpc_news_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetChurchAnnouncementsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetChurchAnnouncementsRequest) ProtoMessage() {}

func (x *GetChurchAnnouncementsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_news_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetChurchAnnouncementsRequest.ProtoReflect.Descriptor instead.
func (*GetChurchAnnouncementsRequest) Descriptor() ([]byte, []int) {
	return file_rpc_news_proto_rawDescGZIP(), []int{2}
}

func (x *GetChurchAnnouncementsRequest) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

func (x *GetChurchAnnouncementsRequest) GetPageNumber() int32 {
	if x != nil {
		return x.PageNumber
	}
	return 0
}

func (x *GetChurchAnnouncementsRequest) GetPageSize() int32 {
	if x != nil {
		return x.PageSize
	}
	return 0
}

type GetAnnouncementsResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Post    *Announcement `protobuf:"bytes,1,opt,name=post,proto3" json:"post,omitempty"`
	HasMore bool          `protobuf:"varint,3,opt,name=has_more,json=hasMore,proto3" json:"has_more,omitempty"`
}

func (x *GetAnnouncementsResponse) Reset() {
	*x = GetAnnouncementsResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rpc_news_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetAnnouncementsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetAnnouncementsResponse) ProtoMessage() {}

func (x *GetAnnouncementsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_news_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetAnnouncementsResponse.ProtoReflect.Descriptor instead.
func (*GetAnnouncementsResponse) Descriptor() ([]byte, []int) {
	return file_rpc_news_proto_rawDescGZIP(), []int{3}
}

func (x *GetAnnouncementsResponse) GetPost() *Announcement {
	if x != nil {
		return x.Post
	}
	return nil
}

func (x *GetAnnouncementsResponse) GetHasMore() bool {
	if x != nil {
		return x.HasMore
	}
	return false
}

type Announcement struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	PostId       string                 `protobuf:"bytes,1,opt,name=post_id,json=postId,proto3" json:"post_id,omitempty"`
	UserId       string                 `protobuf:"bytes,2,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	Title        string                 `protobuf:"bytes,3,opt,name=title,proto3" json:"title,omitempty"`
	Content      string                 `protobuf:"bytes,4,opt,name=content,proto3" json:"content,omitempty"`
	Name         string                 `protobuf:"bytes,5,opt,name=name,proto3" json:"name,omitempty"`
	Username     string                 `protobuf:"bytes,6,opt,name=username,proto3" json:"username,omitempty"`
	ProfileImage string                 `protobuf:"bytes,7,opt,name=profile_image,json=profileImage,proto3" json:"profile_image,omitempty"`
	Timestamp    *timestamppb.Timestamp `protobuf:"bytes,8,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	IsVerified   bool                   `protobuf:"varint,9,opt,name=is_verified,json=isVerified,proto3" json:"is_verified,omitempty"`
	Reason       *string                `protobuf:"bytes,10,opt,name=reason,proto3,oneof" json:"reason,omitempty"`
}

func (x *Announcement) Reset() {
	*x = Announcement{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rpc_news_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Announcement) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Announcement) ProtoMessage() {}

func (x *Announcement) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_news_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Announcement.ProtoReflect.Descriptor instead.
func (*Announcement) Descriptor() ([]byte, []int) {
	return file_rpc_news_proto_rawDescGZIP(), []int{4}
}

func (x *Announcement) GetPostId() string {
	if x != nil {
		return x.PostId
	}
	return ""
}

func (x *Announcement) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

func (x *Announcement) GetTitle() string {
	if x != nil {
		return x.Title
	}
	return ""
}

func (x *Announcement) GetContent() string {
	if x != nil {
		return x.Content
	}
	return ""
}

func (x *Announcement) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Announcement) GetUsername() string {
	if x != nil {
		return x.Username
	}
	return ""
}

func (x *Announcement) GetProfileImage() string {
	if x != nil {
		return x.ProfileImage
	}
	return ""
}

func (x *Announcement) GetTimestamp() *timestamppb.Timestamp {
	if x != nil {
		return x.Timestamp
	}
	return nil
}

func (x *Announcement) GetIsVerified() bool {
	if x != nil {
		return x.IsVerified
	}
	return false
}

func (x *Announcement) GetReason() string {
	if x != nil && x.Reason != nil {
		return *x.Reason
	}
	return ""
}

type GetAnnouncementsForUserRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	PageNumber int32 `protobuf:"varint,1,opt,name=page_number,json=pageNumber,proto3" json:"page_number,omitempty"`
	PageSize   int32 `protobuf:"varint,2,opt,name=page_size,json=pageSize,proto3" json:"page_size,omitempty"`
}

func (x *GetAnnouncementsForUserRequest) Reset() {
	*x = GetAnnouncementsForUserRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rpc_news_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetAnnouncementsForUserRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetAnnouncementsForUserRequest) ProtoMessage() {}

func (x *GetAnnouncementsForUserRequest) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_news_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetAnnouncementsForUserRequest.ProtoReflect.Descriptor instead.
func (*GetAnnouncementsForUserRequest) Descriptor() ([]byte, []int) {
	return file_rpc_news_proto_rawDescGZIP(), []int{5}
}

func (x *GetAnnouncementsForUserRequest) GetPageNumber() int32 {
	if x != nil {
		return x.PageNumber
	}
	return 0
}

func (x *GetAnnouncementsForUserRequest) GetPageSize() int32 {
	if x != nil {
		return x.PageSize
	}
	return 0
}

var File_rpc_news_proto protoreflect.FileDescriptor

var file_rpc_news_proto_rawDesc = []byte{
	0x0a, 0x0e, 0x72, 0x70, 0x63, 0x5f, 0x6e, 0x65, 0x77, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x02, 0x70, 0x62, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x4b, 0x0a, 0x19, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x41,
	0x6e, 0x6e, 0x6f, 0x75, 0x6e, 0x63, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x05, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x63, 0x6f, 0x6e, 0x74,
	0x65, 0x6e, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65,
	0x6e, 0x74, 0x22, 0x47, 0x0a, 0x1a, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x41, 0x6e, 0x6e, 0x6f,
	0x75, 0x6e, 0x63, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x17, 0x0a, 0x07, 0x6e, 0x65, 0x77, 0x73, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x06, 0x6e, 0x65, 0x77, 0x73, 0x49, 0x64, 0x12, 0x10, 0x0a, 0x03, 0x6d, 0x73, 0x67,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6d, 0x73, 0x67, 0x22, 0x76, 0x0a, 0x1d, 0x47,
	0x65, 0x74, 0x43, 0x68, 0x75, 0x72, 0x63, 0x68, 0x41, 0x6e, 0x6e, 0x6f, 0x75, 0x6e, 0x63, 0x65,
	0x6d, 0x65, 0x6e, 0x74, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x17, 0x0a, 0x07,
	0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x75,
	0x73, 0x65, 0x72, 0x49, 0x64, 0x12, 0x1f, 0x0a, 0x0b, 0x70, 0x61, 0x67, 0x65, 0x5f, 0x6e, 0x75,
	0x6d, 0x62, 0x65, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0a, 0x70, 0x61, 0x67, 0x65,
	0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x12, 0x1b, 0x0a, 0x09, 0x70, 0x61, 0x67, 0x65, 0x5f, 0x73,
	0x69, 0x7a, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x08, 0x70, 0x61, 0x67, 0x65, 0x53,
	0x69, 0x7a, 0x65, 0x22, 0x5b, 0x0a, 0x18, 0x47, 0x65, 0x74, 0x41, 0x6e, 0x6e, 0x6f, 0x75, 0x6e,
	0x63, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x24, 0x0a, 0x04, 0x70, 0x6f, 0x73, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x10, 0x2e,
	0x70, 0x62, 0x2e, 0x41, 0x6e, 0x6e, 0x6f, 0x75, 0x6e, 0x63, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x52,
	0x04, 0x70, 0x6f, 0x73, 0x74, 0x12, 0x19, 0x0a, 0x08, 0x68, 0x61, 0x73, 0x5f, 0x6d, 0x6f, 0x72,
	0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x68, 0x61, 0x73, 0x4d, 0x6f, 0x72, 0x65,
	0x22, 0xc8, 0x02, 0x0a, 0x0c, 0x41, 0x6e, 0x6e, 0x6f, 0x75, 0x6e, 0x63, 0x65, 0x6d, 0x65, 0x6e,
	0x74, 0x12, 0x17, 0x0a, 0x07, 0x70, 0x6f, 0x73, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x06, 0x70, 0x6f, 0x73, 0x74, 0x49, 0x64, 0x12, 0x17, 0x0a, 0x07, 0x75, 0x73,
	0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x75, 0x73, 0x65,
	0x72, 0x49, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x05, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x63, 0x6f, 0x6e,
	0x74, 0x65, 0x6e, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x63, 0x6f, 0x6e, 0x74,
	0x65, 0x6e, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x75, 0x73, 0x65, 0x72, 0x6e,
	0x61, 0x6d, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x75, 0x73, 0x65, 0x72, 0x6e,
	0x61, 0x6d, 0x65, 0x12, 0x23, 0x0a, 0x0d, 0x70, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x5f, 0x69,
	0x6d, 0x61, 0x67, 0x65, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x70, 0x72, 0x6f, 0x66,
	0x69, 0x6c, 0x65, 0x49, 0x6d, 0x61, 0x67, 0x65, 0x12, 0x38, 0x0a, 0x09, 0x74, 0x69, 0x6d, 0x65,
	0x73, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x08, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69,
	0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61,
	0x6d, 0x70, 0x12, 0x1f, 0x0a, 0x0b, 0x69, 0x73, 0x5f, 0x76, 0x65, 0x72, 0x69, 0x66, 0x69, 0x65,
	0x64, 0x18, 0x09, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0a, 0x69, 0x73, 0x56, 0x65, 0x72, 0x69, 0x66,
	0x69, 0x65, 0x64, 0x12, 0x1b, 0x0a, 0x06, 0x72, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x18, 0x0a, 0x20,
	0x01, 0x28, 0x09, 0x48, 0x00, 0x52, 0x06, 0x72, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x88, 0x01, 0x01,
	0x42, 0x09, 0x0a, 0x07, 0x5f, 0x72, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x22, 0x5e, 0x0a, 0x1e, 0x47,
	0x65, 0x74, 0x41, 0x6e, 0x6e, 0x6f, 0x75, 0x6e, 0x63, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x46,
	0x6f, 0x72, 0x55, 0x73, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1f, 0x0a,
	0x0b, 0x70, 0x61, 0x67, 0x65, 0x5f, 0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x05, 0x52, 0x0a, 0x70, 0x61, 0x67, 0x65, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x12, 0x1b,
	0x0a, 0x09, 0x70, 0x61, 0x67, 0x65, 0x5f, 0x73, 0x69, 0x7a, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x05, 0x52, 0x08, 0x70, 0x61, 0x67, 0x65, 0x53, 0x69, 0x7a, 0x65, 0x42, 0x28, 0x5a, 0x26, 0x67,
	0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x73, 0x74, 0x65, 0x76, 0x65, 0x2d,
	0x6d, 0x69, 0x72, 0x2f, 0x64, 0x69, 0x69, 0x76, 0x69, 0x78, 0x5f, 0x62, 0x61, 0x63, 0x6b, 0x65,
	0x6e, 0x64, 0x2f, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_rpc_news_proto_rawDescOnce sync.Once
	file_rpc_news_proto_rawDescData = file_rpc_news_proto_rawDesc
)

func file_rpc_news_proto_rawDescGZIP() []byte {
	file_rpc_news_proto_rawDescOnce.Do(func() {
		file_rpc_news_proto_rawDescData = protoimpl.X.CompressGZIP(file_rpc_news_proto_rawDescData)
	})
	return file_rpc_news_proto_rawDescData
}

var file_rpc_news_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_rpc_news_proto_goTypes = []interface{}{
	(*CreateAnnouncementRequest)(nil),      // 0: pb.CreateAnnouncementRequest
	(*CreateAnnouncementResponse)(nil),     // 1: pb.CreateAnnouncementResponse
	(*GetChurchAnnouncementsRequest)(nil),  // 2: pb.GetChurchAnnouncementsRequest
	(*GetAnnouncementsResponse)(nil),       // 3: pb.GetAnnouncementsResponse
	(*Announcement)(nil),                   // 4: pb.Announcement
	(*GetAnnouncementsForUserRequest)(nil), // 5: pb.GetAnnouncementsForUserRequest
	(*timestamppb.Timestamp)(nil),          // 6: google.protobuf.Timestamp
}
var file_rpc_news_proto_depIdxs = []int32{
	4, // 0: pb.GetAnnouncementsResponse.post:type_name -> pb.Announcement
	6, // 1: pb.Announcement.timestamp:type_name -> google.protobuf.Timestamp
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_rpc_news_proto_init() }
func file_rpc_news_proto_init() {
	if File_rpc_news_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_rpc_news_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateAnnouncementRequest); i {
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
		file_rpc_news_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateAnnouncementResponse); i {
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
		file_rpc_news_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetChurchAnnouncementsRequest); i {
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
		file_rpc_news_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetAnnouncementsResponse); i {
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
		file_rpc_news_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Announcement); i {
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
		file_rpc_news_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetAnnouncementsForUserRequest); i {
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
	file_rpc_news_proto_msgTypes[4].OneofWrappers = []interface{}{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_rpc_news_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_rpc_news_proto_goTypes,
		DependencyIndexes: file_rpc_news_proto_depIdxs,
		MessageInfos:      file_rpc_news_proto_msgTypes,
	}.Build()
	File_rpc_news_proto = out.File
	file_rpc_news_proto_rawDesc = nil
	file_rpc_news_proto_goTypes = nil
	file_rpc_news_proto_depIdxs = nil
}
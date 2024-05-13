// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v3.19.6
// source: rpc_church.proto

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

// A message for User's membership data
type Membership struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id       int32                  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	JoinDate *timestamppb.Timestamp `protobuf:"bytes,2,opt,name=join_date,json=joinDate,proto3" json:"join_date,omitempty"`
	// Depending on the context, this could be either a church_id or denomination_id
	EntityId int32 `protobuf:"varint,3,opt,name=entity_id,json=entityId,proto3" json:"entity_id,omitempty"`
}

func (x *Membership) Reset() {
	*x = Membership{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rpc_church_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Membership) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Membership) ProtoMessage() {}

func (x *Membership) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_church_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Membership.ProtoReflect.Descriptor instead.
func (*Membership) Descriptor() ([]byte, []int) {
	return file_rpc_church_proto_rawDescGZIP(), []int{0}
}

func (x *Membership) GetId() int32 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *Membership) GetJoinDate() *timestamppb.Timestamp {
	if x != nil {
		return x.JoinDate
	}
	return nil
}

func (x *Membership) GetEntityId() int32 {
	if x != nil {
		return x.EntityId
	}
	return 0
}

// Request to join or leave a church/denomination
type MembershipChangeRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to Entity:
	//
	//	*MembershipChangeRequest_ChurchId
	//	*MembershipChangeRequest_DenominationId
	Entity isMembershipChangeRequest_Entity `protobuf_oneof:"entity"`
}

func (x *MembershipChangeRequest) Reset() {
	*x = MembershipChangeRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rpc_church_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MembershipChangeRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MembershipChangeRequest) ProtoMessage() {}

func (x *MembershipChangeRequest) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_church_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MembershipChangeRequest.ProtoReflect.Descriptor instead.
func (*MembershipChangeRequest) Descriptor() ([]byte, []int) {
	return file_rpc_church_proto_rawDescGZIP(), []int{1}
}

func (m *MembershipChangeRequest) GetEntity() isMembershipChangeRequest_Entity {
	if m != nil {
		return m.Entity
	}
	return nil
}

func (x *MembershipChangeRequest) GetChurchId() int32 {
	if x, ok := x.GetEntity().(*MembershipChangeRequest_ChurchId); ok {
		return x.ChurchId
	}
	return 0
}

func (x *MembershipChangeRequest) GetDenominationId() int32 {
	if x, ok := x.GetEntity().(*MembershipChangeRequest_DenominationId); ok {
		return x.DenominationId
	}
	return 0
}

type isMembershipChangeRequest_Entity interface {
	isMembershipChangeRequest_Entity()
}

type MembershipChangeRequest_ChurchId struct {
	ChurchId int32 `protobuf:"varint,1,opt,name=church_id,json=churchId,proto3,oneof"`
}

type MembershipChangeRequest_DenominationId struct {
	DenominationId int32 `protobuf:"varint,2,opt,name=denomination_id,json=denominationId,proto3,oneof"`
}

func (*MembershipChangeRequest_ChurchId) isMembershipChangeRequest_Entity() {}

func (*MembershipChangeRequest_DenominationId) isMembershipChangeRequest_Entity() {}

type MembershipChangeResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *MembershipChangeResponse) Reset() {
	*x = MembershipChangeResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rpc_church_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MembershipChangeResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MembershipChangeResponse) ProtoMessage() {}

func (x *MembershipChangeResponse) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_church_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MembershipChangeResponse.ProtoReflect.Descriptor instead.
func (*MembershipChangeResponse) Descriptor() ([]byte, []int) {
	return file_rpc_church_proto_rawDescGZIP(), []int{2}
}

// Search request message
type SearchRequestChurch struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Query        string      `protobuf:"bytes,1,opt,name=query,proto3" json:"query,omitempty"` // The search query string
	PageNumber   int32       `protobuf:"varint,2,opt,name=page_number,json=pageNumber,proto3" json:"page_number,omitempty"`
	PageSize     int32       `protobuf:"varint,3,opt,name=page_size,json=pageSize,proto3" json:"page_size,omitempty"`
	UserLocation *Coordinate `protobuf:"bytes,4,opt,name=user_location,json=userLocation,proto3,oneof" json:"user_location,omitempty"` // Optional for nearby search
}

func (x *SearchRequestChurch) Reset() {
	*x = SearchRequestChurch{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rpc_church_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SearchRequestChurch) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SearchRequestChurch) ProtoMessage() {}

func (x *SearchRequestChurch) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_church_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SearchRequestChurch.ProtoReflect.Descriptor instead.
func (*SearchRequestChurch) Descriptor() ([]byte, []int) {
	return file_rpc_church_proto_rawDescGZIP(), []int{3}
}

func (x *SearchRequestChurch) GetQuery() string {
	if x != nil {
		return x.Query
	}
	return ""
}

func (x *SearchRequestChurch) GetPageNumber() int32 {
	if x != nil {
		return x.PageNumber
	}
	return 0
}

func (x *SearchRequestChurch) GetPageSize() int32 {
	if x != nil {
		return x.PageSize
	}
	return 0
}

func (x *SearchRequestChurch) GetUserLocation() *Coordinate {
	if x != nil {
		return x.UserLocation
	}
	return nil
}

type GetChurchProfileRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	AuthId string `protobuf:"bytes,1,opt,name=auth_id,json=authId,proto3" json:"auth_id,omitempty"`
}

func (x *GetChurchProfileRequest) Reset() {
	*x = GetChurchProfileRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rpc_church_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetChurchProfileRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetChurchProfileRequest) ProtoMessage() {}

func (x *GetChurchProfileRequest) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_church_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetChurchProfileRequest.ProtoReflect.Descriptor instead.
func (*GetChurchProfileRequest) Descriptor() ([]byte, []int) {
	return file_rpc_church_proto_rawDescGZIP(), []int{4}
}

func (x *GetChurchProfileRequest) GetAuthId() string {
	if x != nil {
		return x.AuthId
	}
	return ""
}

type GetChurchProfileResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Church      *Church `protobuf:"bytes,1,opt,name=church,proto3" json:"church,omitempty"`
	IsFollowing bool    `protobuf:"varint,2,opt,name=is_following,json=isFollowing,proto3" json:"is_following,omitempty"`
	IsFollowed  bool    `protobuf:"varint,3,opt,name=is_followed,json=isFollowed,proto3" json:"is_followed,omitempty"`
	IsMember    bool    `protobuf:"varint,4,opt,name=is_member,json=isMember,proto3" json:"is_member,omitempty"`
}

func (x *GetChurchProfileResponse) Reset() {
	*x = GetChurchProfileResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rpc_church_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetChurchProfileResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetChurchProfileResponse) ProtoMessage() {}

func (x *GetChurchProfileResponse) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_church_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetChurchProfileResponse.ProtoReflect.Descriptor instead.
func (*GetChurchProfileResponse) Descriptor() ([]byte, []int) {
	return file_rpc_church_proto_rawDescGZIP(), []int{5}
}

func (x *GetChurchProfileResponse) GetChurch() *Church {
	if x != nil {
		return x.Church
	}
	return nil
}

func (x *GetChurchProfileResponse) GetIsFollowing() bool {
	if x != nil {
		return x.IsFollowing
	}
	return false
}

func (x *GetChurchProfileResponse) GetIsFollowed() bool {
	if x != nil {
		return x.IsFollowed
	}
	return false
}

func (x *GetChurchProfileResponse) GetIsMember() bool {
	if x != nil {
		return x.IsMember
	}
	return false
}

type GetUserChurchRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *GetUserChurchRequest) Reset() {
	*x = GetUserChurchRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rpc_church_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetUserChurchRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetUserChurchRequest) ProtoMessage() {}

func (x *GetUserChurchRequest) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_church_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetUserChurchRequest.ProtoReflect.Descriptor instead.
func (*GetUserChurchRequest) Descriptor() ([]byte, []int) {
	return file_rpc_church_proto_rawDescGZIP(), []int{6}
}

type GetUserChurchResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Church      *Church `protobuf:"bytes,1,opt,name=church,proto3" json:"church,omitempty"`
	IsFollowing bool    `protobuf:"varint,2,opt,name=is_following,json=isFollowing,proto3" json:"is_following,omitempty"`
	IsFollowed  bool    `protobuf:"varint,3,opt,name=is_followed,json=isFollowed,proto3" json:"is_followed,omitempty"`
}

func (x *GetUserChurchResponse) Reset() {
	*x = GetUserChurchResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rpc_church_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetUserChurchResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetUserChurchResponse) ProtoMessage() {}

func (x *GetUserChurchResponse) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_church_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetUserChurchResponse.ProtoReflect.Descriptor instead.
func (*GetUserChurchResponse) Descriptor() ([]byte, []int) {
	return file_rpc_church_proto_rawDescGZIP(), []int{7}
}

func (x *GetUserChurchResponse) GetChurch() *Church {
	if x != nil {
		return x.Church
	}
	return nil
}

func (x *GetUserChurchResponse) GetIsFollowing() bool {
	if x != nil {
		return x.IsFollowing
	}
	return false
}

func (x *GetUserChurchResponse) GetIsFollowed() bool {
	if x != nil {
		return x.IsFollowed
	}
	return false
}

type ChurchMember struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	AuthId     string `protobuf:"bytes,1,opt,name=auth_id,json=authId,proto3" json:"auth_id,omitempty"`
	ImageUrl   string `protobuf:"bytes,2,opt,name=image_url,json=imageUrl,proto3" json:"image_url,omitempty"`
	FirstName  string `protobuf:"bytes,3,opt,name=first_name,json=firstName,proto3" json:"first_name,omitempty"`
	Username   string `protobuf:"bytes,4,opt,name=username,proto3" json:"username,omitempty"`
	IsVerified bool   `protobuf:"varint,5,opt,name=is_verified,json=isVerified,proto3" json:"is_verified,omitempty"`
}

func (x *ChurchMember) Reset() {
	*x = ChurchMember{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rpc_church_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ChurchMember) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ChurchMember) ProtoMessage() {}

func (x *ChurchMember) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_church_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ChurchMember.ProtoReflect.Descriptor instead.
func (*ChurchMember) Descriptor() ([]byte, []int) {
	return file_rpc_church_proto_rawDescGZIP(), []int{8}
}

func (x *ChurchMember) GetAuthId() string {
	if x != nil {
		return x.AuthId
	}
	return ""
}

func (x *ChurchMember) GetImageUrl() string {
	if x != nil {
		return x.ImageUrl
	}
	return ""
}

func (x *ChurchMember) GetFirstName() string {
	if x != nil {
		return x.FirstName
	}
	return ""
}

func (x *ChurchMember) GetUsername() string {
	if x != nil {
		return x.Username
	}
	return ""
}

func (x *ChurchMember) GetIsVerified() bool {
	if x != nil {
		return x.IsVerified
	}
	return false
}

// Get church members
type GetChurchMembersRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ChurchId   int32  `protobuf:"varint,1,opt,name=church_id,json=churchId,proto3" json:"church_id,omitempty"`
	PageNumber int32  `protobuf:"varint,2,opt,name=page_number,json=pageNumber,proto3" json:"page_number,omitempty"`
	PageSize   int32  `protobuf:"varint,3,opt,name=page_size,json=pageSize,proto3" json:"page_size,omitempty"`
	Order      string `protobuf:"bytes,4,opt,name=order,proto3" json:"order,omitempty"`
}

func (x *GetChurchMembersRequest) Reset() {
	*x = GetChurchMembersRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rpc_church_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetChurchMembersRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetChurchMembersRequest) ProtoMessage() {}

func (x *GetChurchMembersRequest) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_church_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetChurchMembersRequest.ProtoReflect.Descriptor instead.
func (*GetChurchMembersRequest) Descriptor() ([]byte, []int) {
	return file_rpc_church_proto_rawDescGZIP(), []int{9}
}

func (x *GetChurchMembersRequest) GetChurchId() int32 {
	if x != nil {
		return x.ChurchId
	}
	return 0
}

func (x *GetChurchMembersRequest) GetPageNumber() int32 {
	if x != nil {
		return x.PageNumber
	}
	return 0
}

func (x *GetChurchMembersRequest) GetPageSize() int32 {
	if x != nil {
		return x.PageSize
	}
	return 0
}

func (x *GetChurchMembersRequest) GetOrder() string {
	if x != nil {
		return x.Order
	}
	return ""
}

type GetChurchMembersResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	IsMember bool          `protobuf:"varint,1,opt,name=is_member,json=isMember,proto3" json:"is_member,omitempty"`
	Members  *ChurchMember `protobuf:"bytes,2,opt,name=members,proto3" json:"members,omitempty"`
}

func (x *GetChurchMembersResponse) Reset() {
	*x = GetChurchMembersResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rpc_church_proto_msgTypes[10]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetChurchMembersResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetChurchMembersResponse) ProtoMessage() {}

func (x *GetChurchMembersResponse) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_church_proto_msgTypes[10]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetChurchMembersResponse.ProtoReflect.Descriptor instead.
func (*GetChurchMembersResponse) Descriptor() ([]byte, []int) {
	return file_rpc_church_proto_rawDescGZIP(), []int{10}
}

func (x *GetChurchMembersResponse) GetIsMember() bool {
	if x != nil {
		return x.IsMember
	}
	return false
}

func (x *GetChurchMembersResponse) GetMembers() *ChurchMember {
	if x != nil {
		return x.Members
	}
	return nil
}

var File_rpc_church_proto protoreflect.FileDescriptor

var file_rpc_church_proto_rawDesc = []byte{
	0x0a, 0x10, 0x72, 0x70, 0x63, 0x5f, 0x63, 0x68, 0x75, 0x72, 0x63, 0x68, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x02, 0x70, 0x62, 0x1a, 0x0e, 0x6c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x0c, 0x63, 0x68, 0x75, 0x72, 0x63, 0x68, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x72, 0x0a, 0x0a, 0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x73,
	0x68, 0x69, 0x70, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52,
	0x02, 0x69, 0x64, 0x12, 0x37, 0x0a, 0x09, 0x6a, 0x6f, 0x69, 0x6e, 0x5f, 0x64, 0x61, 0x74, 0x65,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61,
	0x6d, 0x70, 0x52, 0x08, 0x6a, 0x6f, 0x69, 0x6e, 0x44, 0x61, 0x74, 0x65, 0x12, 0x1b, 0x0a, 0x09,
	0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x5f, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52,
	0x08, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x49, 0x64, 0x22, 0x6d, 0x0a, 0x17, 0x4d, 0x65, 0x6d,
	0x62, 0x65, 0x72, 0x73, 0x68, 0x69, 0x70, 0x43, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x1d, 0x0a, 0x09, 0x63, 0x68, 0x75, 0x72, 0x63, 0x68, 0x5f, 0x69,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x48, 0x00, 0x52, 0x08, 0x63, 0x68, 0x75, 0x72, 0x63,
	0x68, 0x49, 0x64, 0x12, 0x29, 0x0a, 0x0f, 0x64, 0x65, 0x6e, 0x6f, 0x6d, 0x69, 0x6e, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x48, 0x00, 0x52, 0x0e,
	0x64, 0x65, 0x6e, 0x6f, 0x6d, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x64, 0x42, 0x08,
	0x0a, 0x06, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x22, 0x1a, 0x0a, 0x18, 0x4d, 0x65, 0x6d, 0x62,
	0x65, 0x72, 0x73, 0x68, 0x69, 0x70, 0x43, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x22, 0xb5, 0x01, 0x0a, 0x13, 0x53, 0x65, 0x61, 0x72, 0x63, 0x68, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x43, 0x68, 0x75, 0x72, 0x63, 0x68, 0x12, 0x14, 0x0a, 0x05,
	0x71, 0x75, 0x65, 0x72, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x71, 0x75, 0x65,
	0x72, 0x79, 0x12, 0x1f, 0x0a, 0x0b, 0x70, 0x61, 0x67, 0x65, 0x5f, 0x6e, 0x75, 0x6d, 0x62, 0x65,
	0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0a, 0x70, 0x61, 0x67, 0x65, 0x4e, 0x75, 0x6d,
	0x62, 0x65, 0x72, 0x12, 0x1b, 0x0a, 0x09, 0x70, 0x61, 0x67, 0x65, 0x5f, 0x73, 0x69, 0x7a, 0x65,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x08, 0x70, 0x61, 0x67, 0x65, 0x53, 0x69, 0x7a, 0x65,
	0x12, 0x38, 0x0a, 0x0d, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x6c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x70, 0x62, 0x2e, 0x43, 0x6f, 0x6f,
	0x72, 0x64, 0x69, 0x6e, 0x61, 0x74, 0x65, 0x48, 0x00, 0x52, 0x0c, 0x75, 0x73, 0x65, 0x72, 0x4c,
	0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x88, 0x01, 0x01, 0x42, 0x10, 0x0a, 0x0e, 0x5f, 0x75,
	0x73, 0x65, 0x72, 0x5f, 0x6c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x22, 0x32, 0x0a, 0x17,
	0x47, 0x65, 0x74, 0x43, 0x68, 0x75, 0x72, 0x63, 0x68, 0x50, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x17, 0x0a, 0x07, 0x61, 0x75, 0x74, 0x68, 0x5f,
	0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x61, 0x75, 0x74, 0x68, 0x49, 0x64,
	0x22, 0x9f, 0x01, 0x0a, 0x18, 0x47, 0x65, 0x74, 0x43, 0x68, 0x75, 0x72, 0x63, 0x68, 0x50, 0x72,
	0x6f, 0x66, 0x69, 0x6c, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x22, 0x0a,
	0x06, 0x63, 0x68, 0x75, 0x72, 0x63, 0x68, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0a, 0x2e,
	0x70, 0x62, 0x2e, 0x43, 0x68, 0x75, 0x72, 0x63, 0x68, 0x52, 0x06, 0x63, 0x68, 0x75, 0x72, 0x63,
	0x68, 0x12, 0x21, 0x0a, 0x0c, 0x69, 0x73, 0x5f, 0x66, 0x6f, 0x6c, 0x6c, 0x6f, 0x77, 0x69, 0x6e,
	0x67, 0x18, 0x02, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0b, 0x69, 0x73, 0x46, 0x6f, 0x6c, 0x6c, 0x6f,
	0x77, 0x69, 0x6e, 0x67, 0x12, 0x1f, 0x0a, 0x0b, 0x69, 0x73, 0x5f, 0x66, 0x6f, 0x6c, 0x6c, 0x6f,
	0x77, 0x65, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0a, 0x69, 0x73, 0x46, 0x6f, 0x6c,
	0x6c, 0x6f, 0x77, 0x65, 0x64, 0x12, 0x1b, 0x0a, 0x09, 0x69, 0x73, 0x5f, 0x6d, 0x65, 0x6d, 0x62,
	0x65, 0x72, 0x18, 0x04, 0x20, 0x01, 0x28, 0x08, 0x52, 0x08, 0x69, 0x73, 0x4d, 0x65, 0x6d, 0x62,
	0x65, 0x72, 0x22, 0x16, 0x0a, 0x14, 0x47, 0x65, 0x74, 0x55, 0x73, 0x65, 0x72, 0x43, 0x68, 0x75,
	0x72, 0x63, 0x68, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x7f, 0x0a, 0x15, 0x47, 0x65,
	0x74, 0x55, 0x73, 0x65, 0x72, 0x43, 0x68, 0x75, 0x72, 0x63, 0x68, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x22, 0x0a, 0x06, 0x63, 0x68, 0x75, 0x72, 0x63, 0x68, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x0a, 0x2e, 0x70, 0x62, 0x2e, 0x43, 0x68, 0x75, 0x72, 0x63, 0x68, 0x52,
	0x06, 0x63, 0x68, 0x75, 0x72, 0x63, 0x68, 0x12, 0x21, 0x0a, 0x0c, 0x69, 0x73, 0x5f, 0x66, 0x6f,
	0x6c, 0x6c, 0x6f, 0x77, 0x69, 0x6e, 0x67, 0x18, 0x02, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0b, 0x69,
	0x73, 0x46, 0x6f, 0x6c, 0x6c, 0x6f, 0x77, 0x69, 0x6e, 0x67, 0x12, 0x1f, 0x0a, 0x0b, 0x69, 0x73,
	0x5f, 0x66, 0x6f, 0x6c, 0x6c, 0x6f, 0x77, 0x65, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x08, 0x52,
	0x0a, 0x69, 0x73, 0x46, 0x6f, 0x6c, 0x6c, 0x6f, 0x77, 0x65, 0x64, 0x22, 0xa0, 0x01, 0x0a, 0x0c,
	0x43, 0x68, 0x75, 0x72, 0x63, 0x68, 0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x12, 0x17, 0x0a, 0x07,
	0x61, 0x75, 0x74, 0x68, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x61,
	0x75, 0x74, 0x68, 0x49, 0x64, 0x12, 0x1b, 0x0a, 0x09, 0x69, 0x6d, 0x61, 0x67, 0x65, 0x5f, 0x75,
	0x72, 0x6c, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x69, 0x6d, 0x61, 0x67, 0x65, 0x55,
	0x72, 0x6c, 0x12, 0x1d, 0x0a, 0x0a, 0x66, 0x69, 0x72, 0x73, 0x74, 0x5f, 0x6e, 0x61, 0x6d, 0x65,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x66, 0x69, 0x72, 0x73, 0x74, 0x4e, 0x61, 0x6d,
	0x65, 0x12, 0x1a, 0x0a, 0x08, 0x75, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x08, 0x75, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x1f, 0x0a,
	0x0b, 0x69, 0x73, 0x5f, 0x76, 0x65, 0x72, 0x69, 0x66, 0x69, 0x65, 0x64, 0x18, 0x05, 0x20, 0x01,
	0x28, 0x08, 0x52, 0x0a, 0x69, 0x73, 0x56, 0x65, 0x72, 0x69, 0x66, 0x69, 0x65, 0x64, 0x22, 0x8a,
	0x01, 0x0a, 0x17, 0x47, 0x65, 0x74, 0x43, 0x68, 0x75, 0x72, 0x63, 0x68, 0x4d, 0x65, 0x6d, 0x62,
	0x65, 0x72, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1b, 0x0a, 0x09, 0x63, 0x68,
	0x75, 0x72, 0x63, 0x68, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x08, 0x63,
	0x68, 0x75, 0x72, 0x63, 0x68, 0x49, 0x64, 0x12, 0x1f, 0x0a, 0x0b, 0x70, 0x61, 0x67, 0x65, 0x5f,
	0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0a, 0x70, 0x61,
	0x67, 0x65, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x12, 0x1b, 0x0a, 0x09, 0x70, 0x61, 0x67, 0x65,
	0x5f, 0x73, 0x69, 0x7a, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x08, 0x70, 0x61, 0x67,
	0x65, 0x53, 0x69, 0x7a, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x22, 0x63, 0x0a, 0x18, 0x47,
	0x65, 0x74, 0x43, 0x68, 0x75, 0x72, 0x63, 0x68, 0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x73, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x1b, 0x0a, 0x09, 0x69, 0x73, 0x5f, 0x6d, 0x65,
	0x6d, 0x62, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x08, 0x69, 0x73, 0x4d, 0x65,
	0x6d, 0x62, 0x65, 0x72, 0x12, 0x2a, 0x0a, 0x07, 0x6d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x73, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x70, 0x62, 0x2e, 0x43, 0x68, 0x75, 0x72, 0x63,
	0x68, 0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x52, 0x07, 0x6d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x73,
	0x42, 0x28, 0x5a, 0x26, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x73,
	0x74, 0x65, 0x76, 0x65, 0x2d, 0x6d, 0x69, 0x72, 0x2f, 0x64, 0x69, 0x69, 0x76, 0x69, 0x78, 0x5f,
	0x62, 0x61, 0x63, 0x6b, 0x65, 0x6e, 0x64, 0x2f, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_rpc_church_proto_rawDescOnce sync.Once
	file_rpc_church_proto_rawDescData = file_rpc_church_proto_rawDesc
)

func file_rpc_church_proto_rawDescGZIP() []byte {
	file_rpc_church_proto_rawDescOnce.Do(func() {
		file_rpc_church_proto_rawDescData = protoimpl.X.CompressGZIP(file_rpc_church_proto_rawDescData)
	})
	return file_rpc_church_proto_rawDescData
}

var file_rpc_church_proto_msgTypes = make([]protoimpl.MessageInfo, 11)
var file_rpc_church_proto_goTypes = []interface{}{
	(*Membership)(nil),               // 0: pb.Membership
	(*MembershipChangeRequest)(nil),  // 1: pb.MembershipChangeRequest
	(*MembershipChangeResponse)(nil), // 2: pb.MembershipChangeResponse
	(*SearchRequestChurch)(nil),      // 3: pb.SearchRequestChurch
	(*GetChurchProfileRequest)(nil),  // 4: pb.GetChurchProfileRequest
	(*GetChurchProfileResponse)(nil), // 5: pb.GetChurchProfileResponse
	(*GetUserChurchRequest)(nil),     // 6: pb.GetUserChurchRequest
	(*GetUserChurchResponse)(nil),    // 7: pb.GetUserChurchResponse
	(*ChurchMember)(nil),             // 8: pb.ChurchMember
	(*GetChurchMembersRequest)(nil),  // 9: pb.GetChurchMembersRequest
	(*GetChurchMembersResponse)(nil), // 10: pb.GetChurchMembersResponse
	(*timestamppb.Timestamp)(nil),    // 11: google.protobuf.Timestamp
	(*Coordinate)(nil),               // 12: pb.Coordinate
	(*Church)(nil),                   // 13: pb.Church
}
var file_rpc_church_proto_depIdxs = []int32{
	11, // 0: pb.Membership.join_date:type_name -> google.protobuf.Timestamp
	12, // 1: pb.SearchRequestChurch.user_location:type_name -> pb.Coordinate
	13, // 2: pb.GetChurchProfileResponse.church:type_name -> pb.Church
	13, // 3: pb.GetUserChurchResponse.church:type_name -> pb.Church
	8,  // 4: pb.GetChurchMembersResponse.members:type_name -> pb.ChurchMember
	5,  // [5:5] is the sub-list for method output_type
	5,  // [5:5] is the sub-list for method input_type
	5,  // [5:5] is the sub-list for extension type_name
	5,  // [5:5] is the sub-list for extension extendee
	0,  // [0:5] is the sub-list for field type_name
}

func init() { file_rpc_church_proto_init() }
func file_rpc_church_proto_init() {
	if File_rpc_church_proto != nil {
		return
	}
	file_location_proto_init()
	file_church_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_rpc_church_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Membership); i {
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
		file_rpc_church_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MembershipChangeRequest); i {
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
		file_rpc_church_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MembershipChangeResponse); i {
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
		file_rpc_church_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SearchRequestChurch); i {
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
		file_rpc_church_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetChurchProfileRequest); i {
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
		file_rpc_church_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetChurchProfileResponse); i {
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
		file_rpc_church_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetUserChurchRequest); i {
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
		file_rpc_church_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetUserChurchResponse); i {
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
		file_rpc_church_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ChurchMember); i {
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
		file_rpc_church_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetChurchMembersRequest); i {
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
		file_rpc_church_proto_msgTypes[10].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetChurchMembersResponse); i {
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
	file_rpc_church_proto_msgTypes[1].OneofWrappers = []interface{}{
		(*MembershipChangeRequest_ChurchId)(nil),
		(*MembershipChangeRequest_DenominationId)(nil),
	}
	file_rpc_church_proto_msgTypes[3].OneofWrappers = []interface{}{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_rpc_church_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   11,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_rpc_church_proto_goTypes,
		DependencyIndexes: file_rpc_church_proto_depIdxs,
		MessageInfos:      file_rpc_church_proto_msgTypes,
	}.Build()
	File_rpc_church_proto = out.File
	file_rpc_church_proto_rawDesc = nil
	file_rpc_church_proto_goTypes = nil
	file_rpc_church_proto_depIdxs = nil
}
// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v3.19.6
// source: church.proto

package pb

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// A message for Church data
type Church struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	AuthId          string    `protobuf:"bytes,1,opt,name=auth_id,json=authId,proto3" json:"auth_id,omitempty"`
	Id              int32     `protobuf:"varint,2,opt,name=id,proto3" json:"id,omitempty"`
	Email           string    `protobuf:"bytes,3,opt,name=email,proto3" json:"email,omitempty"`
	Username        string    `protobuf:"bytes,4,opt,name=username,proto3" json:"username,omitempty"`
	Phone           string    `protobuf:"bytes,5,opt,name=phone,proto3" json:"phone,omitempty"`
	Password        string    `protobuf:"bytes,6,opt,name=password,proto3" json:"password,omitempty"`
	DenominationId  int32     `protobuf:"varint,7,opt,name=denomination_id,json=denominationId,proto3" json:"denomination_id,omitempty"`
	Name            string    `protobuf:"bytes,8,opt,name=name,proto3" json:"name,omitempty"`
	Location        *Location `protobuf:"bytes,9,opt,name=location,proto3" json:"location,omitempty"`
	ImageUrl        string    `protobuf:"bytes,10,opt,name=image_url,json=imageUrl,proto3" json:"image_url,omitempty"`
	FollowingCount  int64     `protobuf:"varint,11,opt,name=following_count,json=followingCount,proto3" json:"following_count,omitempty"`
	FollowerCount   int64     `protobuf:"varint,12,opt,name=follower_count,json=followerCount,proto3" json:"follower_count,omitempty"`
	MembershipCount int64     `protobuf:"varint,13,opt,name=membership_count,json=membershipCount,proto3" json:"membership_count,omitempty"`
	IsVerified      bool      `protobuf:"varint,14,opt,name=is_verified,json=isVerified,proto3" json:"is_verified,omitempty"`
	PostCount       int64     `protobuf:"varint,15,opt,name=post_count,json=postCount,proto3" json:"post_count,omitempty"`
	About           string    `protobuf:"bytes,16,opt,name=about,proto3" json:"about,omitempty"`
	Website         string    `protobuf:"bytes,17,opt,name=website,proto3" json:"website,omitempty"`
	HeaderImageUrl  string    `protobuf:"bytes,18,opt,name=header_image_url,json=headerImageUrl,proto3" json:"header_image_url,omitempty"`
	AccountName     string    `protobuf:"bytes,19,opt,name=account_name,json=accountName,proto3" json:"account_name,omitempty"`
	AccountNumber   string    `protobuf:"bytes,20,opt,name=account_number,json=accountNumber,proto3" json:"account_number,omitempty"`
	BankName        string    `protobuf:"bytes,21,opt,name=bank_name,json=bankName,proto3" json:"bank_name,omitempty"`
}

func (x *Church) Reset() {
	*x = Church{}
	if protoimpl.UnsafeEnabled {
		mi := &file_church_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Church) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Church) ProtoMessage() {}

func (x *Church) ProtoReflect() protoreflect.Message {
	mi := &file_church_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Church.ProtoReflect.Descriptor instead.
func (*Church) Descriptor() ([]byte, []int) {
	return file_church_proto_rawDescGZIP(), []int{0}
}

func (x *Church) GetAuthId() string {
	if x != nil {
		return x.AuthId
	}
	return ""
}

func (x *Church) GetId() int32 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *Church) GetEmail() string {
	if x != nil {
		return x.Email
	}
	return ""
}

func (x *Church) GetUsername() string {
	if x != nil {
		return x.Username
	}
	return ""
}

func (x *Church) GetPhone() string {
	if x != nil {
		return x.Phone
	}
	return ""
}

func (x *Church) GetPassword() string {
	if x != nil {
		return x.Password
	}
	return ""
}

func (x *Church) GetDenominationId() int32 {
	if x != nil {
		return x.DenominationId
	}
	return 0
}

func (x *Church) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Church) GetLocation() *Location {
	if x != nil {
		return x.Location
	}
	return nil
}

func (x *Church) GetImageUrl() string {
	if x != nil {
		return x.ImageUrl
	}
	return ""
}

func (x *Church) GetFollowingCount() int64 {
	if x != nil {
		return x.FollowingCount
	}
	return 0
}

func (x *Church) GetFollowerCount() int64 {
	if x != nil {
		return x.FollowerCount
	}
	return 0
}

func (x *Church) GetMembershipCount() int64 {
	if x != nil {
		return x.MembershipCount
	}
	return 0
}

func (x *Church) GetIsVerified() bool {
	if x != nil {
		return x.IsVerified
	}
	return false
}

func (x *Church) GetPostCount() int64 {
	if x != nil {
		return x.PostCount
	}
	return 0
}

func (x *Church) GetAbout() string {
	if x != nil {
		return x.About
	}
	return ""
}

func (x *Church) GetWebsite() string {
	if x != nil {
		return x.Website
	}
	return ""
}

func (x *Church) GetHeaderImageUrl() string {
	if x != nil {
		return x.HeaderImageUrl
	}
	return ""
}

func (x *Church) GetAccountName() string {
	if x != nil {
		return x.AccountName
	}
	return ""
}

func (x *Church) GetAccountNumber() string {
	if x != nil {
		return x.AccountNumber
	}
	return ""
}

func (x *Church) GetBankName() string {
	if x != nil {
		return x.BankName
	}
	return ""
}

type CreateChurchResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Success bool   `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
	Message string `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
}

func (x *CreateChurchResponse) Reset() {
	*x = CreateChurchResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_church_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateChurchResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateChurchResponse) ProtoMessage() {}

func (x *CreateChurchResponse) ProtoReflect() protoreflect.Message {
	mi := &file_church_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateChurchResponse.ProtoReflect.Descriptor instead.
func (*CreateChurchResponse) Descriptor() ([]byte, []int) {
	return file_church_proto_rawDescGZIP(), []int{1}
}

func (x *CreateChurchResponse) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

func (x *CreateChurchResponse) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

type SearchChurchResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Church  []*Church `protobuf:"bytes,1,rep,name=church,proto3" json:"church,omitempty"`
	HasMore bool      `protobuf:"varint,2,opt,name=has_more,json=hasMore,proto3" json:"has_more,omitempty"`
}

func (x *SearchChurchResponse) Reset() {
	*x = SearchChurchResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_church_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SearchChurchResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SearchChurchResponse) ProtoMessage() {}

func (x *SearchChurchResponse) ProtoReflect() protoreflect.Message {
	mi := &file_church_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SearchChurchResponse.ProtoReflect.Descriptor instead.
func (*SearchChurchResponse) Descriptor() ([]byte, []int) {
	return file_church_proto_rawDescGZIP(), []int{2}
}

func (x *SearchChurchResponse) GetChurch() []*Church {
	if x != nil {
		return x.Church
	}
	return nil
}

func (x *SearchChurchResponse) GetHasMore() bool {
	if x != nil {
		return x.HasMore
	}
	return false
}

var File_church_proto protoreflect.FileDescriptor

var file_church_proto_rawDesc = []byte{
	0x0a, 0x0c, 0x63, 0x68, 0x75, 0x72, 0x63, 0x68, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x02,
	0x70, 0x62, 0x1a, 0x0e, 0x6c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x22, 0x95, 0x05, 0x0a, 0x06, 0x43, 0x68, 0x75, 0x72, 0x63, 0x68, 0x12, 0x17, 0x0a,
	0x07, 0x61, 0x75, 0x74, 0x68, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06,
	0x61, 0x75, 0x74, 0x68, 0x49, 0x64, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x05, 0x52, 0x02, 0x69, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x12, 0x1a, 0x0a, 0x08,
	0x75, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08,
	0x75, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x70, 0x68, 0x6f, 0x6e,
	0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x70, 0x68, 0x6f, 0x6e, 0x65, 0x12, 0x1a,
	0x0a, 0x08, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x08, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x12, 0x27, 0x0a, 0x0f, 0x64, 0x65,
	0x6e, 0x6f, 0x6d, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x69, 0x64, 0x18, 0x07, 0x20,
	0x01, 0x28, 0x05, 0x52, 0x0e, 0x64, 0x65, 0x6e, 0x6f, 0x6d, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x49, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x08, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x28, 0x0a, 0x08, 0x6c, 0x6f, 0x63, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x18, 0x09, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0c, 0x2e, 0x70, 0x62, 0x2e, 0x4c,
	0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x08, 0x6c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x12, 0x1b, 0x0a, 0x09, 0x69, 0x6d, 0x61, 0x67, 0x65, 0x5f, 0x75, 0x72, 0x6c, 0x18, 0x0a,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x69, 0x6d, 0x61, 0x67, 0x65, 0x55, 0x72, 0x6c, 0x12, 0x27,
	0x0a, 0x0f, 0x66, 0x6f, 0x6c, 0x6c, 0x6f, 0x77, 0x69, 0x6e, 0x67, 0x5f, 0x63, 0x6f, 0x75, 0x6e,
	0x74, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0e, 0x66, 0x6f, 0x6c, 0x6c, 0x6f, 0x77, 0x69,
	0x6e, 0x67, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x25, 0x0a, 0x0e, 0x66, 0x6f, 0x6c, 0x6c, 0x6f,
	0x77, 0x65, 0x72, 0x5f, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x0c, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x0d, 0x66, 0x6f, 0x6c, 0x6c, 0x6f, 0x77, 0x65, 0x72, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x29,
	0x0a, 0x10, 0x6d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x73, 0x68, 0x69, 0x70, 0x5f, 0x63, 0x6f, 0x75,
	0x6e, 0x74, 0x18, 0x0d, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0f, 0x6d, 0x65, 0x6d, 0x62, 0x65, 0x72,
	0x73, 0x68, 0x69, 0x70, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x1f, 0x0a, 0x0b, 0x69, 0x73, 0x5f,
	0x76, 0x65, 0x72, 0x69, 0x66, 0x69, 0x65, 0x64, 0x18, 0x0e, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0a,
	0x69, 0x73, 0x56, 0x65, 0x72, 0x69, 0x66, 0x69, 0x65, 0x64, 0x12, 0x1d, 0x0a, 0x0a, 0x70, 0x6f,
	0x73, 0x74, 0x5f, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x0f, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09,
	0x70, 0x6f, 0x73, 0x74, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x61, 0x62, 0x6f,
	0x75, 0x74, 0x18, 0x10, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x61, 0x62, 0x6f, 0x75, 0x74, 0x12,
	0x18, 0x0a, 0x07, 0x77, 0x65, 0x62, 0x73, 0x69, 0x74, 0x65, 0x18, 0x11, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x07, 0x77, 0x65, 0x62, 0x73, 0x69, 0x74, 0x65, 0x12, 0x28, 0x0a, 0x10, 0x68, 0x65, 0x61,
	0x64, 0x65, 0x72, 0x5f, 0x69, 0x6d, 0x61, 0x67, 0x65, 0x5f, 0x75, 0x72, 0x6c, 0x18, 0x12, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0e, 0x68, 0x65, 0x61, 0x64, 0x65, 0x72, 0x49, 0x6d, 0x61, 0x67, 0x65,
	0x55, 0x72, 0x6c, 0x12, 0x21, 0x0a, 0x0c, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x5f, 0x6e,
	0x61, 0x6d, 0x65, 0x18, 0x13, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x61, 0x63, 0x63, 0x6f, 0x75,
	0x6e, 0x74, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x25, 0x0a, 0x0e, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e,
	0x74, 0x5f, 0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x18, 0x14, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d,
	0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x12, 0x1b, 0x0a,
	0x09, 0x62, 0x61, 0x6e, 0x6b, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x15, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x08, 0x62, 0x61, 0x6e, 0x6b, 0x4e, 0x61, 0x6d, 0x65, 0x22, 0x4a, 0x0a, 0x14, 0x43, 0x72,
	0x65, 0x61, 0x74, 0x65, 0x43, 0x68, 0x75, 0x72, 0x63, 0x68, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x08, 0x52, 0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x12, 0x18, 0x0a, 0x07,
	0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x22, 0x55, 0x0a, 0x14, 0x53, 0x65, 0x61, 0x72, 0x63, 0x68,
	0x43, 0x68, 0x75, 0x72, 0x63, 0x68, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x22,
	0x0a, 0x06, 0x63, 0x68, 0x75, 0x72, 0x63, 0x68, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0a,
	0x2e, 0x70, 0x62, 0x2e, 0x43, 0x68, 0x75, 0x72, 0x63, 0x68, 0x52, 0x06, 0x63, 0x68, 0x75, 0x72,
	0x63, 0x68, 0x12, 0x19, 0x0a, 0x08, 0x68, 0x61, 0x73, 0x5f, 0x6d, 0x6f, 0x72, 0x65, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x68, 0x61, 0x73, 0x4d, 0x6f, 0x72, 0x65, 0x42, 0x28, 0x5a,
	0x26, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x73, 0x74, 0x65, 0x76,
	0x65, 0x2d, 0x6d, 0x69, 0x72, 0x2f, 0x64, 0x69, 0x69, 0x76, 0x69, 0x78, 0x5f, 0x62, 0x61, 0x63,
	0x6b, 0x65, 0x6e, 0x64, 0x2f, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_church_proto_rawDescOnce sync.Once
	file_church_proto_rawDescData = file_church_proto_rawDesc
)

func file_church_proto_rawDescGZIP() []byte {
	file_church_proto_rawDescOnce.Do(func() {
		file_church_proto_rawDescData = protoimpl.X.CompressGZIP(file_church_proto_rawDescData)
	})
	return file_church_proto_rawDescData
}

var file_church_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_church_proto_goTypes = []interface{}{
	(*Church)(nil),               // 0: pb.Church
	(*CreateChurchResponse)(nil), // 1: pb.CreateChurchResponse
	(*SearchChurchResponse)(nil), // 2: pb.SearchChurchResponse
	(*Location)(nil),             // 3: pb.Location
}
var file_church_proto_depIdxs = []int32{
	3, // 0: pb.Church.location:type_name -> pb.Location
	0, // 1: pb.SearchChurchResponse.church:type_name -> pb.Church
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_church_proto_init() }
func file_church_proto_init() {
	if File_church_proto != nil {
		return
	}
	file_location_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_church_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Church); i {
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
		file_church_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateChurchResponse); i {
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
		file_church_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SearchChurchResponse); i {
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
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_church_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_church_proto_goTypes,
		DependencyIndexes: file_church_proto_depIdxs,
		MessageInfos:      file_church_proto_msgTypes,
	}.Build()
	File_church_proto = out.File
	file_church_proto_rawDesc = nil
	file_church_proto_goTypes = nil
	file_church_proto_depIdxs = nil
}

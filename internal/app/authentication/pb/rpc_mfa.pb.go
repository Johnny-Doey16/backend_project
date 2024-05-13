// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v3.19.6
// source: rpc_mfa.proto

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

// Request message for MFA registration.
type RegisterMFARequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Password string `protobuf:"bytes,1,opt,name=password,proto3" json:"password,omitempty"`
}

func (x *RegisterMFARequest) Reset() {
	*x = RegisterMFARequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rpc_mfa_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RegisterMFARequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RegisterMFARequest) ProtoMessage() {}

func (x *RegisterMFARequest) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_mfa_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RegisterMFARequest.ProtoReflect.Descriptor instead.
func (*RegisterMFARequest) Descriptor() ([]byte, []int) {
	return file_rpc_mfa_proto_rawDescGZIP(), []int{0}
}

func (x *RegisterMFARequest) GetPassword() string {
	if x != nil {
		return x.Password
	}
	return ""
}

// Response message for MFA registration.
type RegisterMFAResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Secret string `protobuf:"bytes,1,opt,name=secret,proto3" json:"secret,omitempty"` // The secret key to be stored securely.
	// repeated string recovery_codes = 2;
	QrCode []byte  `protobuf:"bytes,2,opt,name=qrCode,proto3" json:"qrCode,omitempty"`
	Url    *string `protobuf:"bytes,3,opt,name=url,proto3,oneof" json:"url,omitempty"` // URL of the QR code to scan with the Authy app.
}

func (x *RegisterMFAResponse) Reset() {
	*x = RegisterMFAResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rpc_mfa_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RegisterMFAResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RegisterMFAResponse) ProtoMessage() {}

func (x *RegisterMFAResponse) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_mfa_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RegisterMFAResponse.ProtoReflect.Descriptor instead.
func (*RegisterMFAResponse) Descriptor() ([]byte, []int) {
	return file_rpc_mfa_proto_rawDescGZIP(), []int{1}
}

func (x *RegisterMFAResponse) GetSecret() string {
	if x != nil {
		return x.Secret
	}
	return ""
}

func (x *RegisterMFAResponse) GetQrCode() []byte {
	if x != nil {
		return x.QrCode
	}
	return nil
}

func (x *RegisterMFAResponse) GetUrl() string {
	if x != nil && x.Url != nil {
		return *x.Url
	}
	return ""
}

// Used to test the mfa for newly created mfa
type VerifyMFAWorksRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Secret string `protobuf:"bytes,1,opt,name=secret,proto3" json:"secret,omitempty"`
	Token  string `protobuf:"bytes,2,opt,name=token,proto3" json:"token,omitempty"`
}

func (x *VerifyMFAWorksRequest) Reset() {
	*x = VerifyMFAWorksRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rpc_mfa_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *VerifyMFAWorksRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*VerifyMFAWorksRequest) ProtoMessage() {}

func (x *VerifyMFAWorksRequest) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_mfa_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use VerifyMFAWorksRequest.ProtoReflect.Descriptor instead.
func (*VerifyMFAWorksRequest) Descriptor() ([]byte, []int) {
	return file_rpc_mfa_proto_rawDescGZIP(), []int{2}
}

func (x *VerifyMFAWorksRequest) GetSecret() string {
	if x != nil {
		return x.Secret
	}
	return ""
}

func (x *VerifyMFAWorksRequest) GetToken() string {
	if x != nil {
		return x.Token
	}
	return ""
}

// Response message for MFA registration.
type VerifyMFAWorksResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	RecoveryCodes []string `protobuf:"bytes,1,rep,name=recovery_codes,json=recoveryCodes,proto3" json:"recovery_codes,omitempty"`
}

func (x *VerifyMFAWorksResponse) Reset() {
	*x = VerifyMFAWorksResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rpc_mfa_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *VerifyMFAWorksResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*VerifyMFAWorksResponse) ProtoMessage() {}

func (x *VerifyMFAWorksResponse) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_mfa_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use VerifyMFAWorksResponse.ProtoReflect.Descriptor instead.
func (*VerifyMFAWorksResponse) Descriptor() ([]byte, []int) {
	return file_rpc_mfa_proto_rawDescGZIP(), []int{3}
}

func (x *VerifyMFAWorksResponse) GetRecoveryCodes() []string {
	if x != nil {
		return x.RecoveryCodes
	}
	return nil
}

// Request message for verifying a TOTP code.
type VerifyMFARequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId string `protobuf:"bytes,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	Token  string `protobuf:"bytes,2,opt,name=token,proto3" json:"token,omitempty"` // The TOTP code to verify.
}

func (x *VerifyMFARequest) Reset() {
	*x = VerifyMFARequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rpc_mfa_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *VerifyMFARequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*VerifyMFARequest) ProtoMessage() {}

func (x *VerifyMFARequest) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_mfa_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use VerifyMFARequest.ProtoReflect.Descriptor instead.
func (*VerifyMFARequest) Descriptor() ([]byte, []int) {
	return file_rpc_mfa_proto_rawDescGZIP(), []int{4}
}

func (x *VerifyMFARequest) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

func (x *VerifyMFARequest) GetToken() string {
	if x != nil {
		return x.Token
	}
	return ""
}

// Response message for verifying a TOTP code.
type VerifyMFAResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Success  bool           `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
	AuthInfo *AuthTokenInfo `protobuf:"bytes,2,opt,name=authInfo,proto3" json:"authInfo,omitempty"`
}

func (x *VerifyMFAResponse) Reset() {
	*x = VerifyMFAResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rpc_mfa_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *VerifyMFAResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*VerifyMFAResponse) ProtoMessage() {}

func (x *VerifyMFAResponse) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_mfa_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use VerifyMFAResponse.ProtoReflect.Descriptor instead.
func (*VerifyMFAResponse) Descriptor() ([]byte, []int) {
	return file_rpc_mfa_proto_rawDescGZIP(), []int{5}
}

func (x *VerifyMFAResponse) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

func (x *VerifyMFAResponse) GetAuthInfo() *AuthTokenInfo {
	if x != nil {
		return x.AuthInfo
	}
	return nil
}

// Request message for bypassing OTP.
type ByPassOtpRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId       string `protobuf:"bytes,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	RecoveryCode string `protobuf:"bytes,2,opt,name=recovery_code,json=recoveryCode,proto3" json:"recovery_code,omitempty"`
}

func (x *ByPassOtpRequest) Reset() {
	*x = ByPassOtpRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rpc_mfa_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ByPassOtpRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ByPassOtpRequest) ProtoMessage() {}

func (x *ByPassOtpRequest) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_mfa_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ByPassOtpRequest.ProtoReflect.Descriptor instead.
func (*ByPassOtpRequest) Descriptor() ([]byte, []int) {
	return file_rpc_mfa_proto_rawDescGZIP(), []int{6}
}

func (x *ByPassOtpRequest) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

func (x *ByPassOtpRequest) GetRecoveryCode() string {
	if x != nil {
		return x.RecoveryCode
	}
	return ""
}

// Response message for bypassing OTP.
type ByPassOtpResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Success  bool           `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
	AuthInfo *AuthTokenInfo `protobuf:"bytes,2,opt,name=authInfo,proto3" json:"authInfo,omitempty"`
}

func (x *ByPassOtpResponse) Reset() {
	*x = ByPassOtpResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rpc_mfa_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ByPassOtpResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ByPassOtpResponse) ProtoMessage() {}

func (x *ByPassOtpResponse) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_mfa_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ByPassOtpResponse.ProtoReflect.Descriptor instead.
func (*ByPassOtpResponse) Descriptor() ([]byte, []int) {
	return file_rpc_mfa_proto_rawDescGZIP(), []int{7}
}

func (x *ByPassOtpResponse) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

func (x *ByPassOtpResponse) GetAuthInfo() *AuthTokenInfo {
	if x != nil {
		return x.AuthInfo
	}
	return nil
}

var File_rpc_mfa_proto protoreflect.FileDescriptor

var file_rpc_mfa_proto_rawDesc = []byte{
	0x0a, 0x0d, 0x72, 0x70, 0x63, 0x5f, 0x6d, 0x66, 0x61, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x02, 0x70, 0x62, 0x1a, 0x15, 0x61, 0x75, 0x74, 0x68, 0x5f, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x5f,
	0x69, 0x6e, 0x66, 0x6f, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x30, 0x0a, 0x12, 0x52, 0x65,
	0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x4d, 0x46, 0x41, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x1a, 0x0a, 0x08, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x08, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x22, 0x64, 0x0a, 0x13,
	0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x4d, 0x46, 0x41, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x65, 0x63, 0x72, 0x65, 0x74, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x65, 0x63, 0x72, 0x65, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x71,
	0x72, 0x43, 0x6f, 0x64, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x06, 0x71, 0x72, 0x43,
	0x6f, 0x64, 0x65, 0x12, 0x15, 0x0a, 0x03, 0x75, 0x72, 0x6c, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09,
	0x48, 0x00, 0x52, 0x03, 0x75, 0x72, 0x6c, 0x88, 0x01, 0x01, 0x42, 0x06, 0x0a, 0x04, 0x5f, 0x75,
	0x72, 0x6c, 0x22, 0x45, 0x0a, 0x15, 0x56, 0x65, 0x72, 0x69, 0x66, 0x79, 0x4d, 0x46, 0x41, 0x57,
	0x6f, 0x72, 0x6b, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x73,
	0x65, 0x63, 0x72, 0x65, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x65, 0x63,
	0x72, 0x65, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x22, 0x3f, 0x0a, 0x16, 0x56, 0x65, 0x72,
	0x69, 0x66, 0x79, 0x4d, 0x46, 0x41, 0x57, 0x6f, 0x72, 0x6b, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x25, 0x0a, 0x0e, 0x72, 0x65, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x79, 0x5f,
	0x63, 0x6f, 0x64, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x09, 0x52, 0x0d, 0x72, 0x65, 0x63,
	0x6f, 0x76, 0x65, 0x72, 0x79, 0x43, 0x6f, 0x64, 0x65, 0x73, 0x22, 0x41, 0x0a, 0x10, 0x56, 0x65,
	0x72, 0x69, 0x66, 0x79, 0x4d, 0x46, 0x41, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x17,
	0x0a, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x22, 0x5c, 0x0a,
	0x11, 0x56, 0x65, 0x72, 0x69, 0x66, 0x79, 0x4d, 0x46, 0x41, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x08, 0x52, 0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x12, 0x2d, 0x0a, 0x08,
	0x61, 0x75, 0x74, 0x68, 0x49, 0x6e, 0x66, 0x6f, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x11,
	0x2e, 0x70, 0x62, 0x2e, 0x41, 0x75, 0x74, 0x68, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x49, 0x6e, 0x66,
	0x6f, 0x52, 0x08, 0x61, 0x75, 0x74, 0x68, 0x49, 0x6e, 0x66, 0x6f, 0x22, 0x50, 0x0a, 0x10, 0x42,
	0x79, 0x50, 0x61, 0x73, 0x73, 0x4f, 0x74, 0x70, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x17, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12, 0x23, 0x0a, 0x0d, 0x72, 0x65, 0x63, 0x6f,
	0x76, 0x65, 0x72, 0x79, 0x5f, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0c, 0x72, 0x65, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x79, 0x43, 0x6f, 0x64, 0x65, 0x22, 0x5c, 0x0a,
	0x11, 0x42, 0x79, 0x50, 0x61, 0x73, 0x73, 0x4f, 0x74, 0x70, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x08, 0x52, 0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x12, 0x2d, 0x0a, 0x08,
	0x61, 0x75, 0x74, 0x68, 0x49, 0x6e, 0x66, 0x6f, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x11,
	0x2e, 0x70, 0x62, 0x2e, 0x41, 0x75, 0x74, 0x68, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x49, 0x6e, 0x66,
	0x6f, 0x52, 0x08, 0x61, 0x75, 0x74, 0x68, 0x49, 0x6e, 0x66, 0x6f, 0x42, 0x28, 0x5a, 0x26, 0x67,
	0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x73, 0x74, 0x65, 0x76, 0x65, 0x2d,
	0x6d, 0x69, 0x72, 0x2f, 0x64, 0x69, 0x69, 0x76, 0x69, 0x78, 0x5f, 0x62, 0x61, 0x63, 0x6b, 0x65,
	0x6e, 0x64, 0x2f, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_rpc_mfa_proto_rawDescOnce sync.Once
	file_rpc_mfa_proto_rawDescData = file_rpc_mfa_proto_rawDesc
)

func file_rpc_mfa_proto_rawDescGZIP() []byte {
	file_rpc_mfa_proto_rawDescOnce.Do(func() {
		file_rpc_mfa_proto_rawDescData = protoimpl.X.CompressGZIP(file_rpc_mfa_proto_rawDescData)
	})
	return file_rpc_mfa_proto_rawDescData
}

var file_rpc_mfa_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_rpc_mfa_proto_goTypes = []interface{}{
	(*RegisterMFARequest)(nil),     // 0: pb.RegisterMFARequest
	(*RegisterMFAResponse)(nil),    // 1: pb.RegisterMFAResponse
	(*VerifyMFAWorksRequest)(nil),  // 2: pb.VerifyMFAWorksRequest
	(*VerifyMFAWorksResponse)(nil), // 3: pb.VerifyMFAWorksResponse
	(*VerifyMFARequest)(nil),       // 4: pb.VerifyMFARequest
	(*VerifyMFAResponse)(nil),      // 5: pb.VerifyMFAResponse
	(*ByPassOtpRequest)(nil),       // 6: pb.ByPassOtpRequest
	(*ByPassOtpResponse)(nil),      // 7: pb.ByPassOtpResponse
	(*AuthTokenInfo)(nil),          // 8: pb.AuthTokenInfo
}
var file_rpc_mfa_proto_depIdxs = []int32{
	8, // 0: pb.VerifyMFAResponse.authInfo:type_name -> pb.AuthTokenInfo
	8, // 1: pb.ByPassOtpResponse.authInfo:type_name -> pb.AuthTokenInfo
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_rpc_mfa_proto_init() }
func file_rpc_mfa_proto_init() {
	if File_rpc_mfa_proto != nil {
		return
	}
	file_auth_token_info_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_rpc_mfa_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RegisterMFARequest); i {
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
		file_rpc_mfa_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RegisterMFAResponse); i {
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
		file_rpc_mfa_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*VerifyMFAWorksRequest); i {
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
		file_rpc_mfa_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*VerifyMFAWorksResponse); i {
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
		file_rpc_mfa_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*VerifyMFARequest); i {
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
		file_rpc_mfa_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*VerifyMFAResponse); i {
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
		file_rpc_mfa_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ByPassOtpRequest); i {
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
		file_rpc_mfa_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ByPassOtpResponse); i {
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
	file_rpc_mfa_proto_msgTypes[1].OneofWrappers = []interface{}{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_rpc_mfa_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_rpc_mfa_proto_goTypes,
		DependencyIndexes: file_rpc_mfa_proto_depIdxs,
		MessageInfos:      file_rpc_mfa_proto_msgTypes,
	}.Build()
	File_rpc_mfa_proto = out.File
	file_rpc_mfa_proto_rawDesc = nil
	file_rpc_mfa_proto_goTypes = nil
	file_rpc_mfa_proto_depIdxs = nil
}

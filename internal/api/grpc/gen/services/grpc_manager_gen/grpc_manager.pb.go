// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.1
// 	protoc        v3.12.4
// source: internal/api/grpc/proto/grpc_manager.proto

package grpc_manager_gen

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

type GrpcManagerStopRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Message       string                 `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GrpcManagerStopRequest) Reset() {
	*x = GrpcManagerStopRequest{}
	mi := &file_internal_api_grpc_proto_grpc_manager_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GrpcManagerStopRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GrpcManagerStopRequest) ProtoMessage() {}

func (x *GrpcManagerStopRequest) ProtoReflect() protoreflect.Message {
	mi := &file_internal_api_grpc_proto_grpc_manager_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GrpcManagerStopRequest.ProtoReflect.Descriptor instead.
func (*GrpcManagerStopRequest) Descriptor() ([]byte, []int) {
	return file_internal_api_grpc_proto_grpc_manager_proto_rawDescGZIP(), []int{0}
}

func (x *GrpcManagerStopRequest) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

type GrpcManagerStopResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Success       bool                   `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
	Message       string                 `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GrpcManagerStopResponse) Reset() {
	*x = GrpcManagerStopResponse{}
	mi := &file_internal_api_grpc_proto_grpc_manager_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GrpcManagerStopResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GrpcManagerStopResponse) ProtoMessage() {}

func (x *GrpcManagerStopResponse) ProtoReflect() protoreflect.Message {
	mi := &file_internal_api_grpc_proto_grpc_manager_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GrpcManagerStopResponse.ProtoReflect.Descriptor instead.
func (*GrpcManagerStopResponse) Descriptor() ([]byte, []int) {
	return file_internal_api_grpc_proto_grpc_manager_proto_rawDescGZIP(), []int{1}
}

func (x *GrpcManagerStopResponse) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

func (x *GrpcManagerStopResponse) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

var File_internal_api_grpc_proto_grpc_manager_proto protoreflect.FileDescriptor

var file_internal_api_grpc_proto_grpc_manager_proto_rawDesc = []byte{
	0x0a, 0x2a, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x67,
	0x72, 0x70, 0x63, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x67, 0x72, 0x70, 0x63, 0x5f, 0x6d,
	0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x13, 0x73, 0x6c,
	0x6f, 0x67, 0x67, 0x65, 0x72, 0x5f, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x70, 0x6f, 0x72, 0x74, 0x65,
	0x72, 0x22, 0x32, 0x0a, 0x16, 0x47, 0x72, 0x70, 0x63, 0x4d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72,
	0x53, 0x74, 0x6f, 0x70, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x18, 0x0a, 0x07, 0x6d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x22, 0x4d, 0x0a, 0x17, 0x47, 0x72, 0x70, 0x63, 0x4d, 0x61, 0x6e,
	0x61, 0x67, 0x65, 0x72, 0x53, 0x74, 0x6f, 0x70, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x18, 0x0a, 0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x08, 0x52, 0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x32, 0x72, 0x0a, 0x0b, 0x47, 0x72, 0x70, 0x63, 0x4d, 0x61, 0x6e, 0x61,
	0x67, 0x65, 0x72, 0x12, 0x63, 0x0a, 0x04, 0x53, 0x74, 0x6f, 0x70, 0x12, 0x2b, 0x2e, 0x73, 0x6c,
	0x6f, 0x67, 0x67, 0x65, 0x72, 0x5f, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x70, 0x6f, 0x72, 0x74, 0x65,
	0x72, 0x2e, 0x47, 0x72, 0x70, 0x63, 0x4d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x53, 0x74, 0x6f,
	0x70, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x2c, 0x2e, 0x73, 0x6c, 0x6f, 0x67, 0x67,
	0x65, 0x72, 0x5f, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x70, 0x6f, 0x72, 0x74, 0x65, 0x72, 0x2e, 0x47,
	0x72, 0x70, 0x63, 0x4d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x53, 0x74, 0x6f, 0x70, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x42, 0x1b, 0x5a, 0x19, 0x73, 0x65, 0x72, 0x76,
	0x69, 0x63, 0x65, 0x73, 0x2f, 0x67, 0x72, 0x70, 0x63, 0x5f, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65,
	0x72, 0x5f, 0x67, 0x65, 0x6e, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_internal_api_grpc_proto_grpc_manager_proto_rawDescOnce sync.Once
	file_internal_api_grpc_proto_grpc_manager_proto_rawDescData = file_internal_api_grpc_proto_grpc_manager_proto_rawDesc
)

func file_internal_api_grpc_proto_grpc_manager_proto_rawDescGZIP() []byte {
	file_internal_api_grpc_proto_grpc_manager_proto_rawDescOnce.Do(func() {
		file_internal_api_grpc_proto_grpc_manager_proto_rawDescData = protoimpl.X.CompressGZIP(file_internal_api_grpc_proto_grpc_manager_proto_rawDescData)
	})
	return file_internal_api_grpc_proto_grpc_manager_proto_rawDescData
}

var file_internal_api_grpc_proto_grpc_manager_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_internal_api_grpc_proto_grpc_manager_proto_goTypes = []any{
	(*GrpcManagerStopRequest)(nil),  // 0: slogger_transporter.GrpcManagerStopRequest
	(*GrpcManagerStopResponse)(nil), // 1: slogger_transporter.GrpcManagerStopResponse
}
var file_internal_api_grpc_proto_grpc_manager_proto_depIdxs = []int32{
	0, // 0: slogger_transporter.GrpcManager.Stop:input_type -> slogger_transporter.GrpcManagerStopRequest
	1, // 1: slogger_transporter.GrpcManager.Stop:output_type -> slogger_transporter.GrpcManagerStopResponse
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_internal_api_grpc_proto_grpc_manager_proto_init() }
func file_internal_api_grpc_proto_grpc_manager_proto_init() {
	if File_internal_api_grpc_proto_grpc_manager_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_internal_api_grpc_proto_grpc_manager_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_internal_api_grpc_proto_grpc_manager_proto_goTypes,
		DependencyIndexes: file_internal_api_grpc_proto_grpc_manager_proto_depIdxs,
		MessageInfos:      file_internal_api_grpc_proto_grpc_manager_proto_msgTypes,
	}.Build()
	File_internal_api_grpc_proto_grpc_manager_proto = out.File
	file_internal_api_grpc_proto_grpc_manager_proto_rawDesc = nil
	file_internal_api_grpc_proto_grpc_manager_proto_goTypes = nil
	file_internal_api_grpc_proto_grpc_manager_proto_depIdxs = nil
}

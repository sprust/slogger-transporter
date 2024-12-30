// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.1
// 	protoc        v3.12.4
// source: internal/api/grpc/proto/trace_transporter.proto

package trace_transporter_gen

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

type TraceTransporterCreateRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Payload       string                 `protobuf:"bytes,1,opt,name=payload,proto3" json:"payload,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *TraceTransporterCreateRequest) Reset() {
	*x = TraceTransporterCreateRequest{}
	mi := &file_internal_api_grpc_proto_trace_transporter_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *TraceTransporterCreateRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TraceTransporterCreateRequest) ProtoMessage() {}

func (x *TraceTransporterCreateRequest) ProtoReflect() protoreflect.Message {
	mi := &file_internal_api_grpc_proto_trace_transporter_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TraceTransporterCreateRequest.ProtoReflect.Descriptor instead.
func (*TraceTransporterCreateRequest) Descriptor() ([]byte, []int) {
	return file_internal_api_grpc_proto_trace_transporter_proto_rawDescGZIP(), []int{0}
}

func (x *TraceTransporterCreateRequest) GetPayload() string {
	if x != nil {
		return x.Payload
	}
	return ""
}

type TraceTransporterUpdateRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Payload       string                 `protobuf:"bytes,1,opt,name=payload,proto3" json:"payload,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *TraceTransporterUpdateRequest) Reset() {
	*x = TraceTransporterUpdateRequest{}
	mi := &file_internal_api_grpc_proto_trace_transporter_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *TraceTransporterUpdateRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TraceTransporterUpdateRequest) ProtoMessage() {}

func (x *TraceTransporterUpdateRequest) ProtoReflect() protoreflect.Message {
	mi := &file_internal_api_grpc_proto_trace_transporter_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TraceTransporterUpdateRequest.ProtoReflect.Descriptor instead.
func (*TraceTransporterUpdateRequest) Descriptor() ([]byte, []int) {
	return file_internal_api_grpc_proto_trace_transporter_proto_rawDescGZIP(), []int{1}
}

func (x *TraceTransporterUpdateRequest) GetPayload() string {
	if x != nil {
		return x.Payload
	}
	return ""
}

type TraceTransporterResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Success       bool                   `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *TraceTransporterResponse) Reset() {
	*x = TraceTransporterResponse{}
	mi := &file_internal_api_grpc_proto_trace_transporter_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *TraceTransporterResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TraceTransporterResponse) ProtoMessage() {}

func (x *TraceTransporterResponse) ProtoReflect() protoreflect.Message {
	mi := &file_internal_api_grpc_proto_trace_transporter_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TraceTransporterResponse.ProtoReflect.Descriptor instead.
func (*TraceTransporterResponse) Descriptor() ([]byte, []int) {
	return file_internal_api_grpc_proto_trace_transporter_proto_rawDescGZIP(), []int{2}
}

func (x *TraceTransporterResponse) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

var File_internal_api_grpc_proto_trace_transporter_proto protoreflect.FileDescriptor

var file_internal_api_grpc_proto_trace_transporter_proto_rawDesc = []byte{
	0x0a, 0x2f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x67,
	0x72, 0x70, 0x63, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x74, 0x72, 0x61, 0x63, 0x65, 0x5f,
	0x74, 0x72, 0x61, 0x6e, 0x73, 0x70, 0x6f, 0x72, 0x74, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x13, 0x73, 0x6c, 0x6f, 0x67, 0x67, 0x65, 0x72, 0x5f, 0x74, 0x72, 0x61, 0x6e, 0x73,
	0x70, 0x6f, 0x72, 0x74, 0x65, 0x72, 0x22, 0x39, 0x0a, 0x1d, 0x54, 0x72, 0x61, 0x63, 0x65, 0x54,
	0x72, 0x61, 0x6e, 0x73, 0x70, 0x6f, 0x72, 0x74, 0x65, 0x72, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x18, 0x0a, 0x07, 0x70, 0x61, 0x79, 0x6c, 0x6f,
	0x61, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61,
	0x64, 0x22, 0x39, 0x0a, 0x1d, 0x54, 0x72, 0x61, 0x63, 0x65, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x70,
	0x6f, 0x72, 0x74, 0x65, 0x72, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x18, 0x0a, 0x07, 0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x07, 0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x22, 0x34, 0x0a, 0x18,
	0x54, 0x72, 0x61, 0x63, 0x65, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x70, 0x6f, 0x72, 0x74, 0x65, 0x72,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x73, 0x75, 0x63, 0x63,
	0x65, 0x73, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x73, 0x75, 0x63, 0x63, 0x65,
	0x73, 0x73, 0x32, 0xf0, 0x01, 0x0a, 0x10, 0x54, 0x72, 0x61, 0x63, 0x65, 0x54, 0x72, 0x61, 0x6e,
	0x73, 0x70, 0x6f, 0x72, 0x74, 0x65, 0x72, 0x12, 0x6d, 0x0a, 0x06, 0x43, 0x72, 0x65, 0x61, 0x74,
	0x65, 0x12, 0x32, 0x2e, 0x73, 0x6c, 0x6f, 0x67, 0x67, 0x65, 0x72, 0x5f, 0x74, 0x72, 0x61, 0x6e,
	0x73, 0x70, 0x6f, 0x72, 0x74, 0x65, 0x72, 0x2e, 0x54, 0x72, 0x61, 0x63, 0x65, 0x54, 0x72, 0x61,
	0x6e, 0x73, 0x70, 0x6f, 0x72, 0x74, 0x65, 0x72, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x2d, 0x2e, 0x73, 0x6c, 0x6f, 0x67, 0x67, 0x65, 0x72, 0x5f,
	0x74, 0x72, 0x61, 0x6e, 0x73, 0x70, 0x6f, 0x72, 0x74, 0x65, 0x72, 0x2e, 0x54, 0x72, 0x61, 0x63,
	0x65, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x70, 0x6f, 0x72, 0x74, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x6d, 0x0a, 0x06, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65,
	0x12, 0x32, 0x2e, 0x73, 0x6c, 0x6f, 0x67, 0x67, 0x65, 0x72, 0x5f, 0x74, 0x72, 0x61, 0x6e, 0x73,
	0x70, 0x6f, 0x72, 0x74, 0x65, 0x72, 0x2e, 0x54, 0x72, 0x61, 0x63, 0x65, 0x54, 0x72, 0x61, 0x6e,
	0x73, 0x70, 0x6f, 0x72, 0x74, 0x65, 0x72, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x2d, 0x2e, 0x73, 0x6c, 0x6f, 0x67, 0x67, 0x65, 0x72, 0x5f, 0x74,
	0x72, 0x61, 0x6e, 0x73, 0x70, 0x6f, 0x72, 0x74, 0x65, 0x72, 0x2e, 0x54, 0x72, 0x61, 0x63, 0x65,
	0x54, 0x72, 0x61, 0x6e, 0x73, 0x70, 0x6f, 0x72, 0x74, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x22, 0x00, 0x42, 0x20, 0x5a, 0x1e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x73, 0x2f, 0x74, 0x72, 0x61, 0x63, 0x65, 0x5f, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x70, 0x6f, 0x72,
	0x74, 0x65, 0x72, 0x5f, 0x67, 0x65, 0x6e, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_internal_api_grpc_proto_trace_transporter_proto_rawDescOnce sync.Once
	file_internal_api_grpc_proto_trace_transporter_proto_rawDescData = file_internal_api_grpc_proto_trace_transporter_proto_rawDesc
)

func file_internal_api_grpc_proto_trace_transporter_proto_rawDescGZIP() []byte {
	file_internal_api_grpc_proto_trace_transporter_proto_rawDescOnce.Do(func() {
		file_internal_api_grpc_proto_trace_transporter_proto_rawDescData = protoimpl.X.CompressGZIP(file_internal_api_grpc_proto_trace_transporter_proto_rawDescData)
	})
	return file_internal_api_grpc_proto_trace_transporter_proto_rawDescData
}

var file_internal_api_grpc_proto_trace_transporter_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_internal_api_grpc_proto_trace_transporter_proto_goTypes = []any{
	(*TraceTransporterCreateRequest)(nil), // 0: slogger_transporter.TraceTransporterCreateRequest
	(*TraceTransporterUpdateRequest)(nil), // 1: slogger_transporter.TraceTransporterUpdateRequest
	(*TraceTransporterResponse)(nil),      // 2: slogger_transporter.TraceTransporterResponse
}
var file_internal_api_grpc_proto_trace_transporter_proto_depIdxs = []int32{
	0, // 0: slogger_transporter.TraceTransporter.Create:input_type -> slogger_transporter.TraceTransporterCreateRequest
	1, // 1: slogger_transporter.TraceTransporter.Update:input_type -> slogger_transporter.TraceTransporterUpdateRequest
	2, // 2: slogger_transporter.TraceTransporter.Create:output_type -> slogger_transporter.TraceTransporterResponse
	2, // 3: slogger_transporter.TraceTransporter.Update:output_type -> slogger_transporter.TraceTransporterResponse
	2, // [2:4] is the sub-list for method output_type
	0, // [0:2] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_internal_api_grpc_proto_trace_transporter_proto_init() }
func file_internal_api_grpc_proto_trace_transporter_proto_init() {
	if File_internal_api_grpc_proto_trace_transporter_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_internal_api_grpc_proto_trace_transporter_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_internal_api_grpc_proto_trace_transporter_proto_goTypes,
		DependencyIndexes: file_internal_api_grpc_proto_trace_transporter_proto_depIdxs,
		MessageInfos:      file_internal_api_grpc_proto_trace_transporter_proto_msgTypes,
	}.Build()
	File_internal_api_grpc_proto_trace_transporter_proto = out.File
	file_internal_api_grpc_proto_trace_transporter_proto_rawDesc = nil
	file_internal_api_grpc_proto_trace_transporter_proto_goTypes = nil
	file_internal_api_grpc_proto_trace_transporter_proto_depIdxs = nil
}

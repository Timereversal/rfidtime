// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.33.0
// 	protoc        v3.21.12
// source: reader.proto

package reader

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

type ReportRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	TagId      int32                  `protobuf:"varint,1,opt,name=tagId,proto3" json:"tagId,omitempty"`
	EventId    int32                  `protobuf:"varint,2,opt,name=eventId,proto3" json:"eventId,omitempty"`
	RunnerTime *timestamppb.Timestamp `protobuf:"bytes,3,opt,name=runnerTime,proto3" json:"runnerTime,omitempty"`
	Stage      int32                  `protobuf:"varint,4,opt,name=stage,proto3" json:"stage,omitempty"`
}

func (x *ReportRequest) Reset() {
	*x = ReportRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_reader_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReportRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReportRequest) ProtoMessage() {}

func (x *ReportRequest) ProtoReflect() protoreflect.Message {
	mi := &file_reader_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReportRequest.ProtoReflect.Descriptor instead.
func (*ReportRequest) Descriptor() ([]byte, []int) {
	return file_reader_proto_rawDescGZIP(), []int{0}
}

func (x *ReportRequest) GetTagId() int32 {
	if x != nil {
		return x.TagId
	}
	return 0
}

func (x *ReportRequest) GetEventId() int32 {
	if x != nil {
		return x.EventId
	}
	return 0
}

func (x *ReportRequest) GetRunnerTime() *timestamppb.Timestamp {
	if x != nil {
		return x.RunnerTime
	}
	return nil
}

func (x *ReportRequest) GetStage() int32 {
	if x != nil {
		return x.Stage
	}
	return 0
}

type ReportResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Status bool `protobuf:"varint,1,opt,name=status,proto3" json:"status,omitempty"`
}

func (x *ReportResponse) Reset() {
	*x = ReportResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_reader_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReportResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReportResponse) ProtoMessage() {}

func (x *ReportResponse) ProtoReflect() protoreflect.Message {
	mi := &file_reader_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReportResponse.ProtoReflect.Descriptor instead.
func (*ReportResponse) Descriptor() ([]byte, []int) {
	return file_reader_proto_rawDescGZIP(), []int{1}
}

func (x *ReportResponse) GetStatus() bool {
	if x != nil {
		return x.Status
	}
	return false
}

var File_reader_proto protoreflect.FileDescriptor

var file_reader_proto_rawDesc = []byte{
	0x0a, 0x0c, 0x72, 0x65, 0x61, 0x64, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1f,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f,
	0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22,
	0x91, 0x01, 0x0a, 0x0d, 0x52, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x61, 0x67, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05,
	0x52, 0x05, 0x74, 0x61, 0x67, 0x49, 0x64, 0x12, 0x18, 0x0a, 0x07, 0x65, 0x76, 0x65, 0x6e, 0x74,
	0x49, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x07, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x49,
	0x64, 0x12, 0x3a, 0x0a, 0x0a, 0x72, 0x75, 0x6e, 0x6e, 0x65, 0x72, 0x54, 0x69, 0x6d, 0x65, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d,
	0x70, 0x52, 0x0a, 0x72, 0x75, 0x6e, 0x6e, 0x65, 0x72, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x14, 0x0a,
	0x05, 0x73, 0x74, 0x61, 0x67, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x05, 0x52, 0x05, 0x73, 0x74,
	0x61, 0x67, 0x65, 0x22, 0x28, 0x0a, 0x0e, 0x52, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x32, 0x33, 0x0a,
	0x06, 0x52, 0x65, 0x61, 0x64, 0x65, 0x72, 0x12, 0x29, 0x0a, 0x06, 0x52, 0x65, 0x70, 0x6f, 0x72,
	0x74, 0x12, 0x0e, 0x2e, 0x52, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x0f, 0x2e, 0x52, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x42, 0x2b, 0x5a, 0x29, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d,
	0x2f, 0x54, 0x69, 0x6d, 0x65, 0x72, 0x65, 0x76, 0x65, 0x72, 0x73, 0x61, 0x6c, 0x2f, 0x72, 0x66,
	0x69, 0x64, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2f, 0x72, 0x65, 0x61, 0x64, 0x65, 0x72, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_reader_proto_rawDescOnce sync.Once
	file_reader_proto_rawDescData = file_reader_proto_rawDesc
)

func file_reader_proto_rawDescGZIP() []byte {
	file_reader_proto_rawDescOnce.Do(func() {
		file_reader_proto_rawDescData = protoimpl.X.CompressGZIP(file_reader_proto_rawDescData)
	})
	return file_reader_proto_rawDescData
}

var file_reader_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_reader_proto_goTypes = []interface{}{
	(*ReportRequest)(nil),         // 0: ReportRequest
	(*ReportResponse)(nil),        // 1: ReportResponse
	(*timestamppb.Timestamp)(nil), // 2: google.protobuf.Timestamp
}
var file_reader_proto_depIdxs = []int32{
	2, // 0: ReportRequest.runnerTime:type_name -> google.protobuf.Timestamp
	0, // 1: Reader.Report:input_type -> ReportRequest
	1, // 2: Reader.Report:output_type -> ReportResponse
	2, // [2:3] is the sub-list for method output_type
	1, // [1:2] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_reader_proto_init() }
func file_reader_proto_init() {
	if File_reader_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_reader_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ReportRequest); i {
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
		file_reader_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ReportResponse); i {
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
			RawDescriptor: file_reader_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_reader_proto_goTypes,
		DependencyIndexes: file_reader_proto_depIdxs,
		MessageInfos:      file_reader_proto_msgTypes,
	}.Build()
	File_reader_proto = out.File
	file_reader_proto_rawDesc = nil
	file_reader_proto_goTypes = nil
	file_reader_proto_depIdxs = nil
}

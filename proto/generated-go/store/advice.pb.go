// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.0
// 	protoc        (unknown)
// source: store/advice.proto

package store

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

type Advice_Status int32

const (
	// Unspecified.
	Advice_STATUS_UNSPECIFIED Advice_Status = 0
	Advice_SUCCESS            Advice_Status = 1
	Advice_WARNING            Advice_Status = 2
	Advice_ERROR              Advice_Status = 3
)

// Enum value maps for Advice_Status.
var (
	Advice_Status_name = map[int32]string{
		0: "STATUS_UNSPECIFIED",
		1: "SUCCESS",
		2: "WARNING",
		3: "ERROR",
	}
	Advice_Status_value = map[string]int32{
		"STATUS_UNSPECIFIED": 0,
		"SUCCESS":            1,
		"WARNING":            2,
		"ERROR":              3,
	}
)

func (x Advice_Status) Enum() *Advice_Status {
	p := new(Advice_Status)
	*p = x
	return p
}

func (x Advice_Status) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Advice_Status) Descriptor() protoreflect.EnumDescriptor {
	return file_store_advice_proto_enumTypes[0].Descriptor()
}

func (Advice_Status) Type() protoreflect.EnumType {
	return &file_store_advice_proto_enumTypes[0]
}

func (x Advice_Status) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Advice_Status.Descriptor instead.
func (Advice_Status) EnumDescriptor() ([]byte, []int) {
	return file_store_advice_proto_rawDescGZIP(), []int{0, 0}
}

type Advice struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// The advice status.
	Status Advice_Status `protobuf:"varint,1,opt,name=status,proto3,enum=bytebase.store.Advice_Status" json:"status,omitempty"`
	// The advice code.
	Code int32 `protobuf:"varint,2,opt,name=code,proto3" json:"code,omitempty"`
	// The advice title.
	Title string `protobuf:"bytes,3,opt,name=title,proto3" json:"title,omitempty"`
	// The advice content.
	Content string `protobuf:"bytes,4,opt,name=content,proto3" json:"content,omitempty"`
	// The advice detail.
	Detail string `protobuf:"bytes,5,opt,name=detail,proto3" json:"detail,omitempty"`
	// 1-based positions of the sql statment.
	StartPosition *Position `protobuf:"bytes,6,opt,name=start_position,json=startPosition,proto3" json:"start_position,omitempty"`
	EndPosition   *Position `protobuf:"bytes,7,opt,name=end_position,json=endPosition,proto3" json:"end_position,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Advice) Reset() {
	*x = Advice{}
	mi := &file_store_advice_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Advice) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Advice) ProtoMessage() {}

func (x *Advice) ProtoReflect() protoreflect.Message {
	mi := &file_store_advice_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Advice.ProtoReflect.Descriptor instead.
func (*Advice) Descriptor() ([]byte, []int) {
	return file_store_advice_proto_rawDescGZIP(), []int{0}
}

func (x *Advice) GetStatus() Advice_Status {
	if x != nil {
		return x.Status
	}
	return Advice_STATUS_UNSPECIFIED
}

func (x *Advice) GetCode() int32 {
	if x != nil {
		return x.Code
	}
	return 0
}

func (x *Advice) GetTitle() string {
	if x != nil {
		return x.Title
	}
	return ""
}

func (x *Advice) GetContent() string {
	if x != nil {
		return x.Content
	}
	return ""
}

func (x *Advice) GetDetail() string {
	if x != nil {
		return x.Detail
	}
	return ""
}

func (x *Advice) GetStartPosition() *Position {
	if x != nil {
		return x.StartPosition
	}
	return nil
}

func (x *Advice) GetEndPosition() *Position {
	if x != nil {
		return x.EndPosition
	}
	return nil
}

var File_store_advice_proto protoreflect.FileDescriptor

var file_store_advice_proto_rawDesc = []byte{
	0x0a, 0x12, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2f, 0x61, 0x64, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0e, 0x62, 0x79, 0x74, 0x65, 0x62, 0x61, 0x73, 0x65, 0x2e, 0x73,
	0x74, 0x6f, 0x72, 0x65, 0x1a, 0x12, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2f, 0x63, 0x6f, 0x6d, 0x6d,
	0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xe0, 0x02, 0x0a, 0x06, 0x41, 0x64, 0x76,
	0x69, 0x63, 0x65, 0x12, 0x35, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0e, 0x32, 0x1d, 0x2e, 0x62, 0x79, 0x74, 0x65, 0x62, 0x61, 0x73, 0x65, 0x2e, 0x73,
	0x74, 0x6f, 0x72, 0x65, 0x2e, 0x41, 0x64, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x53, 0x74, 0x61, 0x74,
	0x75, 0x73, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x63, 0x6f,
	0x64, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x12, 0x14,
	0x0a, 0x05, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74,
	0x69, 0x74, 0x6c, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x12, 0x16,
	0x0a, 0x06, 0x64, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06,
	0x64, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x12, 0x3f, 0x0a, 0x0e, 0x73, 0x74, 0x61, 0x72, 0x74, 0x5f,
	0x70, 0x6f, 0x73, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x18,
	0x2e, 0x62, 0x79, 0x74, 0x65, 0x62, 0x61, 0x73, 0x65, 0x2e, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2e,
	0x50, 0x6f, 0x73, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x0d, 0x73, 0x74, 0x61, 0x72, 0x74, 0x50,
	0x6f, 0x73, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x3b, 0x0a, 0x0c, 0x65, 0x6e, 0x64, 0x5f, 0x70,
	0x6f, 0x73, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x07, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x18, 0x2e,
	0x62, 0x79, 0x74, 0x65, 0x62, 0x61, 0x73, 0x65, 0x2e, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2e, 0x50,
	0x6f, 0x73, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x0b, 0x65, 0x6e, 0x64, 0x50, 0x6f, 0x73, 0x69,
	0x74, 0x69, 0x6f, 0x6e, 0x22, 0x45, 0x0a, 0x06, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x16,
	0x0a, 0x12, 0x53, 0x54, 0x41, 0x54, 0x55, 0x53, 0x5f, 0x55, 0x4e, 0x53, 0x50, 0x45, 0x43, 0x49,
	0x46, 0x49, 0x45, 0x44, 0x10, 0x00, 0x12, 0x0b, 0x0a, 0x07, 0x53, 0x55, 0x43, 0x43, 0x45, 0x53,
	0x53, 0x10, 0x01, 0x12, 0x0b, 0x0a, 0x07, 0x57, 0x41, 0x52, 0x4e, 0x49, 0x4e, 0x47, 0x10, 0x02,
	0x12, 0x09, 0x0a, 0x05, 0x45, 0x52, 0x52, 0x4f, 0x52, 0x10, 0x03, 0x42, 0x14, 0x5a, 0x12, 0x67,
	0x65, 0x6e, 0x65, 0x72, 0x61, 0x74, 0x65, 0x64, 0x2d, 0x67, 0x6f, 0x2f, 0x73, 0x74, 0x6f, 0x72,
	0x65, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_store_advice_proto_rawDescOnce sync.Once
	file_store_advice_proto_rawDescData = file_store_advice_proto_rawDesc
)

func file_store_advice_proto_rawDescGZIP() []byte {
	file_store_advice_proto_rawDescOnce.Do(func() {
		file_store_advice_proto_rawDescData = protoimpl.X.CompressGZIP(file_store_advice_proto_rawDescData)
	})
	return file_store_advice_proto_rawDescData
}

var file_store_advice_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_store_advice_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_store_advice_proto_goTypes = []any{
	(Advice_Status)(0), // 0: bytebase.store.Advice.Status
	(*Advice)(nil),     // 1: bytebase.store.Advice
	(*Position)(nil),   // 2: bytebase.store.Position
}
var file_store_advice_proto_depIdxs = []int32{
	0, // 0: bytebase.store.Advice.status:type_name -> bytebase.store.Advice.Status
	2, // 1: bytebase.store.Advice.start_position:type_name -> bytebase.store.Position
	2, // 2: bytebase.store.Advice.end_position:type_name -> bytebase.store.Position
	3, // [3:3] is the sub-list for method output_type
	3, // [3:3] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_store_advice_proto_init() }
func file_store_advice_proto_init() {
	if File_store_advice_proto != nil {
		return
	}
	file_store_common_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_store_advice_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_store_advice_proto_goTypes,
		DependencyIndexes: file_store_advice_proto_depIdxs,
		EnumInfos:         file_store_advice_proto_enumTypes,
		MessageInfos:      file_store_advice_proto_msgTypes,
	}.Build()
	File_store_advice_proto = out.File
	file_store_advice_proto_rawDesc = nil
	file_store_advice_proto_goTypes = nil
	file_store_advice_proto_depIdxs = nil
}

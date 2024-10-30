// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.1
// 	protoc        (unknown)
// source: store/release.proto

package store

import (
	_ "google.golang.org/genproto/googleapis/api/annotations"
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

type ReleaseFileType int32

const (
	ReleaseFileType_TYPE_UNSPECIFIED ReleaseFileType = 0
	ReleaseFileType_VERSIONED        ReleaseFileType = 1
)

// Enum value maps for ReleaseFileType.
var (
	ReleaseFileType_name = map[int32]string{
		0: "TYPE_UNSPECIFIED",
		1: "VERSIONED",
	}
	ReleaseFileType_value = map[string]int32{
		"TYPE_UNSPECIFIED": 0,
		"VERSIONED":        1,
	}
)

func (x ReleaseFileType) Enum() *ReleaseFileType {
	p := new(ReleaseFileType)
	*p = x
	return p
}

func (x ReleaseFileType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ReleaseFileType) Descriptor() protoreflect.EnumDescriptor {
	return file_store_release_proto_enumTypes[0].Descriptor()
}

func (ReleaseFileType) Type() protoreflect.EnumType {
	return &file_store_release_proto_enumTypes[0]
}

func (x ReleaseFileType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ReleaseFileType.Descriptor instead.
func (ReleaseFileType) EnumDescriptor() ([]byte, []int) {
	return file_store_release_proto_rawDescGZIP(), []int{0}
}

type ReleasePayload struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Title     string                    `protobuf:"bytes,1,opt,name=title,proto3" json:"title,omitempty"`
	Files     []*ReleasePayload_File    `protobuf:"bytes,2,rep,name=files,proto3" json:"files,omitempty"`
	VcsSource *ReleasePayload_VCSSource `protobuf:"bytes,3,opt,name=vcs_source,json=vcsSource,proto3" json:"vcs_source,omitempty"`
}

func (x *ReleasePayload) Reset() {
	*x = ReleasePayload{}
	mi := &file_store_release_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ReleasePayload) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReleasePayload) ProtoMessage() {}

func (x *ReleasePayload) ProtoReflect() protoreflect.Message {
	mi := &file_store_release_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReleasePayload.ProtoReflect.Descriptor instead.
func (*ReleasePayload) Descriptor() ([]byte, []int) {
	return file_store_release_proto_rawDescGZIP(), []int{0}
}

func (x *ReleasePayload) GetTitle() string {
	if x != nil {
		return x.Title
	}
	return ""
}

func (x *ReleasePayload) GetFiles() []*ReleasePayload_File {
	if x != nil {
		return x.Files
	}
	return nil
}

func (x *ReleasePayload) GetVcsSource() *ReleasePayload_VCSSource {
	if x != nil {
		return x.VcsSource
	}
	return nil
}

type ReleasePayload_File struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The unique identifier for the file.
	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	// The path of the file. e.g. `2.2/V0001_create_table.sql`.
	Path string `protobuf:"bytes,2,opt,name=path,proto3" json:"path,omitempty"`
	// The sheet that holds the content.
	// Format: projects/{project}/sheets/{sheet}
	Sheet string `protobuf:"bytes,3,opt,name=sheet,proto3" json:"sheet,omitempty"`
	// The SHA256 hash value of the sheet.
	SheetSha256 string          `protobuf:"bytes,4,opt,name=sheet_sha256,json=sheetSha256,proto3" json:"sheet_sha256,omitempty"`
	Type        ReleaseFileType `protobuf:"varint,5,opt,name=type,proto3,enum=bytebase.store.ReleaseFileType" json:"type,omitempty"`
	Version     string          `protobuf:"bytes,6,opt,name=version,proto3" json:"version,omitempty"`
}

func (x *ReleasePayload_File) Reset() {
	*x = ReleasePayload_File{}
	mi := &file_store_release_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ReleasePayload_File) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReleasePayload_File) ProtoMessage() {}

func (x *ReleasePayload_File) ProtoReflect() protoreflect.Message {
	mi := &file_store_release_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReleasePayload_File.ProtoReflect.Descriptor instead.
func (*ReleasePayload_File) Descriptor() ([]byte, []int) {
	return file_store_release_proto_rawDescGZIP(), []int{0, 0}
}

func (x *ReleasePayload_File) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *ReleasePayload_File) GetPath() string {
	if x != nil {
		return x.Path
	}
	return ""
}

func (x *ReleasePayload_File) GetSheet() string {
	if x != nil {
		return x.Sheet
	}
	return ""
}

func (x *ReleasePayload_File) GetSheetSha256() string {
	if x != nil {
		return x.SheetSha256
	}
	return ""
}

func (x *ReleasePayload_File) GetType() ReleaseFileType {
	if x != nil {
		return x.Type
	}
	return ReleaseFileType_TYPE_UNSPECIFIED
}

func (x *ReleasePayload_File) GetVersion() string {
	if x != nil {
		return x.Version
	}
	return ""
}

type ReleasePayload_VCSSource struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	VcsType        VCSType `protobuf:"varint,1,opt,name=vcs_type,json=vcsType,proto3,enum=bytebase.store.VCSType" json:"vcs_type,omitempty"`
	PullRequestUrl string  `protobuf:"bytes,2,opt,name=pull_request_url,json=pullRequestUrl,proto3" json:"pull_request_url,omitempty"`
}

func (x *ReleasePayload_VCSSource) Reset() {
	*x = ReleasePayload_VCSSource{}
	mi := &file_store_release_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ReleasePayload_VCSSource) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReleasePayload_VCSSource) ProtoMessage() {}

func (x *ReleasePayload_VCSSource) ProtoReflect() protoreflect.Message {
	mi := &file_store_release_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReleasePayload_VCSSource.ProtoReflect.Descriptor instead.
func (*ReleasePayload_VCSSource) Descriptor() ([]byte, []int) {
	return file_store_release_proto_rawDescGZIP(), []int{0, 1}
}

func (x *ReleasePayload_VCSSource) GetVcsType() VCSType {
	if x != nil {
		return x.VcsType
	}
	return VCSType_VCS_TYPE_UNSPECIFIED
}

func (x *ReleasePayload_VCSSource) GetPullRequestUrl() string {
	if x != nil {
		return x.PullRequestUrl
	}
	return ""
}

var File_store_release_proto protoreflect.FileDescriptor

var file_store_release_proto_rawDesc = []byte{
	0x0a, 0x13, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2f, 0x72, 0x65, 0x6c, 0x65, 0x61, 0x73, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0e, 0x62, 0x79, 0x74, 0x65, 0x62, 0x61, 0x73, 0x65, 0x2e,
	0x73, 0x74, 0x6f, 0x72, 0x65, 0x1a, 0x19, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70,
	0x69, 0x2f, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x1a, 0x12, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x22, 0xe3, 0x03, 0x0a, 0x0e, 0x52, 0x65, 0x6c, 0x65, 0x61, 0x73, 0x65,
	0x50, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x69, 0x74, 0x6c, 0x65,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x12, 0x39, 0x0a,
	0x05, 0x66, 0x69, 0x6c, 0x65, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x23, 0x2e, 0x62,
	0x79, 0x74, 0x65, 0x62, 0x61, 0x73, 0x65, 0x2e, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2e, 0x52, 0x65,
	0x6c, 0x65, 0x61, 0x73, 0x65, 0x50, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x2e, 0x46, 0x69, 0x6c,
	0x65, 0x52, 0x05, 0x66, 0x69, 0x6c, 0x65, 0x73, 0x12, 0x47, 0x0a, 0x0a, 0x76, 0x63, 0x73, 0x5f,
	0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x28, 0x2e, 0x62,
	0x79, 0x74, 0x65, 0x62, 0x61, 0x73, 0x65, 0x2e, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2e, 0x52, 0x65,
	0x6c, 0x65, 0x61, 0x73, 0x65, 0x50, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x2e, 0x56, 0x43, 0x53,
	0x53, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x52, 0x09, 0x76, 0x63, 0x73, 0x53, 0x6f, 0x75, 0x72, 0x63,
	0x65, 0x1a, 0xcb, 0x01, 0x0a, 0x04, 0x46, 0x69, 0x6c, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x70, 0x61,
	0x74, 0x68, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x70, 0x61, 0x74, 0x68, 0x12, 0x2d,
	0x0a, 0x05, 0x73, 0x68, 0x65, 0x65, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x42, 0x17, 0xfa,
	0x41, 0x14, 0x0a, 0x12, 0x62, 0x79, 0x74, 0x65, 0x62, 0x61, 0x73, 0x65, 0x2e, 0x63, 0x6f, 0x6d,
	0x2f, 0x53, 0x68, 0x65, 0x65, 0x74, 0x52, 0x05, 0x73, 0x68, 0x65, 0x65, 0x74, 0x12, 0x21, 0x0a,
	0x0c, 0x73, 0x68, 0x65, 0x65, 0x74, 0x5f, 0x73, 0x68, 0x61, 0x32, 0x35, 0x36, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0b, 0x73, 0x68, 0x65, 0x65, 0x74, 0x53, 0x68, 0x61, 0x32, 0x35, 0x36,
	0x12, 0x33, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x1f,
	0x2e, 0x62, 0x79, 0x74, 0x65, 0x62, 0x61, 0x73, 0x65, 0x2e, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2e,
	0x52, 0x65, 0x6c, 0x65, 0x61, 0x73, 0x65, 0x46, 0x69, 0x6c, 0x65, 0x54, 0x79, 0x70, 0x65, 0x52,
	0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e,
	0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x1a,
	0x69, 0x0a, 0x09, 0x56, 0x43, 0x53, 0x53, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x12, 0x32, 0x0a, 0x08,
	0x76, 0x63, 0x73, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x17,
	0x2e, 0x62, 0x79, 0x74, 0x65, 0x62, 0x61, 0x73, 0x65, 0x2e, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2e,
	0x56, 0x43, 0x53, 0x54, 0x79, 0x70, 0x65, 0x52, 0x07, 0x76, 0x63, 0x73, 0x54, 0x79, 0x70, 0x65,
	0x12, 0x28, 0x0a, 0x10, 0x70, 0x75, 0x6c, 0x6c, 0x5f, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x5f, 0x75, 0x72, 0x6c, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0e, 0x70, 0x75, 0x6c, 0x6c,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x55, 0x72, 0x6c, 0x2a, 0x36, 0x0a, 0x0f, 0x52, 0x65,
	0x6c, 0x65, 0x61, 0x73, 0x65, 0x46, 0x69, 0x6c, 0x65, 0x54, 0x79, 0x70, 0x65, 0x12, 0x14, 0x0a,
	0x10, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x55, 0x4e, 0x53, 0x50, 0x45, 0x43, 0x49, 0x46, 0x49, 0x45,
	0x44, 0x10, 0x00, 0x12, 0x0d, 0x0a, 0x09, 0x56, 0x45, 0x52, 0x53, 0x49, 0x4f, 0x4e, 0x45, 0x44,
	0x10, 0x01, 0x42, 0x14, 0x5a, 0x12, 0x67, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x74, 0x65, 0x64, 0x2d,
	0x67, 0x6f, 0x2f, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_store_release_proto_rawDescOnce sync.Once
	file_store_release_proto_rawDescData = file_store_release_proto_rawDesc
)

func file_store_release_proto_rawDescGZIP() []byte {
	file_store_release_proto_rawDescOnce.Do(func() {
		file_store_release_proto_rawDescData = protoimpl.X.CompressGZIP(file_store_release_proto_rawDescData)
	})
	return file_store_release_proto_rawDescData
}

var file_store_release_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_store_release_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_store_release_proto_goTypes = []any{
	(ReleaseFileType)(0),             // 0: bytebase.store.ReleaseFileType
	(*ReleasePayload)(nil),           // 1: bytebase.store.ReleasePayload
	(*ReleasePayload_File)(nil),      // 2: bytebase.store.ReleasePayload.File
	(*ReleasePayload_VCSSource)(nil), // 3: bytebase.store.ReleasePayload.VCSSource
	(VCSType)(0),                     // 4: bytebase.store.VCSType
}
var file_store_release_proto_depIdxs = []int32{
	2, // 0: bytebase.store.ReleasePayload.files:type_name -> bytebase.store.ReleasePayload.File
	3, // 1: bytebase.store.ReleasePayload.vcs_source:type_name -> bytebase.store.ReleasePayload.VCSSource
	0, // 2: bytebase.store.ReleasePayload.File.type:type_name -> bytebase.store.ReleaseFileType
	4, // 3: bytebase.store.ReleasePayload.VCSSource.vcs_type:type_name -> bytebase.store.VCSType
	4, // [4:4] is the sub-list for method output_type
	4, // [4:4] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_store_release_proto_init() }
func file_store_release_proto_init() {
	if File_store_release_proto != nil {
		return
	}
	file_store_common_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_store_release_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_store_release_proto_goTypes,
		DependencyIndexes: file_store_release_proto_depIdxs,
		EnumInfos:         file_store_release_proto_enumTypes,
		MessageInfos:      file_store_release_proto_msgTypes,
	}.Build()
	File_store_release_proto = out.File
	file_store_release_proto_rawDesc = nil
	file_store_release_proto_goTypes = nil
	file_store_release_proto_depIdxs = nil
}

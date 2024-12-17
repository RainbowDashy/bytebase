// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.0
// 	protoc        (unknown)
// source: store/changelog.proto

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

type ChangelogPayload_Type int32

const (
	ChangelogPayload_TYPE_UNSPECIFIED ChangelogPayload_Type = 0
	ChangelogPayload_BASELINE         ChangelogPayload_Type = 1
	ChangelogPayload_MIGRATE          ChangelogPayload_Type = 2
	ChangelogPayload_MIGRATE_SDL      ChangelogPayload_Type = 3
	ChangelogPayload_MIGRATE_GHOST    ChangelogPayload_Type = 4
	ChangelogPayload_DATA             ChangelogPayload_Type = 6
)

// Enum value maps for ChangelogPayload_Type.
var (
	ChangelogPayload_Type_name = map[int32]string{
		0: "TYPE_UNSPECIFIED",
		1: "BASELINE",
		2: "MIGRATE",
		3: "MIGRATE_SDL",
		4: "MIGRATE_GHOST",
		6: "DATA",
	}
	ChangelogPayload_Type_value = map[string]int32{
		"TYPE_UNSPECIFIED": 0,
		"BASELINE":         1,
		"MIGRATE":          2,
		"MIGRATE_SDL":      3,
		"MIGRATE_GHOST":    4,
		"DATA":             6,
	}
)

func (x ChangelogPayload_Type) Enum() *ChangelogPayload_Type {
	p := new(ChangelogPayload_Type)
	*p = x
	return p
}

func (x ChangelogPayload_Type) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ChangelogPayload_Type) Descriptor() protoreflect.EnumDescriptor {
	return file_store_changelog_proto_enumTypes[0].Descriptor()
}

func (ChangelogPayload_Type) Type() protoreflect.EnumType {
	return &file_store_changelog_proto_enumTypes[0]
}

func (x ChangelogPayload_Type) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ChangelogPayload_Type.Descriptor instead.
func (ChangelogPayload_Type) EnumDescriptor() ([]byte, []int) {
	return file_store_changelog_proto_rawDescGZIP(), []int{0, 0}
}

type ChangelogPayload struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Format: projects/{project}/rollouts/{rollout}/stages/{stage}/tasks/{task}/taskruns/{taskrun}
	TaskRun string `protobuf:"bytes,1,opt,name=task_run,json=taskRun,proto3" json:"task_run,omitempty"`
	// Format: projects/{project}/issues/{issue}
	Issue string `protobuf:"bytes,2,opt,name=issue,proto3" json:"issue,omitempty"`
	// The revision uid.
	// optional
	Revision         int64             `protobuf:"varint,3,opt,name=revision,proto3" json:"revision,omitempty"`
	ChangedResources *ChangedResources `protobuf:"bytes,4,opt,name=changed_resources,json=changedResources,proto3" json:"changed_resources,omitempty"`
	// The sheet that holds the content.
	// Format: projects/{project}/sheets/{sheet}
	Sheet         string                `protobuf:"bytes,5,opt,name=sheet,proto3" json:"sheet,omitempty"`
	Version       string                `protobuf:"bytes,6,opt,name=version,proto3" json:"version,omitempty"`
	Type          ChangelogPayload_Type `protobuf:"varint,7,opt,name=type,proto3,enum=bytebase.store.ChangelogPayload_Type" json:"type,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ChangelogPayload) Reset() {
	*x = ChangelogPayload{}
	mi := &file_store_changelog_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ChangelogPayload) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ChangelogPayload) ProtoMessage() {}

func (x *ChangelogPayload) ProtoReflect() protoreflect.Message {
	mi := &file_store_changelog_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ChangelogPayload.ProtoReflect.Descriptor instead.
func (*ChangelogPayload) Descriptor() ([]byte, []int) {
	return file_store_changelog_proto_rawDescGZIP(), []int{0}
}

func (x *ChangelogPayload) GetTaskRun() string {
	if x != nil {
		return x.TaskRun
	}
	return ""
}

func (x *ChangelogPayload) GetIssue() string {
	if x != nil {
		return x.Issue
	}
	return ""
}

func (x *ChangelogPayload) GetRevision() int64 {
	if x != nil {
		return x.Revision
	}
	return 0
}

func (x *ChangelogPayload) GetChangedResources() *ChangedResources {
	if x != nil {
		return x.ChangedResources
	}
	return nil
}

func (x *ChangelogPayload) GetSheet() string {
	if x != nil {
		return x.Sheet
	}
	return ""
}

func (x *ChangelogPayload) GetVersion() string {
	if x != nil {
		return x.Version
	}
	return ""
}

func (x *ChangelogPayload) GetType() ChangelogPayload_Type {
	if x != nil {
		return x.Type
	}
	return ChangelogPayload_TYPE_UNSPECIFIED
}

var File_store_changelog_proto protoreflect.FileDescriptor

var file_store_changelog_proto_rawDesc = []byte{
	0x0a, 0x15, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2f, 0x63, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x6c, 0x6f,
	0x67, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0e, 0x62, 0x79, 0x74, 0x65, 0x62, 0x61, 0x73,
	0x65, 0x2e, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x1a, 0x23, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2f, 0x69,
	0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x5f, 0x63, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x5f, 0x68,
	0x69, 0x73, 0x74, 0x6f, 0x72, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x80, 0x03, 0x0a,
	0x10, 0x43, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x6c, 0x6f, 0x67, 0x50, 0x61, 0x79, 0x6c, 0x6f, 0x61,
	0x64, 0x12, 0x19, 0x0a, 0x08, 0x74, 0x61, 0x73, 0x6b, 0x5f, 0x72, 0x75, 0x6e, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x07, 0x74, 0x61, 0x73, 0x6b, 0x52, 0x75, 0x6e, 0x12, 0x14, 0x0a, 0x05,
	0x69, 0x73, 0x73, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x69, 0x73, 0x73,
	0x75, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x72, 0x65, 0x76, 0x69, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x03, 0x52, 0x08, 0x72, 0x65, 0x76, 0x69, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x4d,
	0x0a, 0x11, 0x63, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x64, 0x5f, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72,
	0x63, 0x65, 0x73, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x20, 0x2e, 0x62, 0x79, 0x74, 0x65,
	0x62, 0x61, 0x73, 0x65, 0x2e, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2e, 0x43, 0x68, 0x61, 0x6e, 0x67,
	0x65, 0x64, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x52, 0x10, 0x63, 0x68, 0x61,
	0x6e, 0x67, 0x65, 0x64, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x12, 0x14, 0x0a,
	0x05, 0x73, 0x68, 0x65, 0x65, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x73, 0x68,
	0x65, 0x65, 0x74, 0x12, 0x18, 0x0a, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x06,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x39, 0x0a,
	0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x07, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x25, 0x2e, 0x62, 0x79,
	0x74, 0x65, 0x62, 0x61, 0x73, 0x65, 0x2e, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2e, 0x43, 0x68, 0x61,
	0x6e, 0x67, 0x65, 0x6c, 0x6f, 0x67, 0x50, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x2e, 0x54, 0x79,
	0x70, 0x65, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x22, 0x65, 0x0a, 0x04, 0x54, 0x79, 0x70, 0x65,
	0x12, 0x14, 0x0a, 0x10, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x55, 0x4e, 0x53, 0x50, 0x45, 0x43, 0x49,
	0x46, 0x49, 0x45, 0x44, 0x10, 0x00, 0x12, 0x0c, 0x0a, 0x08, 0x42, 0x41, 0x53, 0x45, 0x4c, 0x49,
	0x4e, 0x45, 0x10, 0x01, 0x12, 0x0b, 0x0a, 0x07, 0x4d, 0x49, 0x47, 0x52, 0x41, 0x54, 0x45, 0x10,
	0x02, 0x12, 0x0f, 0x0a, 0x0b, 0x4d, 0x49, 0x47, 0x52, 0x41, 0x54, 0x45, 0x5f, 0x53, 0x44, 0x4c,
	0x10, 0x03, 0x12, 0x11, 0x0a, 0x0d, 0x4d, 0x49, 0x47, 0x52, 0x41, 0x54, 0x45, 0x5f, 0x47, 0x48,
	0x4f, 0x53, 0x54, 0x10, 0x04, 0x12, 0x08, 0x0a, 0x04, 0x44, 0x41, 0x54, 0x41, 0x10, 0x06, 0x42,
	0x14, 0x5a, 0x12, 0x67, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x74, 0x65, 0x64, 0x2d, 0x67, 0x6f, 0x2f,
	0x73, 0x74, 0x6f, 0x72, 0x65, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_store_changelog_proto_rawDescOnce sync.Once
	file_store_changelog_proto_rawDescData = file_store_changelog_proto_rawDesc
)

func file_store_changelog_proto_rawDescGZIP() []byte {
	file_store_changelog_proto_rawDescOnce.Do(func() {
		file_store_changelog_proto_rawDescData = protoimpl.X.CompressGZIP(file_store_changelog_proto_rawDescData)
	})
	return file_store_changelog_proto_rawDescData
}

var file_store_changelog_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_store_changelog_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_store_changelog_proto_goTypes = []any{
	(ChangelogPayload_Type)(0), // 0: bytebase.store.ChangelogPayload.Type
	(*ChangelogPayload)(nil),   // 1: bytebase.store.ChangelogPayload
	(*ChangedResources)(nil),   // 2: bytebase.store.ChangedResources
}
var file_store_changelog_proto_depIdxs = []int32{
	2, // 0: bytebase.store.ChangelogPayload.changed_resources:type_name -> bytebase.store.ChangedResources
	0, // 1: bytebase.store.ChangelogPayload.type:type_name -> bytebase.store.ChangelogPayload.Type
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_store_changelog_proto_init() }
func file_store_changelog_proto_init() {
	if File_store_changelog_proto != nil {
		return
	}
	file_store_instance_change_history_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_store_changelog_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_store_changelog_proto_goTypes,
		DependencyIndexes: file_store_changelog_proto_depIdxs,
		EnumInfos:         file_store_changelog_proto_enumTypes,
		MessageInfos:      file_store_changelog_proto_msgTypes,
	}.Build()
	File_store_changelog_proto = out.File
	file_store_changelog_proto_rawDesc = nil
	file_store_changelog_proto_goTypes = nil
	file_store_changelog_proto_depIdxs = nil
}

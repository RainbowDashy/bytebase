// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.1
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

type ChangedResources struct {
	state         protoimpl.MessageState     `protogen:"open.v1"`
	Databases     []*ChangedResourceDatabase `protobuf:"bytes,1,rep,name=databases,proto3" json:"databases,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ChangedResources) Reset() {
	*x = ChangedResources{}
	mi := &file_store_changelog_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ChangedResources) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ChangedResources) ProtoMessage() {}

func (x *ChangedResources) ProtoReflect() protoreflect.Message {
	mi := &file_store_changelog_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ChangedResources.ProtoReflect.Descriptor instead.
func (*ChangedResources) Descriptor() ([]byte, []int) {
	return file_store_changelog_proto_rawDescGZIP(), []int{1}
}

func (x *ChangedResources) GetDatabases() []*ChangedResourceDatabase {
	if x != nil {
		return x.Databases
	}
	return nil
}

type ChangedResourceDatabase struct {
	state         protoimpl.MessageState   `protogen:"open.v1"`
	Name          string                   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Schemas       []*ChangedResourceSchema `protobuf:"bytes,2,rep,name=schemas,proto3" json:"schemas,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ChangedResourceDatabase) Reset() {
	*x = ChangedResourceDatabase{}
	mi := &file_store_changelog_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ChangedResourceDatabase) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ChangedResourceDatabase) ProtoMessage() {}

func (x *ChangedResourceDatabase) ProtoReflect() protoreflect.Message {
	mi := &file_store_changelog_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ChangedResourceDatabase.ProtoReflect.Descriptor instead.
func (*ChangedResourceDatabase) Descriptor() ([]byte, []int) {
	return file_store_changelog_proto_rawDescGZIP(), []int{2}
}

func (x *ChangedResourceDatabase) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *ChangedResourceDatabase) GetSchemas() []*ChangedResourceSchema {
	if x != nil {
		return x.Schemas
	}
	return nil
}

type ChangedResourceSchema struct {
	state         protoimpl.MessageState      `protogen:"open.v1"`
	Name          string                      `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Tables        []*ChangedResourceTable     `protobuf:"bytes,2,rep,name=tables,proto3" json:"tables,omitempty"`
	Views         []*ChangedResourceView      `protobuf:"bytes,3,rep,name=views,proto3" json:"views,omitempty"`
	Functions     []*ChangedResourceFunction  `protobuf:"bytes,4,rep,name=functions,proto3" json:"functions,omitempty"`
	Procedures    []*ChangedResourceProcedure `protobuf:"bytes,5,rep,name=procedures,proto3" json:"procedures,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ChangedResourceSchema) Reset() {
	*x = ChangedResourceSchema{}
	mi := &file_store_changelog_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ChangedResourceSchema) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ChangedResourceSchema) ProtoMessage() {}

func (x *ChangedResourceSchema) ProtoReflect() protoreflect.Message {
	mi := &file_store_changelog_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ChangedResourceSchema.ProtoReflect.Descriptor instead.
func (*ChangedResourceSchema) Descriptor() ([]byte, []int) {
	return file_store_changelog_proto_rawDescGZIP(), []int{3}
}

func (x *ChangedResourceSchema) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *ChangedResourceSchema) GetTables() []*ChangedResourceTable {
	if x != nil {
		return x.Tables
	}
	return nil
}

func (x *ChangedResourceSchema) GetViews() []*ChangedResourceView {
	if x != nil {
		return x.Views
	}
	return nil
}

func (x *ChangedResourceSchema) GetFunctions() []*ChangedResourceFunction {
	if x != nil {
		return x.Functions
	}
	return nil
}

func (x *ChangedResourceSchema) GetProcedures() []*ChangedResourceProcedure {
	if x != nil {
		return x.Procedures
	}
	return nil
}

type ChangedResourceTable struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	Name  string                 `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	// estimated row count of the table
	TableRows int64 `protobuf:"varint,2,opt,name=table_rows,json=tableRows,proto3" json:"table_rows,omitempty"`
	// The ranges of sub-strings correspond to the statements on the sheet.
	Ranges        []*Range `protobuf:"bytes,3,rep,name=ranges,proto3" json:"ranges,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ChangedResourceTable) Reset() {
	*x = ChangedResourceTable{}
	mi := &file_store_changelog_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ChangedResourceTable) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ChangedResourceTable) ProtoMessage() {}

func (x *ChangedResourceTable) ProtoReflect() protoreflect.Message {
	mi := &file_store_changelog_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ChangedResourceTable.ProtoReflect.Descriptor instead.
func (*ChangedResourceTable) Descriptor() ([]byte, []int) {
	return file_store_changelog_proto_rawDescGZIP(), []int{4}
}

func (x *ChangedResourceTable) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *ChangedResourceTable) GetTableRows() int64 {
	if x != nil {
		return x.TableRows
	}
	return 0
}

func (x *ChangedResourceTable) GetRanges() []*Range {
	if x != nil {
		return x.Ranges
	}
	return nil
}

type ChangedResourceView struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	Name  string                 `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	// The ranges of sub-strings correspond to the statements on the sheet.
	Ranges        []*Range `protobuf:"bytes,2,rep,name=ranges,proto3" json:"ranges,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ChangedResourceView) Reset() {
	*x = ChangedResourceView{}
	mi := &file_store_changelog_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ChangedResourceView) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ChangedResourceView) ProtoMessage() {}

func (x *ChangedResourceView) ProtoReflect() protoreflect.Message {
	mi := &file_store_changelog_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ChangedResourceView.ProtoReflect.Descriptor instead.
func (*ChangedResourceView) Descriptor() ([]byte, []int) {
	return file_store_changelog_proto_rawDescGZIP(), []int{5}
}

func (x *ChangedResourceView) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *ChangedResourceView) GetRanges() []*Range {
	if x != nil {
		return x.Ranges
	}
	return nil
}

type ChangedResourceFunction struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	Name  string                 `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	// The ranges of sub-strings correspond to the statements on the sheet.
	Ranges        []*Range `protobuf:"bytes,2,rep,name=ranges,proto3" json:"ranges,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ChangedResourceFunction) Reset() {
	*x = ChangedResourceFunction{}
	mi := &file_store_changelog_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ChangedResourceFunction) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ChangedResourceFunction) ProtoMessage() {}

func (x *ChangedResourceFunction) ProtoReflect() protoreflect.Message {
	mi := &file_store_changelog_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ChangedResourceFunction.ProtoReflect.Descriptor instead.
func (*ChangedResourceFunction) Descriptor() ([]byte, []int) {
	return file_store_changelog_proto_rawDescGZIP(), []int{6}
}

func (x *ChangedResourceFunction) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *ChangedResourceFunction) GetRanges() []*Range {
	if x != nil {
		return x.Ranges
	}
	return nil
}

type ChangedResourceProcedure struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	Name  string                 `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	// The ranges of sub-strings correspond to the statements on the sheet.
	Ranges        []*Range `protobuf:"bytes,2,rep,name=ranges,proto3" json:"ranges,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ChangedResourceProcedure) Reset() {
	*x = ChangedResourceProcedure{}
	mi := &file_store_changelog_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ChangedResourceProcedure) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ChangedResourceProcedure) ProtoMessage() {}

func (x *ChangedResourceProcedure) ProtoReflect() protoreflect.Message {
	mi := &file_store_changelog_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ChangedResourceProcedure.ProtoReflect.Descriptor instead.
func (*ChangedResourceProcedure) Descriptor() ([]byte, []int) {
	return file_store_changelog_proto_rawDescGZIP(), []int{7}
}

func (x *ChangedResourceProcedure) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *ChangedResourceProcedure) GetRanges() []*Range {
	if x != nil {
		return x.Ranges
	}
	return nil
}

var File_store_changelog_proto protoreflect.FileDescriptor

var file_store_changelog_proto_rawDesc = []byte{
	0x0a, 0x15, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2f, 0x63, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x6c, 0x6f,
	0x67, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0e, 0x62, 0x79, 0x74, 0x65, 0x62, 0x61, 0x73,
	0x65, 0x2e, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x1a, 0x12, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2f, 0x63,
	0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x80, 0x03, 0x0a, 0x10,
	0x43, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x6c, 0x6f, 0x67, 0x50, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64,
	0x12, 0x19, 0x0a, 0x08, 0x74, 0x61, 0x73, 0x6b, 0x5f, 0x72, 0x75, 0x6e, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x07, 0x74, 0x61, 0x73, 0x6b, 0x52, 0x75, 0x6e, 0x12, 0x14, 0x0a, 0x05, 0x69,
	0x73, 0x73, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x69, 0x73, 0x73, 0x75,
	0x65, 0x12, 0x1a, 0x0a, 0x08, 0x72, 0x65, 0x76, 0x69, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x03, 0x52, 0x08, 0x72, 0x65, 0x76, 0x69, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x4d, 0x0a,
	0x11, 0x63, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x64, 0x5f, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63,
	0x65, 0x73, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x20, 0x2e, 0x62, 0x79, 0x74, 0x65, 0x62,
	0x61, 0x73, 0x65, 0x2e, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2e, 0x43, 0x68, 0x61, 0x6e, 0x67, 0x65,
	0x64, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x52, 0x10, 0x63, 0x68, 0x61, 0x6e,
	0x67, 0x65, 0x64, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x12, 0x14, 0x0a, 0x05,
	0x73, 0x68, 0x65, 0x65, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x73, 0x68, 0x65,
	0x65, 0x74, 0x12, 0x18, 0x0a, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x06, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x39, 0x0a, 0x04,
	0x74, 0x79, 0x70, 0x65, 0x18, 0x07, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x25, 0x2e, 0x62, 0x79, 0x74,
	0x65, 0x62, 0x61, 0x73, 0x65, 0x2e, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2e, 0x43, 0x68, 0x61, 0x6e,
	0x67, 0x65, 0x6c, 0x6f, 0x67, 0x50, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x2e, 0x54, 0x79, 0x70,
	0x65, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x22, 0x65, 0x0a, 0x04, 0x54, 0x79, 0x70, 0x65, 0x12,
	0x14, 0x0a, 0x10, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x55, 0x4e, 0x53, 0x50, 0x45, 0x43, 0x49, 0x46,
	0x49, 0x45, 0x44, 0x10, 0x00, 0x12, 0x0c, 0x0a, 0x08, 0x42, 0x41, 0x53, 0x45, 0x4c, 0x49, 0x4e,
	0x45, 0x10, 0x01, 0x12, 0x0b, 0x0a, 0x07, 0x4d, 0x49, 0x47, 0x52, 0x41, 0x54, 0x45, 0x10, 0x02,
	0x12, 0x0f, 0x0a, 0x0b, 0x4d, 0x49, 0x47, 0x52, 0x41, 0x54, 0x45, 0x5f, 0x53, 0x44, 0x4c, 0x10,
	0x03, 0x12, 0x11, 0x0a, 0x0d, 0x4d, 0x49, 0x47, 0x52, 0x41, 0x54, 0x45, 0x5f, 0x47, 0x48, 0x4f,
	0x53, 0x54, 0x10, 0x04, 0x12, 0x08, 0x0a, 0x04, 0x44, 0x41, 0x54, 0x41, 0x10, 0x06, 0x22, 0x59,
	0x0a, 0x10, 0x43, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x64, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63,
	0x65, 0x73, 0x12, 0x45, 0x0a, 0x09, 0x64, 0x61, 0x74, 0x61, 0x62, 0x61, 0x73, 0x65, 0x73, 0x18,
	0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x27, 0x2e, 0x62, 0x79, 0x74, 0x65, 0x62, 0x61, 0x73, 0x65,
	0x2e, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2e, 0x43, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x64, 0x52, 0x65,
	0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x44, 0x61, 0x74, 0x61, 0x62, 0x61, 0x73, 0x65, 0x52, 0x09,
	0x64, 0x61, 0x74, 0x61, 0x62, 0x61, 0x73, 0x65, 0x73, 0x22, 0x6e, 0x0a, 0x17, 0x43, 0x68, 0x61,
	0x6e, 0x67, 0x65, 0x64, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x44, 0x61, 0x74, 0x61,
	0x62, 0x61, 0x73, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x3f, 0x0a, 0x07, 0x73, 0x63, 0x68, 0x65,
	0x6d, 0x61, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x25, 0x2e, 0x62, 0x79, 0x74, 0x65,
	0x62, 0x61, 0x73, 0x65, 0x2e, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2e, 0x43, 0x68, 0x61, 0x6e, 0x67,
	0x65, 0x64, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x53, 0x63, 0x68, 0x65, 0x6d, 0x61,
	0x52, 0x07, 0x73, 0x63, 0x68, 0x65, 0x6d, 0x61, 0x73, 0x22, 0xb5, 0x02, 0x0a, 0x15, 0x43, 0x68,
	0x61, 0x6e, 0x67, 0x65, 0x64, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x53, 0x63, 0x68,
	0x65, 0x6d, 0x61, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x3c, 0x0a, 0x06, 0x74, 0x61, 0x62, 0x6c, 0x65,
	0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x24, 0x2e, 0x62, 0x79, 0x74, 0x65, 0x62, 0x61,
	0x73, 0x65, 0x2e, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2e, 0x43, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x64,
	0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x54, 0x61, 0x62, 0x6c, 0x65, 0x52, 0x06, 0x74,
	0x61, 0x62, 0x6c, 0x65, 0x73, 0x12, 0x39, 0x0a, 0x05, 0x76, 0x69, 0x65, 0x77, 0x73, 0x18, 0x03,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x23, 0x2e, 0x62, 0x79, 0x74, 0x65, 0x62, 0x61, 0x73, 0x65, 0x2e,
	0x73, 0x74, 0x6f, 0x72, 0x65, 0x2e, 0x43, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x64, 0x52, 0x65, 0x73,
	0x6f, 0x75, 0x72, 0x63, 0x65, 0x56, 0x69, 0x65, 0x77, 0x52, 0x05, 0x76, 0x69, 0x65, 0x77, 0x73,
	0x12, 0x45, 0x0a, 0x09, 0x66, 0x75, 0x6e, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x04, 0x20,
	0x03, 0x28, 0x0b, 0x32, 0x27, 0x2e, 0x62, 0x79, 0x74, 0x65, 0x62, 0x61, 0x73, 0x65, 0x2e, 0x73,
	0x74, 0x6f, 0x72, 0x65, 0x2e, 0x43, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x64, 0x52, 0x65, 0x73, 0x6f,
	0x75, 0x72, 0x63, 0x65, 0x46, 0x75, 0x6e, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x09, 0x66, 0x75,
	0x6e, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x12, 0x48, 0x0a, 0x0a, 0x70, 0x72, 0x6f, 0x63, 0x65,
	0x64, 0x75, 0x72, 0x65, 0x73, 0x18, 0x05, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x28, 0x2e, 0x62, 0x79,
	0x74, 0x65, 0x62, 0x61, 0x73, 0x65, 0x2e, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2e, 0x43, 0x68, 0x61,
	0x6e, 0x67, 0x65, 0x64, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x50, 0x72, 0x6f, 0x63,
	0x65, 0x64, 0x75, 0x72, 0x65, 0x52, 0x0a, 0x70, 0x72, 0x6f, 0x63, 0x65, 0x64, 0x75, 0x72, 0x65,
	0x73, 0x22, 0x78, 0x0a, 0x14, 0x43, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x64, 0x52, 0x65, 0x73, 0x6f,
	0x75, 0x72, 0x63, 0x65, 0x54, 0x61, 0x62, 0x6c, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d,
	0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x1d, 0x0a,
	0x0a, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x5f, 0x72, 0x6f, 0x77, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x03, 0x52, 0x09, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x52, 0x6f, 0x77, 0x73, 0x12, 0x2d, 0x0a, 0x06,
	0x72, 0x61, 0x6e, 0x67, 0x65, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x62,
	0x79, 0x74, 0x65, 0x62, 0x61, 0x73, 0x65, 0x2e, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2e, 0x52, 0x61,
	0x6e, 0x67, 0x65, 0x52, 0x06, 0x72, 0x61, 0x6e, 0x67, 0x65, 0x73, 0x22, 0x58, 0x0a, 0x13, 0x43,
	0x68, 0x61, 0x6e, 0x67, 0x65, 0x64, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x56, 0x69,
	0x65, 0x77, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x2d, 0x0a, 0x06, 0x72, 0x61, 0x6e, 0x67, 0x65, 0x73,
	0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x62, 0x79, 0x74, 0x65, 0x62, 0x61, 0x73,
	0x65, 0x2e, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2e, 0x52, 0x61, 0x6e, 0x67, 0x65, 0x52, 0x06, 0x72,
	0x61, 0x6e, 0x67, 0x65, 0x73, 0x22, 0x5c, 0x0a, 0x17, 0x43, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x64,
	0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x46, 0x75, 0x6e, 0x63, 0x74, 0x69, 0x6f, 0x6e,
	0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04,
	0x6e, 0x61, 0x6d, 0x65, 0x12, 0x2d, 0x0a, 0x06, 0x72, 0x61, 0x6e, 0x67, 0x65, 0x73, 0x18, 0x02,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x62, 0x79, 0x74, 0x65, 0x62, 0x61, 0x73, 0x65, 0x2e,
	0x73, 0x74, 0x6f, 0x72, 0x65, 0x2e, 0x52, 0x61, 0x6e, 0x67, 0x65, 0x52, 0x06, 0x72, 0x61, 0x6e,
	0x67, 0x65, 0x73, 0x22, 0x5d, 0x0a, 0x18, 0x43, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x64, 0x52, 0x65,
	0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x50, 0x72, 0x6f, 0x63, 0x65, 0x64, 0x75, 0x72, 0x65, 0x12,
	0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e,
	0x61, 0x6d, 0x65, 0x12, 0x2d, 0x0a, 0x06, 0x72, 0x61, 0x6e, 0x67, 0x65, 0x73, 0x18, 0x02, 0x20,
	0x03, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x62, 0x79, 0x74, 0x65, 0x62, 0x61, 0x73, 0x65, 0x2e, 0x73,
	0x74, 0x6f, 0x72, 0x65, 0x2e, 0x52, 0x61, 0x6e, 0x67, 0x65, 0x52, 0x06, 0x72, 0x61, 0x6e, 0x67,
	0x65, 0x73, 0x42, 0x14, 0x5a, 0x12, 0x67, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x74, 0x65, 0x64, 0x2d,
	0x67, 0x6f, 0x2f, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
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
var file_store_changelog_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_store_changelog_proto_goTypes = []any{
	(ChangelogPayload_Type)(0),       // 0: bytebase.store.ChangelogPayload.Type
	(*ChangelogPayload)(nil),         // 1: bytebase.store.ChangelogPayload
	(*ChangedResources)(nil),         // 2: bytebase.store.ChangedResources
	(*ChangedResourceDatabase)(nil),  // 3: bytebase.store.ChangedResourceDatabase
	(*ChangedResourceSchema)(nil),    // 4: bytebase.store.ChangedResourceSchema
	(*ChangedResourceTable)(nil),     // 5: bytebase.store.ChangedResourceTable
	(*ChangedResourceView)(nil),      // 6: bytebase.store.ChangedResourceView
	(*ChangedResourceFunction)(nil),  // 7: bytebase.store.ChangedResourceFunction
	(*ChangedResourceProcedure)(nil), // 8: bytebase.store.ChangedResourceProcedure
	(*Range)(nil),                    // 9: bytebase.store.Range
}
var file_store_changelog_proto_depIdxs = []int32{
	2,  // 0: bytebase.store.ChangelogPayload.changed_resources:type_name -> bytebase.store.ChangedResources
	0,  // 1: bytebase.store.ChangelogPayload.type:type_name -> bytebase.store.ChangelogPayload.Type
	3,  // 2: bytebase.store.ChangedResources.databases:type_name -> bytebase.store.ChangedResourceDatabase
	4,  // 3: bytebase.store.ChangedResourceDatabase.schemas:type_name -> bytebase.store.ChangedResourceSchema
	5,  // 4: bytebase.store.ChangedResourceSchema.tables:type_name -> bytebase.store.ChangedResourceTable
	6,  // 5: bytebase.store.ChangedResourceSchema.views:type_name -> bytebase.store.ChangedResourceView
	7,  // 6: bytebase.store.ChangedResourceSchema.functions:type_name -> bytebase.store.ChangedResourceFunction
	8,  // 7: bytebase.store.ChangedResourceSchema.procedures:type_name -> bytebase.store.ChangedResourceProcedure
	9,  // 8: bytebase.store.ChangedResourceTable.ranges:type_name -> bytebase.store.Range
	9,  // 9: bytebase.store.ChangedResourceView.ranges:type_name -> bytebase.store.Range
	9,  // 10: bytebase.store.ChangedResourceFunction.ranges:type_name -> bytebase.store.Range
	9,  // 11: bytebase.store.ChangedResourceProcedure.ranges:type_name -> bytebase.store.Range
	12, // [12:12] is the sub-list for method output_type
	12, // [12:12] is the sub-list for method input_type
	12, // [12:12] is the sub-list for extension type_name
	12, // [12:12] is the sub-list for extension extendee
	0,  // [0:12] is the sub-list for field type_name
}

func init() { file_store_changelog_proto_init() }
func file_store_changelog_proto_init() {
	if File_store_changelog_proto != nil {
		return
	}
	file_store_common_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_store_changelog_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   8,
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

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.0
// 	protoc        (unknown)
// source: v1/instance_role_service.proto

package v1

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

type GetInstanceRoleRequest struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// The name of the role to retrieve.
	// Format: instances/{instance}/roles/{role name}
	// The role name is the unique name for the role.
	Name          string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetInstanceRoleRequest) Reset() {
	*x = GetInstanceRoleRequest{}
	mi := &file_v1_instance_role_service_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetInstanceRoleRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetInstanceRoleRequest) ProtoMessage() {}

func (x *GetInstanceRoleRequest) ProtoReflect() protoreflect.Message {
	mi := &file_v1_instance_role_service_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetInstanceRoleRequest.ProtoReflect.Descriptor instead.
func (*GetInstanceRoleRequest) Descriptor() ([]byte, []int) {
	return file_v1_instance_role_service_proto_rawDescGZIP(), []int{0}
}

func (x *GetInstanceRoleRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

type ListInstanceRolesRequest struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// The parent, which owns this collection of roles.
	// Format: instances/{instance}
	Parent string `protobuf:"bytes,1,opt,name=parent,proto3" json:"parent,omitempty"`
	// Not used.
	// The maximum number of roles to return. The service may return fewer than
	// this value.
	// If unspecified, at most 10 roles will be returned.
	// The maximum value is 1000; values above 1000 will be coerced to 1000.
	PageSize int32 `protobuf:"varint,2,opt,name=page_size,json=pageSize,proto3" json:"page_size,omitempty"`
	// Not used.
	// A page token, received from a previous `ListInstanceRoles` call.
	// Provide this to retrieve the subsequent page.
	//
	// When paginating, all other parameters provided to `ListInstanceRoles` must match
	// the call that provided the page token.
	PageToken string `protobuf:"bytes,3,opt,name=page_token,json=pageToken,proto3" json:"page_token,omitempty"`
	// Refresh will refresh and return the latest data.
	Refresh       bool `protobuf:"varint,4,opt,name=refresh,proto3" json:"refresh,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ListInstanceRolesRequest) Reset() {
	*x = ListInstanceRolesRequest{}
	mi := &file_v1_instance_role_service_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListInstanceRolesRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListInstanceRolesRequest) ProtoMessage() {}

func (x *ListInstanceRolesRequest) ProtoReflect() protoreflect.Message {
	mi := &file_v1_instance_role_service_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListInstanceRolesRequest.ProtoReflect.Descriptor instead.
func (*ListInstanceRolesRequest) Descriptor() ([]byte, []int) {
	return file_v1_instance_role_service_proto_rawDescGZIP(), []int{1}
}

func (x *ListInstanceRolesRequest) GetParent() string {
	if x != nil {
		return x.Parent
	}
	return ""
}

func (x *ListInstanceRolesRequest) GetPageSize() int32 {
	if x != nil {
		return x.PageSize
	}
	return 0
}

func (x *ListInstanceRolesRequest) GetPageToken() string {
	if x != nil {
		return x.PageToken
	}
	return ""
}

func (x *ListInstanceRolesRequest) GetRefresh() bool {
	if x != nil {
		return x.Refresh
	}
	return false
}

type ListInstanceRolesResponse struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// The roles from the specified request.
	Roles []*InstanceRole `protobuf:"bytes,1,rep,name=roles,proto3" json:"roles,omitempty"`
	// A token, which can be sent as `page_token` to retrieve the next page.
	// If this field is omitted, there are no subsequent pages.
	NextPageToken string `protobuf:"bytes,2,opt,name=next_page_token,json=nextPageToken,proto3" json:"next_page_token,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ListInstanceRolesResponse) Reset() {
	*x = ListInstanceRolesResponse{}
	mi := &file_v1_instance_role_service_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListInstanceRolesResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListInstanceRolesResponse) ProtoMessage() {}

func (x *ListInstanceRolesResponse) ProtoReflect() protoreflect.Message {
	mi := &file_v1_instance_role_service_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListInstanceRolesResponse.ProtoReflect.Descriptor instead.
func (*ListInstanceRolesResponse) Descriptor() ([]byte, []int) {
	return file_v1_instance_role_service_proto_rawDescGZIP(), []int{2}
}

func (x *ListInstanceRolesResponse) GetRoles() []*InstanceRole {
	if x != nil {
		return x.Roles
	}
	return nil
}

func (x *ListInstanceRolesResponse) GetNextPageToken() string {
	if x != nil {
		return x.NextPageToken
	}
	return ""
}

// InstanceRole is the API message for instance role.
type InstanceRole struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// The name of the role.
	// Format: instances/{instance}/roles/{role}
	// The role name is the unique name for the role.
	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	// The role name. It's unique within the instance.
	RoleName string `protobuf:"bytes,2,opt,name=role_name,json=roleName,proto3" json:"role_name,omitempty"`
	// The role password.
	Password *string `protobuf:"bytes,3,opt,name=password,proto3,oneof" json:"password,omitempty"`
	// The connection count limit for this role.
	ConnectionLimit *int32 `protobuf:"varint,4,opt,name=connection_limit,json=connectionLimit,proto3,oneof" json:"connection_limit,omitempty"`
	// The expiration for the role's password.
	ValidUntil *string `protobuf:"bytes,5,opt,name=valid_until,json=validUntil,proto3,oneof" json:"valid_until,omitempty"`
	// The role attribute.
	// For PostgreSQL, it containt super_user, no_inherit, create_role, create_db, can_login, replication and bypass_rls. Docs: https://www.postgresql.org/docs/current/role-attributes.html
	// For MySQL, it's the global privileges as GRANT statements, which means it only contains "GRANT ... ON *.* TO ...". Docs: https://dev.mysql.com/doc/refman/8.0/en/grant.html
	Attribute     *string `protobuf:"bytes,6,opt,name=attribute,proto3,oneof" json:"attribute,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *InstanceRole) Reset() {
	*x = InstanceRole{}
	mi := &file_v1_instance_role_service_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *InstanceRole) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InstanceRole) ProtoMessage() {}

func (x *InstanceRole) ProtoReflect() protoreflect.Message {
	mi := &file_v1_instance_role_service_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use InstanceRole.ProtoReflect.Descriptor instead.
func (*InstanceRole) Descriptor() ([]byte, []int) {
	return file_v1_instance_role_service_proto_rawDescGZIP(), []int{3}
}

func (x *InstanceRole) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *InstanceRole) GetRoleName() string {
	if x != nil {
		return x.RoleName
	}
	return ""
}

func (x *InstanceRole) GetPassword() string {
	if x != nil && x.Password != nil {
		return *x.Password
	}
	return ""
}

func (x *InstanceRole) GetConnectionLimit() int32 {
	if x != nil && x.ConnectionLimit != nil {
		return *x.ConnectionLimit
	}
	return 0
}

func (x *InstanceRole) GetValidUntil() string {
	if x != nil && x.ValidUntil != nil {
		return *x.ValidUntil
	}
	return ""
}

func (x *InstanceRole) GetAttribute() string {
	if x != nil && x.Attribute != nil {
		return *x.Attribute
	}
	return ""
}

var File_v1_instance_role_service_proto protoreflect.FileDescriptor

var file_v1_instance_role_service_proto_rawDesc = []byte{
	0x0a, 0x1e, 0x76, 0x31, 0x2f, 0x69, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x5f, 0x72, 0x6f,
	0x6c, 0x65, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x0b, 0x62, 0x79, 0x74, 0x65, 0x62, 0x61, 0x73, 0x65, 0x2e, 0x76, 0x31, 0x1a, 0x1c, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x17, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69,
	0x2f, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x5f, 0x62, 0x65, 0x68, 0x61, 0x76, 0x69, 0x6f, 0x72, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x19, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70,
	0x69, 0x2f, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x1a, 0x13, 0x76, 0x31, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x50, 0x0a, 0x16, 0x47, 0x65, 0x74, 0x49, 0x6e, 0x73, 0x74,
	0x61, 0x6e, 0x63, 0x65, 0x52, 0x6f, 0x6c, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x36, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x22, 0xe2,
	0x41, 0x01, 0x02, 0xfa, 0x41, 0x1b, 0x0a, 0x19, 0x62, 0x79, 0x74, 0x65, 0x62, 0x61, 0x73, 0x65,
	0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x49, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x52, 0x6f, 0x6c,
	0x65, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x22, 0xa8, 0x01, 0x0a, 0x18, 0x4c, 0x69, 0x73, 0x74,
	0x49, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x52, 0x6f, 0x6c, 0x65, 0x73, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x36, 0x0a, 0x06, 0x70, 0x61, 0x72, 0x65, 0x6e, 0x74, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x42, 0x1e, 0xe2, 0x41, 0x01, 0x02, 0xfa, 0x41, 0x17, 0x0a, 0x15, 0x62,
	0x79, 0x74, 0x65, 0x62, 0x61, 0x73, 0x65, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x49, 0x6e, 0x73, 0x74,
	0x61, 0x6e, 0x63, 0x65, 0x52, 0x06, 0x70, 0x61, 0x72, 0x65, 0x6e, 0x74, 0x12, 0x1b, 0x0a, 0x09,
	0x70, 0x61, 0x67, 0x65, 0x5f, 0x73, 0x69, 0x7a, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52,
	0x08, 0x70, 0x61, 0x67, 0x65, 0x53, 0x69, 0x7a, 0x65, 0x12, 0x1d, 0x0a, 0x0a, 0x70, 0x61, 0x67,
	0x65, 0x5f, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x70,
	0x61, 0x67, 0x65, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x18, 0x0a, 0x07, 0x72, 0x65, 0x66, 0x72,
	0x65, 0x73, 0x68, 0x18, 0x04, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x72, 0x65, 0x66, 0x72, 0x65,
	0x73, 0x68, 0x22, 0x74, 0x0a, 0x19, 0x4c, 0x69, 0x73, 0x74, 0x49, 0x6e, 0x73, 0x74, 0x61, 0x6e,
	0x63, 0x65, 0x52, 0x6f, 0x6c, 0x65, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x2f, 0x0a, 0x05, 0x72, 0x6f, 0x6c, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x19,
	0x2e, 0x62, 0x79, 0x74, 0x65, 0x62, 0x61, 0x73, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x49, 0x6e, 0x73,
	0x74, 0x61, 0x6e, 0x63, 0x65, 0x52, 0x6f, 0x6c, 0x65, 0x52, 0x05, 0x72, 0x6f, 0x6c, 0x65, 0x73,
	0x12, 0x26, 0x0a, 0x0f, 0x6e, 0x65, 0x78, 0x74, 0x5f, 0x70, 0x61, 0x67, 0x65, 0x5f, 0x74, 0x6f,
	0x6b, 0x65, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x6e, 0x65, 0x78, 0x74, 0x50,
	0x61, 0x67, 0x65, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x22, 0xe2, 0x02, 0x0a, 0x0c, 0x49, 0x6e, 0x73,
	0x74, 0x61, 0x6e, 0x63, 0x65, 0x52, 0x6f, 0x6c, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d,
	0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x1b, 0x0a,
	0x09, 0x72, 0x6f, 0x6c, 0x65, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x08, 0x72, 0x6f, 0x6c, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x25, 0x0a, 0x08, 0x70, 0x61,
	0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x42, 0x04, 0xe2, 0x41,
	0x01, 0x04, 0x48, 0x00, 0x52, 0x08, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x88, 0x01,
	0x01, 0x12, 0x2e, 0x0a, 0x10, 0x63, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x5f,
	0x6c, 0x69, 0x6d, 0x69, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x05, 0x48, 0x01, 0x52, 0x0f, 0x63,
	0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x4c, 0x69, 0x6d, 0x69, 0x74, 0x88, 0x01,
	0x01, 0x12, 0x24, 0x0a, 0x0b, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x5f, 0x75, 0x6e, 0x74, 0x69, 0x6c,
	0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x48, 0x02, 0x52, 0x0a, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x55,
	0x6e, 0x74, 0x69, 0x6c, 0x88, 0x01, 0x01, 0x12, 0x21, 0x0a, 0x09, 0x61, 0x74, 0x74, 0x72, 0x69,
	0x62, 0x75, 0x74, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x48, 0x03, 0x52, 0x09, 0x61, 0x74,
	0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x65, 0x88, 0x01, 0x01, 0x3a, 0x41, 0xea, 0x41, 0x3e, 0x0a,
	0x19, 0x62, 0x79, 0x74, 0x65, 0x62, 0x61, 0x73, 0x65, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x49, 0x6e,
	0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x52, 0x6f, 0x6c, 0x65, 0x12, 0x21, 0x69, 0x6e, 0x73, 0x74,
	0x61, 0x6e, 0x63, 0x65, 0x73, 0x2f, 0x7b, 0x69, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x7d,
	0x2f, 0x72, 0x6f, 0x6c, 0x65, 0x73, 0x2f, 0x7b, 0x72, 0x6f, 0x6c, 0x65, 0x7d, 0x42, 0x0b, 0x0a,
	0x09, 0x5f, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x42, 0x13, 0x0a, 0x11, 0x5f, 0x63,
	0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x42,
	0x0e, 0x0a, 0x0c, 0x5f, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x5f, 0x75, 0x6e, 0x74, 0x69, 0x6c, 0x42,
	0x0c, 0x0a, 0x0a, 0x5f, 0x61, 0x74, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x65, 0x32, 0xe6, 0x02,
	0x0a, 0x13, 0x49, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x52, 0x6f, 0x6c, 0x65, 0x53, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x9c, 0x01, 0x0a, 0x0f, 0x47, 0x65, 0x74, 0x49, 0x6e, 0x73,
	0x74, 0x61, 0x6e, 0x63, 0x65, 0x52, 0x6f, 0x6c, 0x65, 0x12, 0x23, 0x2e, 0x62, 0x79, 0x74, 0x65,
	0x62, 0x61, 0x73, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x47, 0x65, 0x74, 0x49, 0x6e, 0x73, 0x74, 0x61,
	0x6e, 0x63, 0x65, 0x52, 0x6f, 0x6c, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x19,
	0x2e, 0x62, 0x79, 0x74, 0x65, 0x62, 0x61, 0x73, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x49, 0x6e, 0x73,
	0x74, 0x61, 0x6e, 0x63, 0x65, 0x52, 0x6f, 0x6c, 0x65, 0x22, 0x49, 0xda, 0x41, 0x04, 0x6e, 0x61,
	0x6d, 0x65, 0x8a, 0xea, 0x30, 0x14, 0x62, 0x62, 0x2e, 0x69, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63,
	0x65, 0x52, 0x6f, 0x6c, 0x65, 0x73, 0x2e, 0x67, 0x65, 0x74, 0x90, 0xea, 0x30, 0x01, 0x82, 0xd3,
	0xe4, 0x93, 0x02, 0x20, 0x12, 0x1e, 0x2f, 0x76, 0x31, 0x2f, 0x7b, 0x6e, 0x61, 0x6d, 0x65, 0x3d,
	0x69, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x73, 0x2f, 0x2a, 0x2f, 0x72, 0x6f, 0x6c, 0x65,
	0x73, 0x2f, 0x2a, 0x7d, 0x12, 0xaf, 0x01, 0x0a, 0x11, 0x4c, 0x69, 0x73, 0x74, 0x49, 0x6e, 0x73,
	0x74, 0x61, 0x6e, 0x63, 0x65, 0x52, 0x6f, 0x6c, 0x65, 0x73, 0x12, 0x25, 0x2e, 0x62, 0x79, 0x74,
	0x65, 0x62, 0x61, 0x73, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x49, 0x6e, 0x73,
	0x74, 0x61, 0x6e, 0x63, 0x65, 0x52, 0x6f, 0x6c, 0x65, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x26, 0x2e, 0x62, 0x79, 0x74, 0x65, 0x62, 0x61, 0x73, 0x65, 0x2e, 0x76, 0x31, 0x2e,
	0x4c, 0x69, 0x73, 0x74, 0x49, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x52, 0x6f, 0x6c, 0x65,
	0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x4b, 0xda, 0x41, 0x06, 0x70, 0x61,
	0x72, 0x65, 0x6e, 0x74, 0x8a, 0xea, 0x30, 0x14, 0x62, 0x62, 0x2e, 0x69, 0x6e, 0x73, 0x74, 0x61,
	0x6e, 0x63, 0x65, 0x52, 0x6f, 0x6c, 0x65, 0x73, 0x2e, 0x67, 0x65, 0x74, 0x90, 0xea, 0x30, 0x01,
	0x82, 0xd3, 0xe4, 0x93, 0x02, 0x20, 0x12, 0x1e, 0x2f, 0x76, 0x31, 0x2f, 0x7b, 0x70, 0x61, 0x72,
	0x65, 0x6e, 0x74, 0x3d, 0x69, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x73, 0x2f, 0x2a, 0x7d,
	0x2f, 0x72, 0x6f, 0x6c, 0x65, 0x73, 0x42, 0x11, 0x5a, 0x0f, 0x67, 0x65, 0x6e, 0x65, 0x72, 0x61,
	0x74, 0x65, 0x64, 0x2d, 0x67, 0x6f, 0x2f, 0x76, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_v1_instance_role_service_proto_rawDescOnce sync.Once
	file_v1_instance_role_service_proto_rawDescData = file_v1_instance_role_service_proto_rawDesc
)

func file_v1_instance_role_service_proto_rawDescGZIP() []byte {
	file_v1_instance_role_service_proto_rawDescOnce.Do(func() {
		file_v1_instance_role_service_proto_rawDescData = protoimpl.X.CompressGZIP(file_v1_instance_role_service_proto_rawDescData)
	})
	return file_v1_instance_role_service_proto_rawDescData
}

var file_v1_instance_role_service_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_v1_instance_role_service_proto_goTypes = []any{
	(*GetInstanceRoleRequest)(nil),    // 0: bytebase.v1.GetInstanceRoleRequest
	(*ListInstanceRolesRequest)(nil),  // 1: bytebase.v1.ListInstanceRolesRequest
	(*ListInstanceRolesResponse)(nil), // 2: bytebase.v1.ListInstanceRolesResponse
	(*InstanceRole)(nil),              // 3: bytebase.v1.InstanceRole
}
var file_v1_instance_role_service_proto_depIdxs = []int32{
	3, // 0: bytebase.v1.ListInstanceRolesResponse.roles:type_name -> bytebase.v1.InstanceRole
	0, // 1: bytebase.v1.InstanceRoleService.GetInstanceRole:input_type -> bytebase.v1.GetInstanceRoleRequest
	1, // 2: bytebase.v1.InstanceRoleService.ListInstanceRoles:input_type -> bytebase.v1.ListInstanceRolesRequest
	3, // 3: bytebase.v1.InstanceRoleService.GetInstanceRole:output_type -> bytebase.v1.InstanceRole
	2, // 4: bytebase.v1.InstanceRoleService.ListInstanceRoles:output_type -> bytebase.v1.ListInstanceRolesResponse
	3, // [3:5] is the sub-list for method output_type
	1, // [1:3] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_v1_instance_role_service_proto_init() }
func file_v1_instance_role_service_proto_init() {
	if File_v1_instance_role_service_proto != nil {
		return
	}
	file_v1_annotation_proto_init()
	file_v1_instance_role_service_proto_msgTypes[3].OneofWrappers = []any{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_v1_instance_role_service_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_v1_instance_role_service_proto_goTypes,
		DependencyIndexes: file_v1_instance_role_service_proto_depIdxs,
		MessageInfos:      file_v1_instance_role_service_proto_msgTypes,
	}.Build()
	File_v1_instance_role_service_proto = out.File
	file_v1_instance_role_service_proto_rawDesc = nil
	file_v1_instance_role_service_proto_goTypes = nil
	file_v1_instance_role_service_proto_depIdxs = nil
}

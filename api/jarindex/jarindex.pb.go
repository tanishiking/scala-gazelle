// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.14.0
// source: api/jarindex/jarindex.proto

package jarindex

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

type JarIndex struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	JarFile    []*JarFile `protobuf:"bytes,2,rep,name=jar_file,json=jarFiles,proto3" json:"jar_file,omitempty"`
	Predefined []string   `protobuf:"bytes,3,rep,name=predefined,proto3" json:"predefined,omitempty"`
	Preferred  []string   `protobuf:"bytes,4,rep,name=preferred,proto3" json:"preferred,omitempty"`
}

func (x *JarIndex) Reset() {
	*x = JarIndex{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_jarindex_jarindex_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *JarIndex) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*JarIndex) ProtoMessage() {}

func (x *JarIndex) ProtoReflect() protoreflect.Message {
	mi := &file_api_jarindex_jarindex_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use JarIndex.ProtoReflect.Descriptor instead.
func (*JarIndex) Descriptor() ([]byte, []int) {
	return file_api_jarindex_jarindex_proto_rawDescGZIP(), []int{0}
}

func (x *JarIndex) GetJarFile() []*JarFile {
	if x != nil {
		return x.JarFile
	}
	return nil
}

func (x *JarIndex) GetPredefined() []string {
	if x != nil {
		return x.Predefined
	}
	return nil
}

func (x *JarIndex) GetPreferred() []string {
	if x != nil {
		return x.Preferred
	}
	return nil
}

type JarFile struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Filename    string            `protobuf:"bytes,1,opt,name=filename,proto3" json:"filename,omitempty"`
	Label       string            `protobuf:"bytes,2,opt,name=label,proto3" json:"label,omitempty"`
	ClassFile   []*ClassFile      `protobuf:"bytes,3,rep,name=class_file,json=files,proto3" json:"class_file,omitempty"`
	Symbols     []string          `protobuf:"bytes,4,rep,name=symbols,proto3" json:"symbols,omitempty"`
	ClassName   []string          `protobuf:"bytes,5,rep,name=class_name,json=classes,proto3" json:"class_name,omitempty"`
	PackageName []string          `protobuf:"bytes,6,rep,name=package_name,json=packages,proto3" json:"package_name,omitempty"`
	Extends     map[string]string `protobuf:"bytes,7,rep,name=extends,proto3" json:"extends,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *JarFile) Reset() {
	*x = JarFile{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_jarindex_jarindex_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *JarFile) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*JarFile) ProtoMessage() {}

func (x *JarFile) ProtoReflect() protoreflect.Message {
	mi := &file_api_jarindex_jarindex_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use JarFile.ProtoReflect.Descriptor instead.
func (*JarFile) Descriptor() ([]byte, []int) {
	return file_api_jarindex_jarindex_proto_rawDescGZIP(), []int{1}
}

func (x *JarFile) GetFilename() string {
	if x != nil {
		return x.Filename
	}
	return ""
}

func (x *JarFile) GetLabel() string {
	if x != nil {
		return x.Label
	}
	return ""
}

func (x *JarFile) GetClassFile() []*ClassFile {
	if x != nil {
		return x.ClassFile
	}
	return nil
}

func (x *JarFile) GetSymbols() []string {
	if x != nil {
		return x.Symbols
	}
	return nil
}

func (x *JarFile) GetClassName() []string {
	if x != nil {
		return x.ClassName
	}
	return nil
}

func (x *JarFile) GetPackageName() []string {
	if x != nil {
		return x.PackageName
	}
	return nil
}

func (x *JarFile) GetExtends() map[string]string {
	if x != nil {
		return x.Extends
	}
	return nil
}

type ClassFile struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name         string         `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Classes      []int32        `protobuf:"varint,2,rep,packed,name=classes,proto3" json:"classes,omitempty"`
	Symbols      []string       `protobuf:"bytes,3,rep,name=symbols,proto3" json:"symbols,omitempty"`
	Superclasses []string       `protobuf:"bytes,4,rep,name=superclasses,proto3" json:"superclasses,omitempty"`
	Interfaces   []string       `protobuf:"bytes,5,rep,name=interfaces,proto3" json:"interfaces,omitempty"`
	Fields       []*ClassField  `protobuf:"bytes,6,rep,name=fields,proto3" json:"fields,omitempty"`
	Methods      []*ClassMethod `protobuf:"bytes,7,rep,name=methods,proto3" json:"methods,omitempty"`
}

func (x *ClassFile) Reset() {
	*x = ClassFile{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_jarindex_jarindex_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ClassFile) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ClassFile) ProtoMessage() {}

func (x *ClassFile) ProtoReflect() protoreflect.Message {
	mi := &file_api_jarindex_jarindex_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ClassFile.ProtoReflect.Descriptor instead.
func (*ClassFile) Descriptor() ([]byte, []int) {
	return file_api_jarindex_jarindex_proto_rawDescGZIP(), []int{2}
}

func (x *ClassFile) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *ClassFile) GetClasses() []int32 {
	if x != nil {
		return x.Classes
	}
	return nil
}

func (x *ClassFile) GetSymbols() []string {
	if x != nil {
		return x.Symbols
	}
	return nil
}

func (x *ClassFile) GetSuperclasses() []string {
	if x != nil {
		return x.Superclasses
	}
	return nil
}

func (x *ClassFile) GetInterfaces() []string {
	if x != nil {
		return x.Interfaces
	}
	return nil
}

func (x *ClassFile) GetFields() []*ClassField {
	if x != nil {
		return x.Fields
	}
	return nil
}

func (x *ClassFile) GetMethods() []*ClassMethod {
	if x != nil {
		return x.Methods
	}
	return nil
}

type ClassField struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name string     `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Type *ClassType `protobuf:"bytes,2,opt,name=type,proto3" json:"type,omitempty"`
}

func (x *ClassField) Reset() {
	*x = ClassField{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_jarindex_jarindex_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ClassField) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ClassField) ProtoMessage() {}

func (x *ClassField) ProtoReflect() protoreflect.Message {
	mi := &file_api_jarindex_jarindex_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ClassField.ProtoReflect.Descriptor instead.
func (*ClassField) Descriptor() ([]byte, []int) {
	return file_api_jarindex_jarindex_proto_rawDescGZIP(), []int{3}
}

func (x *ClassField) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *ClassField) GetType() *ClassType {
	if x != nil {
		return x.Type
	}
	return nil
}

type ClassType struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Kind  string `protobuf:"bytes,1,opt,name=kind,proto3" json:"kind,omitempty"`
	Value string `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *ClassType) Reset() {
	*x = ClassType{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_jarindex_jarindex_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ClassType) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ClassType) ProtoMessage() {}

func (x *ClassType) ProtoReflect() protoreflect.Message {
	mi := &file_api_jarindex_jarindex_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ClassType.ProtoReflect.Descriptor instead.
func (*ClassType) Descriptor() ([]byte, []int) {
	return file_api_jarindex_jarindex_proto_rawDescGZIP(), []int{4}
}

func (x *ClassType) GetKind() string {
	if x != nil {
		return x.Kind
	}
	return ""
}

func (x *ClassType) GetValue() string {
	if x != nil {
		return x.Value
	}
	return ""
}

type ClassMethod struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name    string              `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Returns *ClassType          `protobuf:"bytes,2,opt,name=returns,proto3" json:"returns,omitempty"`
	Params  []*ClassMethodParam `protobuf:"bytes,3,rep,name=params,proto3" json:"params,omitempty"`
	Types   []*ClassType        `protobuf:"bytes,4,rep,name=types,proto3" json:"types,omitempty"`
	Throws  []*ClassType        `protobuf:"bytes,5,rep,name=throws,proto3" json:"throws,omitempty"`
}

func (x *ClassMethod) Reset() {
	*x = ClassMethod{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_jarindex_jarindex_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ClassMethod) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ClassMethod) ProtoMessage() {}

func (x *ClassMethod) ProtoReflect() protoreflect.Message {
	mi := &file_api_jarindex_jarindex_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ClassMethod.ProtoReflect.Descriptor instead.
func (*ClassMethod) Descriptor() ([]byte, []int) {
	return file_api_jarindex_jarindex_proto_rawDescGZIP(), []int{5}
}

func (x *ClassMethod) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *ClassMethod) GetReturns() *ClassType {
	if x != nil {
		return x.Returns
	}
	return nil
}

func (x *ClassMethod) GetParams() []*ClassMethodParam {
	if x != nil {
		return x.Params
	}
	return nil
}

func (x *ClassMethod) GetTypes() []*ClassType {
	if x != nil {
		return x.Types
	}
	return nil
}

func (x *ClassMethod) GetThrows() []*ClassType {
	if x != nil {
		return x.Throws
	}
	return nil
}

type ClassMethodParam struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Returns *ClassType `protobuf:"bytes,1,opt,name=returns,proto3" json:"returns,omitempty"`
}

func (x *ClassMethodParam) Reset() {
	*x = ClassMethodParam{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_jarindex_jarindex_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ClassMethodParam) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ClassMethodParam) ProtoMessage() {}

func (x *ClassMethodParam) ProtoReflect() protoreflect.Message {
	mi := &file_api_jarindex_jarindex_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ClassMethodParam.ProtoReflect.Descriptor instead.
func (*ClassMethodParam) Descriptor() ([]byte, []int) {
	return file_api_jarindex_jarindex_proto_rawDescGZIP(), []int{6}
}

func (x *ClassMethodParam) GetReturns() *ClassType {
	if x != nil {
		return x.Returns
	}
	return nil
}

var File_api_jarindex_jarindex_proto protoreflect.FileDescriptor

var file_api_jarindex_jarindex_proto_rawDesc = []byte{
	0x0a, 0x1b, 0x61, 0x70, 0x69, 0x2f, 0x6a, 0x61, 0x72, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x2f, 0x6a,
	0x61, 0x72, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0c, 0x61,
	0x70, 0x69, 0x2e, 0x6a, 0x61, 0x72, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x22, 0x7b, 0x0a, 0x08, 0x4a,
	0x61, 0x72, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x12, 0x31, 0x0a, 0x08, 0x6a, 0x61, 0x72, 0x5f, 0x66,
	0x69, 0x6c, 0x65, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x61, 0x70, 0x69, 0x2e,
	0x6a, 0x61, 0x72, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x2e, 0x4a, 0x61, 0x72, 0x46, 0x69, 0x6c, 0x65,
	0x52, 0x08, 0x6a, 0x61, 0x72, 0x46, 0x69, 0x6c, 0x65, 0x73, 0x12, 0x1e, 0x0a, 0x0a, 0x70, 0x72,
	0x65, 0x64, 0x65, 0x66, 0x69, 0x6e, 0x65, 0x64, 0x18, 0x03, 0x20, 0x03, 0x28, 0x09, 0x52, 0x0a,
	0x70, 0x72, 0x65, 0x64, 0x65, 0x66, 0x69, 0x6e, 0x65, 0x64, 0x12, 0x1c, 0x0a, 0x09, 0x70, 0x72,
	0x65, 0x66, 0x65, 0x72, 0x72, 0x65, 0x64, 0x18, 0x04, 0x20, 0x03, 0x28, 0x09, 0x52, 0x09, 0x70,
	0x72, 0x65, 0x66, 0x65, 0x72, 0x72, 0x65, 0x64, 0x22, 0xc0, 0x02, 0x0a, 0x07, 0x4a, 0x61, 0x72,
	0x46, 0x69, 0x6c, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x66, 0x69, 0x6c, 0x65, 0x6e, 0x61, 0x6d, 0x65,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x66, 0x69, 0x6c, 0x65, 0x6e, 0x61, 0x6d, 0x65,
	0x12, 0x14, 0x0a, 0x05, 0x6c, 0x61, 0x62, 0x65, 0x6c, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x05, 0x6c, 0x61, 0x62, 0x65, 0x6c, 0x12, 0x32, 0x0a, 0x0a, 0x63, 0x6c, 0x61, 0x73, 0x73, 0x5f,
	0x66, 0x69, 0x6c, 0x65, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x61, 0x70, 0x69,
	0x2e, 0x6a, 0x61, 0x72, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x2e, 0x43, 0x6c, 0x61, 0x73, 0x73, 0x46,
	0x69, 0x6c, 0x65, 0x52, 0x05, 0x66, 0x69, 0x6c, 0x65, 0x73, 0x12, 0x18, 0x0a, 0x07, 0x73, 0x79,
	0x6d, 0x62, 0x6f, 0x6c, 0x73, 0x18, 0x04, 0x20, 0x03, 0x28, 0x09, 0x52, 0x07, 0x73, 0x79, 0x6d,
	0x62, 0x6f, 0x6c, 0x73, 0x12, 0x1b, 0x0a, 0x0a, 0x63, 0x6c, 0x61, 0x73, 0x73, 0x5f, 0x6e, 0x61,
	0x6d, 0x65, 0x18, 0x05, 0x20, 0x03, 0x28, 0x09, 0x52, 0x07, 0x63, 0x6c, 0x61, 0x73, 0x73, 0x65,
	0x73, 0x12, 0x1e, 0x0a, 0x0c, 0x70, 0x61, 0x63, 0x6b, 0x61, 0x67, 0x65, 0x5f, 0x6e, 0x61, 0x6d,
	0x65, 0x18, 0x06, 0x20, 0x03, 0x28, 0x09, 0x52, 0x08, 0x70, 0x61, 0x63, 0x6b, 0x61, 0x67, 0x65,
	0x73, 0x12, 0x3c, 0x0a, 0x07, 0x65, 0x78, 0x74, 0x65, 0x6e, 0x64, 0x73, 0x18, 0x07, 0x20, 0x03,
	0x28, 0x0b, 0x32, 0x22, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x6a, 0x61, 0x72, 0x69, 0x6e, 0x64, 0x65,
	0x78, 0x2e, 0x4a, 0x61, 0x72, 0x46, 0x69, 0x6c, 0x65, 0x2e, 0x45, 0x78, 0x74, 0x65, 0x6e, 0x64,
	0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x07, 0x65, 0x78, 0x74, 0x65, 0x6e, 0x64, 0x73, 0x1a,
	0x3a, 0x0a, 0x0c, 0x45, 0x78, 0x74, 0x65, 0x6e, 0x64, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12,
	0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65,
	0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x22, 0xfe, 0x01, 0x0a, 0x09,
	0x43, 0x6c, 0x61, 0x73, 0x73, 0x46, 0x69, 0x6c, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d,
	0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x18, 0x0a,
	0x07, 0x63, 0x6c, 0x61, 0x73, 0x73, 0x65, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x05, 0x52, 0x07,
	0x63, 0x6c, 0x61, 0x73, 0x73, 0x65, 0x73, 0x12, 0x18, 0x0a, 0x07, 0x73, 0x79, 0x6d, 0x62, 0x6f,
	0x6c, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x09, 0x52, 0x07, 0x73, 0x79, 0x6d, 0x62, 0x6f, 0x6c,
	0x73, 0x12, 0x22, 0x0a, 0x0c, 0x73, 0x75, 0x70, 0x65, 0x72, 0x63, 0x6c, 0x61, 0x73, 0x73, 0x65,
	0x73, 0x18, 0x04, 0x20, 0x03, 0x28, 0x09, 0x52, 0x0c, 0x73, 0x75, 0x70, 0x65, 0x72, 0x63, 0x6c,
	0x61, 0x73, 0x73, 0x65, 0x73, 0x12, 0x1e, 0x0a, 0x0a, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x66, 0x61,
	0x63, 0x65, 0x73, 0x18, 0x05, 0x20, 0x03, 0x28, 0x09, 0x52, 0x0a, 0x69, 0x6e, 0x74, 0x65, 0x72,
	0x66, 0x61, 0x63, 0x65, 0x73, 0x12, 0x30, 0x0a, 0x06, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x73, 0x18,
	0x06, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x18, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x6a, 0x61, 0x72, 0x69,
	0x6e, 0x64, 0x65, 0x78, 0x2e, 0x43, 0x6c, 0x61, 0x73, 0x73, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x52,
	0x06, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x73, 0x12, 0x33, 0x0a, 0x07, 0x6d, 0x65, 0x74, 0x68, 0x6f,
	0x64, 0x73, 0x18, 0x07, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x6a,
	0x61, 0x72, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x2e, 0x43, 0x6c, 0x61, 0x73, 0x73, 0x4d, 0x65, 0x74,
	0x68, 0x6f, 0x64, 0x52, 0x07, 0x6d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x73, 0x22, 0x4d, 0x0a, 0x0a,
	0x43, 0x6c, 0x61, 0x73, 0x73, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61,
	0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x2b,
	0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x61,
	0x70, 0x69, 0x2e, 0x6a, 0x61, 0x72, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x2e, 0x43, 0x6c, 0x61, 0x73,
	0x73, 0x54, 0x79, 0x70, 0x65, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x22, 0x35, 0x0a, 0x09, 0x43,
	0x6c, 0x61, 0x73, 0x73, 0x54, 0x79, 0x70, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x6b, 0x69, 0x6e, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6b, 0x69, 0x6e, 0x64, 0x12, 0x14, 0x0a, 0x05,
	0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c,
	0x75, 0x65, 0x22, 0xec, 0x01, 0x0a, 0x0b, 0x43, 0x6c, 0x61, 0x73, 0x73, 0x4d, 0x65, 0x74, 0x68,
	0x6f, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x31, 0x0a, 0x07, 0x72, 0x65, 0x74, 0x75, 0x72, 0x6e,
	0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x6a, 0x61,
	0x72, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x2e, 0x43, 0x6c, 0x61, 0x73, 0x73, 0x54, 0x79, 0x70, 0x65,
	0x52, 0x07, 0x72, 0x65, 0x74, 0x75, 0x72, 0x6e, 0x73, 0x12, 0x36, 0x0a, 0x06, 0x70, 0x61, 0x72,
	0x61, 0x6d, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x61, 0x70, 0x69, 0x2e,
	0x6a, 0x61, 0x72, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x2e, 0x43, 0x6c, 0x61, 0x73, 0x73, 0x4d, 0x65,
	0x74, 0x68, 0x6f, 0x64, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x52, 0x06, 0x70, 0x61, 0x72, 0x61, 0x6d,
	0x73, 0x12, 0x2d, 0x0a, 0x05, 0x74, 0x79, 0x70, 0x65, 0x73, 0x18, 0x04, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x17, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x6a, 0x61, 0x72, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x2e,
	0x43, 0x6c, 0x61, 0x73, 0x73, 0x54, 0x79, 0x70, 0x65, 0x52, 0x05, 0x74, 0x79, 0x70, 0x65, 0x73,
	0x12, 0x2f, 0x0a, 0x06, 0x74, 0x68, 0x72, 0x6f, 0x77, 0x73, 0x18, 0x05, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x17, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x6a, 0x61, 0x72, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x2e,
	0x43, 0x6c, 0x61, 0x73, 0x73, 0x54, 0x79, 0x70, 0x65, 0x52, 0x06, 0x74, 0x68, 0x72, 0x6f, 0x77,
	0x73, 0x22, 0x45, 0x0a, 0x10, 0x43, 0x6c, 0x61, 0x73, 0x73, 0x4d, 0x65, 0x74, 0x68, 0x6f, 0x64,
	0x50, 0x61, 0x72, 0x61, 0x6d, 0x12, 0x31, 0x0a, 0x07, 0x72, 0x65, 0x74, 0x75, 0x72, 0x6e, 0x73,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x6a, 0x61, 0x72,
	0x69, 0x6e, 0x64, 0x65, 0x78, 0x2e, 0x43, 0x6c, 0x61, 0x73, 0x73, 0x54, 0x79, 0x70, 0x65, 0x52,
	0x07, 0x72, 0x65, 0x74, 0x75, 0x72, 0x6e, 0x73, 0x42, 0x58, 0x0a, 0x26, 0x62, 0x75, 0x69, 0x6c,
	0x64, 0x2e, 0x73, 0x74, 0x61, 0x63, 0x6b, 0x2e, 0x73, 0x63, 0x61, 0x6c, 0x61, 0x2e, 0x67, 0x61,
	0x7a, 0x65, 0x6c, 0x6c, 0x65, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x6a, 0x61, 0x72, 0x69, 0x6e, 0x64,
	0x65, 0x78, 0x50, 0x01, 0x5a, 0x2c, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d,
	0x2f, 0x73, 0x74, 0x61, 0x63, 0x6b, 0x62, 0x2f, 0x73, 0x63, 0x61, 0x6c, 0x61, 0x2d, 0x67, 0x61,
	0x7a, 0x65, 0x6c, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x6a, 0x61, 0x72, 0x69, 0x6e, 0x64,
	0x65, 0x78, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_api_jarindex_jarindex_proto_rawDescOnce sync.Once
	file_api_jarindex_jarindex_proto_rawDescData = file_api_jarindex_jarindex_proto_rawDesc
)

func file_api_jarindex_jarindex_proto_rawDescGZIP() []byte {
	file_api_jarindex_jarindex_proto_rawDescOnce.Do(func() {
		file_api_jarindex_jarindex_proto_rawDescData = protoimpl.X.CompressGZIP(file_api_jarindex_jarindex_proto_rawDescData)
	})
	return file_api_jarindex_jarindex_proto_rawDescData
}

var file_api_jarindex_jarindex_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_api_jarindex_jarindex_proto_goTypes = []interface{}{
	(*JarIndex)(nil),         // 0: api.jarindex.JarIndex
	(*JarFile)(nil),          // 1: api.jarindex.JarFile
	(*ClassFile)(nil),        // 2: api.jarindex.ClassFile
	(*ClassField)(nil),       // 3: api.jarindex.ClassField
	(*ClassType)(nil),        // 4: api.jarindex.ClassType
	(*ClassMethod)(nil),      // 5: api.jarindex.ClassMethod
	(*ClassMethodParam)(nil), // 6: api.jarindex.ClassMethodParam
	nil,                      // 7: api.jarindex.JarFile.ExtendsEntry
}
var file_api_jarindex_jarindex_proto_depIdxs = []int32{
	1,  // 0: api.jarindex.JarIndex.jar_file:type_name -> api.jarindex.JarFile
	2,  // 1: api.jarindex.JarFile.class_file:type_name -> api.jarindex.ClassFile
	7,  // 2: api.jarindex.JarFile.extends:type_name -> api.jarindex.JarFile.ExtendsEntry
	3,  // 3: api.jarindex.ClassFile.fields:type_name -> api.jarindex.ClassField
	5,  // 4: api.jarindex.ClassFile.methods:type_name -> api.jarindex.ClassMethod
	4,  // 5: api.jarindex.ClassField.type:type_name -> api.jarindex.ClassType
	4,  // 6: api.jarindex.ClassMethod.returns:type_name -> api.jarindex.ClassType
	6,  // 7: api.jarindex.ClassMethod.params:type_name -> api.jarindex.ClassMethodParam
	4,  // 8: api.jarindex.ClassMethod.types:type_name -> api.jarindex.ClassType
	4,  // 9: api.jarindex.ClassMethod.throws:type_name -> api.jarindex.ClassType
	4,  // 10: api.jarindex.ClassMethodParam.returns:type_name -> api.jarindex.ClassType
	11, // [11:11] is the sub-list for method output_type
	11, // [11:11] is the sub-list for method input_type
	11, // [11:11] is the sub-list for extension type_name
	11, // [11:11] is the sub-list for extension extendee
	0,  // [0:11] is the sub-list for field type_name
}

func init() { file_api_jarindex_jarindex_proto_init() }
func file_api_jarindex_jarindex_proto_init() {
	if File_api_jarindex_jarindex_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_api_jarindex_jarindex_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*JarIndex); i {
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
		file_api_jarindex_jarindex_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*JarFile); i {
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
		file_api_jarindex_jarindex_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ClassFile); i {
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
		file_api_jarindex_jarindex_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ClassField); i {
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
		file_api_jarindex_jarindex_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ClassType); i {
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
		file_api_jarindex_jarindex_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ClassMethod); i {
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
		file_api_jarindex_jarindex_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ClassMethodParam); i {
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
			RawDescriptor: file_api_jarindex_jarindex_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_api_jarindex_jarindex_proto_goTypes,
		DependencyIndexes: file_api_jarindex_jarindex_proto_depIdxs,
		MessageInfos:      file_api_jarindex_jarindex_proto_msgTypes,
	}.Build()
	File_api_jarindex_jarindex_proto = out.File
	file_api_jarindex_jarindex_proto_rawDesc = nil
	file_api_jarindex_jarindex_proto_goTypes = nil
	file_api_jarindex_jarindex_proto_depIdxs = nil
}

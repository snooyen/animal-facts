// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.17.3
// source: facts.proto

package pb

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type CreateFactRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Animal string `protobuf:"bytes,1,opt,name=Animal,proto3" json:"Animal,omitempty"`
	Fact   string `protobuf:"bytes,2,opt,name=Fact,proto3" json:"Fact,omitempty"`
}

func (x *CreateFactRequest) Reset() {
	*x = CreateFactRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_facts_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateFactRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateFactRequest) ProtoMessage() {}

func (x *CreateFactRequest) ProtoReflect() protoreflect.Message {
	mi := &file_facts_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateFactRequest.ProtoReflect.Descriptor instead.
func (*CreateFactRequest) Descriptor() ([]byte, []int) {
	return file_facts_proto_rawDescGZIP(), []int{0}
}

func (x *CreateFactRequest) GetAnimal() string {
	if x != nil {
		return x.Animal
	}
	return ""
}

func (x *CreateFactRequest) GetFact() string {
	if x != nil {
		return x.Fact
	}
	return ""
}

type CreateFactReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ID  int64  `protobuf:"varint,1,opt,name=ID,proto3" json:"ID,omitempty"`
	Err string `protobuf:"bytes,2,opt,name=Err,proto3" json:"Err,omitempty"`
}

func (x *CreateFactReply) Reset() {
	*x = CreateFactReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_facts_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateFactReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateFactReply) ProtoMessage() {}

func (x *CreateFactReply) ProtoReflect() protoreflect.Message {
	mi := &file_facts_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateFactReply.ProtoReflect.Descriptor instead.
func (*CreateFactReply) Descriptor() ([]byte, []int) {
	return file_facts_proto_rawDescGZIP(), []int{1}
}

func (x *CreateFactReply) GetID() int64 {
	if x != nil {
		return x.ID
	}
	return 0
}

func (x *CreateFactReply) GetErr() string {
	if x != nil {
		return x.Err
	}
	return ""
}

type GetFactRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ID int64 `protobuf:"varint,1,opt,name=ID,proto3" json:"ID,omitempty"`
}

func (x *GetFactRequest) Reset() {
	*x = GetFactRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_facts_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetFactRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetFactRequest) ProtoMessage() {}

func (x *GetFactRequest) ProtoReflect() protoreflect.Message {
	mi := &file_facts_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetFactRequest.ProtoReflect.Descriptor instead.
func (*GetFactRequest) Descriptor() ([]byte, []int) {
	return file_facts_proto_rawDescGZIP(), []int{2}
}

func (x *GetFactRequest) GetID() int64 {
	if x != nil {
		return x.ID
	}
	return 0
}

type GetFactReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Animal  string `protobuf:"bytes,1,opt,name=Animal,proto3" json:"Animal,omitempty"`
	Fact    string `protobuf:"bytes,2,opt,name=Fact,proto3" json:"Fact,omitempty"`
	ID      int64  `protobuf:"varint,3,opt,name=ID,proto3" json:"ID,omitempty"`
	Deleted bool   `protobuf:"varint,4,opt,name=Deleted,proto3" json:"Deleted,omitempty"`
	Err     string `protobuf:"bytes,5,opt,name=Err,proto3" json:"Err,omitempty"`
}

func (x *GetFactReply) Reset() {
	*x = GetFactReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_facts_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetFactReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetFactReply) ProtoMessage() {}

func (x *GetFactReply) ProtoReflect() protoreflect.Message {
	mi := &file_facts_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetFactReply.ProtoReflect.Descriptor instead.
func (*GetFactReply) Descriptor() ([]byte, []int) {
	return file_facts_proto_rawDescGZIP(), []int{3}
}

func (x *GetFactReply) GetAnimal() string {
	if x != nil {
		return x.Animal
	}
	return ""
}

func (x *GetFactReply) GetFact() string {
	if x != nil {
		return x.Fact
	}
	return ""
}

func (x *GetFactReply) GetID() int64 {
	if x != nil {
		return x.ID
	}
	return 0
}

func (x *GetFactReply) GetDeleted() bool {
	if x != nil {
		return x.Deleted
	}
	return false
}

func (x *GetFactReply) GetErr() string {
	if x != nil {
		return x.Err
	}
	return ""
}

type DeleteFactRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ID int64 `protobuf:"varint,1,opt,name=ID,proto3" json:"ID,omitempty"`
}

func (x *DeleteFactRequest) Reset() {
	*x = DeleteFactRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_facts_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeleteFactRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteFactRequest) ProtoMessage() {}

func (x *DeleteFactRequest) ProtoReflect() protoreflect.Message {
	mi := &file_facts_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteFactRequest.ProtoReflect.Descriptor instead.
func (*DeleteFactRequest) Descriptor() ([]byte, []int) {
	return file_facts_proto_rawDescGZIP(), []int{4}
}

func (x *DeleteFactRequest) GetID() int64 {
	if x != nil {
		return x.ID
	}
	return 0
}

type DeleteFactReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Err string `protobuf:"bytes,1,opt,name=Err,proto3" json:"Err,omitempty"`
}

func (x *DeleteFactReply) Reset() {
	*x = DeleteFactReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_facts_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeleteFactReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteFactReply) ProtoMessage() {}

func (x *DeleteFactReply) ProtoReflect() protoreflect.Message {
	mi := &file_facts_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteFactReply.ProtoReflect.Descriptor instead.
func (*DeleteFactReply) Descriptor() ([]byte, []int) {
	return file_facts_proto_rawDescGZIP(), []int{5}
}

func (x *DeleteFactReply) GetErr() string {
	if x != nil {
		return x.Err
	}
	return ""
}

type GetAnimalsReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Animals []string `protobuf:"bytes,1,rep,name=Animals,proto3" json:"Animals,omitempty"`
	Err     string   `protobuf:"bytes,2,opt,name=Err,proto3" json:"Err,omitempty"`
}

func (x *GetAnimalsReply) Reset() {
	*x = GetAnimalsReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_facts_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetAnimalsReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetAnimalsReply) ProtoMessage() {}

func (x *GetAnimalsReply) ProtoReflect() protoreflect.Message {
	mi := &file_facts_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetAnimalsReply.ProtoReflect.Descriptor instead.
func (*GetAnimalsReply) Descriptor() ([]byte, []int) {
	return file_facts_proto_rawDescGZIP(), []int{6}
}

func (x *GetAnimalsReply) GetAnimals() []string {
	if x != nil {
		return x.Animals
	}
	return nil
}

func (x *GetAnimalsReply) GetErr() string {
	if x != nil {
		return x.Err
	}
	return ""
}

type GetRandAnimalFactRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Animal string `protobuf:"bytes,1,opt,name=Animal,proto3" json:"Animal,omitempty"`
}

func (x *GetRandAnimalFactRequest) Reset() {
	*x = GetRandAnimalFactRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_facts_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetRandAnimalFactRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetRandAnimalFactRequest) ProtoMessage() {}

func (x *GetRandAnimalFactRequest) ProtoReflect() protoreflect.Message {
	mi := &file_facts_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetRandAnimalFactRequest.ProtoReflect.Descriptor instead.
func (*GetRandAnimalFactRequest) Descriptor() ([]byte, []int) {
	return file_facts_proto_rawDescGZIP(), []int{7}
}

func (x *GetRandAnimalFactRequest) GetAnimal() string {
	if x != nil {
		return x.Animal
	}
	return ""
}

type GetRandAnimalFactReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Fact string `protobuf:"bytes,1,opt,name=Fact,proto3" json:"Fact,omitempty"`
	ID   int64  `protobuf:"varint,2,opt,name=ID,proto3" json:"ID,omitempty"`
	Err  string `protobuf:"bytes,3,opt,name=Err,proto3" json:"Err,omitempty"`
}

func (x *GetRandAnimalFactReply) Reset() {
	*x = GetRandAnimalFactReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_facts_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetRandAnimalFactReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetRandAnimalFactReply) ProtoMessage() {}

func (x *GetRandAnimalFactReply) ProtoReflect() protoreflect.Message {
	mi := &file_facts_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetRandAnimalFactReply.ProtoReflect.Descriptor instead.
func (*GetRandAnimalFactReply) Descriptor() ([]byte, []int) {
	return file_facts_proto_rawDescGZIP(), []int{8}
}

func (x *GetRandAnimalFactReply) GetFact() string {
	if x != nil {
		return x.Fact
	}
	return ""
}

func (x *GetRandAnimalFactReply) GetID() int64 {
	if x != nil {
		return x.ID
	}
	return 0
}

func (x *GetRandAnimalFactReply) GetErr() string {
	if x != nil {
		return x.Err
	}
	return ""
}

type PublishFactRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Animal string `protobuf:"bytes,1,opt,name=Animal,proto3" json:"Animal,omitempty"`
}

func (x *PublishFactRequest) Reset() {
	*x = PublishFactRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_facts_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PublishFactRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PublishFactRequest) ProtoMessage() {}

func (x *PublishFactRequest) ProtoReflect() protoreflect.Message {
	mi := &file_facts_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PublishFactRequest.ProtoReflect.Descriptor instead.
func (*PublishFactRequest) Descriptor() ([]byte, []int) {
	return file_facts_proto_rawDescGZIP(), []int{9}
}

func (x *PublishFactRequest) GetAnimal() string {
	if x != nil {
		return x.Animal
	}
	return ""
}

type PublishFactReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Fact string `protobuf:"bytes,1,opt,name=Fact,proto3" json:"Fact,omitempty"`
	ID   int64  `protobuf:"varint,2,opt,name=ID,proto3" json:"ID,omitempty"`
	Err  string `protobuf:"bytes,3,opt,name=Err,proto3" json:"Err,omitempty"`
}

func (x *PublishFactReply) Reset() {
	*x = PublishFactReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_facts_proto_msgTypes[10]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PublishFactReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PublishFactReply) ProtoMessage() {}

func (x *PublishFactReply) ProtoReflect() protoreflect.Message {
	mi := &file_facts_proto_msgTypes[10]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PublishFactReply.ProtoReflect.Descriptor instead.
func (*PublishFactReply) Descriptor() ([]byte, []int) {
	return file_facts_proto_rawDescGZIP(), []int{10}
}

func (x *PublishFactReply) GetFact() string {
	if x != nil {
		return x.Fact
	}
	return ""
}

func (x *PublishFactReply) GetID() int64 {
	if x != nil {
		return x.ID
	}
	return 0
}

func (x *PublishFactReply) GetErr() string {
	if x != nil {
		return x.Err
	}
	return ""
}

var File_facts_proto protoreflect.FileDescriptor

var file_facts_proto_rawDesc = []byte{
	0x0a, 0x0b, 0x66, 0x61, 0x63, 0x74, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1b, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65,
	0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x3f, 0x0a, 0x11, 0x43, 0x72,
	0x65, 0x61, 0x74, 0x65, 0x46, 0x61, 0x63, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x16, 0x0a, 0x06, 0x41, 0x6e, 0x69, 0x6d, 0x61, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x06, 0x41, 0x6e, 0x69, 0x6d, 0x61, 0x6c, 0x12, 0x12, 0x0a, 0x04, 0x46, 0x61, 0x63, 0x74, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x46, 0x61, 0x63, 0x74, 0x22, 0x33, 0x0a, 0x0f, 0x43,
	0x72, 0x65, 0x61, 0x74, 0x65, 0x46, 0x61, 0x63, 0x74, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x0e,
	0x0a, 0x02, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x02, 0x49, 0x44, 0x12, 0x10,
	0x0a, 0x03, 0x45, 0x72, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x45, 0x72, 0x72,
	0x22, 0x20, 0x0a, 0x0e, 0x47, 0x65, 0x74, 0x46, 0x61, 0x63, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x02,
	0x49, 0x44, 0x22, 0x76, 0x0a, 0x0c, 0x47, 0x65, 0x74, 0x46, 0x61, 0x63, 0x74, 0x52, 0x65, 0x70,
	0x6c, 0x79, 0x12, 0x16, 0x0a, 0x06, 0x41, 0x6e, 0x69, 0x6d, 0x61, 0x6c, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x06, 0x41, 0x6e, 0x69, 0x6d, 0x61, 0x6c, 0x12, 0x12, 0x0a, 0x04, 0x46, 0x61,
	0x63, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x46, 0x61, 0x63, 0x74, 0x12, 0x0e,
	0x0a, 0x02, 0x49, 0x44, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x02, 0x49, 0x44, 0x12, 0x18,
	0x0a, 0x07, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x08, 0x52,
	0x07, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x64, 0x12, 0x10, 0x0a, 0x03, 0x45, 0x72, 0x72, 0x18,
	0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x45, 0x72, 0x72, 0x22, 0x23, 0x0a, 0x11, 0x44, 0x65,
	0x6c, 0x65, 0x74, 0x65, 0x46, 0x61, 0x63, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x0e, 0x0a, 0x02, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x02, 0x49, 0x44, 0x22,
	0x23, 0x0a, 0x0f, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x46, 0x61, 0x63, 0x74, 0x52, 0x65, 0x70,
	0x6c, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x45, 0x72, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x03, 0x45, 0x72, 0x72, 0x22, 0x3d, 0x0a, 0x0f, 0x47, 0x65, 0x74, 0x41, 0x6e, 0x69, 0x6d, 0x61,
	0x6c, 0x73, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x18, 0x0a, 0x07, 0x41, 0x6e, 0x69, 0x6d, 0x61,
	0x6c, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x09, 0x52, 0x07, 0x41, 0x6e, 0x69, 0x6d, 0x61, 0x6c,
	0x73, 0x12, 0x10, 0x0a, 0x03, 0x45, 0x72, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03,
	0x45, 0x72, 0x72, 0x22, 0x32, 0x0a, 0x18, 0x47, 0x65, 0x74, 0x52, 0x61, 0x6e, 0x64, 0x41, 0x6e,
	0x69, 0x6d, 0x61, 0x6c, 0x46, 0x61, 0x63, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x16, 0x0a, 0x06, 0x41, 0x6e, 0x69, 0x6d, 0x61, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x06, 0x41, 0x6e, 0x69, 0x6d, 0x61, 0x6c, 0x22, 0x4e, 0x0a, 0x16, 0x47, 0x65, 0x74, 0x52, 0x61,
	0x6e, 0x64, 0x41, 0x6e, 0x69, 0x6d, 0x61, 0x6c, 0x46, 0x61, 0x63, 0x74, 0x52, 0x65, 0x70, 0x6c,
	0x79, 0x12, 0x12, 0x0a, 0x04, 0x46, 0x61, 0x63, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x04, 0x46, 0x61, 0x63, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x49, 0x44, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x03, 0x52, 0x02, 0x49, 0x44, 0x12, 0x10, 0x0a, 0x03, 0x45, 0x72, 0x72, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x03, 0x45, 0x72, 0x72, 0x22, 0x2c, 0x0a, 0x12, 0x50, 0x75, 0x62, 0x6c, 0x69,
	0x73, 0x68, 0x46, 0x61, 0x63, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x16, 0x0a,
	0x06, 0x41, 0x6e, 0x69, 0x6d, 0x61, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x41,
	0x6e, 0x69, 0x6d, 0x61, 0x6c, 0x22, 0x48, 0x0a, 0x10, 0x50, 0x75, 0x62, 0x6c, 0x69, 0x73, 0x68,
	0x46, 0x61, 0x63, 0x74, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x12, 0x0a, 0x04, 0x46, 0x61, 0x63,
	0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x46, 0x61, 0x63, 0x74, 0x12, 0x0e, 0x0a,
	0x02, 0x49, 0x44, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x02, 0x49, 0x44, 0x12, 0x10, 0x0a,
	0x03, 0x45, 0x72, 0x72, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x45, 0x72, 0x72, 0x32,
	0xd2, 0x02, 0x0a, 0x05, 0x46, 0x61, 0x63, 0x74, 0x73, 0x12, 0x32, 0x0a, 0x0a, 0x43, 0x72, 0x65,
	0x61, 0x74, 0x65, 0x46, 0x61, 0x63, 0x74, 0x12, 0x12, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65,
	0x46, 0x61, 0x63, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x10, 0x2e, 0x43, 0x72,
	0x65, 0x61, 0x74, 0x65, 0x46, 0x61, 0x63, 0x74, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x29, 0x0a,
	0x07, 0x47, 0x65, 0x74, 0x46, 0x61, 0x63, 0x74, 0x12, 0x0f, 0x2e, 0x47, 0x65, 0x74, 0x46, 0x61,
	0x63, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x0d, 0x2e, 0x47, 0x65, 0x74, 0x46,
	0x61, 0x63, 0x74, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x32, 0x0a, 0x0a, 0x44, 0x65, 0x6c, 0x65,
	0x74, 0x65, 0x46, 0x61, 0x63, 0x74, 0x12, 0x12, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x46,
	0x61, 0x63, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x10, 0x2e, 0x44, 0x65, 0x6c,
	0x65, 0x74, 0x65, 0x46, 0x61, 0x63, 0x74, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x36, 0x0a, 0x0a,
	0x47, 0x65, 0x74, 0x41, 0x6e, 0x69, 0x6d, 0x61, 0x6c, 0x73, 0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70,
	0x74, 0x79, 0x1a, 0x10, 0x2e, 0x47, 0x65, 0x74, 0x41, 0x6e, 0x69, 0x6d, 0x61, 0x6c, 0x73, 0x52,
	0x65, 0x70, 0x6c, 0x79, 0x12, 0x47, 0x0a, 0x11, 0x47, 0x65, 0x74, 0x52, 0x61, 0x6e, 0x64, 0x41,
	0x6e, 0x69, 0x6d, 0x61, 0x6c, 0x46, 0x61, 0x63, 0x74, 0x12, 0x19, 0x2e, 0x47, 0x65, 0x74, 0x52,
	0x61, 0x6e, 0x64, 0x41, 0x6e, 0x69, 0x6d, 0x61, 0x6c, 0x46, 0x61, 0x63, 0x74, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x17, 0x2e, 0x47, 0x65, 0x74, 0x52, 0x61, 0x6e, 0x64, 0x41, 0x6e,
	0x69, 0x6d, 0x61, 0x6c, 0x46, 0x61, 0x63, 0x74, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x35, 0x0a,
	0x0b, 0x50, 0x75, 0x62, 0x6c, 0x69, 0x73, 0x68, 0x46, 0x61, 0x63, 0x74, 0x12, 0x13, 0x2e, 0x50,
	0x75, 0x62, 0x6c, 0x69, 0x73, 0x68, 0x46, 0x61, 0x63, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x11, 0x2e, 0x50, 0x75, 0x62, 0x6c, 0x69, 0x73, 0x68, 0x46, 0x61, 0x63, 0x74, 0x52,
	0x65, 0x70, 0x6c, 0x79, 0x42, 0x31, 0x5a, 0x2f, 0x6d, 0x6f, 0x64, 0x75, 0x6c, 0x65, 0x20, 0x67,
	0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x73, 0x6e, 0x6f, 0x6f, 0x79, 0x65,
	0x6e, 0x2f, 0x61, 0x6e, 0x69, 0x6d, 0x61, 0x6c, 0x2d, 0x66, 0x61, 0x63, 0x74, 0x73, 0x2f, 0x66,
	0x61, 0x63, 0x74, 0x73, 0x2f, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_facts_proto_rawDescOnce sync.Once
	file_facts_proto_rawDescData = file_facts_proto_rawDesc
)

func file_facts_proto_rawDescGZIP() []byte {
	file_facts_proto_rawDescOnce.Do(func() {
		file_facts_proto_rawDescData = protoimpl.X.CompressGZIP(file_facts_proto_rawDescData)
	})
	return file_facts_proto_rawDescData
}

var file_facts_proto_msgTypes = make([]protoimpl.MessageInfo, 11)
var file_facts_proto_goTypes = []interface{}{
	(*CreateFactRequest)(nil),        // 0: CreateFactRequest
	(*CreateFactReply)(nil),          // 1: CreateFactReply
	(*GetFactRequest)(nil),           // 2: GetFactRequest
	(*GetFactReply)(nil),             // 3: GetFactReply
	(*DeleteFactRequest)(nil),        // 4: DeleteFactRequest
	(*DeleteFactReply)(nil),          // 5: DeleteFactReply
	(*GetAnimalsReply)(nil),          // 6: GetAnimalsReply
	(*GetRandAnimalFactRequest)(nil), // 7: GetRandAnimalFactRequest
	(*GetRandAnimalFactReply)(nil),   // 8: GetRandAnimalFactReply
	(*PublishFactRequest)(nil),       // 9: PublishFactRequest
	(*PublishFactReply)(nil),         // 10: PublishFactReply
	(*emptypb.Empty)(nil),            // 11: google.protobuf.Empty
}
var file_facts_proto_depIdxs = []int32{
	0,  // 0: Facts.CreateFact:input_type -> CreateFactRequest
	2,  // 1: Facts.GetFact:input_type -> GetFactRequest
	4,  // 2: Facts.DeleteFact:input_type -> DeleteFactRequest
	11, // 3: Facts.GetAnimals:input_type -> google.protobuf.Empty
	7,  // 4: Facts.GetRandAnimalFact:input_type -> GetRandAnimalFactRequest
	9,  // 5: Facts.PublishFact:input_type -> PublishFactRequest
	1,  // 6: Facts.CreateFact:output_type -> CreateFactReply
	3,  // 7: Facts.GetFact:output_type -> GetFactReply
	5,  // 8: Facts.DeleteFact:output_type -> DeleteFactReply
	6,  // 9: Facts.GetAnimals:output_type -> GetAnimalsReply
	8,  // 10: Facts.GetRandAnimalFact:output_type -> GetRandAnimalFactReply
	10, // 11: Facts.PublishFact:output_type -> PublishFactReply
	6,  // [6:12] is the sub-list for method output_type
	0,  // [0:6] is the sub-list for method input_type
	0,  // [0:0] is the sub-list for extension type_name
	0,  // [0:0] is the sub-list for extension extendee
	0,  // [0:0] is the sub-list for field type_name
}

func init() { file_facts_proto_init() }
func file_facts_proto_init() {
	if File_facts_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_facts_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateFactRequest); i {
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
		file_facts_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateFactReply); i {
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
		file_facts_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetFactRequest); i {
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
		file_facts_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetFactReply); i {
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
		file_facts_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeleteFactRequest); i {
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
		file_facts_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeleteFactReply); i {
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
		file_facts_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetAnimalsReply); i {
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
		file_facts_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetRandAnimalFactRequest); i {
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
		file_facts_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetRandAnimalFactReply); i {
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
		file_facts_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PublishFactRequest); i {
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
		file_facts_proto_msgTypes[10].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PublishFactReply); i {
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
			RawDescriptor: file_facts_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   11,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_facts_proto_goTypes,
		DependencyIndexes: file_facts_proto_depIdxs,
		MessageInfos:      file_facts_proto_msgTypes,
	}.Build()
	File_facts_proto = out.File
	file_facts_proto_rawDesc = nil
	file_facts_proto_goTypes = nil
	file_facts_proto_depIdxs = nil
}

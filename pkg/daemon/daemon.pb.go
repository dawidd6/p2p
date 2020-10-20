// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.23.0
// 	protoc        v3.13.0
// source: pkg/daemon/daemon.proto

package daemon

import (
	torrent "github.com/dawidd6/p2p/pkg/torrent"
	proto "github.com/golang/protobuf/proto"
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

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

type SeedRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	FileHash    string `protobuf:"bytes,1,opt,name=file_hash,json=fileHash,proto3" json:"file_hash,omitempty"`
	PieceNumber int64  `protobuf:"varint,2,opt,name=piece_number,json=pieceNumber,proto3" json:"piece_number,omitempty"`
}

func (x *SeedRequest) Reset() {
	*x = SeedRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_daemon_daemon_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SeedRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SeedRequest) ProtoMessage() {}

func (x *SeedRequest) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_daemon_daemon_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SeedRequest.ProtoReflect.Descriptor instead.
func (*SeedRequest) Descriptor() ([]byte, []int) {
	return file_pkg_daemon_daemon_proto_rawDescGZIP(), []int{0}
}

func (x *SeedRequest) GetFileHash() string {
	if x != nil {
		return x.FileHash
	}
	return ""
}

func (x *SeedRequest) GetPieceNumber() int64 {
	if x != nil {
		return x.PieceNumber
	}
	return 0
}

type SeedResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	PieceData []byte `protobuf:"bytes,2,opt,name=piece_data,json=pieceData,proto3" json:"piece_data,omitempty"`
}

func (x *SeedResponse) Reset() {
	*x = SeedResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_daemon_daemon_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SeedResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SeedResponse) ProtoMessage() {}

func (x *SeedResponse) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_daemon_daemon_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SeedResponse.ProtoReflect.Descriptor instead.
func (*SeedResponse) Descriptor() ([]byte, []int) {
	return file_pkg_daemon_daemon_proto_rawDescGZIP(), []int{1}
}

func (x *SeedResponse) GetPieceData() []byte {
	if x != nil {
		return x.PieceData
	}
	return nil
}

type AddRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Torrent *torrent.Torrent `protobuf:"bytes,1,opt,name=torrent,proto3" json:"torrent,omitempty"`
}

func (x *AddRequest) Reset() {
	*x = AddRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_daemon_daemon_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AddRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AddRequest) ProtoMessage() {}

func (x *AddRequest) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_daemon_daemon_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AddRequest.ProtoReflect.Descriptor instead.
func (*AddRequest) Descriptor() ([]byte, []int) {
	return file_pkg_daemon_daemon_proto_rawDescGZIP(), []int{2}
}

func (x *AddRequest) GetTorrent() *torrent.Torrent {
	if x != nil {
		return x.Torrent
	}
	return nil
}

type AddResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *AddResponse) Reset() {
	*x = AddResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_daemon_daemon_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AddResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AddResponse) ProtoMessage() {}

func (x *AddResponse) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_daemon_daemon_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AddResponse.ProtoReflect.Descriptor instead.
func (*AddResponse) Descriptor() ([]byte, []int) {
	return file_pkg_daemon_daemon_proto_rawDescGZIP(), []int{3}
}

type DeleteRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	FileHash string `protobuf:"bytes,1,opt,name=file_hash,json=fileHash,proto3" json:"file_hash,omitempty"`
	WithData bool   `protobuf:"varint,2,opt,name=with_data,json=withData,proto3" json:"with_data,omitempty"`
}

func (x *DeleteRequest) Reset() {
	*x = DeleteRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_daemon_daemon_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeleteRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteRequest) ProtoMessage() {}

func (x *DeleteRequest) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_daemon_daemon_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteRequest.ProtoReflect.Descriptor instead.
func (*DeleteRequest) Descriptor() ([]byte, []int) {
	return file_pkg_daemon_daemon_proto_rawDescGZIP(), []int{4}
}

func (x *DeleteRequest) GetFileHash() string {
	if x != nil {
		return x.FileHash
	}
	return ""
}

func (x *DeleteRequest) GetWithData() bool {
	if x != nil {
		return x.WithData
	}
	return false
}

type DeleteResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *DeleteResponse) Reset() {
	*x = DeleteResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_daemon_daemon_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeleteResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteResponse) ProtoMessage() {}

func (x *DeleteResponse) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_daemon_daemon_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteResponse.ProtoReflect.Descriptor instead.
func (*DeleteResponse) Descriptor() ([]byte, []int) {
	return file_pkg_daemon_daemon_proto_rawDescGZIP(), []int{5}
}

type StatusRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	FileHash string `protobuf:"bytes,1,opt,name=file_hash,json=fileHash,proto3" json:"file_hash,omitempty"`
}

func (x *StatusRequest) Reset() {
	*x = StatusRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_daemon_daemon_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StatusRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StatusRequest) ProtoMessage() {}

func (x *StatusRequest) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_daemon_daemon_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StatusRequest.ProtoReflect.Descriptor instead.
func (*StatusRequest) Descriptor() ([]byte, []int) {
	return file_pkg_daemon_daemon_proto_rawDescGZIP(), []int{6}
}

func (x *StatusRequest) GetFileHash() string {
	if x != nil {
		return x.FileHash
	}
	return ""
}

type StatusResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Torrents []*torrent.Torrent `protobuf:"bytes,1,rep,name=torrents,proto3" json:"torrents,omitempty"`
}

func (x *StatusResponse) Reset() {
	*x = StatusResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_daemon_daemon_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StatusResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StatusResponse) ProtoMessage() {}

func (x *StatusResponse) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_daemon_daemon_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StatusResponse.ProtoReflect.Descriptor instead.
func (*StatusResponse) Descriptor() ([]byte, []int) {
	return file_pkg_daemon_daemon_proto_rawDescGZIP(), []int{7}
}

func (x *StatusResponse) GetTorrents() []*torrent.Torrent {
	if x != nil {
		return x.Torrents
	}
	return nil
}

var File_pkg_daemon_daemon_proto protoreflect.FileDescriptor

var file_pkg_daemon_daemon_proto_rawDesc = []byte{
	0x0a, 0x17, 0x70, 0x6b, 0x67, 0x2f, 0x64, 0x61, 0x65, 0x6d, 0x6f, 0x6e, 0x2f, 0x64, 0x61, 0x65,
	0x6d, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x19, 0x70, 0x6b, 0x67, 0x2f, 0x74,
	0x6f, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x2f, 0x74, 0x6f, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x22, 0x4d, 0x0a, 0x0b, 0x53, 0x65, 0x65, 0x64, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x1b, 0x0a, 0x09, 0x66, 0x69, 0x6c, 0x65, 0x5f, 0x68, 0x61, 0x73, 0x68,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x66, 0x69, 0x6c, 0x65, 0x48, 0x61, 0x73, 0x68,
	0x12, 0x21, 0x0a, 0x0c, 0x70, 0x69, 0x65, 0x63, 0x65, 0x5f, 0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0b, 0x70, 0x69, 0x65, 0x63, 0x65, 0x4e, 0x75, 0x6d,
	0x62, 0x65, 0x72, 0x22, 0x2d, 0x0a, 0x0c, 0x53, 0x65, 0x65, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x1d, 0x0a, 0x0a, 0x70, 0x69, 0x65, 0x63, 0x65, 0x5f, 0x64, 0x61, 0x74,
	0x61, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x09, 0x70, 0x69, 0x65, 0x63, 0x65, 0x44, 0x61,
	0x74, 0x61, 0x22, 0x30, 0x0a, 0x0a, 0x41, 0x64, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x22, 0x0a, 0x07, 0x74, 0x6f, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x08, 0x2e, 0x54, 0x6f, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x52, 0x07, 0x74, 0x6f, 0x72,
	0x72, 0x65, 0x6e, 0x74, 0x22, 0x0d, 0x0a, 0x0b, 0x41, 0x64, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x22, 0x49, 0x0a, 0x0d, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x1b, 0x0a, 0x09, 0x66, 0x69, 0x6c, 0x65, 0x5f, 0x68, 0x61, 0x73,
	0x68, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x66, 0x69, 0x6c, 0x65, 0x48, 0x61, 0x73,
	0x68, 0x12, 0x1b, 0x0a, 0x09, 0x77, 0x69, 0x74, 0x68, 0x5f, 0x64, 0x61, 0x74, 0x61, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x08, 0x52, 0x08, 0x77, 0x69, 0x74, 0x68, 0x44, 0x61, 0x74, 0x61, 0x22, 0x10,
	0x0a, 0x0e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x22, 0x2c, 0x0a, 0x0d, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x1b, 0x0a, 0x09, 0x66, 0x69, 0x6c, 0x65, 0x5f, 0x68, 0x61, 0x73, 0x68, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x66, 0x69, 0x6c, 0x65, 0x48, 0x61, 0x73, 0x68, 0x22, 0x36,
	0x0a, 0x0e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x24, 0x0a, 0x08, 0x74, 0x6f, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x73, 0x18, 0x01, 0x20, 0x03,
	0x28, 0x0b, 0x32, 0x08, 0x2e, 0x54, 0x6f, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x52, 0x08, 0x74, 0x6f,
	0x72, 0x72, 0x65, 0x6e, 0x74, 0x73, 0x32, 0x80, 0x01, 0x0a, 0x06, 0x44, 0x61, 0x65, 0x6d, 0x6f,
	0x6e, 0x12, 0x20, 0x0a, 0x03, 0x41, 0x64, 0x64, 0x12, 0x0b, 0x2e, 0x41, 0x64, 0x64, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x0c, 0x2e, 0x41, 0x64, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x29, 0x0a, 0x06, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x12, 0x0e, 0x2e,
	0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x0f, 0x2e,
	0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x29,
	0x0a, 0x06, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x0e, 0x2e, 0x53, 0x74, 0x61, 0x74, 0x75,
	0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x0f, 0x2e, 0x53, 0x74, 0x61, 0x74, 0x75,
	0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x32, 0x2d, 0x0a, 0x06, 0x53, 0x65, 0x65,
	0x64, 0x65, 0x72, 0x12, 0x23, 0x0a, 0x04, 0x53, 0x65, 0x65, 0x64, 0x12, 0x0c, 0x2e, 0x53, 0x65,
	0x65, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x0d, 0x2e, 0x53, 0x65, 0x65, 0x64,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x23, 0x5a, 0x21, 0x67, 0x69, 0x74, 0x68,
	0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x64, 0x61, 0x77, 0x69, 0x64, 0x64, 0x36, 0x2f, 0x70,
	0x32, 0x70, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x64, 0x61, 0x65, 0x6d, 0x6f, 0x6e, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_pkg_daemon_daemon_proto_rawDescOnce sync.Once
	file_pkg_daemon_daemon_proto_rawDescData = file_pkg_daemon_daemon_proto_rawDesc
)

func file_pkg_daemon_daemon_proto_rawDescGZIP() []byte {
	file_pkg_daemon_daemon_proto_rawDescOnce.Do(func() {
		file_pkg_daemon_daemon_proto_rawDescData = protoimpl.X.CompressGZIP(file_pkg_daemon_daemon_proto_rawDescData)
	})
	return file_pkg_daemon_daemon_proto_rawDescData
}

var file_pkg_daemon_daemon_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_pkg_daemon_daemon_proto_goTypes = []interface{}{
	(*SeedRequest)(nil),     // 0: SeedRequest
	(*SeedResponse)(nil),    // 1: SeedResponse
	(*AddRequest)(nil),      // 2: AddRequest
	(*AddResponse)(nil),     // 3: AddResponse
	(*DeleteRequest)(nil),   // 4: DeleteRequest
	(*DeleteResponse)(nil),  // 5: DeleteResponse
	(*StatusRequest)(nil),   // 6: StatusRequest
	(*StatusResponse)(nil),  // 7: StatusResponse
	(*torrent.Torrent)(nil), // 8: Torrent
}
var file_pkg_daemon_daemon_proto_depIdxs = []int32{
	8, // 0: AddRequest.torrent:type_name -> Torrent
	8, // 1: StatusResponse.torrents:type_name -> Torrent
	2, // 2: Daemon.Add:input_type -> AddRequest
	4, // 3: Daemon.Delete:input_type -> DeleteRequest
	6, // 4: Daemon.Status:input_type -> StatusRequest
	0, // 5: Seeder.Seed:input_type -> SeedRequest
	3, // 6: Daemon.Add:output_type -> AddResponse
	5, // 7: Daemon.Delete:output_type -> DeleteResponse
	7, // 8: Daemon.Status:output_type -> StatusResponse
	1, // 9: Seeder.Seed:output_type -> SeedResponse
	6, // [6:10] is the sub-list for method output_type
	2, // [2:6] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_pkg_daemon_daemon_proto_init() }
func file_pkg_daemon_daemon_proto_init() {
	if File_pkg_daemon_daemon_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_pkg_daemon_daemon_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SeedRequest); i {
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
		file_pkg_daemon_daemon_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SeedResponse); i {
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
		file_pkg_daemon_daemon_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AddRequest); i {
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
		file_pkg_daemon_daemon_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AddResponse); i {
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
		file_pkg_daemon_daemon_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeleteRequest); i {
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
		file_pkg_daemon_daemon_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeleteResponse); i {
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
		file_pkg_daemon_daemon_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StatusRequest); i {
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
		file_pkg_daemon_daemon_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StatusResponse); i {
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
			RawDescriptor: file_pkg_daemon_daemon_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   2,
		},
		GoTypes:           file_pkg_daemon_daemon_proto_goTypes,
		DependencyIndexes: file_pkg_daemon_daemon_proto_depIdxs,
		MessageInfos:      file_pkg_daemon_daemon_proto_msgTypes,
	}.Build()
	File_pkg_daemon_daemon_proto = out.File
	file_pkg_daemon_daemon_proto_rawDesc = nil
	file_pkg_daemon_daemon_proto_goTypes = nil
	file_pkg_daemon_daemon_proto_depIdxs = nil
}

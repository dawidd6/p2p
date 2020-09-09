// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0-devel
// 	protoc        v3.13.0
// source: torrent/torrent.proto

package torrent

import (
	file "github.com/dawidd6/p2p/file"
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

type Torrent struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name      string       `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Timestamp int64        `protobuf:"varint,2,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	Files     []*file.File `protobuf:"bytes,3,rep,name=files,proto3" json:"files,omitempty"`
}

func (x *Torrent) Reset() {
	*x = Torrent{}
	if protoimpl.UnsafeEnabled {
		mi := &file_torrent_torrent_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Torrent) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Torrent) ProtoMessage() {}

func (x *Torrent) ProtoReflect() protoreflect.Message {
	mi := &file_torrent_torrent_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Torrent.ProtoReflect.Descriptor instead.
func (*Torrent) Descriptor() ([]byte, []int) {
	return file_torrent_torrent_proto_rawDescGZIP(), []int{0}
}

func (x *Torrent) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Torrent) GetTimestamp() int64 {
	if x != nil {
		return x.Timestamp
	}
	return 0
}

func (x *Torrent) GetFiles() []*file.File {
	if x != nil {
		return x.Files
	}
	return nil
}

var File_torrent_torrent_proto protoreflect.FileDescriptor

var file_torrent_torrent_proto_rawDesc = []byte{
	0x0a, 0x15, 0x74, 0x6f, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x2f, 0x74, 0x6f, 0x72, 0x72, 0x65, 0x6e,
	0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x0f, 0x66, 0x69, 0x6c, 0x65, 0x2f, 0x66, 0x69,
	0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x58, 0x0a, 0x07, 0x54, 0x6f, 0x72, 0x72,
	0x65, 0x6e, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x1c, 0x0a, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73,
	0x74, 0x61, 0x6d, 0x70, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x74, 0x69, 0x6d, 0x65,
	0x73, 0x74, 0x61, 0x6d, 0x70, 0x12, 0x1b, 0x0a, 0x05, 0x66, 0x69, 0x6c, 0x65, 0x73, 0x18, 0x03,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x05, 0x2e, 0x46, 0x69, 0x6c, 0x65, 0x52, 0x05, 0x66, 0x69, 0x6c,
	0x65, 0x73, 0x42, 0x20, 0x5a, 0x1e, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d,
	0x2f, 0x64, 0x61, 0x77, 0x69, 0x64, 0x64, 0x36, 0x2f, 0x70, 0x32, 0x70, 0x2f, 0x74, 0x6f, 0x72,
	0x72, 0x65, 0x6e, 0x74, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_torrent_torrent_proto_rawDescOnce sync.Once
	file_torrent_torrent_proto_rawDescData = file_torrent_torrent_proto_rawDesc
)

func file_torrent_torrent_proto_rawDescGZIP() []byte {
	file_torrent_torrent_proto_rawDescOnce.Do(func() {
		file_torrent_torrent_proto_rawDescData = protoimpl.X.CompressGZIP(file_torrent_torrent_proto_rawDescData)
	})
	return file_torrent_torrent_proto_rawDescData
}

var file_torrent_torrent_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_torrent_torrent_proto_goTypes = []interface{}{
	(*Torrent)(nil),   // 0: Torrent
	(*file.File)(nil), // 1: File
}
var file_torrent_torrent_proto_depIdxs = []int32{
	1, // 0: Torrent.files:type_name -> File
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_torrent_torrent_proto_init() }
func file_torrent_torrent_proto_init() {
	if File_torrent_torrent_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_torrent_torrent_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Torrent); i {
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
			RawDescriptor: file_torrent_torrent_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_torrent_torrent_proto_goTypes,
		DependencyIndexes: file_torrent_torrent_proto_depIdxs,
		MessageInfos:      file_torrent_torrent_proto_msgTypes,
	}.Build()
	File_torrent_torrent_proto = out.File
	file_torrent_torrent_proto_rawDesc = nil
	file_torrent_torrent_proto_goTypes = nil
	file_torrent_torrent_proto_depIdxs = nil
}

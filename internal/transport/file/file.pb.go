// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.5
// 	protoc        v3.6.1
// source: api/file.proto

package file

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type DocumentID struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *DocumentID) Reset() {
	*x = DocumentID{}
	mi := &file_api_file_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DocumentID) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DocumentID) ProtoMessage() {}

func (x *DocumentID) ProtoReflect() protoreflect.Message {
	mi := &file_api_file_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DocumentID.ProtoReflect.Descriptor instead.
func (*DocumentID) Descriptor() ([]byte, []int) {
	return file_api_file_proto_rawDescGZIP(), []int{0}
}

func (x *DocumentID) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

type FileStatus struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Status        string                 `protobuf:"bytes,1,opt,name=status,proto3" json:"status,omitempty"`
	Message       string                 `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *FileStatus) Reset() {
	*x = FileStatus{}
	mi := &file_api_file_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *FileStatus) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FileStatus) ProtoMessage() {}

func (x *FileStatus) ProtoReflect() protoreflect.Message {
	mi := &file_api_file_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FileStatus.ProtoReflect.Descriptor instead.
func (*FileStatus) Descriptor() ([]byte, []int) {
	return file_api_file_proto_rawDescGZIP(), []int{1}
}

func (x *FileStatus) GetStatus() string {
	if x != nil {
		return x.Status
	}
	return ""
}

func (x *FileStatus) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

type File struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Name          string                 `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Chunks        []*File_Chunk          `protobuf:"bytes,2,rep,name=chunks,proto3" json:"chunks,omitempty"`
	Ext           string                 `protobuf:"bytes,3,opt,name=ext,proto3" json:"ext,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *File) Reset() {
	*x = File{}
	mi := &file_api_file_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *File) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*File) ProtoMessage() {}

func (x *File) ProtoReflect() protoreflect.Message {
	mi := &file_api_file_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use File.ProtoReflect.Descriptor instead.
func (*File) Descriptor() ([]byte, []int) {
	return file_api_file_proto_rawDescGZIP(), []int{2}
}

func (x *File) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *File) GetChunks() []*File_Chunk {
	if x != nil {
		return x.Chunks
	}
	return nil
}

func (x *File) GetExt() string {
	if x != nil {
		return x.Ext
	}
	return ""
}

type File_Chunk struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            int32                  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Data          []byte                 `protobuf:"bytes,2,opt,name=data,proto3" json:"data,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *File_Chunk) Reset() {
	*x = File_Chunk{}
	mi := &file_api_file_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *File_Chunk) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*File_Chunk) ProtoMessage() {}

func (x *File_Chunk) ProtoReflect() protoreflect.Message {
	mi := &file_api_file_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use File_Chunk.ProtoReflect.Descriptor instead.
func (*File_Chunk) Descriptor() ([]byte, []int) {
	return file_api_file_proto_rawDescGZIP(), []int{2, 0}
}

func (x *File_Chunk) GetId() int32 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *File_Chunk) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

var File_api_file_proto protoreflect.FileDescriptor

var file_api_file_proto_rawDesc = string([]byte{
	0x0a, 0x0e, 0x61, 0x70, 0x69, 0x2f, 0x66, 0x69, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x08, 0x61, 0x70, 0x69, 0x2e, 0x66, 0x69, 0x6c, 0x65, 0x22, 0x1c, 0x0a, 0x0a, 0x44, 0x6f,
	0x63, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x49, 0x44, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x22, 0x3e, 0x0a, 0x0a, 0x46, 0x69, 0x6c, 0x65,
	0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x18,
	0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x22, 0x87, 0x01, 0x0a, 0x04, 0x46, 0x69, 0x6c,
	0x65, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x2c, 0x0a, 0x06, 0x63, 0x68, 0x75, 0x6e, 0x6b, 0x73, 0x18,
	0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x66, 0x69, 0x6c, 0x65,
	0x2e, 0x46, 0x69, 0x6c, 0x65, 0x2e, 0x43, 0x68, 0x75, 0x6e, 0x6b, 0x52, 0x06, 0x63, 0x68, 0x75,
	0x6e, 0x6b, 0x73, 0x12, 0x10, 0x0a, 0x03, 0x65, 0x78, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x03, 0x65, 0x78, 0x74, 0x1a, 0x2b, 0x0a, 0x05, 0x43, 0x68, 0x75, 0x6e, 0x6b, 0x12, 0x0e,
	0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x02, 0x69, 0x64, 0x12, 0x12,
	0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x64, 0x61,
	0x74, 0x61, 0x32, 0x79, 0x0a, 0x0b, 0x46, 0x69, 0x6c, 0x65, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x12, 0x34, 0x0a, 0x0a, 0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x46, 0x69, 0x6c, 0x65, 0x12,
	0x0e, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x66, 0x69, 0x6c, 0x65, 0x2e, 0x46, 0x69, 0x6c, 0x65, 0x1a,
	0x14, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x66, 0x69, 0x6c, 0x65, 0x2e, 0x46, 0x69, 0x6c, 0x65, 0x53,
	0x74, 0x61, 0x74, 0x75, 0x73, 0x28, 0x01, 0x12, 0x34, 0x0a, 0x0c, 0x44, 0x6f, 0x77, 0x6e, 0x6c,
	0x6f, 0x61, 0x64, 0x46, 0x69, 0x6c, 0x65, 0x12, 0x14, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x66, 0x69,
	0x6c, 0x65, 0x2e, 0x44, 0x6f, 0x63, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x49, 0x44, 0x1a, 0x0e, 0x2e,
	0x61, 0x70, 0x69, 0x2e, 0x66, 0x69, 0x6c, 0x65, 0x2e, 0x46, 0x69, 0x6c, 0x65, 0x42, 0x36, 0x5a,
	0x34, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x43, 0x53, 0x38, 0x30,
	0x2d, 0x54, 0x65, 0x61, 0x6d, 0x2f, 0x47, 0x6f, 0x6f, 0x6c, 0x65, 0x61, 0x6e, 0x2f, 0x69, 0x6e,
	0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x70, 0x6f, 0x72, 0x74,
	0x2f, 0x66, 0x69, 0x6c, 0x65, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
})

var (
	file_api_file_proto_rawDescOnce sync.Once
	file_api_file_proto_rawDescData []byte
)

func file_api_file_proto_rawDescGZIP() []byte {
	file_api_file_proto_rawDescOnce.Do(func() {
		file_api_file_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_api_file_proto_rawDesc), len(file_api_file_proto_rawDesc)))
	})
	return file_api_file_proto_rawDescData
}

var file_api_file_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_api_file_proto_goTypes = []any{
	(*DocumentID)(nil), // 0: api.file.DocumentID
	(*FileStatus)(nil), // 1: api.file.FileStatus
	(*File)(nil),       // 2: api.file.File
	(*File_Chunk)(nil), // 3: api.file.File.Chunk
}
var file_api_file_proto_depIdxs = []int32{
	3, // 0: api.file.File.chunks:type_name -> api.file.File.Chunk
	2, // 1: api.file.FileService.UploadFile:input_type -> api.file.File
	0, // 2: api.file.FileService.DownloadFile:input_type -> api.file.DocumentID
	1, // 3: api.file.FileService.UploadFile:output_type -> api.file.FileStatus
	2, // 4: api.file.FileService.DownloadFile:output_type -> api.file.File
	3, // [3:5] is the sub-list for method output_type
	1, // [1:3] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_api_file_proto_init() }
func file_api_file_proto_init() {
	if File_api_file_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_api_file_proto_rawDesc), len(file_api_file_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_api_file_proto_goTypes,
		DependencyIndexes: file_api_file_proto_depIdxs,
		MessageInfos:      file_api_file_proto_msgTypes,
	}.Build()
	File_api_file_proto = out.File
	file_api_file_proto_goTypes = nil
	file_api_file_proto_depIdxs = nil
}

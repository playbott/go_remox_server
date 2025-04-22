package remotemouse

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

type InputState struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Dx            int32                  `protobuf:"varint,1,opt,name=dx,proto3" json:"dx,omitempty"`
	Dy            int32                  `protobuf:"varint,2,opt,name=dy,proto3" json:"dy,omitempty"`
	Buttons       map[string]bool        `protobuf:"bytes,3,rep,name=buttons,proto3" json:"buttons,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"varint,2,opt,name=value"`
	ScrollY       int32                  `protobuf:"varint,4,opt,name=scroll_y,json=scrollY,proto3" json:"scroll_y,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *InputState) Reset() {
	*x = InputState{}
	mi := &file_input_state_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *InputState) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InputState) ProtoMessage() {}

func (x *InputState) ProtoReflect() protoreflect.Message {
	mi := &file_input_state_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use InputState.ProtoReflect.Descriptor instead.
func (*InputState) Descriptor() ([]byte, []int) {
	return file_input_state_proto_rawDescGZIP(), []int{0}
}

func (x *InputState) GetDx() int32 {
	if x != nil {
		return x.Dx
	}
	return 0
}

func (x *InputState) GetDy() int32 {
	if x != nil {
		return x.Dy
	}
	return 0
}

func (x *InputState) GetButtons() map[string]bool {
	if x != nil {
		return x.Buttons
	}
	return nil
}

func (x *InputState) GetScrollY() int32 {
	if x != nil {
		return x.ScrollY
	}
	return 0
}

var File_input_state_proto protoreflect.FileDescriptor

const file_input_state_proto_rawDesc = "" +
	"\n" +
	"\x11input_state.proto\x12\vremotemouse\"\xc3\x01\n" +
	"\n" +
	"InputState\x12\x0e\n" +
	"\x02dx\x18\x01 \x01(\x05R\x02dx\x12\x0e\n" +
	"\x02dy\x18\x02 \x01(\x05R\x02dy\x12>\n" +
	"\abuttons\x18\x03 \x03(\v2$.remotemouse.InputState.ButtonsEntryR\abuttons\x12\x19\n" +
	"\bscroll_y\x18\x04 \x01(\x05R\ascrollY\x1a:\n" +
	"\fButtonsEntry\x12\x10\n" +
	"\x03key\x18\x01 \x01(\tR\x03key\x12\x14\n" +
	"\x05value\x18\x02 \x01(\bR\x05value:\x028\x01B\x15Z\x13./proto;remotemouseb\x06proto3"

var (
	file_input_state_proto_rawDescOnce sync.Once
	file_input_state_proto_rawDescData []byte
)

func file_input_state_proto_rawDescGZIP() []byte {
	file_input_state_proto_rawDescOnce.Do(func() {
		file_input_state_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_input_state_proto_rawDesc), len(file_input_state_proto_rawDesc)))
	})
	return file_input_state_proto_rawDescData
}

var file_input_state_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_input_state_proto_goTypes = []any{
	(*InputState)(nil), // 0: remotemouse.InputState
	nil,                // 1: remotemouse.InputState.ButtonsEntry
}
var file_input_state_proto_depIdxs = []int32{
	1, // 0: remotemouse.InputState.buttons:type_name -> remotemouse.InputState.ButtonsEntry
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_input_state_proto_init() }
func file_input_state_proto_init() {
	if File_input_state_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_input_state_proto_rawDesc), len(file_input_state_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_input_state_proto_goTypes,
		DependencyIndexes: file_input_state_proto_depIdxs,
		MessageInfos:      file_input_state_proto_msgTypes,
	}.Build()
	File_input_state_proto = out.File
	file_input_state_proto_goTypes = nil
	file_input_state_proto_depIdxs = nil
}

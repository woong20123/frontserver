// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.21.0-devel
// 	protoc        v3.11.4
// source: PacketCmd.proto

package examchatserverPacket

import (
	reflect "reflect"
	sync "sync"

	proto "github.com/golang/protobuf/proto"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
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

type Cmd int32

const (
	Cmd_packetCommandStart  Cmd = 0
	Cmd_F2CSServerRegistReq Cmd = 1
	Cmd_CS2FServerRegistRes Cmd = 2
	Cmd_LogicStart          Cmd = 10000
	Cmd_C2SLoginUserReq     Cmd = 10001
	Cmd_S2CLoginUserRes     Cmd = 10002
	Cmd_C2SLogOutUserReq    Cmd = 10003
	Cmd_S2CLogOutUserRes    Cmd = 10004
	Cmd_C2SLobbyMsgReq      Cmd = 10005
	Cmd_S2CLobbyMsgRes      Cmd = 10006
	Cmd_C2SRoomCreateReq    Cmd = 10007
	Cmd_S2CRoomCreateRes    Cmd = 10008
	Cmd_C2SRoomEnterReq     Cmd = 10009
	Cmd_S2CRoomEnterRes     Cmd = 10010
	Cmd_C2SRoomLeaveReq     Cmd = 10011
	Cmd_S2CRoomLeaveRes     Cmd = 10012
	Cmd_C2SRoomMsgReq       Cmd = 10013
	Cmd_S2CRoomMsgRes       Cmd = 10014
	Cmd_S2CSystemMsgSend    Cmd = 10015
)

// Enum value maps for Cmd.
var (
	Cmd_name = map[int32]string{
		0:     "packetCommandStart",
		1:     "F2CSPacketCmdSysServerRegistReq",
		2:     "C2FPacketCmdSysServerRegistRes",
		10000: "packetLogicCommandStart",
		10001: "C2SPacketCommandLoginUserReq",
		10002: "S2CPacketCommandLoginUserRes",
		10003: "C2SPacketCommandLogOutUserReq",
		10004: "S2CPacketCommandLogOutUserRes",
		10005: "C2SPacketCommandLobbyMsgReq",
		10006: "S2CPacketCommandLobbyMsgRes",
		10007: "C2SPacketCommandRoomCreateReq",
		10008: "S2CPacketCommandRoomCreateRes",
		10009: "C2SPacketCommandRoomEnterReq",
		10010: "S2CPacketCommandRoomEnterRes",
		10011: "C2SPacketCommandRoomLeaveReq",
		10012: "S2CPacketCommandRoomLeaveRes",
		10013: "C2SPacketCommandRoomMsgReq",
		10014: "S2CPacketCommandRoomMsgRes",
		10015: "S2CPacketCommandSystemMsgSend",
	}
	Cmd_value = map[string]int32{
		"packetCommandStart":              0,
		"F2CSPacketCmdSysServerRegistReq": 1,
		"C2FPacketCmdSysServerRegistRes":  2,
		"packetLogicCommandStart":         10000,
		"C2SPacketCommandLoginUserReq":    10001,
		"S2CPacketCommandLoginUserRes":    10002,
		"C2SPacketCommandLogOutUserReq":   10003,
		"S2CPacketCommandLogOutUserRes":   10004,
		"C2SPacketCommandLobbyMsgReq":     10005,
		"S2CPacketCommandLobbyMsgRes":     10006,
		"C2SPacketCommandRoomCreateReq":   10007,
		"S2CPacketCommandRoomCreateRes":   10008,
		"C2SPacketCommandRoomEnterReq":    10009,
		"S2CPacketCommandRoomEnterRes":    10010,
		"C2SPacketCommandRoomLeaveReq":    10011,
		"S2CPacketCommandRoomLeaveRes":    10012,
		"C2SPacketCommandRoomMsgReq":      10013,
		"S2CPacketCommandRoomMsgRes":      10014,
		"S2CPacketCommandSystemMsgSend":   10015,
	}
)

func (x Cmd) Enum() *Cmd {
	p := new(Cmd)
	*p = x
	return p
}

func (x Cmd) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Cmd) Descriptor() protoreflect.EnumDescriptor {
	return file_PacketCmd_proto_enumTypes[0].Descriptor()
}

func (Cmd) Type() protoreflect.EnumType {
	return &file_PacketCmd_proto_enumTypes[0]
}

func (x Cmd) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Cmd.Descriptor instead.
func (Cmd) EnumDescriptor() ([]byte, []int) {
	return file_PacketCmd_proto_rawDescGZIP(), []int{0}
}

var File_PacketCmd_proto protoreflect.FileDescriptor

var file_PacketCmd_proto_rawDesc = []byte{
	0x0a, 0x0f, 0x50, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x43, 0x6d, 0x64, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x14, 0x65, 0x78, 0x61, 0x6d, 0x63, 0x68, 0x61, 0x74, 0x73, 0x65, 0x72, 0x76, 0x65,
	0x72, 0x50, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x2a, 0x90, 0x05, 0x0a, 0x03, 0x43, 0x6d, 0x64, 0x12,
	0x16, 0x0a, 0x12, 0x70, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64,
	0x53, 0x74, 0x61, 0x72, 0x74, 0x10, 0x00, 0x12, 0x23, 0x0a, 0x1f, 0x46, 0x32, 0x43, 0x53, 0x50,
	0x61, 0x63, 0x6b, 0x65, 0x74, 0x43, 0x6d, 0x64, 0x53, 0x79, 0x73, 0x53, 0x65, 0x72, 0x76, 0x65,
	0x72, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x52, 0x65, 0x71, 0x10, 0x01, 0x12, 0x22, 0x0a, 0x1e,
	0x43, 0x32, 0x46, 0x50, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x43, 0x6d, 0x64, 0x53, 0x79, 0x73, 0x53,
	0x65, 0x72, 0x76, 0x65, 0x72, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x52, 0x65, 0x73, 0x10, 0x02,
	0x12, 0x1c, 0x0a, 0x17, 0x70, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x4c, 0x6f, 0x67, 0x69, 0x63, 0x43,
	0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x53, 0x74, 0x61, 0x72, 0x74, 0x10, 0x90, 0x4e, 0x12, 0x21,
	0x0a, 0x1c, 0x43, 0x32, 0x53, 0x50, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x43, 0x6f, 0x6d, 0x6d, 0x61,
	0x6e, 0x64, 0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x55, 0x73, 0x65, 0x72, 0x52, 0x65, 0x71, 0x10, 0x91,
	0x4e, 0x12, 0x21, 0x0a, 0x1c, 0x53, 0x32, 0x43, 0x50, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x43, 0x6f,
	0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x55, 0x73, 0x65, 0x72, 0x52, 0x65,
	0x73, 0x10, 0x92, 0x4e, 0x12, 0x22, 0x0a, 0x1d, 0x43, 0x32, 0x53, 0x50, 0x61, 0x63, 0x6b, 0x65,
	0x74, 0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x4c, 0x6f, 0x67, 0x4f, 0x75, 0x74, 0x55, 0x73,
	0x65, 0x72, 0x52, 0x65, 0x71, 0x10, 0x93, 0x4e, 0x12, 0x22, 0x0a, 0x1d, 0x53, 0x32, 0x43, 0x50,
	0x61, 0x63, 0x6b, 0x65, 0x74, 0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x4c, 0x6f, 0x67, 0x4f,
	0x75, 0x74, 0x55, 0x73, 0x65, 0x72, 0x52, 0x65, 0x73, 0x10, 0x94, 0x4e, 0x12, 0x20, 0x0a, 0x1b,
	0x43, 0x32, 0x53, 0x50, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64,
	0x4c, 0x6f, 0x62, 0x62, 0x79, 0x4d, 0x73, 0x67, 0x52, 0x65, 0x71, 0x10, 0x95, 0x4e, 0x12, 0x20,
	0x0a, 0x1b, 0x53, 0x32, 0x43, 0x50, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x43, 0x6f, 0x6d, 0x6d, 0x61,
	0x6e, 0x64, 0x4c, 0x6f, 0x62, 0x62, 0x79, 0x4d, 0x73, 0x67, 0x52, 0x65, 0x73, 0x10, 0x96, 0x4e,
	0x12, 0x22, 0x0a, 0x1d, 0x43, 0x32, 0x53, 0x50, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x43, 0x6f, 0x6d,
	0x6d, 0x61, 0x6e, 0x64, 0x52, 0x6f, 0x6f, 0x6d, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x52, 0x65,
	0x71, 0x10, 0x97, 0x4e, 0x12, 0x22, 0x0a, 0x1d, 0x53, 0x32, 0x43, 0x50, 0x61, 0x63, 0x6b, 0x65,
	0x74, 0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x52, 0x6f, 0x6f, 0x6d, 0x43, 0x72, 0x65, 0x61,
	0x74, 0x65, 0x52, 0x65, 0x73, 0x10, 0x98, 0x4e, 0x12, 0x21, 0x0a, 0x1c, 0x43, 0x32, 0x53, 0x50,
	0x61, 0x63, 0x6b, 0x65, 0x74, 0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x52, 0x6f, 0x6f, 0x6d,
	0x45, 0x6e, 0x74, 0x65, 0x72, 0x52, 0x65, 0x71, 0x10, 0x99, 0x4e, 0x12, 0x21, 0x0a, 0x1c, 0x53,
	0x32, 0x43, 0x50, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x52,
	0x6f, 0x6f, 0x6d, 0x45, 0x6e, 0x74, 0x65, 0x72, 0x52, 0x65, 0x73, 0x10, 0x9a, 0x4e, 0x12, 0x21,
	0x0a, 0x1c, 0x43, 0x32, 0x53, 0x50, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x43, 0x6f, 0x6d, 0x6d, 0x61,
	0x6e, 0x64, 0x52, 0x6f, 0x6f, 0x6d, 0x4c, 0x65, 0x61, 0x76, 0x65, 0x52, 0x65, 0x71, 0x10, 0x9b,
	0x4e, 0x12, 0x21, 0x0a, 0x1c, 0x53, 0x32, 0x43, 0x50, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x43, 0x6f,
	0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x52, 0x6f, 0x6f, 0x6d, 0x4c, 0x65, 0x61, 0x76, 0x65, 0x52, 0x65,
	0x73, 0x10, 0x9c, 0x4e, 0x12, 0x1f, 0x0a, 0x1a, 0x43, 0x32, 0x53, 0x50, 0x61, 0x63, 0x6b, 0x65,
	0x74, 0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x52, 0x6f, 0x6f, 0x6d, 0x4d, 0x73, 0x67, 0x52,
	0x65, 0x71, 0x10, 0x9d, 0x4e, 0x12, 0x1f, 0x0a, 0x1a, 0x53, 0x32, 0x43, 0x50, 0x61, 0x63, 0x6b,
	0x65, 0x74, 0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x52, 0x6f, 0x6f, 0x6d, 0x4d, 0x73, 0x67,
	0x52, 0x65, 0x73, 0x10, 0x9e, 0x4e, 0x12, 0x22, 0x0a, 0x1d, 0x53, 0x32, 0x43, 0x50, 0x61, 0x63,
	0x6b, 0x65, 0x74, 0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x53, 0x79, 0x73, 0x74, 0x65, 0x6d,
	0x4d, 0x73, 0x67, 0x53, 0x65, 0x6e, 0x64, 0x10, 0x9f, 0x4e, 0x42, 0x18, 0x5a, 0x16, 0x2e, 0x3b,
	0x65, 0x78, 0x61, 0x6d, 0x63, 0x68, 0x61, 0x74, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x50, 0x61,
	0x63, 0x6b, 0x65, 0x74, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_PacketCmd_proto_rawDescOnce sync.Once
	file_PacketCmd_proto_rawDescData = file_PacketCmd_proto_rawDesc
)

func file_PacketCmd_proto_rawDescGZIP() []byte {
	file_PacketCmd_proto_rawDescOnce.Do(func() {
		file_PacketCmd_proto_rawDescData = protoimpl.X.CompressGZIP(file_PacketCmd_proto_rawDescData)
	})
	return file_PacketCmd_proto_rawDescData
}

var file_PacketCmd_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_PacketCmd_proto_goTypes = []interface{}{
	(Cmd)(0), // 0: examchatserverPacket.Cmd
}
var file_PacketCmd_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_PacketCmd_proto_init() }
func file_PacketCmd_proto_init() {
	if File_PacketCmd_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_PacketCmd_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_PacketCmd_proto_goTypes,
		DependencyIndexes: file_PacketCmd_proto_depIdxs,
		EnumInfos:         file_PacketCmd_proto_enumTypes,
	}.Build()
	File_PacketCmd_proto = out.File
	file_PacketCmd_proto_rawDesc = nil
	file_PacketCmd_proto_goTypes = nil
	file_PacketCmd_proto_depIdxs = nil
}
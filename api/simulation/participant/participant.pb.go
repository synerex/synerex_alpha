// Code generated by protoc-gen-go. DO NOT EDIT.
// source: simulation/participant/participant.proto

package participant

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type StatusType int32

const (
	StatusType_OK   StatusType = 0
	StatusType_NG   StatusType = 1
	StatusType_NONE StatusType = 2
)

var StatusType_name = map[int32]string{
	0: "OK",
	1: "NG",
	2: "NONE",
}

var StatusType_value = map[string]int32{
	"OK":   0,
	"NG":   1,
	"NONE": 2,
}

func (x StatusType) String() string {
	return proto.EnumName(StatusType_name, int32(x))
}

func (StatusType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_6631733b0a9fa50e, []int{0}
}

type AgentType int32

const (
	AgentType_PEDESTRIAN AgentType = 0
	AgentType_CAR        AgentType = 1
	AgentType_TRAIN      AgentType = 2
	AgentType_BICYCLE    AgentType = 3
)

var AgentType_name = map[int32]string{
	0: "PEDESTRIAN",
	1: "CAR",
	2: "TRAIN",
	3: "BICYCLE",
}

var AgentType_value = map[string]int32{
	"PEDESTRIAN": 0,
	"CAR":        1,
	"TRAIN":      2,
	"BICYCLE":    3,
}

func (x AgentType) String() string {
	return proto.EnumName(AgentType_name, int32(x))
}

func (AgentType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_6631733b0a9fa50e, []int{1}
}

type ClientType int32

const (
	ClientType_AREA    ClientType = 0
	ClientType_LOG     ClientType = 1
	ClientType_CARAREA ClientType = 2
	ClientType_PEDAREA ClientType = 3
)

var ClientType_name = map[int32]string{
	0: "AREA",
	1: "LOG",
	2: "CARAREA",
	3: "PEDAREA",
}

var ClientType_value = map[string]int32{
	"AREA":    0,
	"LOG":     1,
	"CARAREA": 2,
	"PEDAREA": 3,
}

func (x ClientType) String() string {
	return proto.EnumName(ClientType_name, int32(x))
}

func (ClientType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_6631733b0a9fa50e, []int{2}
}

type DemandType int32

const (
	DemandType_GET DemandType = 0
)

var DemandType_name = map[int32]string{
	0: "GET",
}

var DemandType_value = map[string]int32{
	"GET": 0,
}

func (x DemandType) String() string {
	return proto.EnumName(DemandType_name, int32(x))
}

func (DemandType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_6631733b0a9fa50e, []int{3}
}

type ParticipantInfo struct {
	// participant info
	ClientId   uint64     `protobuf:"varint,1,opt,name=client_id,json=clientId,proto3" json:"client_id,omitempty"`
	ClientType ClientType `protobuf:"varint,2,opt,name=client_type,json=clientType,proto3,enum=api.participant.ClientType" json:"client_type,omitempty"`
	AreaId     uint32     `protobuf:"varint,3,opt,name=area_id,json=areaId,proto3" json:"area_id,omitempty"`
	AgentType  AgentType  `protobuf:"varint,4,opt,name=agent_type,json=agentType,proto3,enum=api.participant.AgentType" json:"agent_type,omitempty"`
	// meta data
	StatusType           StatusType `protobuf:"varint,5,opt,name=status_type,json=statusType,proto3,enum=api.participant.StatusType" json:"status_type,omitempty"`
	Meta                 string     `protobuf:"bytes,6,opt,name=meta,proto3" json:"meta,omitempty"`
	XXX_NoUnkeyedLiteral struct{}   `json:"-"`
	XXX_unrecognized     []byte     `json:"-"`
	XXX_sizecache        int32      `json:"-"`
}

func (m *ParticipantInfo) Reset()         { *m = ParticipantInfo{} }
func (m *ParticipantInfo) String() string { return proto.CompactTextString(m) }
func (*ParticipantInfo) ProtoMessage()    {}
func (*ParticipantInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_6631733b0a9fa50e, []int{0}
}

func (m *ParticipantInfo) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ParticipantInfo.Unmarshal(m, b)
}
func (m *ParticipantInfo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ParticipantInfo.Marshal(b, m, deterministic)
}
func (m *ParticipantInfo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ParticipantInfo.Merge(m, src)
}
func (m *ParticipantInfo) XXX_Size() int {
	return xxx_messageInfo_ParticipantInfo.Size(m)
}
func (m *ParticipantInfo) XXX_DiscardUnknown() {
	xxx_messageInfo_ParticipantInfo.DiscardUnknown(m)
}

var xxx_messageInfo_ParticipantInfo proto.InternalMessageInfo

func (m *ParticipantInfo) GetClientId() uint64 {
	if m != nil {
		return m.ClientId
	}
	return 0
}

func (m *ParticipantInfo) GetClientType() ClientType {
	if m != nil {
		return m.ClientType
	}
	return ClientType_AREA
}

func (m *ParticipantInfo) GetAreaId() uint32 {
	if m != nil {
		return m.AreaId
	}
	return 0
}

func (m *ParticipantInfo) GetAgentType() AgentType {
	if m != nil {
		return m.AgentType
	}
	return AgentType_PEDESTRIAN
}

func (m *ParticipantInfo) GetStatusType() StatusType {
	if m != nil {
		return m.StatusType
	}
	return StatusType_OK
}

func (m *ParticipantInfo) GetMeta() string {
	if m != nil {
		return m.Meta
	}
	return ""
}

type ParticipantDemand struct {
	// demand info
	ClientId uint64 `protobuf:"varint,1,opt,name=client_id,json=clientId,proto3" json:"client_id,omitempty"`
	// demand type
	DemandType DemandType `protobuf:"varint,2,opt,name=demand_type,json=demandType,proto3,enum=api.participant.DemandType" json:"demand_type,omitempty"`
	// meta data
	StatusType           StatusType `protobuf:"varint,3,opt,name=status_type,json=statusType,proto3,enum=api.participant.StatusType" json:"status_type,omitempty"`
	Meta                 string     `protobuf:"bytes,4,opt,name=meta,proto3" json:"meta,omitempty"`
	XXX_NoUnkeyedLiteral struct{}   `json:"-"`
	XXX_unrecognized     []byte     `json:"-"`
	XXX_sizecache        int32      `json:"-"`
}

func (m *ParticipantDemand) Reset()         { *m = ParticipantDemand{} }
func (m *ParticipantDemand) String() string { return proto.CompactTextString(m) }
func (*ParticipantDemand) ProtoMessage()    {}
func (*ParticipantDemand) Descriptor() ([]byte, []int) {
	return fileDescriptor_6631733b0a9fa50e, []int{1}
}

func (m *ParticipantDemand) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ParticipantDemand.Unmarshal(m, b)
}
func (m *ParticipantDemand) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ParticipantDemand.Marshal(b, m, deterministic)
}
func (m *ParticipantDemand) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ParticipantDemand.Merge(m, src)
}
func (m *ParticipantDemand) XXX_Size() int {
	return xxx_messageInfo_ParticipantDemand.Size(m)
}
func (m *ParticipantDemand) XXX_DiscardUnknown() {
	xxx_messageInfo_ParticipantDemand.DiscardUnknown(m)
}

var xxx_messageInfo_ParticipantDemand proto.InternalMessageInfo

func (m *ParticipantDemand) GetClientId() uint64 {
	if m != nil {
		return m.ClientId
	}
	return 0
}

func (m *ParticipantDemand) GetDemandType() DemandType {
	if m != nil {
		return m.DemandType
	}
	return DemandType_GET
}

func (m *ParticipantDemand) GetStatusType() StatusType {
	if m != nil {
		return m.StatusType
	}
	return StatusType_OK
}

func (m *ParticipantDemand) GetMeta() string {
	if m != nil {
		return m.Meta
	}
	return ""
}

func init() {
	proto.RegisterEnum("api.participant.StatusType", StatusType_name, StatusType_value)
	proto.RegisterEnum("api.participant.AgentType", AgentType_name, AgentType_value)
	proto.RegisterEnum("api.participant.ClientType", ClientType_name, ClientType_value)
	proto.RegisterEnum("api.participant.DemandType", DemandType_name, DemandType_value)
	proto.RegisterType((*ParticipantInfo)(nil), "api.participant.ParticipantInfo")
	proto.RegisterType((*ParticipantDemand)(nil), "api.participant.ParticipantDemand")
}

func init() {
	proto.RegisterFile("simulation/participant/participant.proto", fileDescriptor_6631733b0a9fa50e)
}

var fileDescriptor_6631733b0a9fa50e = []byte{
	// 412 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x92, 0xcf, 0x6e, 0x9b, 0x40,
	0x10, 0xc6, 0xbd, 0x40, 0xec, 0x30, 0x56, 0x93, 0xed, 0x4a, 0x55, 0xad, 0xe6, 0x62, 0xe5, 0x50,
	0x21, 0x0e, 0x20, 0xb5, 0xa7, 0xa8, 0xe9, 0x81, 0x60, 0x14, 0xa1, 0x46, 0xd8, 0xda, 0x70, 0x69,
	0x2f, 0xd1, 0x06, 0xb6, 0xce, 0x4a, 0x06, 0x56, 0xb0, 0x96, 0xea, 0xd7, 0xe8, 0x13, 0xf5, 0xd1,
	0x2a, 0x96, 0x1a, 0x2c, 0xf7, 0x9f, 0x94, 0xd3, 0x7e, 0x33, 0xcc, 0xfc, 0x98, 0xf9, 0x76, 0xc1,
	0x69, 0x44, 0xb1, 0xdd, 0x30, 0x25, 0xaa, 0xd2, 0x97, 0xac, 0x56, 0x22, 0x13, 0x92, 0x95, 0xea,
	0x50, 0x7b, 0xb2, 0xae, 0x54, 0x45, 0xce, 0x99, 0x14, 0xde, 0x41, 0xfa, 0xf2, 0xbb, 0x01, 0xe7,
	0xab, 0x21, 0x8e, 0xcb, 0xaf, 0x15, 0xb9, 0x00, 0x3b, 0xdb, 0x08, 0x5e, 0xaa, 0x07, 0x91, 0xcf,
	0xd0, 0x1c, 0x39, 0x16, 0x3d, 0xed, 0x12, 0x71, 0x4e, 0xae, 0x61, 0xfa, 0xeb, 0xa3, 0xda, 0x49,
	0x3e, 0x33, 0xe6, 0xc8, 0x39, 0x7b, 0x77, 0xe1, 0x1d, 0x71, 0xbd, 0x50, 0xd7, 0xa4, 0x3b, 0xc9,
	0x29, 0x64, 0xbd, 0x26, 0xaf, 0x61, 0xc2, 0x6a, 0xce, 0x5a, 0xb0, 0x39, 0x47, 0xce, 0x0b, 0x3a,
	0x6e, 0xc3, 0x38, 0x27, 0x57, 0x00, 0x6c, 0xdd, 0x53, 0x2d, 0x4d, 0x7d, 0xf3, 0x1b, 0x35, 0x58,
	0xef, 0xa1, 0x36, 0xdb, 0xcb, 0x76, 0xa2, 0x46, 0x31, 0xb5, 0x6d, 0xba, 0xde, 0x93, 0xbf, 0x4c,
	0x74, 0xaf, 0x6b, 0xba, 0x89, 0x9a, 0x5e, 0x13, 0x02, 0x56, 0xc1, 0x15, 0x9b, 0x8d, 0xe7, 0xc8,
	0xb1, 0xa9, 0xd6, 0x97, 0x3f, 0x10, 0xbc, 0x3c, 0x30, 0x65, 0xc1, 0x0b, 0x56, 0xe6, 0xff, 0xb5,
	0x25, 0xd7, 0x65, 0xff, 0xb6, 0xa5, 0x43, 0x75, 0x43, 0xe4, 0xbd, 0x3e, 0x5e, 0xc1, 0x7c, 0xde,
	0x0a, 0xd6, 0xb0, 0x82, 0xfb, 0x16, 0x60, 0xa8, 0x26, 0x63, 0x30, 0x96, 0x9f, 0xf0, 0xa8, 0x3d,
	0x93, 0x5b, 0x8c, 0xc8, 0x29, 0x58, 0xc9, 0x32, 0x89, 0xb0, 0xe1, 0x5e, 0x83, 0xdd, 0x9b, 0x4a,
	0xce, 0x00, 0x56, 0xd1, 0x22, 0xba, 0x4f, 0x69, 0x1c, 0x24, 0x78, 0x44, 0x26, 0x60, 0x86, 0x01,
	0xc5, 0x88, 0xd8, 0x70, 0x92, 0xd2, 0x20, 0x4e, 0xb0, 0x41, 0xa6, 0x30, 0xb9, 0x89, 0xc3, 0xcf,
	0xe1, 0x5d, 0x84, 0x4d, 0xf7, 0x0a, 0x60, 0xb8, 0xe8, 0x96, 0x1a, 0xd0, 0x28, 0xe8, 0x1a, 0xef,
	0x96, 0xed, 0x8f, 0xa6, 0x30, 0x09, 0x03, 0xaa, 0xb3, 0xba, 0x75, 0x15, 0x2d, 0x74, 0x60, 0xba,
	0xaf, 0x00, 0x06, 0x33, 0xda, 0x86, 0xdb, 0x28, 0xc5, 0xa3, 0x9b, 0x8f, 0x5f, 0x3e, 0xac, 0x85,
	0x7a, 0xda, 0x3e, 0x7a, 0x59, 0x55, 0xf8, 0xcd, 0xae, 0xe4, 0x35, 0xff, 0xb6, 0x3f, 0x1f, 0xd8,
	0x46, 0x3e, 0x31, 0x9f, 0x49, 0xe1, 0xff, 0xf9, 0xc5, 0x3f, 0x8e, 0xf5, 0x33, 0x7f, 0xff, 0x33,
	0x00, 0x00, 0xff, 0xff, 0x9c, 0x41, 0xb4, 0x1e, 0x12, 0x03, 0x00, 0x00,
}

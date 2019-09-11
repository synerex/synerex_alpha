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
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

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
	ClientParticipantId uint64     `protobuf:"varint,1,opt,name=client_participant_id,json=clientParticipantId,proto3" json:"client_participant_id,omitempty"`
	ClientAreaId        uint64     `protobuf:"varint,2,opt,name=client_area_id,json=clientAreaId,proto3" json:"client_area_id,omitempty"`
	ClientAgentId       uint64     `protobuf:"varint,3,opt,name=client_agent_id,json=clientAgentId,proto3" json:"client_agent_id,omitempty"`
	ClientClockId       uint64     `protobuf:"varint,4,opt,name=client_clock_id,json=clientClockId,proto3" json:"client_clock_id,omitempty"`
	ClientType          ClientType `protobuf:"varint,5,opt,name=client_type,json=clientType,proto3,enum=api.participant.ClientType" json:"client_type,omitempty"`
	AreaId              uint32     `protobuf:"varint,6,opt,name=area_id,json=areaId,proto3" json:"area_id,omitempty"`
	AgentType           AgentType  `protobuf:"varint,7,opt,name=agent_type,json=agentType,proto3,enum=api.participant.AgentType" json:"agent_type,omitempty"`
	// meta data
	StatusType           StatusType `protobuf:"varint,8,opt,name=status_type,json=statusType,proto3,enum=api.participant.StatusType" json:"status_type,omitempty"`
	Meta                 string     `protobuf:"bytes,9,opt,name=meta,proto3" json:"meta,omitempty"`
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

func (m *ParticipantInfo) GetClientParticipantId() uint64 {
	if m != nil {
		return m.ClientParticipantId
	}
	return 0
}

func (m *ParticipantInfo) GetClientAreaId() uint64 {
	if m != nil {
		return m.ClientAreaId
	}
	return 0
}

func (m *ParticipantInfo) GetClientAgentId() uint64 {
	if m != nil {
		return m.ClientAgentId
	}
	return 0
}

func (m *ParticipantInfo) GetClientClockId() uint64 {
	if m != nil {
		return m.ClientClockId
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
	// 473 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x93, 0xcf, 0x6e, 0xda, 0x40,
	0x10, 0xc6, 0xf1, 0x9f, 0x60, 0x3c, 0x34, 0xe0, 0x6e, 0x15, 0x15, 0x35, 0x17, 0x14, 0x55, 0x91,
	0xc5, 0xc1, 0x48, 0xe9, 0x29, 0x6a, 0x7a, 0x70, 0x8c, 0x15, 0x59, 0x8d, 0x0c, 0xda, 0x70, 0x69,
	0x2f, 0x68, 0x63, 0x6f, 0x93, 0x55, 0xc1, 0xb6, 0xec, 0x45, 0x2a, 0x4f, 0xd3, 0x57, 0xe9, 0xa3,
	0x55, 0xbb, 0x0b, 0xb6, 0x9b, 0x96, 0x4b, 0x4f, 0x9e, 0x1d, 0x7e, 0xdf, 0xc7, 0x37, 0xb3, 0x5a,
	0x70, 0x2b, 0xb6, 0xd9, 0xae, 0x09, 0x67, 0x79, 0x36, 0x2d, 0x48, 0xc9, 0x59, 0xc2, 0x0a, 0x92,
	0xf1, 0x76, 0xed, 0x15, 0x65, 0xce, 0x73, 0x34, 0x24, 0x05, 0xf3, 0x5a, 0xed, 0x8b, 0x9f, 0x06,
	0x0c, 0x17, 0xcd, 0x39, 0xca, 0xbe, 0xe5, 0xe8, 0x0a, 0xce, 0x92, 0x35, 0xa3, 0x19, 0x5f, 0xb5,
	0xc8, 0x15, 0x4b, 0x47, 0xda, 0x58, 0x73, 0x4d, 0xfc, 0x46, 0xfd, 0xd8, 0x56, 0xa5, 0xe8, 0x3d,
	0x0c, 0xf6, 0x1a, 0x52, 0x52, 0x22, 0x60, 0x5d, 0xc2, 0xaf, 0x54, 0xd7, 0x2f, 0x29, 0x89, 0x52,
	0x74, 0x09, 0xc3, 0x03, 0xf5, 0x44, 0x95, 0xa7, 0x21, 0xb1, 0xd3, 0x3d, 0x26, 0xba, 0x7f, 0x70,
	0xc9, 0x3a, 0x4f, 0xbe, 0x0b, 0xce, 0x6c, 0x73, 0x81, 0xe8, 0x46, 0x29, 0xba, 0x81, 0xfe, 0x9e,
	0xe3, 0xbb, 0x82, 0x8e, 0x4e, 0xc6, 0x9a, 0x3b, 0xb8, 0x3a, 0xf7, 0x5e, 0x0c, 0xe9, 0x05, 0x92,
	0x59, 0xee, 0x0a, 0x8a, 0x21, 0xa9, 0x6b, 0xf4, 0x16, 0xac, 0x43, 0xd8, 0xee, 0x58, 0x73, 0x4f,
	0x71, 0x97, 0xa8, 0x98, 0xd7, 0x00, 0x2a, 0x9f, 0x74, 0xb5, 0xa4, 0xeb, 0xbb, 0xbf, 0x5c, 0x65,
	0x58, 0x69, 0x6a, 0x93, 0x43, 0x29, 0x12, 0x55, 0x9c, 0xf0, 0x6d, 0xa5, 0xb4, 0xbd, 0x23, 0x89,
	0x1e, 0x24, 0xa3, 0x12, 0x55, 0x75, 0x8d, 0x10, 0x98, 0x1b, 0xca, 0xc9, 0xc8, 0x1e, 0x6b, 0xae,
	0x8d, 0x65, 0x7d, 0xf1, 0x4b, 0x83, 0xd7, 0xad, 0x5d, 0xcf, 0xe8, 0x86, 0x64, 0x29, 0x3a, 0x07,
	0x7b, 0x3f, 0x79, 0x7d, 0x2f, 0x3d, 0xd5, 0x50, 0x6b, 0x49, 0x25, 0xa6, 0x42, 0xe8, 0x47, 0x42,
	0x28, 0x2b, 0x15, 0x22, 0xad, 0xeb, 0x97, 0x23, 0x18, 0xff, 0x37, 0x82, 0xd9, 0x8c, 0x30, 0xb9,
	0x04, 0x68, 0x68, 0xd4, 0x05, 0x7d, 0xfe, 0xd9, 0xe9, 0x88, 0x6f, 0x7c, 0xe7, 0x68, 0xa8, 0x07,
	0x66, 0x3c, 0x8f, 0x43, 0x47, 0x9f, 0xdc, 0x80, 0x5d, 0x2f, 0x15, 0x0d, 0x00, 0x16, 0xe1, 0x2c,
	0x7c, 0x58, 0xe2, 0xc8, 0x8f, 0x9d, 0x0e, 0xb2, 0xc0, 0x08, 0x7c, 0xec, 0x68, 0xc8, 0x86, 0x93,
	0x25, 0xf6, 0xa3, 0xd8, 0xd1, 0x51, 0x1f, 0xac, 0xdb, 0x28, 0xf8, 0x12, 0xdc, 0x87, 0x8e, 0x31,
	0xb9, 0x06, 0x68, 0x2e, 0x5a, 0xb8, 0xfa, 0x38, 0xf4, 0x95, 0xf0, 0x7e, 0x2e, 0xfe, 0xa8, 0x0f,
	0x56, 0xe0, 0x63, 0xd9, 0x95, 0xd2, 0x45, 0x38, 0x93, 0x07, 0x63, 0x72, 0x06, 0xd0, 0x2c, 0x43,
	0x08, 0xee, 0xc2, 0xa5, 0xd3, 0xb9, 0xfd, 0xf4, 0xf5, 0xe3, 0x13, 0xe3, 0xcf, 0xdb, 0x47, 0x2f,
	0xc9, 0x37, 0xd3, 0x6a, 0x97, 0xd1, 0x92, 0xfe, 0x38, 0x7c, 0x57, 0x64, 0x5d, 0x3c, 0x93, 0x29,
	0x29, 0xd8, 0xf4, 0xdf, 0xcf, 0xef, 0xb1, 0x2b, 0xdf, 0xdc, 0x87, 0xdf, 0x01, 0x00, 0x00, 0xff,
	0xff, 0xd1, 0xda, 0xbd, 0x1e, 0x9f, 0x03, 0x00, 0x00,
}

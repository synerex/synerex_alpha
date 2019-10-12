// Code generated by protoc-gen-go. DO NOT EDIT.
// source: simulation/clock/clock.proto

package clock

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
	return fileDescriptor_e96fed1976809896, []int{0}
}

type DemandType int32

const (
	DemandType_FORWARD DemandType = 0
	DemandType_BACK    DemandType = 1
	DemandType_SET     DemandType = 2
	DemandType_START   DemandType = 3
	DemandType_STOP    DemandType = 4
)

var DemandType_name = map[int32]string{
	0: "FORWARD",
	1: "BACK",
	2: "SET",
	3: "START",
	4: "STOP",
}

var DemandType_value = map[string]int32{
	"FORWARD": 0,
	"BACK":    1,
	"SET":     2,
	"START":   3,
	"STOP":    4,
}

func (x DemandType) String() string {
	return proto.EnumName(DemandType_name, int32(x))
}

func (DemandType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_e96fed1976809896, []int{1}
}

type SupplyType int32

const (
	SupplyType_RES_FORWARD SupplyType = 0
	SupplyType_RES_BACK    SupplyType = 1
	SupplyType_RES_SET     SupplyType = 2
	SupplyType_RES_START   SupplyType = 3
	SupplyType_RES_STOP    SupplyType = 4
)

var SupplyType_name = map[int32]string{
	0: "RES_FORWARD",
	1: "RES_BACK",
	2: "RES_SET",
	3: "RES_START",
	4: "RES_STOP",
}

var SupplyType_value = map[string]int32{
	"RES_FORWARD": 0,
	"RES_BACK":    1,
	"RES_SET":     2,
	"RES_START":   3,
	"RES_STOP":    4,
}

func (x SupplyType) String() string {
	return proto.EnumName(SupplyType_name, int32(x))
}

func (SupplyType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_e96fed1976809896, []int{2}
}

type ClockInfo struct {
	// clock info
	Time uint32 `protobuf:"varint,1,opt,name=time,proto3" json:"time,omitempty"`
	// supply type
	SupplyType SupplyType `protobuf:"varint,2,opt,name=supply_type,json=supplyType,proto3,enum=api.clock.SupplyType" json:"supply_type,omitempty"`
	// meta data
	StatusType           StatusType `protobuf:"varint,3,opt,name=status_type,json=statusType,proto3,enum=api.clock.StatusType" json:"status_type,omitempty"`
	Meta                 string     `protobuf:"bytes,4,opt,name=meta,proto3" json:"meta,omitempty"`
	XXX_NoUnkeyedLiteral struct{}   `json:"-"`
	XXX_unrecognized     []byte     `json:"-"`
	XXX_sizecache        int32      `json:"-"`
}

func (m *ClockInfo) Reset()         { *m = ClockInfo{} }
func (m *ClockInfo) String() string { return proto.CompactTextString(m) }
func (*ClockInfo) ProtoMessage()    {}
func (*ClockInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_e96fed1976809896, []int{0}
}

func (m *ClockInfo) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ClockInfo.Unmarshal(m, b)
}
func (m *ClockInfo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ClockInfo.Marshal(b, m, deterministic)
}
func (m *ClockInfo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ClockInfo.Merge(m, src)
}
func (m *ClockInfo) XXX_Size() int {
	return xxx_messageInfo_ClockInfo.Size(m)
}
func (m *ClockInfo) XXX_DiscardUnknown() {
	xxx_messageInfo_ClockInfo.DiscardUnknown(m)
}

var xxx_messageInfo_ClockInfo proto.InternalMessageInfo

func (m *ClockInfo) GetTime() uint32 {
	if m != nil {
		return m.Time
	}
	return 0
}

func (m *ClockInfo) GetSupplyType() SupplyType {
	if m != nil {
		return m.SupplyType
	}
	return SupplyType_RES_FORWARD
}

func (m *ClockInfo) GetStatusType() StatusType {
	if m != nil {
		return m.StatusType
	}
	return StatusType_OK
}

func (m *ClockInfo) GetMeta() string {
	if m != nil {
		return m.Meta
	}
	return ""
}

type ClockDemand struct {
	// demand info
	Time          uint32 `protobuf:"varint,1,opt,name=time,proto3" json:"time,omitempty"`
	CycleNum      uint32 `protobuf:"varint,2,opt,name=cycle_num,json=cycleNum,proto3" json:"cycle_num,omitempty"`
	CycleDuration uint32 `protobuf:"varint,3,opt,name=cycle_duration,json=cycleDuration,proto3" json:"cycle_duration,omitempty"`
	CycleInterval uint32 `protobuf:"varint,4,opt,name=cycle_interval,json=cycleInterval,proto3" json:"cycle_interval,omitempty"`
	// demand type
	DemandType DemandType `protobuf:"varint,5,opt,name=demand_type,json=demandType,proto3,enum=api.clock.DemandType" json:"demand_type,omitempty"`
	// meta data
	StatusType           StatusType `protobuf:"varint,6,opt,name=status_type,json=statusType,proto3,enum=api.clock.StatusType" json:"status_type,omitempty"`
	Meta                 string     `protobuf:"bytes,7,opt,name=meta,proto3" json:"meta,omitempty"`
	XXX_NoUnkeyedLiteral struct{}   `json:"-"`
	XXX_unrecognized     []byte     `json:"-"`
	XXX_sizecache        int32      `json:"-"`
}

func (m *ClockDemand) Reset()         { *m = ClockDemand{} }
func (m *ClockDemand) String() string { return proto.CompactTextString(m) }
func (*ClockDemand) ProtoMessage()    {}
func (*ClockDemand) Descriptor() ([]byte, []int) {
	return fileDescriptor_e96fed1976809896, []int{1}
}

func (m *ClockDemand) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ClockDemand.Unmarshal(m, b)
}
func (m *ClockDemand) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ClockDemand.Marshal(b, m, deterministic)
}
func (m *ClockDemand) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ClockDemand.Merge(m, src)
}
func (m *ClockDemand) XXX_Size() int {
	return xxx_messageInfo_ClockDemand.Size(m)
}
func (m *ClockDemand) XXX_DiscardUnknown() {
	xxx_messageInfo_ClockDemand.DiscardUnknown(m)
}

var xxx_messageInfo_ClockDemand proto.InternalMessageInfo

func (m *ClockDemand) GetTime() uint32 {
	if m != nil {
		return m.Time
	}
	return 0
}

func (m *ClockDemand) GetCycleNum() uint32 {
	if m != nil {
		return m.CycleNum
	}
	return 0
}

func (m *ClockDemand) GetCycleDuration() uint32 {
	if m != nil {
		return m.CycleDuration
	}
	return 0
}

func (m *ClockDemand) GetCycleInterval() uint32 {
	if m != nil {
		return m.CycleInterval
	}
	return 0
}

func (m *ClockDemand) GetDemandType() DemandType {
	if m != nil {
		return m.DemandType
	}
	return DemandType_FORWARD
}

func (m *ClockDemand) GetStatusType() StatusType {
	if m != nil {
		return m.StatusType
	}
	return StatusType_OK
}

func (m *ClockDemand) GetMeta() string {
	if m != nil {
		return m.Meta
	}
	return ""
}

func init() {
	proto.RegisterEnum("api.clock.StatusType", StatusType_name, StatusType_value)
	proto.RegisterEnum("api.clock.DemandType", DemandType_name, DemandType_value)
	proto.RegisterEnum("api.clock.SupplyType", SupplyType_name, SupplyType_value)
	proto.RegisterType((*ClockInfo)(nil), "api.clock.ClockInfo")
	proto.RegisterType((*ClockDemand)(nil), "api.clock.ClockDemand")
}

func init() { proto.RegisterFile("simulation/clock/clock.proto", fileDescriptor_e96fed1976809896) }

var fileDescriptor_e96fed1976809896 = []byte{
	// 411 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x52, 0x5d, 0x8b, 0xd3, 0x40,
	0x14, 0xdd, 0xa4, 0xdd, 0xb6, 0xb9, 0x31, 0xeb, 0x30, 0x20, 0x14, 0xf4, 0xa1, 0x2c, 0x28, 0xa5,
	0x0f, 0x29, 0x28, 0xea, 0x73, 0x77, 0x5b, 0x65, 0x59, 0x68, 0x64, 0x12, 0x11, 0x7c, 0x09, 0xb3,
	0xc9, 0xe8, 0x06, 0xf3, 0x31, 0x24, 0x13, 0x31, 0xff, 0xc2, 0xff, 0xe0, 0x1f, 0x95, 0xb9, 0x93,
	0x4d, 0x60, 0xeb, 0x8b, 0x2f, 0x99, 0x33, 0xe7, 0x9e, 0x7b, 0x38, 0xf7, 0x66, 0xe0, 0x45, 0x93,
	0x15, 0x6d, 0xce, 0x55, 0x56, 0x95, 0xdb, 0x24, 0xaf, 0x92, 0x1f, 0xe6, 0xeb, 0xcb, 0xba, 0x52,
	0x15, 0x75, 0xb8, 0xcc, 0x7c, 0x24, 0x2e, 0xff, 0x58, 0xe0, 0x5c, 0x6b, 0x74, 0x53, 0x7e, 0xab,
	0x28, 0x85, 0xa9, 0xca, 0x0a, 0xb1, 0xb4, 0x56, 0xd6, 0xda, 0x63, 0x88, 0xe9, 0x3b, 0x70, 0x9b,
	0x56, 0xca, 0xbc, 0x8b, 0x55, 0x27, 0xc5, 0xd2, 0x5e, 0x59, 0xeb, 0x8b, 0xd7, 0xcf, 0xfc, 0xc1,
	0xc2, 0x0f, 0xb1, 0x1a, 0x75, 0x52, 0x30, 0x68, 0x06, 0x8c, 0x7d, 0x8a, 0xab, 0xb6, 0x31, 0x7d,
	0x93, 0xd3, 0x3e, 0xac, 0xf6, 0x7d, 0x03, 0xd6, 0x19, 0x0a, 0xa1, 0xf8, 0x72, 0xba, 0xb2, 0xd6,
	0x0e, 0x43, 0x7c, 0xf9, 0xdb, 0x06, 0x17, 0x53, 0xee, 0x45, 0xc1, 0xcb, 0xf4, 0x9f, 0x39, 0x9f,
	0x83, 0x93, 0x74, 0x49, 0x2e, 0xe2, 0xb2, 0x2d, 0x30, 0xa5, 0xc7, 0x16, 0x48, 0x1c, 0xdb, 0x82,
	0xbe, 0x84, 0x0b, 0x53, 0x4c, 0xdb, 0x1a, 0xb7, 0x82, 0x79, 0x3c, 0xe6, 0x21, 0xbb, 0xef, 0xc9,
	0x51, 0x96, 0x95, 0x4a, 0xd4, 0x3f, 0x79, 0x8e, 0x29, 0x1e, 0x64, 0x37, 0x3d, 0xa9, 0x47, 0x4b,
	0x31, 0x88, 0x19, 0xed, 0xfc, 0x64, 0x34, 0x13, 0xd3, 0x8c, 0x96, 0x0e, 0xf8, 0xf1, 0x4a, 0x66,
	0xff, 0xbb, 0x92, 0xf9, 0xb8, 0x92, 0xcd, 0x2b, 0x80, 0x51, 0x4d, 0x67, 0x60, 0x07, 0xb7, 0xe4,
	0x4c, 0x9f, 0xc7, 0x8f, 0xc4, 0xa2, 0x0b, 0x98, 0x1e, 0x83, 0xe3, 0x81, 0xd8, 0x9b, 0x1d, 0xc0,
	0x98, 0x86, 0xba, 0x30, 0xff, 0x10, 0xb0, 0x2f, 0x3b, 0xb6, 0x27, 0x67, 0x5a, 0x74, 0xb5, 0xbb,
	0xbe, 0x25, 0x16, 0x9d, 0xc3, 0x24, 0x3c, 0x44, 0xc4, 0xa6, 0x0e, 0x9c, 0x87, 0xd1, 0x8e, 0x45,
	0x64, 0xa2, 0xab, 0x61, 0x14, 0x7c, 0x22, 0xd3, 0xcd, 0x67, 0x80, 0xf1, 0x1f, 0xd3, 0xa7, 0xe0,
	0xb2, 0x43, 0x18, 0x8f, 0x36, 0x4f, 0x60, 0xa1, 0x89, 0xde, 0xca, 0x85, 0xb9, 0xbe, 0x19, 0x3b,
	0x0f, 0x1c, 0xbc, 0xf4, 0x96, 0xbd, 0xd2, 0xd8, 0x5e, 0xbd, 0xff, 0xfa, 0xf6, 0x7b, 0xa6, 0xee,
	0xdb, 0x3b, 0x3f, 0xa9, 0x8a, 0x6d, 0xd3, 0x95, 0xa2, 0x16, 0xbf, 0x1e, 0xce, 0x98, 0xe7, 0xf2,
	0x9e, 0x6f, 0xb9, 0xcc, 0xb6, 0x8f, 0x9f, 0xf2, 0xdd, 0x0c, 0x5f, 0xf1, 0x9b, 0xbf, 0x01, 0x00,
	0x00, 0xff, 0xff, 0x30, 0x4a, 0xd1, 0x18, 0xe5, 0x02, 0x00, 0x00,
}
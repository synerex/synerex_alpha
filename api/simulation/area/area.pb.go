// Code generated by protoc-gen-go. DO NOT EDIT.
// source: simulation/area/area.proto

package area

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

type DemandType int32

const (
	DemandType_SET DemandType = 0
	DemandType_GET DemandType = 1
)

var DemandType_name = map[int32]string{
	0: "SET",
	1: "GET",
}

var DemandType_value = map[string]int32{
	"SET": 0,
	"GET": 1,
}

func (x DemandType) String() string {
	return proto.EnumName(DemandType_name, int32(x))
}

func (DemandType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_c8212774c97f163d, []int{0}
}

type SupplyType int32

const (
	SupplyType_RES_SET SupplyType = 0
	SupplyType_RES_GET SupplyType = 1
)

var SupplyType_name = map[int32]string{
	0: "RES_SET",
	1: "RES_GET",
}

var SupplyType_value = map[string]int32{
	"RES_SET": 0,
	"RES_GET": 1,
}

func (x SupplyType) String() string {
	return proto.EnumName(SupplyType_name, int32(x))
}

func (SupplyType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_c8212774c97f163d, []int{1}
}

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
	return fileDescriptor_c8212774c97f163d, []int{2}
}

type AreaInfo struct {
	// area info
	Time     uint32   `protobuf:"varint,1,opt,name=time,proto3" json:"time,omitempty"`
	AreaId   uint32   `protobuf:"varint,2,opt,name=area_id,json=areaId,proto3" json:"area_id,omitempty"`
	AreaName string   `protobuf:"bytes,3,opt,name=area_name,json=areaName,proto3" json:"area_name,omitempty"`
	Signal   *Signal  `protobuf:"bytes,4,opt,name=signal,proto3" json:"signal,omitempty"`
	Map      *Map     `protobuf:"bytes,5,opt,name=map,proto3" json:"map,omitempty"`
	Climate  *Climate `protobuf:"bytes,6,opt,name=climate,proto3" json:"climate,omitempty"`
	// supply type
	SupplyType SupplyType `protobuf:"varint,7,opt,name=supply_type,json=supplyType,proto3,enum=api.area.SupplyType" json:"supply_type,omitempty"`
	// meta data
	StatusType           StatusType `protobuf:"varint,8,opt,name=status_type,json=statusType,proto3,enum=api.area.StatusType" json:"status_type,omitempty"`
	Meta                 string     `protobuf:"bytes,9,opt,name=meta,proto3" json:"meta,omitempty"`
	XXX_NoUnkeyedLiteral struct{}   `json:"-"`
	XXX_unrecognized     []byte     `json:"-"`
	XXX_sizecache        int32      `json:"-"`
}

func (m *AreaInfo) Reset()         { *m = AreaInfo{} }
func (m *AreaInfo) String() string { return proto.CompactTextString(m) }
func (*AreaInfo) ProtoMessage()    {}
func (*AreaInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_c8212774c97f163d, []int{0}
}

func (m *AreaInfo) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AreaInfo.Unmarshal(m, b)
}
func (m *AreaInfo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AreaInfo.Marshal(b, m, deterministic)
}
func (m *AreaInfo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AreaInfo.Merge(m, src)
}
func (m *AreaInfo) XXX_Size() int {
	return xxx_messageInfo_AreaInfo.Size(m)
}
func (m *AreaInfo) XXX_DiscardUnknown() {
	xxx_messageInfo_AreaInfo.DiscardUnknown(m)
}

var xxx_messageInfo_AreaInfo proto.InternalMessageInfo

func (m *AreaInfo) GetTime() uint32 {
	if m != nil {
		return m.Time
	}
	return 0
}

func (m *AreaInfo) GetAreaId() uint32 {
	if m != nil {
		return m.AreaId
	}
	return 0
}

func (m *AreaInfo) GetAreaName() string {
	if m != nil {
		return m.AreaName
	}
	return ""
}

func (m *AreaInfo) GetSignal() *Signal {
	if m != nil {
		return m.Signal
	}
	return nil
}

func (m *AreaInfo) GetMap() *Map {
	if m != nil {
		return m.Map
	}
	return nil
}

func (m *AreaInfo) GetClimate() *Climate {
	if m != nil {
		return m.Climate
	}
	return nil
}

func (m *AreaInfo) GetSupplyType() SupplyType {
	if m != nil {
		return m.SupplyType
	}
	return SupplyType_RES_SET
}

func (m *AreaInfo) GetStatusType() StatusType {
	if m != nil {
		return m.StatusType
	}
	return StatusType_OK
}

func (m *AreaInfo) GetMeta() string {
	if m != nil {
		return m.Meta
	}
	return ""
}

type AreaDemand struct {
	// demand info
	Time   uint32 `protobuf:"varint,1,opt,name=time,proto3" json:"time,omitempty"`
	AreaId uint32 `protobuf:"varint,2,opt,name=area_id,json=areaId,proto3" json:"area_id,omitempty"`
	// demand type
	DemandType DemandType `protobuf:"varint,3,opt,name=demand_type,json=demandType,proto3,enum=api.area.DemandType" json:"demand_type,omitempty"`
	// meta data
	StatusType           StatusType `protobuf:"varint,4,opt,name=status_type,json=statusType,proto3,enum=api.area.StatusType" json:"status_type,omitempty"`
	Meta                 string     `protobuf:"bytes,5,opt,name=meta,proto3" json:"meta,omitempty"`
	XXX_NoUnkeyedLiteral struct{}   `json:"-"`
	XXX_unrecognized     []byte     `json:"-"`
	XXX_sizecache        int32      `json:"-"`
}

func (m *AreaDemand) Reset()         { *m = AreaDemand{} }
func (m *AreaDemand) String() string { return proto.CompactTextString(m) }
func (*AreaDemand) ProtoMessage()    {}
func (*AreaDemand) Descriptor() ([]byte, []int) {
	return fileDescriptor_c8212774c97f163d, []int{1}
}

func (m *AreaDemand) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AreaDemand.Unmarshal(m, b)
}
func (m *AreaDemand) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AreaDemand.Marshal(b, m, deterministic)
}
func (m *AreaDemand) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AreaDemand.Merge(m, src)
}
func (m *AreaDemand) XXX_Size() int {
	return xxx_messageInfo_AreaDemand.Size(m)
}
func (m *AreaDemand) XXX_DiscardUnknown() {
	xxx_messageInfo_AreaDemand.DiscardUnknown(m)
}

var xxx_messageInfo_AreaDemand proto.InternalMessageInfo

func (m *AreaDemand) GetTime() uint32 {
	if m != nil {
		return m.Time
	}
	return 0
}

func (m *AreaDemand) GetAreaId() uint32 {
	if m != nil {
		return m.AreaId
	}
	return 0
}

func (m *AreaDemand) GetDemandType() DemandType {
	if m != nil {
		return m.DemandType
	}
	return DemandType_SET
}

func (m *AreaDemand) GetStatusType() StatusType {
	if m != nil {
		return m.StatusType
	}
	return StatusType_OK
}

func (m *AreaDemand) GetMeta() string {
	if m != nil {
		return m.Meta
	}
	return ""
}

type Signal struct {
	SignalInfo           uint32   `protobuf:"varint,1,opt,name=signal_info,json=signalInfo,proto3" json:"signal_info,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Signal) Reset()         { *m = Signal{} }
func (m *Signal) String() string { return proto.CompactTextString(m) }
func (*Signal) ProtoMessage()    {}
func (*Signal) Descriptor() ([]byte, []int) {
	return fileDescriptor_c8212774c97f163d, []int{2}
}

func (m *Signal) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Signal.Unmarshal(m, b)
}
func (m *Signal) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Signal.Marshal(b, m, deterministic)
}
func (m *Signal) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Signal.Merge(m, src)
}
func (m *Signal) XXX_Size() int {
	return xxx_messageInfo_Signal.Size(m)
}
func (m *Signal) XXX_DiscardUnknown() {
	xxx_messageInfo_Signal.DiscardUnknown(m)
}

var xxx_messageInfo_Signal proto.InternalMessageInfo

func (m *Signal) GetSignalInfo() uint32 {
	if m != nil {
		return m.SignalInfo
	}
	return 0
}

type Map struct {
	Coord                *Map_Coord `protobuf:"bytes,1,opt,name=coord,proto3" json:"coord,omitempty"`
	Neighbor             []uint32   `protobuf:"varint,2,rep,packed,name=neighbor,proto3" json:"neighbor,omitempty"`
	Controlled           *Map_Coord `protobuf:"bytes,3,opt,name=controlled,proto3" json:"controlled,omitempty"`
	MapInfo              uint32     `protobuf:"varint,4,opt,name=map_info,json=mapInfo,proto3" json:"map_info,omitempty"`
	XXX_NoUnkeyedLiteral struct{}   `json:"-"`
	XXX_unrecognized     []byte     `json:"-"`
	XXX_sizecache        int32      `json:"-"`
}

func (m *Map) Reset()         { *m = Map{} }
func (m *Map) String() string { return proto.CompactTextString(m) }
func (*Map) ProtoMessage()    {}
func (*Map) Descriptor() ([]byte, []int) {
	return fileDescriptor_c8212774c97f163d, []int{3}
}

func (m *Map) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Map.Unmarshal(m, b)
}
func (m *Map) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Map.Marshal(b, m, deterministic)
}
func (m *Map) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Map.Merge(m, src)
}
func (m *Map) XXX_Size() int {
	return xxx_messageInfo_Map.Size(m)
}
func (m *Map) XXX_DiscardUnknown() {
	xxx_messageInfo_Map.DiscardUnknown(m)
}

var xxx_messageInfo_Map proto.InternalMessageInfo

func (m *Map) GetCoord() *Map_Coord {
	if m != nil {
		return m.Coord
	}
	return nil
}

func (m *Map) GetNeighbor() []uint32 {
	if m != nil {
		return m.Neighbor
	}
	return nil
}

func (m *Map) GetControlled() *Map_Coord {
	if m != nil {
		return m.Controlled
	}
	return nil
}

func (m *Map) GetMapInfo() uint32 {
	if m != nil {
		return m.MapInfo
	}
	return 0
}

type Map_Coord struct {
	StartLat             float32  `protobuf:"fixed32,1,opt,name=start_lat,json=startLat,proto3" json:"start_lat,omitempty"`
	StartLon             float32  `protobuf:"fixed32,2,opt,name=start_lon,json=startLon,proto3" json:"start_lon,omitempty"`
	EndLat               float32  `protobuf:"fixed32,3,opt,name=end_lat,json=endLat,proto3" json:"end_lat,omitempty"`
	EndLon               float32  `protobuf:"fixed32,4,opt,name=end_lon,json=endLon,proto3" json:"end_lon,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Map_Coord) Reset()         { *m = Map_Coord{} }
func (m *Map_Coord) String() string { return proto.CompactTextString(m) }
func (*Map_Coord) ProtoMessage()    {}
func (*Map_Coord) Descriptor() ([]byte, []int) {
	return fileDescriptor_c8212774c97f163d, []int{3, 0}
}

func (m *Map_Coord) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Map_Coord.Unmarshal(m, b)
}
func (m *Map_Coord) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Map_Coord.Marshal(b, m, deterministic)
}
func (m *Map_Coord) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Map_Coord.Merge(m, src)
}
func (m *Map_Coord) XXX_Size() int {
	return xxx_messageInfo_Map_Coord.Size(m)
}
func (m *Map_Coord) XXX_DiscardUnknown() {
	xxx_messageInfo_Map_Coord.DiscardUnknown(m)
}

var xxx_messageInfo_Map_Coord proto.InternalMessageInfo

func (m *Map_Coord) GetStartLat() float32 {
	if m != nil {
		return m.StartLat
	}
	return 0
}

func (m *Map_Coord) GetStartLon() float32 {
	if m != nil {
		return m.StartLon
	}
	return 0
}

func (m *Map_Coord) GetEndLat() float32 {
	if m != nil {
		return m.EndLat
	}
	return 0
}

func (m *Map_Coord) GetEndLon() float32 {
	if m != nil {
		return m.EndLon
	}
	return 0
}

type Climate struct {
	Weather              uint32   `protobuf:"varint,1,opt,name=Weather,proto3" json:"Weather,omitempty"`
	Temperature          uint32   `protobuf:"varint,2,opt,name=Temperature,proto3" json:"Temperature,omitempty"`
	Humidity             uint32   `protobuf:"varint,3,opt,name=Humidity,proto3" json:"Humidity,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Climate) Reset()         { *m = Climate{} }
func (m *Climate) String() string { return proto.CompactTextString(m) }
func (*Climate) ProtoMessage()    {}
func (*Climate) Descriptor() ([]byte, []int) {
	return fileDescriptor_c8212774c97f163d, []int{4}
}

func (m *Climate) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Climate.Unmarshal(m, b)
}
func (m *Climate) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Climate.Marshal(b, m, deterministic)
}
func (m *Climate) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Climate.Merge(m, src)
}
func (m *Climate) XXX_Size() int {
	return xxx_messageInfo_Climate.Size(m)
}
func (m *Climate) XXX_DiscardUnknown() {
	xxx_messageInfo_Climate.DiscardUnknown(m)
}

var xxx_messageInfo_Climate proto.InternalMessageInfo

func (m *Climate) GetWeather() uint32 {
	if m != nil {
		return m.Weather
	}
	return 0
}

func (m *Climate) GetTemperature() uint32 {
	if m != nil {
		return m.Temperature
	}
	return 0
}

func (m *Climate) GetHumidity() uint32 {
	if m != nil {
		return m.Humidity
	}
	return 0
}

func init() {
	proto.RegisterEnum("api.area.DemandType", DemandType_name, DemandType_value)
	proto.RegisterEnum("api.area.SupplyType", SupplyType_name, SupplyType_value)
	proto.RegisterEnum("api.area.StatusType", StatusType_name, StatusType_value)
	proto.RegisterType((*AreaInfo)(nil), "api.area.AreaInfo")
	proto.RegisterType((*AreaDemand)(nil), "api.area.AreaDemand")
	proto.RegisterType((*Signal)(nil), "api.area.Signal")
	proto.RegisterType((*Map)(nil), "api.area.Map")
	proto.RegisterType((*Map_Coord)(nil), "api.area.Map.Coord")
	proto.RegisterType((*Climate)(nil), "api.area.Climate")
}

func init() { proto.RegisterFile("simulation/area/area.proto", fileDescriptor_c8212774c97f163d) }

var fileDescriptor_c8212774c97f163d = []byte{
	// 587 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x54, 0xdd, 0x6a, 0xdb, 0x4c,
	0x10, 0x8d, 0x24, 0x5b, 0x52, 0x46, 0xf8, 0xc3, 0xdf, 0xb6, 0x50, 0x35, 0x85, 0xd6, 0xf8, 0xa2,
	0x38, 0x29, 0xd8, 0x90, 0x34, 0xbd, 0x6f, 0x53, 0x93, 0x86, 0x36, 0x0e, 0xc8, 0x86, 0x42, 0x6f,
	0xcc, 0xc4, 0xda, 0xc4, 0x0b, 0xda, 0x1f, 0xa4, 0x35, 0xd4, 0x8f, 0xd1, 0xe7, 0xe9, 0x3b, 0xf4,
	0x99, 0xca, 0x8e, 0xe4, 0x9f, 0xa6, 0xf4, 0x22, 0x37, 0xd2, 0xcc, 0x9c, 0x33, 0x3b, 0x67, 0x67,
	0x86, 0x85, 0xa3, 0x4a, 0xc8, 0x55, 0x81, 0x56, 0x68, 0x35, 0xc2, 0x92, 0x23, 0x7d, 0x86, 0xa6,
	0xd4, 0x56, 0xb3, 0x18, 0x8d, 0x18, 0x3a, 0xbf, 0xff, 0xcb, 0x87, 0xf8, 0x7d, 0xc9, 0xf1, 0x4a,
	0xdd, 0x69, 0xc6, 0xa0, 0x65, 0x85, 0xe4, 0xa9, 0xd7, 0xf3, 0x06, 0x9d, 0x8c, 0x6c, 0xf6, 0x0c,
	0x22, 0x47, 0x9c, 0x8b, 0x3c, 0xf5, 0x29, 0x1c, 0x3a, 0xf7, 0x2a, 0x67, 0x2f, 0xe0, 0x90, 0x00,
	0x85, 0x92, 0xa7, 0x41, 0xcf, 0x1b, 0x1c, 0x66, 0xb1, 0x0b, 0x4c, 0x50, 0x72, 0x36, 0x80, 0xb0,
	0x12, 0xf7, 0x0a, 0x8b, 0xb4, 0xd5, 0xf3, 0x06, 0xc9, 0x69, 0x77, 0xb8, 0xa9, 0x38, 0x9c, 0x52,
	0x3c, 0x6b, 0x70, 0xf6, 0x0a, 0x02, 0x89, 0x26, 0x6d, 0x13, 0xad, 0xb3, 0xa3, 0x5d, 0xa3, 0xc9,
	0x1c, 0xc2, 0xde, 0x40, 0xb4, 0x28, 0x84, 0x44, 0xcb, 0xd3, 0x90, 0x48, 0xff, 0xef, 0x48, 0x17,
	0x35, 0x90, 0x6d, 0x18, 0xec, 0x1c, 0x92, 0x6a, 0x65, 0x4c, 0xb1, 0x9e, 0xdb, 0xb5, 0xe1, 0x69,
	0xd4, 0xf3, 0x06, 0xff, 0x9d, 0x3e, 0xdd, 0x2b, 0x4e, 0xe0, 0x6c, 0x6d, 0x78, 0x06, 0xd5, 0xd6,
	0xa6, 0x34, 0x8b, 0x76, 0x55, 0xd5, 0x69, 0xf1, 0x5f, 0x69, 0x04, 0x36, 0x69, 0x5b, 0xdb, 0xf5,
	0x4b, 0x72, 0x8b, 0xe9, 0x21, 0xdd, 0x9e, 0xec, 0xfe, 0x4f, 0x0f, 0xc0, 0x35, 0xf4, 0x23, 0x97,
	0xa8, 0xf2, 0xc7, 0xb5, 0xf4, 0x1c, 0x92, 0x9c, 0xd2, 0x6a, 0x19, 0xc1, 0x43, 0x19, 0xf5, 0x99,
	0xb5, 0x8c, 0x7c, 0x6b, 0x3f, 0x54, 0xdf, 0x7a, 0xa4, 0xfa, 0xf6, 0x9e, 0xfa, 0x63, 0x08, 0xa7,
	0x9b, 0xb9, 0x24, 0xf5, 0x84, 0xe6, 0x42, 0xdd, 0xe9, 0x46, 0x3f, 0xd4, 0x21, 0xb7, 0x2c, 0xfd,
	0x1f, 0x3e, 0x04, 0xd7, 0x68, 0xd8, 0x31, 0xb4, 0x17, 0x5a, 0x97, 0x39, 0x51, 0x92, 0xd3, 0x27,
	0x7f, 0x8c, 0x70, 0x78, 0xe1, 0xa0, 0xac, 0x66, 0xb0, 0x23, 0x88, 0x15, 0x17, 0xf7, 0xcb, 0x5b,
	0x5d, 0xa6, 0x7e, 0x2f, 0x18, 0x74, 0xb2, 0xad, 0xcf, 0xce, 0x00, 0x16, 0x5a, 0xd9, 0x52, 0x17,
	0x05, 0xcf, 0xe9, 0xea, 0xff, 0x38, 0x6b, 0x8f, 0xc6, 0x9e, 0x43, 0x2c, 0xd1, 0xd4, 0x0a, 0x5b,
	0xa4, 0x30, 0x92, 0x68, 0x9c, 0xbc, 0xa3, 0x0a, 0xda, 0xc4, 0x77, 0x7b, 0x5a, 0x59, 0x2c, 0xed,
	0xbc, 0x40, 0x4b, 0x1a, 0xfd, 0x2c, 0xa6, 0xc0, 0x17, 0xb4, 0x7b, 0xa0, 0x56, 0x34, 0x8c, 0x2d,
	0xa8, 0x95, 0x9b, 0x13, 0x57, 0x39, 0xe5, 0x05, 0x04, 0x85, 0x5c, 0xe5, 0x2e, 0x6b, 0x03, 0x68,
	0x45, 0x55, 0x1b, 0x40, 0xab, 0x3e, 0x42, 0xd4, 0xac, 0x24, 0x4b, 0x21, 0xfa, 0xca, 0xd1, 0x2e,
	0x79, 0xd9, 0xf4, 0x6e, 0xe3, 0xb2, 0x1e, 0x24, 0x33, 0x2e, 0x0d, 0x2f, 0xd1, 0xae, 0x4a, 0xde,
	0xac, 0xc0, 0x7e, 0xc8, 0xf5, 0xe9, 0xd3, 0x4a, 0x8a, 0x5c, 0xd8, 0x35, 0x55, 0xee, 0x64, 0x5b,
	0xff, 0xe4, 0x25, 0xc0, 0x6e, 0x0d, 0x58, 0x04, 0xc1, 0x74, 0x3c, 0xeb, 0x1e, 0x38, 0xe3, 0x72,
	0x3c, 0xeb, 0x7a, 0x27, 0xaf, 0x01, 0x76, 0x4b, 0xce, 0x12, 0x88, 0xb2, 0xf1, 0x74, 0x5e, 0x73,
	0x1a, 0x67, 0xc7, 0xdb, 0xed, 0x42, 0x08, 0xfe, 0xcd, 0xe7, 0xee, 0x81, 0xfb, 0x4f, 0x2e, 0xbb,
	0x1e, 0x8b, 0xa1, 0x35, 0xb9, 0x99, 0x8c, 0xbb, 0xfe, 0x87, 0x77, 0xdf, 0xde, 0xde, 0x0b, 0xbb,
	0x5c, 0xdd, 0x0e, 0x17, 0x5a, 0x8e, 0xaa, 0xb5, 0xe2, 0x25, 0xff, 0xbe, 0xf9, 0xcf, 0xb1, 0x30,
	0x4b, 0x1c, 0xa1, 0x11, 0xa3, 0x07, 0xaf, 0xcd, 0x6d, 0x48, 0x2f, 0xcd, 0xd9, 0xef, 0x00, 0x00,
	0x00, 0xff, 0xff, 0xdf, 0xca, 0x49, 0xd1, 0x87, 0x04, 0x00, 0x00,
}

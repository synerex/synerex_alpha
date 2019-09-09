// Code generated by protoc-gen-go. DO NOT EDIT.
// source: simulation/agent/agent.proto

package agent

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
	return fileDescriptor_fce67ac898dc274e, []int{0}
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
	return fileDescriptor_fce67ac898dc274e, []int{1}
}

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
	return fileDescriptor_fce67ac898dc274e, []int{2}
}

type AgentInfo struct {
	// agent info
	Time        uint32       `protobuf:"varint,1,opt,name=time,proto3" json:"time,omitempty"`
	AgentId     uint32       `protobuf:"varint,2,opt,name=agent_id,json=agentId,proto3" json:"agent_id,omitempty"`
	AgentName   string       `protobuf:"bytes,3,opt,name=agent_name,json=agentName,proto3" json:"agent_name,omitempty"`
	AgentStatus *AgentStatus `protobuf:"bytes,4,opt,name=agent_status,json=agentStatus,proto3" json:"agent_status,omitempty"`
	AgentType   AgentType    `protobuf:"varint,5,opt,name=agent_type,json=agentType,proto3,enum=api.agent.AgentType" json:"agent_type,omitempty"`
	Route       *Route       `protobuf:"bytes,6,opt,name=route,proto3" json:"route,omitempty"`
	Rule        *Rule        `protobuf:"bytes,7,opt,name=rule,proto3" json:"rule,omitempty"`
	// meta data
	StatusType           StatusType `protobuf:"varint,8,opt,name=status_type,json=statusType,proto3,enum=api.agent.StatusType" json:"status_type,omitempty"`
	Meta                 string     `protobuf:"bytes,9,opt,name=meta,proto3" json:"meta,omitempty"`
	XXX_NoUnkeyedLiteral struct{}   `json:"-"`
	XXX_unrecognized     []byte     `json:"-"`
	XXX_sizecache        int32      `json:"-"`
}

func (m *AgentInfo) Reset()         { *m = AgentInfo{} }
func (m *AgentInfo) String() string { return proto.CompactTextString(m) }
func (*AgentInfo) ProtoMessage()    {}
func (*AgentInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_fce67ac898dc274e, []int{0}
}

func (m *AgentInfo) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AgentInfo.Unmarshal(m, b)
}
func (m *AgentInfo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AgentInfo.Marshal(b, m, deterministic)
}
func (m *AgentInfo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AgentInfo.Merge(m, src)
}
func (m *AgentInfo) XXX_Size() int {
	return xxx_messageInfo_AgentInfo.Size(m)
}
func (m *AgentInfo) XXX_DiscardUnknown() {
	xxx_messageInfo_AgentInfo.DiscardUnknown(m)
}

var xxx_messageInfo_AgentInfo proto.InternalMessageInfo

func (m *AgentInfo) GetTime() uint32 {
	if m != nil {
		return m.Time
	}
	return 0
}

func (m *AgentInfo) GetAgentId() uint32 {
	if m != nil {
		return m.AgentId
	}
	return 0
}

func (m *AgentInfo) GetAgentName() string {
	if m != nil {
		return m.AgentName
	}
	return ""
}

func (m *AgentInfo) GetAgentStatus() *AgentStatus {
	if m != nil {
		return m.AgentStatus
	}
	return nil
}

func (m *AgentInfo) GetAgentType() AgentType {
	if m != nil {
		return m.AgentType
	}
	return AgentType_PEDESTRIAN
}

func (m *AgentInfo) GetRoute() *Route {
	if m != nil {
		return m.Route
	}
	return nil
}

func (m *AgentInfo) GetRule() *Rule {
	if m != nil {
		return m.Rule
	}
	return nil
}

func (m *AgentInfo) GetStatusType() StatusType {
	if m != nil {
		return m.StatusType
	}
	return StatusType_OK
}

func (m *AgentInfo) GetMeta() string {
	if m != nil {
		return m.Meta
	}
	return ""
}

type AgentsInfo struct {
	// agents info
	Time      uint32       `protobuf:"varint,1,opt,name=time,proto3" json:"time,omitempty"`
	AgentType AgentType    `protobuf:"varint,2,opt,name=agent_type,json=agentType,proto3,enum=api.agent.AgentType" json:"agent_type,omitempty"`
	AgentInfo []*AgentInfo `protobuf:"bytes,3,rep,name=agent_info,json=agentInfo,proto3" json:"agent_info,omitempty"`
	// meta data
	StatusType           StatusType `protobuf:"varint,4,opt,name=status_type,json=statusType,proto3,enum=api.agent.StatusType" json:"status_type,omitempty"`
	Meta                 string     `protobuf:"bytes,5,opt,name=meta,proto3" json:"meta,omitempty"`
	XXX_NoUnkeyedLiteral struct{}   `json:"-"`
	XXX_unrecognized     []byte     `json:"-"`
	XXX_sizecache        int32      `json:"-"`
}

func (m *AgentsInfo) Reset()         { *m = AgentsInfo{} }
func (m *AgentsInfo) String() string { return proto.CompactTextString(m) }
func (*AgentsInfo) ProtoMessage()    {}
func (*AgentsInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_fce67ac898dc274e, []int{1}
}

func (m *AgentsInfo) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AgentsInfo.Unmarshal(m, b)
}
func (m *AgentsInfo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AgentsInfo.Marshal(b, m, deterministic)
}
func (m *AgentsInfo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AgentsInfo.Merge(m, src)
}
func (m *AgentsInfo) XXX_Size() int {
	return xxx_messageInfo_AgentsInfo.Size(m)
}
func (m *AgentsInfo) XXX_DiscardUnknown() {
	xxx_messageInfo_AgentsInfo.DiscardUnknown(m)
}

var xxx_messageInfo_AgentsInfo proto.InternalMessageInfo

func (m *AgentsInfo) GetTime() uint32 {
	if m != nil {
		return m.Time
	}
	return 0
}

func (m *AgentsInfo) GetAgentType() AgentType {
	if m != nil {
		return m.AgentType
	}
	return AgentType_PEDESTRIAN
}

func (m *AgentsInfo) GetAgentInfo() []*AgentInfo {
	if m != nil {
		return m.AgentInfo
	}
	return nil
}

func (m *AgentsInfo) GetStatusType() StatusType {
	if m != nil {
		return m.StatusType
	}
	return StatusType_OK
}

func (m *AgentsInfo) GetMeta() string {
	if m != nil {
		return m.Meta
	}
	return ""
}

type AgentDemand struct {
	// demand info
	Time        uint32       `protobuf:"varint,1,opt,name=time,proto3" json:"time,omitempty"`
	AgentId     uint32       `protobuf:"varint,2,opt,name=agent_id,json=agentId,proto3" json:"agent_id,omitempty"`
	AgentName   string       `protobuf:"bytes,3,opt,name=agent_name,json=agentName,proto3" json:"agent_name,omitempty"`
	AgentType   AgentType    `protobuf:"varint,4,opt,name=agent_type,json=agentType,proto3,enum=api.agent.AgentType" json:"agent_type,omitempty"`
	AgentStatus *AgentStatus `protobuf:"bytes,5,opt,name=agent_status,json=agentStatus,proto3" json:"agent_status,omitempty"`
	Route       *Route       `protobuf:"bytes,6,opt,name=route,proto3" json:"route,omitempty"`
	Rule        *Rule        `protobuf:"bytes,7,opt,name=rule,proto3" json:"rule,omitempty"`
	// demand info
	DemandType DemandType `protobuf:"varint,8,opt,name=demand_type,json=demandType,proto3,enum=api.agent.DemandType" json:"demand_type,omitempty"`
	// meta data
	StatusType           StatusType `protobuf:"varint,9,opt,name=status_type,json=statusType,proto3,enum=api.agent.StatusType" json:"status_type,omitempty"`
	Meta                 string     `protobuf:"bytes,10,opt,name=meta,proto3" json:"meta,omitempty"`
	XXX_NoUnkeyedLiteral struct{}   `json:"-"`
	XXX_unrecognized     []byte     `json:"-"`
	XXX_sizecache        int32      `json:"-"`
}

func (m *AgentDemand) Reset()         { *m = AgentDemand{} }
func (m *AgentDemand) String() string { return proto.CompactTextString(m) }
func (*AgentDemand) ProtoMessage()    {}
func (*AgentDemand) Descriptor() ([]byte, []int) {
	return fileDescriptor_fce67ac898dc274e, []int{2}
}

func (m *AgentDemand) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AgentDemand.Unmarshal(m, b)
}
func (m *AgentDemand) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AgentDemand.Marshal(b, m, deterministic)
}
func (m *AgentDemand) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AgentDemand.Merge(m, src)
}
func (m *AgentDemand) XXX_Size() int {
	return xxx_messageInfo_AgentDemand.Size(m)
}
func (m *AgentDemand) XXX_DiscardUnknown() {
	xxx_messageInfo_AgentDemand.DiscardUnknown(m)
}

var xxx_messageInfo_AgentDemand proto.InternalMessageInfo

func (m *AgentDemand) GetTime() uint32 {
	if m != nil {
		return m.Time
	}
	return 0
}

func (m *AgentDemand) GetAgentId() uint32 {
	if m != nil {
		return m.AgentId
	}
	return 0
}

func (m *AgentDemand) GetAgentName() string {
	if m != nil {
		return m.AgentName
	}
	return ""
}

func (m *AgentDemand) GetAgentType() AgentType {
	if m != nil {
		return m.AgentType
	}
	return AgentType_PEDESTRIAN
}

func (m *AgentDemand) GetAgentStatus() *AgentStatus {
	if m != nil {
		return m.AgentStatus
	}
	return nil
}

func (m *AgentDemand) GetRoute() *Route {
	if m != nil {
		return m.Route
	}
	return nil
}

func (m *AgentDemand) GetRule() *Rule {
	if m != nil {
		return m.Rule
	}
	return nil
}

func (m *AgentDemand) GetDemandType() DemandType {
	if m != nil {
		return m.DemandType
	}
	return DemandType_SET
}

func (m *AgentDemand) GetStatusType() StatusType {
	if m != nil {
		return m.StatusType
	}
	return StatusType_OK
}

func (m *AgentDemand) GetMeta() string {
	if m != nil {
		return m.Meta
	}
	return ""
}

type AgentsDemand struct {
	// demand info
	AreaId    uint32    `protobuf:"varint,1,opt,name=area_id,json=areaId,proto3" json:"area_id,omitempty"`
	AgentType AgentType `protobuf:"varint,2,opt,name=agent_type,json=agentType,proto3,enum=api.agent.AgentType" json:"agent_type,omitempty"`
	// demand info
	DemandType DemandType `protobuf:"varint,3,opt,name=demand_type,json=demandType,proto3,enum=api.agent.DemandType" json:"demand_type,omitempty"`
	// meta data
	StatusType           StatusType `protobuf:"varint,4,opt,name=status_type,json=statusType,proto3,enum=api.agent.StatusType" json:"status_type,omitempty"`
	Meta                 string     `protobuf:"bytes,5,opt,name=meta,proto3" json:"meta,omitempty"`
	XXX_NoUnkeyedLiteral struct{}   `json:"-"`
	XXX_unrecognized     []byte     `json:"-"`
	XXX_sizecache        int32      `json:"-"`
}

func (m *AgentsDemand) Reset()         { *m = AgentsDemand{} }
func (m *AgentsDemand) String() string { return proto.CompactTextString(m) }
func (*AgentsDemand) ProtoMessage()    {}
func (*AgentsDemand) Descriptor() ([]byte, []int) {
	return fileDescriptor_fce67ac898dc274e, []int{3}
}

func (m *AgentsDemand) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AgentsDemand.Unmarshal(m, b)
}
func (m *AgentsDemand) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AgentsDemand.Marshal(b, m, deterministic)
}
func (m *AgentsDemand) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AgentsDemand.Merge(m, src)
}
func (m *AgentsDemand) XXX_Size() int {
	return xxx_messageInfo_AgentsDemand.Size(m)
}
func (m *AgentsDemand) XXX_DiscardUnknown() {
	xxx_messageInfo_AgentsDemand.DiscardUnknown(m)
}

var xxx_messageInfo_AgentsDemand proto.InternalMessageInfo

func (m *AgentsDemand) GetAreaId() uint32 {
	if m != nil {
		return m.AreaId
	}
	return 0
}

func (m *AgentsDemand) GetAgentType() AgentType {
	if m != nil {
		return m.AgentType
	}
	return AgentType_PEDESTRIAN
}

func (m *AgentsDemand) GetDemandType() DemandType {
	if m != nil {
		return m.DemandType
	}
	return DemandType_SET
}

func (m *AgentsDemand) GetStatusType() StatusType {
	if m != nil {
		return m.StatusType
	}
	return StatusType_OK
}

func (m *AgentsDemand) GetMeta() string {
	if m != nil {
		return m.Meta
	}
	return ""
}

type Route struct {
	Coord                *Route_Coord `protobuf:"bytes,1,opt,name=coord,proto3" json:"coord,omitempty"`
	Direction            float32      `protobuf:"fixed32,2,opt,name=direction,proto3" json:"direction,omitempty"`
	Speed                float32      `protobuf:"fixed32,3,opt,name=speed,proto3" json:"speed,omitempty"`
	Destination          float32      `protobuf:"fixed32,4,opt,name=destination,proto3" json:"destination,omitempty"`
	Departure            float32      `protobuf:"fixed32,5,opt,name=departure,proto3" json:"departure,omitempty"`
	XXX_NoUnkeyedLiteral struct{}     `json:"-"`
	XXX_unrecognized     []byte       `json:"-"`
	XXX_sizecache        int32        `json:"-"`
}

func (m *Route) Reset()         { *m = Route{} }
func (m *Route) String() string { return proto.CompactTextString(m) }
func (*Route) ProtoMessage()    {}
func (*Route) Descriptor() ([]byte, []int) {
	return fileDescriptor_fce67ac898dc274e, []int{4}
}

func (m *Route) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Route.Unmarshal(m, b)
}
func (m *Route) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Route.Marshal(b, m, deterministic)
}
func (m *Route) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Route.Merge(m, src)
}
func (m *Route) XXX_Size() int {
	return xxx_messageInfo_Route.Size(m)
}
func (m *Route) XXX_DiscardUnknown() {
	xxx_messageInfo_Route.DiscardUnknown(m)
}

var xxx_messageInfo_Route proto.InternalMessageInfo

func (m *Route) GetCoord() *Route_Coord {
	if m != nil {
		return m.Coord
	}
	return nil
}

func (m *Route) GetDirection() float32 {
	if m != nil {
		return m.Direction
	}
	return 0
}

func (m *Route) GetSpeed() float32 {
	if m != nil {
		return m.Speed
	}
	return 0
}

func (m *Route) GetDestination() float32 {
	if m != nil {
		return m.Destination
	}
	return 0
}

func (m *Route) GetDeparture() float32 {
	if m != nil {
		return m.Departure
	}
	return 0
}

type Route_Coord struct {
	Lat                  float32  `protobuf:"fixed32,1,opt,name=lat,proto3" json:"lat,omitempty"`
	Lon                  float32  `protobuf:"fixed32,2,opt,name=lon,proto3" json:"lon,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Route_Coord) Reset()         { *m = Route_Coord{} }
func (m *Route_Coord) String() string { return proto.CompactTextString(m) }
func (*Route_Coord) ProtoMessage()    {}
func (*Route_Coord) Descriptor() ([]byte, []int) {
	return fileDescriptor_fce67ac898dc274e, []int{4, 0}
}

func (m *Route_Coord) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Route_Coord.Unmarshal(m, b)
}
func (m *Route_Coord) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Route_Coord.Marshal(b, m, deterministic)
}
func (m *Route_Coord) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Route_Coord.Merge(m, src)
}
func (m *Route_Coord) XXX_Size() int {
	return xxx_messageInfo_Route_Coord.Size(m)
}
func (m *Route_Coord) XXX_DiscardUnknown() {
	xxx_messageInfo_Route_Coord.DiscardUnknown(m)
}

var xxx_messageInfo_Route_Coord proto.InternalMessageInfo

func (m *Route_Coord) GetLat() float32 {
	if m != nil {
		return m.Lat
	}
	return 0
}

func (m *Route_Coord) GetLon() float32 {
	if m != nil {
		return m.Lon
	}
	return 0
}

type Rule struct {
	RuleInfo             string   `protobuf:"bytes,1,opt,name=rule_info,json=ruleInfo,proto3" json:"rule_info,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Rule) Reset()         { *m = Rule{} }
func (m *Rule) String() string { return proto.CompactTextString(m) }
func (*Rule) ProtoMessage()    {}
func (*Rule) Descriptor() ([]byte, []int) {
	return fileDescriptor_fce67ac898dc274e, []int{5}
}

func (m *Rule) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Rule.Unmarshal(m, b)
}
func (m *Rule) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Rule.Marshal(b, m, deterministic)
}
func (m *Rule) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Rule.Merge(m, src)
}
func (m *Rule) XXX_Size() int {
	return xxx_messageInfo_Rule.Size(m)
}
func (m *Rule) XXX_DiscardUnknown() {
	xxx_messageInfo_Rule.DiscardUnknown(m)
}

var xxx_messageInfo_Rule proto.InternalMessageInfo

func (m *Rule) GetRuleInfo() string {
	if m != nil {
		return m.RuleInfo
	}
	return ""
}

type AgentStatus struct {
	Age                  string   `protobuf:"bytes,1,opt,name=age,proto3" json:"age,omitempty"`
	Sex                  string   `protobuf:"bytes,2,opt,name=sex,proto3" json:"sex,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *AgentStatus) Reset()         { *m = AgentStatus{} }
func (m *AgentStatus) String() string { return proto.CompactTextString(m) }
func (*AgentStatus) ProtoMessage()    {}
func (*AgentStatus) Descriptor() ([]byte, []int) {
	return fileDescriptor_fce67ac898dc274e, []int{6}
}

func (m *AgentStatus) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AgentStatus.Unmarshal(m, b)
}
func (m *AgentStatus) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AgentStatus.Marshal(b, m, deterministic)
}
func (m *AgentStatus) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AgentStatus.Merge(m, src)
}
func (m *AgentStatus) XXX_Size() int {
	return xxx_messageInfo_AgentStatus.Size(m)
}
func (m *AgentStatus) XXX_DiscardUnknown() {
	xxx_messageInfo_AgentStatus.DiscardUnknown(m)
}

var xxx_messageInfo_AgentStatus proto.InternalMessageInfo

func (m *AgentStatus) GetAge() string {
	if m != nil {
		return m.Age
	}
	return ""
}

func (m *AgentStatus) GetSex() string {
	if m != nil {
		return m.Sex
	}
	return ""
}

func init() {
	proto.RegisterEnum("api.agent.StatusType", StatusType_name, StatusType_value)
	proto.RegisterEnum("api.agent.AgentType", AgentType_name, AgentType_value)
	proto.RegisterEnum("api.agent.DemandType", DemandType_name, DemandType_value)
	proto.RegisterType((*AgentInfo)(nil), "api.agent.AgentInfo")
	proto.RegisterType((*AgentsInfo)(nil), "api.agent.AgentsInfo")
	proto.RegisterType((*AgentDemand)(nil), "api.agent.AgentDemand")
	proto.RegisterType((*AgentsDemand)(nil), "api.agent.AgentsDemand")
	proto.RegisterType((*Route)(nil), "api.agent.Route")
	proto.RegisterType((*Route_Coord)(nil), "api.agent.Route.Coord")
	proto.RegisterType((*Rule)(nil), "api.agent.Rule")
	proto.RegisterType((*AgentStatus)(nil), "api.agent.AgentStatus")
}

func init() { proto.RegisterFile("simulation/agent/agent.proto", fileDescriptor_fce67ac898dc274e) }

var fileDescriptor_fce67ac898dc274e = []byte{
	// 660 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xac, 0x55, 0x5d, 0x4b, 0x1b, 0x4d,
	0x14, 0x76, 0xf6, 0x23, 0xc9, 0x9e, 0xf8, 0xfa, 0x2e, 0x83, 0x6d, 0xb7, 0xad, 0x2d, 0x21, 0x82,
	0x04, 0x5b, 0x12, 0xaa, 0xb4, 0xa5, 0xd0, 0x9b, 0x18, 0x83, 0x84, 0x96, 0xb5, 0x8c, 0xb9, 0x69,
	0x6f, 0x64, 0x74, 0x47, 0x5d, 0xc8, 0x7e, 0xb0, 0x3b, 0x0b, 0xfa, 0x5f, 0xfa, 0xbb, 0xbc, 0x2e,
	0xfd, 0x25, 0x65, 0xce, 0xac, 0xc9, 0xba, 0x45, 0x88, 0xd6, 0x9b, 0xdd, 0x33, 0xe7, 0x63, 0xce,
	0x79, 0xce, 0x79, 0x66, 0x06, 0x36, 0xf2, 0x30, 0x2a, 0x66, 0x5c, 0x86, 0x49, 0x3c, 0xe0, 0xe7,
	0x22, 0x96, 0xfa, 0xdb, 0x4f, 0xb3, 0x44, 0x26, 0xd4, 0xe1, 0x69, 0xd8, 0x47, 0x45, 0xf7, 0xb7,
	0x01, 0xce, 0x50, 0x49, 0x93, 0xf8, 0x2c, 0xa1, 0x14, 0x2c, 0x19, 0x46, 0xc2, 0x23, 0x1d, 0xd2,
	0xfb, 0x8f, 0xa1, 0x4c, 0x9f, 0x43, 0x0b, 0x5d, 0x8f, 0xc3, 0xc0, 0x33, 0x50, 0xdf, 0xc4, 0xf5,
	0x24, 0xa0, 0xaf, 0x00, 0xb4, 0x29, 0xe6, 0x91, 0xf0, 0xcc, 0x0e, 0xe9, 0x39, 0xcc, 0x41, 0x8d,
	0xcf, 0x23, 0x41, 0x3f, 0xc1, 0xaa, 0x36, 0xe7, 0x92, 0xcb, 0x22, 0xf7, 0xac, 0x0e, 0xe9, 0xb5,
	0x77, 0x9e, 0xf6, 0xe7, 0xd9, 0xfb, 0x98, 0xf9, 0x08, 0xad, 0xac, 0xcd, 0x17, 0x0b, 0xba, 0x7b,
	0xb3, 0xb3, 0xbc, 0x4a, 0x85, 0x67, 0x77, 0x48, 0x6f, 0x6d, 0x67, 0xbd, 0x1e, 0x38, 0xbd, 0x4a,
	0x45, 0x99, 0x4f, 0x89, 0x74, 0x0b, 0xec, 0x2c, 0x29, 0xa4, 0xf0, 0x1a, 0x98, 0xc8, 0xad, 0xf8,
	0x33, 0xa5, 0x67, 0xda, 0x4c, 0x37, 0xc1, 0xca, 0x8a, 0x99, 0xf0, 0x9a, 0xe8, 0xf6, 0x7f, 0xd5,
	0xad, 0x98, 0x09, 0x86, 0x46, 0xfa, 0x01, 0xda, 0xba, 0x6c, 0x5d, 0x42, 0x0b, 0x4b, 0x78, 0x52,
	0xf1, 0xd5, 0x95, 0x62, 0x0d, 0x90, 0xcf, 0x65, 0xd5, 0xc2, 0x48, 0x48, 0xee, 0x39, 0xd8, 0x0d,
	0x94, 0xbb, 0xd7, 0x04, 0x00, 0x2b, 0xce, 0xef, 0xec, 0xf2, 0x6d, 0xc0, 0xc6, 0x72, 0x80, 0xe7,
	0x41, 0x61, 0x7c, 0x96, 0x78, 0x66, 0xc7, 0xec, 0xb5, 0xff, 0x0e, 0x52, 0x29, 0xcb, 0x20, 0xcc,
	0x5e, 0x03, 0x66, 0xdd, 0x17, 0x98, 0x5d, 0x01, 0xf6, 0xd3, 0x84, 0x36, 0x26, 0xd9, 0x17, 0x11,
	0x8f, 0x83, 0x47, 0xe6, 0xcf, 0xed, 0x9e, 0x58, 0xcb, 0xf5, 0xa4, 0x4e, 0x3a, 0x7b, 0x79, 0xd2,
	0x3d, 0x36, 0x7f, 0x02, 0x6c, 0xca, 0x5d, 0xfc, 0xd1, 0x2d, 0xd3, 0x6d, 0x0e, 0xe6, 0x72, 0x7d,
	0x3c, 0xce, 0x7d, 0xc7, 0x03, 0x95, 0xf1, 0xfc, 0x22, 0xb0, 0xaa, 0x79, 0x57, 0xce, 0xe7, 0x19,
	0x34, 0x79, 0x26, 0xb8, 0x1a, 0x85, 0x1e, 0x51, 0x43, 0x2d, 0x27, 0xc1, 0xc3, 0xe8, 0x57, 0x83,
	0x68, 0x3e, 0x10, 0xe2, 0x3f, 0x31, 0xf0, 0x9a, 0x80, 0x8d, 0xc3, 0xa1, 0x6f, 0xc1, 0x3e, 0x4d,
	0x92, 0x4c, 0x23, 0xbb, 0x3d, 0x71, 0x74, 0xe8, 0x8f, 0x94, 0x95, 0x69, 0x27, 0xba, 0x01, 0x4e,
	0x10, 0x66, 0xe2, 0x54, 0xdd, 0x90, 0x88, 0xd7, 0x60, 0x0b, 0x05, 0x5d, 0x07, 0x3b, 0x4f, 0x85,
	0x08, 0x10, 0x93, 0xc1, 0xf4, 0x82, 0x76, 0x14, 0xde, 0x5c, 0x86, 0x31, 0xde, 0xab, 0x58, 0xb7,
	0xc1, 0xaa, 0x2a, 0xdc, 0x55, 0xa4, 0x3c, 0x93, 0x45, 0xa6, 0x6f, 0x2d, 0xb5, 0xeb, 0x8d, 0xe2,
	0xc5, 0x1b, 0xb0, 0xb1, 0x06, 0xea, 0x82, 0x39, 0xe3, 0x12, 0x0b, 0x35, 0x98, 0x12, 0x51, 0x33,
	0x2f, 0x44, 0x89, 0xdd, 0x4d, 0xb0, 0x14, 0x9b, 0xe8, 0x4b, 0x70, 0x14, 0x9f, 0xf4, 0x11, 0x27,
	0x88, 0xbc, 0xa5, 0x14, 0xea, 0x2c, 0x77, 0xdf, 0x95, 0xc7, 0xaf, 0x24, 0xb0, 0x0b, 0x26, 0x3f,
	0x17, 0xa5, 0x97, 0x12, 0x95, 0x26, 0x17, 0x97, 0xb8, 0xaf, 0xc3, 0x94, 0xb8, 0xbd, 0x05, 0xb0,
	0x68, 0x2f, 0x6d, 0x80, 0x71, 0xf8, 0xc5, 0x5d, 0x51, 0x7f, 0xff, 0xc0, 0x25, 0xb4, 0x05, 0x96,
	0x7f, 0xe8, 0x8f, 0x5d, 0x63, 0xfb, 0x73, 0xf9, 0x2e, 0xa0, 0xdb, 0x1a, 0xc0, 0xb7, 0xf1, 0xfe,
	0xf8, 0x68, 0xca, 0x26, 0x43, 0xdf, 0x5d, 0xa1, 0x4d, 0x30, 0x47, 0x43, 0xe6, 0x12, 0xea, 0x80,
	0x3d, 0x65, 0xc3, 0x89, 0xef, 0x1a, 0xb4, 0x0d, 0xcd, 0xbd, 0xc9, 0xe8, 0xfb, 0xe8, 0xeb, 0xd8,
	0x35, 0xb7, 0x5f, 0x03, 0x2c, 0x86, 0xaf, 0xdc, 0x8f, 0xc6, 0x53, 0x1d, 0x77, 0x30, 0x9e, 0xba,
	0x64, 0xef, 0xe3, 0x8f, 0xf7, 0xe7, 0xa1, 0xbc, 0x28, 0x4e, 0xfa, 0xa7, 0x49, 0x34, 0xc8, 0xaf,
	0x62, 0x91, 0x89, 0xcb, 0x9b, 0xff, 0x31, 0x9f, 0xa5, 0x17, 0x7c, 0xc0, 0xd3, 0x70, 0x50, 0x7f,
	0xc6, 0x4e, 0x1a, 0xf8, 0x82, 0xed, 0xfe, 0x09, 0x00, 0x00, 0xff, 0xff, 0x0f, 0x99, 0x9b, 0x7d,
	0xe1, 0x06, 0x00, 0x00,
}

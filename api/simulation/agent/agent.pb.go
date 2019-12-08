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

type AgentType int32

const (
	AgentType_PEDESTRIAN AgentType = 0
	AgentType_CAR        AgentType = 1
	AgentType_TRAIN      AgentType = 2
	AgentType_SIGNAL     AgentType = 3
)

var AgentType_name = map[int32]string{
	0: "PEDESTRIAN",
	1: "CAR",
	2: "TRAIN",
	3: "SIGNAL",
}

var AgentType_value = map[string]int32{
	"PEDESTRIAN": 0,
	"CAR":        1,
	"TRAIN":      2,
	"SIGNAL":     3,
}

func (x AgentType) String() string {
	return proto.EnumName(AgentType_name, int32(x))
}

func (AgentType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_fce67ac898dc274e, []int{0}
}

type GetAgentsDemand struct {
	Time                 uint64    `protobuf:"varint,1,opt,name=time,proto3" json:"time,omitempty"`
	AreaId               uint64    `protobuf:"varint,2,opt,name=area_id,json=areaId,proto3" json:"area_id,omitempty"`
	AgentType            AgentType `protobuf:"varint,3,opt,name=agent_type,json=agentType,proto3,enum=api.agent.AgentType" json:"agent_type,omitempty"`
	XXX_NoUnkeyedLiteral struct{}  `json:"-"`
	XXX_unrecognized     []byte    `json:"-"`
	XXX_sizecache        int32     `json:"-"`
}

func (m *GetAgentsDemand) Reset()         { *m = GetAgentsDemand{} }
func (m *GetAgentsDemand) String() string { return proto.CompactTextString(m) }
func (*GetAgentsDemand) ProtoMessage()    {}
func (*GetAgentsDemand) Descriptor() ([]byte, []int) {
	return fileDescriptor_fce67ac898dc274e, []int{0}
}

func (m *GetAgentsDemand) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetAgentsDemand.Unmarshal(m, b)
}
func (m *GetAgentsDemand) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetAgentsDemand.Marshal(b, m, deterministic)
}
func (m *GetAgentsDemand) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetAgentsDemand.Merge(m, src)
}
func (m *GetAgentsDemand) XXX_Size() int {
	return xxx_messageInfo_GetAgentsDemand.Size(m)
}
func (m *GetAgentsDemand) XXX_DiscardUnknown() {
	xxx_messageInfo_GetAgentsDemand.DiscardUnknown(m)
}

var xxx_messageInfo_GetAgentsDemand proto.InternalMessageInfo

func (m *GetAgentsDemand) GetTime() uint64 {
	if m != nil {
		return m.Time
	}
	return 0
}

func (m *GetAgentsDemand) GetAreaId() uint64 {
	if m != nil {
		return m.AreaId
	}
	return 0
}

func (m *GetAgentsDemand) GetAgentType() AgentType {
	if m != nil {
		return m.AgentType
	}
	return AgentType_PEDESTRIAN
}

type GetAgentsSupply struct {
	Time                 uint64    `protobuf:"varint,1,opt,name=time,proto3" json:"time,omitempty"`
	AreaId               uint64    `protobuf:"varint,2,opt,name=area_id,json=areaId,proto3" json:"area_id,omitempty"`
	AgentType            AgentType `protobuf:"varint,3,opt,name=agent_type,json=agentType,proto3,enum=api.agent.AgentType" json:"agent_type,omitempty"`
	AgentsInfo           []*Agent  `protobuf:"bytes,4,rep,name=agents_info,json=agentsInfo,proto3" json:"agents_info,omitempty"`
	XXX_NoUnkeyedLiteral struct{}  `json:"-"`
	XXX_unrecognized     []byte    `json:"-"`
	XXX_sizecache        int32     `json:"-"`
}

func (m *GetAgentsSupply) Reset()         { *m = GetAgentsSupply{} }
func (m *GetAgentsSupply) String() string { return proto.CompactTextString(m) }
func (*GetAgentsSupply) ProtoMessage()    {}
func (*GetAgentsSupply) Descriptor() ([]byte, []int) {
	return fileDescriptor_fce67ac898dc274e, []int{1}
}

func (m *GetAgentsSupply) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetAgentsSupply.Unmarshal(m, b)
}
func (m *GetAgentsSupply) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetAgentsSupply.Marshal(b, m, deterministic)
}
func (m *GetAgentsSupply) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetAgentsSupply.Merge(m, src)
}
func (m *GetAgentsSupply) XXX_Size() int {
	return xxx_messageInfo_GetAgentsSupply.Size(m)
}
func (m *GetAgentsSupply) XXX_DiscardUnknown() {
	xxx_messageInfo_GetAgentsSupply.DiscardUnknown(m)
}

var xxx_messageInfo_GetAgentsSupply proto.InternalMessageInfo

func (m *GetAgentsSupply) GetTime() uint64 {
	if m != nil {
		return m.Time
	}
	return 0
}

func (m *GetAgentsSupply) GetAreaId() uint64 {
	if m != nil {
		return m.AreaId
	}
	return 0
}

func (m *GetAgentsSupply) GetAgentType() AgentType {
	if m != nil {
		return m.AgentType
	}
	return AgentType_PEDESTRIAN
}

func (m *GetAgentsSupply) GetAgentsInfo() []*Agent {
	if m != nil {
		return m.AgentsInfo
	}
	return nil
}

type SetAgentsDemand struct {
	AgentsInfo           []*Agent `protobuf:"bytes,1,rep,name=agents_info,json=agentsInfo,proto3" json:"agents_info,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SetAgentsDemand) Reset()         { *m = SetAgentsDemand{} }
func (m *SetAgentsDemand) String() string { return proto.CompactTextString(m) }
func (*SetAgentsDemand) ProtoMessage()    {}
func (*SetAgentsDemand) Descriptor() ([]byte, []int) {
	return fileDescriptor_fce67ac898dc274e, []int{2}
}

func (m *SetAgentsDemand) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SetAgentsDemand.Unmarshal(m, b)
}
func (m *SetAgentsDemand) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SetAgentsDemand.Marshal(b, m, deterministic)
}
func (m *SetAgentsDemand) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SetAgentsDemand.Merge(m, src)
}
func (m *SetAgentsDemand) XXX_Size() int {
	return xxx_messageInfo_SetAgentsDemand.Size(m)
}
func (m *SetAgentsDemand) XXX_DiscardUnknown() {
	xxx_messageInfo_SetAgentsDemand.DiscardUnknown(m)
}

var xxx_messageInfo_SetAgentsDemand proto.InternalMessageInfo

func (m *SetAgentsDemand) GetAgentsInfo() []*Agent {
	if m != nil {
		return m.AgentsInfo
	}
	return nil
}

type SetAgentsSupply struct {
	Time                 uint64    `protobuf:"varint,1,opt,name=time,proto3" json:"time,omitempty"`
	AreaId               uint64    `protobuf:"varint,2,opt,name=area_id,json=areaId,proto3" json:"area_id,omitempty"`
	AgentType            AgentType `protobuf:"varint,3,opt,name=agent_type,json=agentType,proto3,enum=api.agent.AgentType" json:"agent_type,omitempty"`
	XXX_NoUnkeyedLiteral struct{}  `json:"-"`
	XXX_unrecognized     []byte    `json:"-"`
	XXX_sizecache        int32     `json:"-"`
}

func (m *SetAgentsSupply) Reset()         { *m = SetAgentsSupply{} }
func (m *SetAgentsSupply) String() string { return proto.CompactTextString(m) }
func (*SetAgentsSupply) ProtoMessage()    {}
func (*SetAgentsSupply) Descriptor() ([]byte, []int) {
	return fileDescriptor_fce67ac898dc274e, []int{3}
}

func (m *SetAgentsSupply) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SetAgentsSupply.Unmarshal(m, b)
}
func (m *SetAgentsSupply) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SetAgentsSupply.Marshal(b, m, deterministic)
}
func (m *SetAgentsSupply) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SetAgentsSupply.Merge(m, src)
}
func (m *SetAgentsSupply) XXX_Size() int {
	return xxx_messageInfo_SetAgentsSupply.Size(m)
}
func (m *SetAgentsSupply) XXX_DiscardUnknown() {
	xxx_messageInfo_SetAgentsSupply.DiscardUnknown(m)
}

var xxx_messageInfo_SetAgentsSupply proto.InternalMessageInfo

func (m *SetAgentsSupply) GetTime() uint64 {
	if m != nil {
		return m.Time
	}
	return 0
}

func (m *SetAgentsSupply) GetAreaId() uint64 {
	if m != nil {
		return m.AreaId
	}
	return 0
}

func (m *SetAgentsSupply) GetAgentType() AgentType {
	if m != nil {
		return m.AgentType
	}
	return AgentType_PEDESTRIAN
}

type ForwardAgentsSupply struct {
	Time                 uint64    `protobuf:"varint,1,opt,name=time,proto3" json:"time,omitempty"`
	AreaId               uint64    `protobuf:"varint,2,opt,name=area_id,json=areaId,proto3" json:"area_id,omitempty"`
	AgentType            AgentType `protobuf:"varint,3,opt,name=agent_type,json=agentType,proto3,enum=api.agent.AgentType" json:"agent_type,omitempty"`
	AgentsInfo           []*Agent  `protobuf:"bytes,4,rep,name=agents_info,json=agentsInfo,proto3" json:"agents_info,omitempty"`
	XXX_NoUnkeyedLiteral struct{}  `json:"-"`
	XXX_unrecognized     []byte    `json:"-"`
	XXX_sizecache        int32     `json:"-"`
}

func (m *ForwardAgentsSupply) Reset()         { *m = ForwardAgentsSupply{} }
func (m *ForwardAgentsSupply) String() string { return proto.CompactTextString(m) }
func (*ForwardAgentsSupply) ProtoMessage()    {}
func (*ForwardAgentsSupply) Descriptor() ([]byte, []int) {
	return fileDescriptor_fce67ac898dc274e, []int{4}
}

func (m *ForwardAgentsSupply) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ForwardAgentsSupply.Unmarshal(m, b)
}
func (m *ForwardAgentsSupply) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ForwardAgentsSupply.Marshal(b, m, deterministic)
}
func (m *ForwardAgentsSupply) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ForwardAgentsSupply.Merge(m, src)
}
func (m *ForwardAgentsSupply) XXX_Size() int {
	return xxx_messageInfo_ForwardAgentsSupply.Size(m)
}
func (m *ForwardAgentsSupply) XXX_DiscardUnknown() {
	xxx_messageInfo_ForwardAgentsSupply.DiscardUnknown(m)
}

var xxx_messageInfo_ForwardAgentsSupply proto.InternalMessageInfo

func (m *ForwardAgentsSupply) GetTime() uint64 {
	if m != nil {
		return m.Time
	}
	return 0
}

func (m *ForwardAgentsSupply) GetAreaId() uint64 {
	if m != nil {
		return m.AreaId
	}
	return 0
}

func (m *ForwardAgentsSupply) GetAgentType() AgentType {
	if m != nil {
		return m.AgentType
	}
	return AgentType_PEDESTRIAN
}

func (m *ForwardAgentsSupply) GetAgentsInfo() []*Agent {
	if m != nil {
		return m.AgentsInfo
	}
	return nil
}

type GetAgentRouteDemand struct {
	AgentsInfo           *Agent   `protobuf:"bytes,1,opt,name=agents_info,json=agentsInfo,proto3" json:"agents_info,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetAgentRouteDemand) Reset()         { *m = GetAgentRouteDemand{} }
func (m *GetAgentRouteDemand) String() string { return proto.CompactTextString(m) }
func (*GetAgentRouteDemand) ProtoMessage()    {}
func (*GetAgentRouteDemand) Descriptor() ([]byte, []int) {
	return fileDescriptor_fce67ac898dc274e, []int{5}
}

func (m *GetAgentRouteDemand) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetAgentRouteDemand.Unmarshal(m, b)
}
func (m *GetAgentRouteDemand) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetAgentRouteDemand.Marshal(b, m, deterministic)
}
func (m *GetAgentRouteDemand) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetAgentRouteDemand.Merge(m, src)
}
func (m *GetAgentRouteDemand) XXX_Size() int {
	return xxx_messageInfo_GetAgentRouteDemand.Size(m)
}
func (m *GetAgentRouteDemand) XXX_DiscardUnknown() {
	xxx_messageInfo_GetAgentRouteDemand.DiscardUnknown(m)
}

var xxx_messageInfo_GetAgentRouteDemand proto.InternalMessageInfo

func (m *GetAgentRouteDemand) GetAgentsInfo() *Agent {
	if m != nil {
		return m.AgentsInfo
	}
	return nil
}

type GetAgentRouteSupply struct {
	AgentsInfo           *Agent   `protobuf:"bytes,1,opt,name=agents_info,json=agentsInfo,proto3" json:"agents_info,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetAgentRouteSupply) Reset()         { *m = GetAgentRouteSupply{} }
func (m *GetAgentRouteSupply) String() string { return proto.CompactTextString(m) }
func (*GetAgentRouteSupply) ProtoMessage()    {}
func (*GetAgentRouteSupply) Descriptor() ([]byte, []int) {
	return fileDescriptor_fce67ac898dc274e, []int{6}
}

func (m *GetAgentRouteSupply) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetAgentRouteSupply.Unmarshal(m, b)
}
func (m *GetAgentRouteSupply) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetAgentRouteSupply.Marshal(b, m, deterministic)
}
func (m *GetAgentRouteSupply) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetAgentRouteSupply.Merge(m, src)
}
func (m *GetAgentRouteSupply) XXX_Size() int {
	return xxx_messageInfo_GetAgentRouteSupply.Size(m)
}
func (m *GetAgentRouteSupply) XXX_DiscardUnknown() {
	xxx_messageInfo_GetAgentRouteSupply.DiscardUnknown(m)
}

var xxx_messageInfo_GetAgentRouteSupply proto.InternalMessageInfo

func (m *GetAgentRouteSupply) GetAgentsInfo() *Agent {
	if m != nil {
		return m.AgentsInfo
	}
	return nil
}

type GetAgentsRouteDemand struct {
	AgentsInfo           []*Agent `protobuf:"bytes,1,rep,name=agents_info,json=agentsInfo,proto3" json:"agents_info,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetAgentsRouteDemand) Reset()         { *m = GetAgentsRouteDemand{} }
func (m *GetAgentsRouteDemand) String() string { return proto.CompactTextString(m) }
func (*GetAgentsRouteDemand) ProtoMessage()    {}
func (*GetAgentsRouteDemand) Descriptor() ([]byte, []int) {
	return fileDescriptor_fce67ac898dc274e, []int{7}
}

func (m *GetAgentsRouteDemand) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetAgentsRouteDemand.Unmarshal(m, b)
}
func (m *GetAgentsRouteDemand) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetAgentsRouteDemand.Marshal(b, m, deterministic)
}
func (m *GetAgentsRouteDemand) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetAgentsRouteDemand.Merge(m, src)
}
func (m *GetAgentsRouteDemand) XXX_Size() int {
	return xxx_messageInfo_GetAgentsRouteDemand.Size(m)
}
func (m *GetAgentsRouteDemand) XXX_DiscardUnknown() {
	xxx_messageInfo_GetAgentsRouteDemand.DiscardUnknown(m)
}

var xxx_messageInfo_GetAgentsRouteDemand proto.InternalMessageInfo

func (m *GetAgentsRouteDemand) GetAgentsInfo() []*Agent {
	if m != nil {
		return m.AgentsInfo
	}
	return nil
}

type GetAgentsRouteSupply struct {
	AgentsInfo           []*Agent `protobuf:"bytes,1,rep,name=agents_info,json=agentsInfo,proto3" json:"agents_info,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetAgentsRouteSupply) Reset()         { *m = GetAgentsRouteSupply{} }
func (m *GetAgentsRouteSupply) String() string { return proto.CompactTextString(m) }
func (*GetAgentsRouteSupply) ProtoMessage()    {}
func (*GetAgentsRouteSupply) Descriptor() ([]byte, []int) {
	return fileDescriptor_fce67ac898dc274e, []int{8}
}

func (m *GetAgentsRouteSupply) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetAgentsRouteSupply.Unmarshal(m, b)
}
func (m *GetAgentsRouteSupply) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetAgentsRouteSupply.Marshal(b, m, deterministic)
}
func (m *GetAgentsRouteSupply) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetAgentsRouteSupply.Merge(m, src)
}
func (m *GetAgentsRouteSupply) XXX_Size() int {
	return xxx_messageInfo_GetAgentsRouteSupply.Size(m)
}
func (m *GetAgentsRouteSupply) XXX_DiscardUnknown() {
	xxx_messageInfo_GetAgentsRouteSupply.DiscardUnknown(m)
}

var xxx_messageInfo_GetAgentsRouteSupply proto.InternalMessageInfo

func (m *GetAgentsRouteSupply) GetAgentsInfo() []*Agent {
	if m != nil {
		return m.AgentsInfo
	}
	return nil
}

type Agent struct {
	Id   uint64    `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Type AgentType `protobuf:"varint,2,opt,name=type,proto3,enum=api.agent.AgentType" json:"type,omitempty"`
	// Types that are valid to be assigned to Data:
	//	*Agent_Pedestrian
	//	*Agent_Car
	//	*Agent_Train
	//	*Agent_Signal
	Data                 isAgent_Data `protobuf_oneof:"data"`
	XXX_NoUnkeyedLiteral struct{}     `json:"-"`
	XXX_unrecognized     []byte       `json:"-"`
	XXX_sizecache        int32        `json:"-"`
}

func (m *Agent) Reset()         { *m = Agent{} }
func (m *Agent) String() string { return proto.CompactTextString(m) }
func (*Agent) ProtoMessage()    {}
func (*Agent) Descriptor() ([]byte, []int) {
	return fileDescriptor_fce67ac898dc274e, []int{9}
}

func (m *Agent) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Agent.Unmarshal(m, b)
}
func (m *Agent) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Agent.Marshal(b, m, deterministic)
}
func (m *Agent) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Agent.Merge(m, src)
}
func (m *Agent) XXX_Size() int {
	return xxx_messageInfo_Agent.Size(m)
}
func (m *Agent) XXX_DiscardUnknown() {
	xxx_messageInfo_Agent.DiscardUnknown(m)
}

var xxx_messageInfo_Agent proto.InternalMessageInfo

func (m *Agent) GetId() uint64 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *Agent) GetType() AgentType {
	if m != nil {
		return m.Type
	}
	return AgentType_PEDESTRIAN
}

type isAgent_Data interface {
	isAgent_Data()
}

type Agent_Pedestrian struct {
	Pedestrian *Pedestrian `protobuf:"bytes,3,opt,name=pedestrian,proto3,oneof"`
}

type Agent_Car struct {
	Car *Car `protobuf:"bytes,4,opt,name=car,proto3,oneof"`
}

type Agent_Train struct {
	Train *Train `protobuf:"bytes,5,opt,name=train,proto3,oneof"`
}

type Agent_Signal struct {
	Signal *Signal `protobuf:"bytes,6,opt,name=signal,proto3,oneof"`
}

func (*Agent_Pedestrian) isAgent_Data() {}

func (*Agent_Car) isAgent_Data() {}

func (*Agent_Train) isAgent_Data() {}

func (*Agent_Signal) isAgent_Data() {}

func (m *Agent) GetData() isAgent_Data {
	if m != nil {
		return m.Data
	}
	return nil
}

func (m *Agent) GetPedestrian() *Pedestrian {
	if x, ok := m.GetData().(*Agent_Pedestrian); ok {
		return x.Pedestrian
	}
	return nil
}

func (m *Agent) GetCar() *Car {
	if x, ok := m.GetData().(*Agent_Car); ok {
		return x.Car
	}
	return nil
}

func (m *Agent) GetTrain() *Train {
	if x, ok := m.GetData().(*Agent_Train); ok {
		return x.Train
	}
	return nil
}

func (m *Agent) GetSignal() *Signal {
	if x, ok := m.GetData().(*Agent_Signal); ok {
		return x.Signal
	}
	return nil
}

// XXX_OneofWrappers is for the internal use of the proto package.
func (*Agent) XXX_OneofWrappers() []interface{} {
	return []interface{}{
		(*Agent_Pedestrian)(nil),
		(*Agent_Car)(nil),
		(*Agent_Train)(nil),
		(*Agent_Signal)(nil),
	}
}

func init() {
	proto.RegisterEnum("api.agent.AgentType", AgentType_name, AgentType_value)
	proto.RegisterType((*GetAgentsDemand)(nil), "api.agent.GetAgentsDemand")
	proto.RegisterType((*GetAgentsSupply)(nil), "api.agent.GetAgentsSupply")
	proto.RegisterType((*SetAgentsDemand)(nil), "api.agent.SetAgentsDemand")
	proto.RegisterType((*SetAgentsSupply)(nil), "api.agent.SetAgentsSupply")
	proto.RegisterType((*ForwardAgentsSupply)(nil), "api.agent.ForwardAgentsSupply")
	proto.RegisterType((*GetAgentRouteDemand)(nil), "api.agent.GetAgentRouteDemand")
	proto.RegisterType((*GetAgentRouteSupply)(nil), "api.agent.GetAgentRouteSupply")
	proto.RegisterType((*GetAgentsRouteDemand)(nil), "api.agent.GetAgentsRouteDemand")
	proto.RegisterType((*GetAgentsRouteSupply)(nil), "api.agent.GetAgentsRouteSupply")
	proto.RegisterType((*Agent)(nil), "api.agent.Agent")
}

func init() { proto.RegisterFile("simulation/agent/agent.proto", fileDescriptor_fce67ac898dc274e) }

var fileDescriptor_fce67ac898dc274e = []byte{
	// 502 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xcc, 0x54, 0x4d, 0x6b, 0xdb, 0x40,
	0x10, 0xd5, 0x97, 0x15, 0x3c, 0x2e, 0x8e, 0xd8, 0x04, 0x2a, 0x4c, 0x0b, 0xae, 0x4f, 0xa2, 0x14,
	0x99, 0x3a, 0x94, 0x1e, 0xda, 0x8b, 0x12, 0xa7, 0x91, 0xa0, 0x98, 0xb0, 0xf2, 0xa9, 0x17, 0x33,
	0xb1, 0x36, 0xce, 0x82, 0x2d, 0x89, 0xd5, 0x9a, 0xd6, 0xc7, 0xfe, 0x99, 0x1e, 0xfa, 0x2b, 0x8b,
	0x56, 0xb2, 0x2b, 0x6c, 0x92, 0x36, 0x85, 0x96, 0x5e, 0xf6, 0x63, 0xe6, 0xcd, 0x9b, 0x9d, 0x99,
	0xe5, 0xc1, 0xb3, 0x82, 0xaf, 0xd6, 0x4b, 0x94, 0x3c, 0x4b, 0x87, 0xb8, 0x60, 0xa9, 0xac, 0x56,
	0x3f, 0x17, 0x99, 0xcc, 0x48, 0x1b, 0x73, 0xee, 0x2b, 0x43, 0xef, 0xc5, 0x01, 0x30, 0x67, 0x09,
	0x2b, 0xa4, 0xe0, 0x98, 0x56, 0xe8, 0x5e, 0xef, 0x00, 0x32, 0x47, 0x51, 0xfb, 0x0e, 0xf3, 0x48,
	0x81, 0x7c, 0x1b, 0xf9, 0xfc, 0xc0, 0x5b, 0xf0, 0x45, 0x8a, 0xcb, 0xca, 0x3d, 0x28, 0xe0, 0xf8,
	0x8a, 0xc9, 0xa0, 0x74, 0x14, 0x63, 0xb6, 0xc2, 0x34, 0x21, 0x04, 0x2c, 0xc9, 0x57, 0xcc, 0xd5,
	0xfb, 0xba, 0x67, 0x51, 0x75, 0x26, 0x4f, 0xe1, 0x08, 0x05, 0xc3, 0x19, 0x4f, 0x5c, 0x43, 0x99,
	0xed, 0xf2, 0x1a, 0x25, 0xe4, 0x0c, 0x40, 0xb1, 0xce, 0xe4, 0x26, 0x67, 0xae, 0xd9, 0xd7, 0xbd,
	0xee, 0xe8, 0xd4, 0xdf, 0xd5, 0xe6, 0x2b, 0xe6, 0xe9, 0x26, 0x67, 0xb4, 0x8d, 0xdb, 0xe3, 0xe0,
	0x9b, 0xde, 0xc8, 0x1a, 0xaf, 0xf3, 0x7c, 0xb9, 0xf9, 0xfb, 0x59, 0xc9, 0x6b, 0xe8, 0xa8, 0x4b,
	0x31, 0xe3, 0xe9, 0x6d, 0xe6, 0x5a, 0x7d, 0xd3, 0xeb, 0x8c, 0x9c, 0xfd, 0x28, 0x5a, 0x31, 0x17,
	0x51, 0x7a, 0x9b, 0x0d, 0xc6, 0x70, 0x1c, 0xef, 0x75, 0x67, 0x8f, 0x45, 0xff, 0x0d, 0x96, 0xa2,
	0xc1, 0xf2, 0xaf, 0xaa, 0x1d, 0x7c, 0xd7, 0xe1, 0xe4, 0x43, 0x26, 0x3e, 0xa3, 0x48, 0xfe, 0xff,
	0x3e, 0x87, 0x70, 0xb2, 0xfd, 0x0f, 0x34, 0x5b, 0x4b, 0x76, 0x5f, 0xaf, 0xf5, 0x47, 0x33, 0xd5,
	0x55, 0xff, 0x01, 0x53, 0x04, 0xa7, 0xbb, 0x3f, 0xfa, 0xe0, 0xa3, 0xcc, 0xc7, 0x53, 0xdd, 0xf7,
	0xaa, 0x5f, 0x53, 0x7d, 0x35, 0xa0, 0xa5, 0xac, 0xa4, 0x0b, 0x06, 0x4f, 0xea, 0x31, 0x1a, 0x3c,
	0x21, 0x1e, 0x58, 0x6a, 0x4a, 0xc6, 0x03, 0x53, 0x52, 0x08, 0xf2, 0x1e, 0xe0, 0xa7, 0xc0, 0xa8,
	0xa9, 0x76, 0x46, 0x3d, 0x85, 0x6f, 0xe8, 0xce, 0xf5, 0xee, 0x18, 0x6a, 0xb4, 0x81, 0x27, 0x7d,
	0x30, 0xe7, 0x28, 0x5c, 0x4b, 0x85, 0x3d, 0x51, 0x61, 0xa5, 0x16, 0x5d, 0xa0, 0x08, 0x35, 0x5a,
	0xba, 0x88, 0x07, 0x2d, 0xa5, 0x40, 0x6e, 0xab, 0xd1, 0xe6, 0x4a, 0x93, 0xa6, 0xe5, 0x1a, 0x6a,
	0xb4, 0x02, 0x90, 0x57, 0x60, 0x57, 0x6a, 0xe4, 0xda, 0x0a, 0x4a, 0x14, 0xb4, 0x16, 0xa8, 0x58,
	0x6d, 0xa1, 0x46, 0x6b, 0xcc, 0xb9, 0x0d, 0x56, 0x82, 0x12, 0x5f, 0xbe, 0x83, 0xf6, 0xae, 0x24,
	0xd2, 0x05, 0xb8, 0xbe, 0x1c, 0x5f, 0xc6, 0x53, 0x1a, 0x05, 0x13, 0x47, 0x23, 0x47, 0x60, 0x5e,
	0x04, 0xd4, 0xd1, 0x49, 0x1b, 0x5a, 0x53, 0x1a, 0x44, 0x13, 0xc7, 0x20, 0x00, 0x76, 0x1c, 0x5d,
	0x4d, 0x82, 0x8f, 0x8e, 0x79, 0xfe, 0xf6, 0xd3, 0x9b, 0x05, 0x97, 0x77, 0xeb, 0x1b, 0x7f, 0x9e,
	0xad, 0x86, 0xc5, 0x26, 0x65, 0x82, 0x7d, 0xd9, 0xee, 0x33, 0x5c, 0xe6, 0x77, 0x38, 0xc4, 0x9c,
	0x0f, 0xf7, 0x65, 0xf3, 0xc6, 0x56, 0x82, 0x79, 0xf6, 0x23, 0x00, 0x00, 0xff, 0xff, 0x91, 0x09,
	0x62, 0x4e, 0xd7, 0x05, 0x00, 0x00,
}

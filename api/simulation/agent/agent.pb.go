// Code generated by protoc-gen-go. DO NOT EDIT.
// source: simulation/agent/agent.proto

package agent

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	common "github.com/synerex/synerex_alpha/api/simulation/common"
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

type GetSameAreaAgentsRequest struct {
	AreaId               uint64           `protobuf:"varint,1,opt,name=area_id,json=areaId,proto3" json:"area_id,omitempty"`
	AgentType            common.AgentType `protobuf:"varint,2,opt,name=agent_type,json=agentType,proto3,enum=api.common.AgentType" json:"agent_type,omitempty"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_unrecognized     []byte           `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *GetSameAreaAgentsRequest) Reset()         { *m = GetSameAreaAgentsRequest{} }
func (m *GetSameAreaAgentsRequest) String() string { return proto.CompactTextString(m) }
func (*GetSameAreaAgentsRequest) ProtoMessage()    {}
func (*GetSameAreaAgentsRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_fce67ac898dc274e, []int{0}
}

func (m *GetSameAreaAgentsRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetSameAreaAgentsRequest.Unmarshal(m, b)
}
func (m *GetSameAreaAgentsRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetSameAreaAgentsRequest.Marshal(b, m, deterministic)
}
func (m *GetSameAreaAgentsRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetSameAreaAgentsRequest.Merge(m, src)
}
func (m *GetSameAreaAgentsRequest) XXX_Size() int {
	return xxx_messageInfo_GetSameAreaAgentsRequest.Size(m)
}
func (m *GetSameAreaAgentsRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_GetSameAreaAgentsRequest.DiscardUnknown(m)
}

var xxx_messageInfo_GetSameAreaAgentsRequest proto.InternalMessageInfo

func (m *GetSameAreaAgentsRequest) GetAreaId() uint64 {
	if m != nil {
		return m.AreaId
	}
	return 0
}

func (m *GetSameAreaAgentsRequest) GetAgentType() common.AgentType {
	if m != nil {
		return m.AgentType
	}
	return common.AgentType_NONE
}

type GetSameAreaAgentsResponse struct {
	Agents               []*Agent `protobuf:"bytes,1,rep,name=agents,proto3" json:"agents,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetSameAreaAgentsResponse) Reset()         { *m = GetSameAreaAgentsResponse{} }
func (m *GetSameAreaAgentsResponse) String() string { return proto.CompactTextString(m) }
func (*GetSameAreaAgentsResponse) ProtoMessage()    {}
func (*GetSameAreaAgentsResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_fce67ac898dc274e, []int{1}
}

func (m *GetSameAreaAgentsResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetSameAreaAgentsResponse.Unmarshal(m, b)
}
func (m *GetSameAreaAgentsResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetSameAreaAgentsResponse.Marshal(b, m, deterministic)
}
func (m *GetSameAreaAgentsResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetSameAreaAgentsResponse.Merge(m, src)
}
func (m *GetSameAreaAgentsResponse) XXX_Size() int {
	return xxx_messageInfo_GetSameAreaAgentsResponse.Size(m)
}
func (m *GetSameAreaAgentsResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_GetSameAreaAgentsResponse.DiscardUnknown(m)
}

var xxx_messageInfo_GetSameAreaAgentsResponse proto.InternalMessageInfo

func (m *GetSameAreaAgentsResponse) GetAgents() []*Agent {
	if m != nil {
		return m.Agents
	}
	return nil
}

type GetNeighborAreaAgentsResponse struct {
	Agents               []*Agent `protobuf:"bytes,1,rep,name=agents,proto3" json:"agents,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetNeighborAreaAgentsResponse) Reset()         { *m = GetNeighborAreaAgentsResponse{} }
func (m *GetNeighborAreaAgentsResponse) String() string { return proto.CompactTextString(m) }
func (*GetNeighborAreaAgentsResponse) ProtoMessage()    {}
func (*GetNeighborAreaAgentsResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_fce67ac898dc274e, []int{2}
}

func (m *GetNeighborAreaAgentsResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetNeighborAreaAgentsResponse.Unmarshal(m, b)
}
func (m *GetNeighborAreaAgentsResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetNeighborAreaAgentsResponse.Marshal(b, m, deterministic)
}
func (m *GetNeighborAreaAgentsResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetNeighborAreaAgentsResponse.Merge(m, src)
}
func (m *GetNeighborAreaAgentsResponse) XXX_Size() int {
	return xxx_messageInfo_GetNeighborAreaAgentsResponse.Size(m)
}
func (m *GetNeighborAreaAgentsResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_GetNeighborAreaAgentsResponse.DiscardUnknown(m)
}

var xxx_messageInfo_GetNeighborAreaAgentsResponse proto.InternalMessageInfo

func (m *GetNeighborAreaAgentsResponse) GetAgents() []*Agent {
	if m != nil {
		return m.Agents
	}
	return nil
}

type SetAgentsRequest struct {
	Agents               []*Agent `protobuf:"bytes,1,rep,name=agents,proto3" json:"agents,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SetAgentsRequest) Reset()         { *m = SetAgentsRequest{} }
func (m *SetAgentsRequest) String() string { return proto.CompactTextString(m) }
func (*SetAgentsRequest) ProtoMessage()    {}
func (*SetAgentsRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_fce67ac898dc274e, []int{3}
}

func (m *SetAgentsRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SetAgentsRequest.Unmarshal(m, b)
}
func (m *SetAgentsRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SetAgentsRequest.Marshal(b, m, deterministic)
}
func (m *SetAgentsRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SetAgentsRequest.Merge(m, src)
}
func (m *SetAgentsRequest) XXX_Size() int {
	return xxx_messageInfo_SetAgentsRequest.Size(m)
}
func (m *SetAgentsRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_SetAgentsRequest.DiscardUnknown(m)
}

var xxx_messageInfo_SetAgentsRequest proto.InternalMessageInfo

func (m *SetAgentsRequest) GetAgents() []*Agent {
	if m != nil {
		return m.Agents
	}
	return nil
}

type SetAgentsResponse struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SetAgentsResponse) Reset()         { *m = SetAgentsResponse{} }
func (m *SetAgentsResponse) String() string { return proto.CompactTextString(m) }
func (*SetAgentsResponse) ProtoMessage()    {}
func (*SetAgentsResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_fce67ac898dc274e, []int{4}
}

func (m *SetAgentsResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SetAgentsResponse.Unmarshal(m, b)
}
func (m *SetAgentsResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SetAgentsResponse.Marshal(b, m, deterministic)
}
func (m *SetAgentsResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SetAgentsResponse.Merge(m, src)
}
func (m *SetAgentsResponse) XXX_Size() int {
	return xxx_messageInfo_SetAgentsResponse.Size(m)
}
func (m *SetAgentsResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_SetAgentsResponse.DiscardUnknown(m)
}

var xxx_messageInfo_SetAgentsResponse proto.InternalMessageInfo

type ClearAgentsRequest struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ClearAgentsRequest) Reset()         { *m = ClearAgentsRequest{} }
func (m *ClearAgentsRequest) String() string { return proto.CompactTextString(m) }
func (*ClearAgentsRequest) ProtoMessage()    {}
func (*ClearAgentsRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_fce67ac898dc274e, []int{5}
}

func (m *ClearAgentsRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ClearAgentsRequest.Unmarshal(m, b)
}
func (m *ClearAgentsRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ClearAgentsRequest.Marshal(b, m, deterministic)
}
func (m *ClearAgentsRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ClearAgentsRequest.Merge(m, src)
}
func (m *ClearAgentsRequest) XXX_Size() int {
	return xxx_messageInfo_ClearAgentsRequest.Size(m)
}
func (m *ClearAgentsRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_ClearAgentsRequest.DiscardUnknown(m)
}

var xxx_messageInfo_ClearAgentsRequest proto.InternalMessageInfo

type ClearAgentsResponse struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ClearAgentsResponse) Reset()         { *m = ClearAgentsResponse{} }
func (m *ClearAgentsResponse) String() string { return proto.CompactTextString(m) }
func (*ClearAgentsResponse) ProtoMessage()    {}
func (*ClearAgentsResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_fce67ac898dc274e, []int{6}
}

func (m *ClearAgentsResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ClearAgentsResponse.Unmarshal(m, b)
}
func (m *ClearAgentsResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ClearAgentsResponse.Marshal(b, m, deterministic)
}
func (m *ClearAgentsResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ClearAgentsResponse.Merge(m, src)
}
func (m *ClearAgentsResponse) XXX_Size() int {
	return xxx_messageInfo_ClearAgentsResponse.Size(m)
}
func (m *ClearAgentsResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_ClearAgentsResponse.DiscardUnknown(m)
}

var xxx_messageInfo_ClearAgentsResponse proto.InternalMessageInfo

type VisualizeAgentsResponse struct {
	AreaId               uint64           `protobuf:"varint,1,opt,name=area_id,json=areaId,proto3" json:"area_id,omitempty"`
	AgentType            common.AgentType `protobuf:"varint,2,opt,name=agent_type,json=agentType,proto3,enum=api.common.AgentType" json:"agent_type,omitempty"`
	Agents               []*Agent         `protobuf:"bytes,3,rep,name=agents,proto3" json:"agents,omitempty"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_unrecognized     []byte           `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *VisualizeAgentsResponse) Reset()         { *m = VisualizeAgentsResponse{} }
func (m *VisualizeAgentsResponse) String() string { return proto.CompactTextString(m) }
func (*VisualizeAgentsResponse) ProtoMessage()    {}
func (*VisualizeAgentsResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_fce67ac898dc274e, []int{7}
}

func (m *VisualizeAgentsResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_VisualizeAgentsResponse.Unmarshal(m, b)
}
func (m *VisualizeAgentsResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_VisualizeAgentsResponse.Marshal(b, m, deterministic)
}
func (m *VisualizeAgentsResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_VisualizeAgentsResponse.Merge(m, src)
}
func (m *VisualizeAgentsResponse) XXX_Size() int {
	return xxx_messageInfo_VisualizeAgentsResponse.Size(m)
}
func (m *VisualizeAgentsResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_VisualizeAgentsResponse.DiscardUnknown(m)
}

var xxx_messageInfo_VisualizeAgentsResponse proto.InternalMessageInfo

func (m *VisualizeAgentsResponse) GetAreaId() uint64 {
	if m != nil {
		return m.AreaId
	}
	return 0
}

func (m *VisualizeAgentsResponse) GetAgentType() common.AgentType {
	if m != nil {
		return m.AgentType
	}
	return common.AgentType_NONE
}

func (m *VisualizeAgentsResponse) GetAgents() []*Agent {
	if m != nil {
		return m.Agents
	}
	return nil
}

type Agent struct {
	Id   uint64           `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Type common.AgentType `protobuf:"varint,2,opt,name=type,proto3,enum=api.common.AgentType" json:"type,omitempty"`
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
	return fileDescriptor_fce67ac898dc274e, []int{8}
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

func (m *Agent) GetType() common.AgentType {
	if m != nil {
		return m.Type
	}
	return common.AgentType_NONE
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
	proto.RegisterType((*GetSameAreaAgentsRequest)(nil), "api.agent.GetSameAreaAgentsRequest")
	proto.RegisterType((*GetSameAreaAgentsResponse)(nil), "api.agent.GetSameAreaAgentsResponse")
	proto.RegisterType((*GetNeighborAreaAgentsResponse)(nil), "api.agent.GetNeighborAreaAgentsResponse")
	proto.RegisterType((*SetAgentsRequest)(nil), "api.agent.SetAgentsRequest")
	proto.RegisterType((*SetAgentsResponse)(nil), "api.agent.SetAgentsResponse")
	proto.RegisterType((*ClearAgentsRequest)(nil), "api.agent.ClearAgentsRequest")
	proto.RegisterType((*ClearAgentsResponse)(nil), "api.agent.ClearAgentsResponse")
	proto.RegisterType((*VisualizeAgentsResponse)(nil), "api.agent.VisualizeAgentsResponse")
	proto.RegisterType((*Agent)(nil), "api.agent.Agent")
}

func init() { proto.RegisterFile("simulation/agent/agent.proto", fileDescriptor_fce67ac898dc274e) }

var fileDescriptor_fce67ac898dc274e = []byte{
	// 453 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xac, 0x93, 0xcf, 0x6b, 0x13, 0x41,
	0x14, 0xc7, 0xb3, 0xf9, 0xb1, 0xd2, 0x57, 0x29, 0x75, 0x6a, 0xe9, 0x18, 0xac, 0xac, 0x7b, 0x5a,
	0x41, 0x36, 0x10, 0x15, 0x2f, 0xbd, 0xb4, 0x45, 0x9a, 0x5e, 0x44, 0x26, 0xc5, 0x83, 0x97, 0xf0,
	0x92, 0x7d, 0x24, 0x03, 0x9b, 0xdd, 0x71, 0x66, 0x02, 0xc6, 0xa3, 0x7f, 0x82, 0x7f, 0xb1, 0xec,
	0xcc, 0x26, 0x6e, 0x1a, 0xa1, 0x82, 0x5e, 0x76, 0x66, 0xe7, 0xfb, 0x79, 0xef, 0xcb, 0x77, 0x98,
	0x07, 0xcf, 0x8d, 0x5c, 0xae, 0x72, 0xb4, 0xb2, 0x2c, 0x06, 0x38, 0xa7, 0xc2, 0xfa, 0x6f, 0xaa,
	0x74, 0x69, 0x4b, 0x76, 0x80, 0x4a, 0xa6, 0xee, 0xa0, 0xff, 0x72, 0x0f, 0x54, 0x94, 0x91, 0xb1,
	0x5a, 0x62, 0xe1, 0xe9, 0x7e, 0x7f, 0x0f, 0x99, 0xa1, 0xae, 0xb5, 0x7d, 0x1f, 0xab, 0x51, 0x6e,
	0x2a, 0xcf, 0xf7, 0x54, 0x23, 0xe7, 0x05, 0xe6, 0xb5, 0xfc, 0xa2, 0x21, 0xcf, 0xca, 0xe5, 0x72,
	0xbb, 0x78, 0x3d, 0x96, 0xc0, 0x6f, 0xc8, 0x8e, 0x71, 0x49, 0x97, 0x9a, 0xf0, 0xb2, 0x6a, 0x60,
	0x04, 0x7d, 0x5d, 0x91, 0xb1, 0xec, 0x0c, 0x1e, 0xa1, 0x26, 0x9c, 0xc8, 0x8c, 0x07, 0x51, 0x90,
	0x74, 0x45, 0x58, 0xfd, 0xde, 0x66, 0xec, 0x2d, 0x80, 0xb3, 0x9a, 0xd8, 0xb5, 0x22, 0xde, 0x8e,
	0x82, 0xe4, 0x68, 0x78, 0x9a, 0x56, 0x81, 0xeb, 0xde, 0xae, 0xcf, 0xdd, 0x5a, 0x91, 0x38, 0xc0,
	0xcd, 0x36, 0xfe, 0x00, 0xcf, 0xfe, 0x60, 0x65, 0x54, 0x59, 0x18, 0x62, 0x09, 0x84, 0x8e, 0x34,
	0x3c, 0x88, 0x3a, 0xc9, 0xe1, 0xf0, 0x38, 0xdd, 0xde, 0x9f, 0xef, 0x26, 0x6a, 0x3d, 0xbe, 0x85,
	0xf3, 0x1b, 0xb2, 0x1f, 0x49, 0xce, 0x17, 0xd3, 0x52, 0xff, 0x53, 0xab, 0x0b, 0x38, 0x1e, 0x93,
	0xdd, 0x0d, 0xfd, 0xf7, 0xd5, 0x27, 0xf0, 0xa4, 0x51, 0xed, 0xcd, 0xe3, 0xa7, 0xc0, 0xae, 0x73,
	0x42, 0xbd, 0xd3, 0x34, 0x3e, 0x85, 0x93, 0x9d, 0xd3, 0x1a, 0xfe, 0x19, 0xc0, 0xd9, 0x67, 0x69,
	0x56, 0x98, 0xcb, 0xef, 0x74, 0x2f, 0xc5, 0xff, 0xbd, 0xfc, 0x46, 0xac, 0xce, 0x03, 0xb1, 0x7e,
	0xb4, 0xa1, 0xe7, 0x4e, 0xd8, 0x11, 0xb4, 0xb7, 0xee, 0x6d, 0x99, 0xb1, 0x57, 0xd0, 0x7d, 0xd8,
	0xd3, 0x21, 0xec, 0x02, 0xe0, 0xf7, 0x1b, 0xe7, 0x9d, 0x28, 0x48, 0x0e, 0x87, 0x7d, 0x57, 0xd0,
	0x78, 0xfa, 0x9f, 0xb6, 0xdb, 0x51, 0x4b, 0x34, 0x78, 0x16, 0x41, 0x67, 0x86, 0x9a, 0x77, 0x5d,
	0xd9, 0x63, 0xef, 0x83, 0x3a, 0xbd, 0x46, 0x3d, 0x6a, 0x89, 0x4a, 0x62, 0x09, 0xf4, 0xdc, 0x10,
	0xf0, 0x9e, 0x63, 0x7c, 0x1a, 0x3f, 0x16, 0x77, 0xd5, 0x77, 0xd4, 0x12, 0x1e, 0x60, 0xaf, 0x21,
	0xf4, 0x03, 0xc1, 0x43, 0x87, 0x32, 0x87, 0xd6, 0x33, 0x32, 0x76, 0xcb, 0xa8, 0x25, 0x6a, 0xe6,
	0x2a, 0x84, 0x6e, 0x86, 0x16, 0xaf, 0xde, 0x7f, 0x79, 0x37, 0x97, 0x76, 0xb1, 0x9a, 0x56, 0x01,
	0x07, 0x66, 0x5d, 0x90, 0xa6, 0x6f, 0x9b, 0x75, 0x82, 0xb9, 0x5a, 0xe0, 0x00, 0x95, 0x1c, 0xdc,
	0x1f, 0xbe, 0x69, 0xe8, 0xc6, 0xea, 0xcd, 0xaf, 0x00, 0x00, 0x00, 0xff, 0xff, 0xb1, 0x36, 0xa2,
	0xf0, 0x1d, 0x04, 0x00, 0x00,
}

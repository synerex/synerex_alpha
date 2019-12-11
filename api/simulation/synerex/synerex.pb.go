// Code generated by protoc-gen-go. DO NOT EDIT.
// source: simulation/synerex/synerex.proto

package synerex

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	agent "github.com/synerex/synerex_alpha/api/simulation/agent"
	area "github.com/synerex/synerex_alpha/api/simulation/area"
	clock "github.com/synerex/synerex_alpha/api/simulation/clock"
	participant "github.com/synerex/synerex_alpha/api/simulation/participant"
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
	return fileDescriptor_81d77a6394be03b4, []int{0}
}

type DemandType int32

const (
	DemandType_GET_SAME_AREA_AGENTS_REQUEST DemandType = 0
	DemandType_SET_AGENTS_REQUEST           DemandType = 1
	DemandType_VISUALIZE_AGENTS_REQUEST     DemandType = 2
	DemandType_GET_AREA_REQUEST             DemandType = 3
	DemandType_RESIST_PARTICIPANT_REQUEST   DemandType = 4
	DemandType_SET_PARTICIPANTS_REQUEST     DemandType = 5
	DemandType_GET_CLOCK_REQUEST            DemandType = 6
	DemandType_SET_CLOCK_REQUEST            DemandType = 7
	DemandType_START_CLOCK_REQUEST          DemandType = 8
	DemandType_STOP_CLOCK_REQUEST           DemandType = 9
	DemandType_FORWARD_CLOCK_REQUEST        DemandType = 10
	DemandType_BACK_CLOCK_REQUEST           DemandType = 11
)

var DemandType_name = map[int32]string{
	0:  "GET_SAME_AREA_AGENTS_REQUEST",
	1:  "SET_AGENTS_REQUEST",
	2:  "VISUALIZE_AGENTS_REQUEST",
	3:  "GET_AREA_REQUEST",
	4:  "RESIST_PARTICIPANT_REQUEST",
	5:  "SET_PARTICIPANTS_REQUEST",
	6:  "GET_CLOCK_REQUEST",
	7:  "SET_CLOCK_REQUEST",
	8:  "START_CLOCK_REQUEST",
	9:  "STOP_CLOCK_REQUEST",
	10: "FORWARD_CLOCK_REQUEST",
	11: "BACK_CLOCK_REQUEST",
}

var DemandType_value = map[string]int32{
	"GET_SAME_AREA_AGENTS_REQUEST": 0,
	"SET_AGENTS_REQUEST":           1,
	"VISUALIZE_AGENTS_REQUEST":     2,
	"GET_AREA_REQUEST":             3,
	"RESIST_PARTICIPANT_REQUEST":   4,
	"SET_PARTICIPANTS_REQUEST":     5,
	"GET_CLOCK_REQUEST":            6,
	"SET_CLOCK_REQUEST":            7,
	"START_CLOCK_REQUEST":          8,
	"STOP_CLOCK_REQUEST":           9,
	"FORWARD_CLOCK_REQUEST":        10,
	"BACK_CLOCK_REQUEST":           11,
}

func (x DemandType) String() string {
	return proto.EnumName(DemandType_name, int32(x))
}

func (DemandType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_81d77a6394be03b4, []int{1}
}

type SupplyType int32

const (
	SupplyType_GET_SAME_AREA_AGENTS_RESPONSE     SupplyType = 0
	SupplyType_GET_NEIGHBOR_AREA_AGENTS_RESPONSE SupplyType = 1
	SupplyType_SET_AGENTS_RESPONSE               SupplyType = 2
	SupplyType_GET_AREA_RESPONSE                 SupplyType = 3
	SupplyType_RESIST_PARTICIPANT_RESPONSE       SupplyType = 4
	SupplyType_SET_PARTICIPANTS_RESPONSE         SupplyType = 5
	SupplyType_GET_CLOCK_RESPONSE                SupplyType = 6
	SupplyType_SET_CLOCK_RESPONSE                SupplyType = 7
	SupplyType_START_CLOCK_RESPONSE              SupplyType = 8
	SupplyType_STOP_CLOCK_RESPONSE               SupplyType = 9
	SupplyType_FORWARD_CLOCK_RESPONSE            SupplyType = 10
	SupplyType_BACK_CLOCK_RESPONSE               SupplyType = 11
)

var SupplyType_name = map[int32]string{
	0:  "GET_SAME_AREA_AGENTS_RESPONSE",
	1:  "GET_NEIGHBOR_AREA_AGENTS_RESPONSE",
	2:  "SET_AGENTS_RESPONSE",
	3:  "GET_AREA_RESPONSE",
	4:  "RESIST_PARTICIPANT_RESPONSE",
	5:  "SET_PARTICIPANTS_RESPONSE",
	6:  "GET_CLOCK_RESPONSE",
	7:  "SET_CLOCK_RESPONSE",
	8:  "START_CLOCK_RESPONSE",
	9:  "STOP_CLOCK_RESPONSE",
	10: "FORWARD_CLOCK_RESPONSE",
	11: "BACK_CLOCK_RESPONSE",
}

var SupplyType_value = map[string]int32{
	"GET_SAME_AREA_AGENTS_RESPONSE":     0,
	"GET_NEIGHBOR_AREA_AGENTS_RESPONSE": 1,
	"SET_AGENTS_RESPONSE":               2,
	"GET_AREA_RESPONSE":                 3,
	"RESIST_PARTICIPANT_RESPONSE":       4,
	"SET_PARTICIPANTS_RESPONSE":         5,
	"GET_CLOCK_RESPONSE":                6,
	"SET_CLOCK_RESPONSE":                7,
	"START_CLOCK_RESPONSE":              8,
	"STOP_CLOCK_RESPONSE":               9,
	"FORWARD_CLOCK_RESPONSE":            10,
	"BACK_CLOCK_RESPONSE":               11,
}

func (x SupplyType) String() string {
	return proto.EnumName(SupplyType_name, int32(x))
}

func (SupplyType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_81d77a6394be03b4, []int{2}
}

type SimDemand struct {
	// participant info
	DemandType DemandType `protobuf:"varint,1,opt,name=demand_type,json=demandType,proto3,enum=api.synerex.DemandType" json:"demand_type,omitempty"`
	// meta data
	StatusType StatusType `protobuf:"varint,2,opt,name=status_type,json=statusType,proto3,enum=api.synerex.StatusType" json:"status_type,omitempty"`
	Meta       string     `protobuf:"bytes,3,opt,name=meta,proto3" json:"meta,omitempty"`
	// demand data
	//
	// Types that are valid to be assigned to Data:
	//	*SimDemand_GetSameAreaAgentsRequest
	//	*SimDemand_SetAgentsRequest
	//	*SimDemand_VisualiseAgentsRequest
	//	*SimDemand_GetAreaRequest
	//	*SimDemand_ResistParticipantRequest
	//	*SimDemand_SetParticipantsRequest
	//	*SimDemand_GetClockRequest
	//	*SimDemand_SetClockRequest
	//	*SimDemand_StartClockRequest
	//	*SimDemand_StopClockRequest
	//	*SimDemand_ForwardClockRequest
	//	*SimDemand_BackClockRequest
	Data                 isSimDemand_Data `protobuf_oneof:"data"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_unrecognized     []byte           `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *SimDemand) Reset()         { *m = SimDemand{} }
func (m *SimDemand) String() string { return proto.CompactTextString(m) }
func (*SimDemand) ProtoMessage()    {}
func (*SimDemand) Descriptor() ([]byte, []int) {
	return fileDescriptor_81d77a6394be03b4, []int{0}
}

func (m *SimDemand) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SimDemand.Unmarshal(m, b)
}
func (m *SimDemand) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SimDemand.Marshal(b, m, deterministic)
}
func (m *SimDemand) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SimDemand.Merge(m, src)
}
func (m *SimDemand) XXX_Size() int {
	return xxx_messageInfo_SimDemand.Size(m)
}
func (m *SimDemand) XXX_DiscardUnknown() {
	xxx_messageInfo_SimDemand.DiscardUnknown(m)
}

var xxx_messageInfo_SimDemand proto.InternalMessageInfo

func (m *SimDemand) GetDemandType() DemandType {
	if m != nil {
		return m.DemandType
	}
	return DemandType_GET_SAME_AREA_AGENTS_REQUEST
}

func (m *SimDemand) GetStatusType() StatusType {
	if m != nil {
		return m.StatusType
	}
	return StatusType_OK
}

func (m *SimDemand) GetMeta() string {
	if m != nil {
		return m.Meta
	}
	return ""
}

type isSimDemand_Data interface {
	isSimDemand_Data()
}

type SimDemand_GetSameAreaAgentsRequest struct {
	GetSameAreaAgentsRequest *agent.GetSameAreaAgentsRequest `protobuf:"bytes,4,opt,name=get_same_area_agents_request,json=getSameAreaAgentsRequest,proto3,oneof"`
}

type SimDemand_SetAgentsRequest struct {
	SetAgentsRequest *agent.SetAgentsRequest `protobuf:"bytes,5,opt,name=set_agents_request,json=setAgentsRequest,proto3,oneof"`
}

type SimDemand_VisualiseAgentsRequest struct {
	VisualiseAgentsRequest *agent.VisualizeAgentsRequest `protobuf:"bytes,6,opt,name=visualise_agents_request,json=visualiseAgentsRequest,proto3,oneof"`
}

type SimDemand_GetAreaRequest struct {
	GetAreaRequest *area.GetAreaRequest `protobuf:"bytes,7,opt,name=get_area_request,json=getAreaRequest,proto3,oneof"`
}

type SimDemand_ResistParticipantRequest struct {
	ResistParticipantRequest *participant.RegistParticipantRequest `protobuf:"bytes,8,opt,name=resist_participant_request,json=resistParticipantRequest,proto3,oneof"`
}

type SimDemand_SetParticipantsRequest struct {
	SetParticipantsRequest *participant.SetParticipantsRequest `protobuf:"bytes,9,opt,name=set_participants_request,json=setParticipantsRequest,proto3,oneof"`
}

type SimDemand_GetClockRequest struct {
	GetClockRequest *clock.GetClockRequest `protobuf:"bytes,10,opt,name=get_clock_request,json=getClockRequest,proto3,oneof"`
}

type SimDemand_SetClockRequest struct {
	SetClockRequest *clock.SetClockRequest `protobuf:"bytes,11,opt,name=set_clock_request,json=setClockRequest,proto3,oneof"`
}

type SimDemand_StartClockRequest struct {
	StartClockRequest *clock.StartClockRequest `protobuf:"bytes,12,opt,name=start_clock_request,json=startClockRequest,proto3,oneof"`
}

type SimDemand_StopClockRequest struct {
	StopClockRequest *clock.StopClockRequest `protobuf:"bytes,13,opt,name=stop_clock_request,json=stopClockRequest,proto3,oneof"`
}

type SimDemand_ForwardClockRequest struct {
	ForwardClockRequest *clock.ForwardClockRequest `protobuf:"bytes,14,opt,name=forward_clock_request,json=forwardClockRequest,proto3,oneof"`
}

type SimDemand_BackClockRequest struct {
	BackClockRequest *clock.BackClockRequest `protobuf:"bytes,15,opt,name=back_clock_request,json=backClockRequest,proto3,oneof"`
}

func (*SimDemand_GetSameAreaAgentsRequest) isSimDemand_Data() {}

func (*SimDemand_SetAgentsRequest) isSimDemand_Data() {}

func (*SimDemand_VisualiseAgentsRequest) isSimDemand_Data() {}

func (*SimDemand_GetAreaRequest) isSimDemand_Data() {}

func (*SimDemand_ResistParticipantRequest) isSimDemand_Data() {}

func (*SimDemand_SetParticipantsRequest) isSimDemand_Data() {}

func (*SimDemand_GetClockRequest) isSimDemand_Data() {}

func (*SimDemand_SetClockRequest) isSimDemand_Data() {}

func (*SimDemand_StartClockRequest) isSimDemand_Data() {}

func (*SimDemand_StopClockRequest) isSimDemand_Data() {}

func (*SimDemand_ForwardClockRequest) isSimDemand_Data() {}

func (*SimDemand_BackClockRequest) isSimDemand_Data() {}

func (m *SimDemand) GetData() isSimDemand_Data {
	if m != nil {
		return m.Data
	}
	return nil
}

func (m *SimDemand) GetGetSameAreaAgentsRequest() *agent.GetSameAreaAgentsRequest {
	if x, ok := m.GetData().(*SimDemand_GetSameAreaAgentsRequest); ok {
		return x.GetSameAreaAgentsRequest
	}
	return nil
}

func (m *SimDemand) GetSetAgentsRequest() *agent.SetAgentsRequest {
	if x, ok := m.GetData().(*SimDemand_SetAgentsRequest); ok {
		return x.SetAgentsRequest
	}
	return nil
}

func (m *SimDemand) GetVisualiseAgentsRequest() *agent.VisualizeAgentsRequest {
	if x, ok := m.GetData().(*SimDemand_VisualiseAgentsRequest); ok {
		return x.VisualiseAgentsRequest
	}
	return nil
}

func (m *SimDemand) GetGetAreaRequest() *area.GetAreaRequest {
	if x, ok := m.GetData().(*SimDemand_GetAreaRequest); ok {
		return x.GetAreaRequest
	}
	return nil
}

func (m *SimDemand) GetResistParticipantRequest() *participant.RegistParticipantRequest {
	if x, ok := m.GetData().(*SimDemand_ResistParticipantRequest); ok {
		return x.ResistParticipantRequest
	}
	return nil
}

func (m *SimDemand) GetSetParticipantsRequest() *participant.SetParticipantsRequest {
	if x, ok := m.GetData().(*SimDemand_SetParticipantsRequest); ok {
		return x.SetParticipantsRequest
	}
	return nil
}

func (m *SimDemand) GetGetClockRequest() *clock.GetClockRequest {
	if x, ok := m.GetData().(*SimDemand_GetClockRequest); ok {
		return x.GetClockRequest
	}
	return nil
}

func (m *SimDemand) GetSetClockRequest() *clock.SetClockRequest {
	if x, ok := m.GetData().(*SimDemand_SetClockRequest); ok {
		return x.SetClockRequest
	}
	return nil
}

func (m *SimDemand) GetStartClockRequest() *clock.StartClockRequest {
	if x, ok := m.GetData().(*SimDemand_StartClockRequest); ok {
		return x.StartClockRequest
	}
	return nil
}

func (m *SimDemand) GetStopClockRequest() *clock.StopClockRequest {
	if x, ok := m.GetData().(*SimDemand_StopClockRequest); ok {
		return x.StopClockRequest
	}
	return nil
}

func (m *SimDemand) GetForwardClockRequest() *clock.ForwardClockRequest {
	if x, ok := m.GetData().(*SimDemand_ForwardClockRequest); ok {
		return x.ForwardClockRequest
	}
	return nil
}

func (m *SimDemand) GetBackClockRequest() *clock.BackClockRequest {
	if x, ok := m.GetData().(*SimDemand_BackClockRequest); ok {
		return x.BackClockRequest
	}
	return nil
}

// XXX_OneofWrappers is for the internal use of the proto package.
func (*SimDemand) XXX_OneofWrappers() []interface{} {
	return []interface{}{
		(*SimDemand_GetSameAreaAgentsRequest)(nil),
		(*SimDemand_SetAgentsRequest)(nil),
		(*SimDemand_VisualiseAgentsRequest)(nil),
		(*SimDemand_GetAreaRequest)(nil),
		(*SimDemand_ResistParticipantRequest)(nil),
		(*SimDemand_SetParticipantsRequest)(nil),
		(*SimDemand_GetClockRequest)(nil),
		(*SimDemand_SetClockRequest)(nil),
		(*SimDemand_StartClockRequest)(nil),
		(*SimDemand_StopClockRequest)(nil),
		(*SimDemand_ForwardClockRequest)(nil),
		(*SimDemand_BackClockRequest)(nil),
	}
}

type SimSupply struct {
	// demand type
	SupplyType SupplyType `protobuf:"varint,1,opt,name=supply_type,json=supplyType,proto3,enum=api.synerex.SupplyType" json:"supply_type,omitempty"`
	// meta data
	StatusType StatusType `protobuf:"varint,2,opt,name=status_type,json=statusType,proto3,enum=api.synerex.StatusType" json:"status_type,omitempty"`
	Meta       string     `protobuf:"bytes,3,opt,name=meta,proto3" json:"meta,omitempty"`
	// supply data
	//
	// Types that are valid to be assigned to Data:
	//	*SimSupply_GetSameAreaAgentsResponse
	//	*SimSupply_GetNeighborAreaAgentsResponse
	//	*SimSupply_SetAgentsResponse
	//	*SimSupply_GetAreaResponse
	//	*SimSupply_RegistParticipantResponse
	//	*SimSupply_SetParticipantsResponse
	//	*SimSupply_GetClockResponse
	//	*SimSupply_SetClockResponse
	//	*SimSupply_StartClockResponse
	//	*SimSupply_StopClockResponse
	//	*SimSupply_ForwardClockResponse
	//	*SimSupply_BackClockResponse
	Data                 isSimSupply_Data `protobuf_oneof:"data"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_unrecognized     []byte           `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *SimSupply) Reset()         { *m = SimSupply{} }
func (m *SimSupply) String() string { return proto.CompactTextString(m) }
func (*SimSupply) ProtoMessage()    {}
func (*SimSupply) Descriptor() ([]byte, []int) {
	return fileDescriptor_81d77a6394be03b4, []int{1}
}

func (m *SimSupply) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SimSupply.Unmarshal(m, b)
}
func (m *SimSupply) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SimSupply.Marshal(b, m, deterministic)
}
func (m *SimSupply) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SimSupply.Merge(m, src)
}
func (m *SimSupply) XXX_Size() int {
	return xxx_messageInfo_SimSupply.Size(m)
}
func (m *SimSupply) XXX_DiscardUnknown() {
	xxx_messageInfo_SimSupply.DiscardUnknown(m)
}

var xxx_messageInfo_SimSupply proto.InternalMessageInfo

func (m *SimSupply) GetSupplyType() SupplyType {
	if m != nil {
		return m.SupplyType
	}
	return SupplyType_GET_SAME_AREA_AGENTS_RESPONSE
}

func (m *SimSupply) GetStatusType() StatusType {
	if m != nil {
		return m.StatusType
	}
	return StatusType_OK
}

func (m *SimSupply) GetMeta() string {
	if m != nil {
		return m.Meta
	}
	return ""
}

type isSimSupply_Data interface {
	isSimSupply_Data()
}

type SimSupply_GetSameAreaAgentsResponse struct {
	GetSameAreaAgentsResponse *agent.GetSameAreaAgentsResponse `protobuf:"bytes,4,opt,name=get_same_area_agents_response,json=getSameAreaAgentsResponse,proto3,oneof"`
}

type SimSupply_GetNeighborAreaAgentsResponse struct {
	GetNeighborAreaAgentsResponse *agent.GetNeighborAreaAgentsResponse `protobuf:"bytes,5,opt,name=get_neighbor_area_agents_response,json=getNeighborAreaAgentsResponse,proto3,oneof"`
}

type SimSupply_SetAgentsResponse struct {
	SetAgentsResponse *agent.SetAgentsResponse `protobuf:"bytes,6,opt,name=set_agents_response,json=setAgentsResponse,proto3,oneof"`
}

type SimSupply_GetAreaResponse struct {
	GetAreaResponse *area.GetAreaResponse `protobuf:"bytes,7,opt,name=get_area_response,json=getAreaResponse,proto3,oneof"`
}

type SimSupply_RegistParticipantResponse struct {
	RegistParticipantResponse *participant.RegistParticipantResponse `protobuf:"bytes,8,opt,name=regist_participant_response,json=registParticipantResponse,proto3,oneof"`
}

type SimSupply_SetParticipantsResponse struct {
	SetParticipantsResponse *participant.SetParticipantsResponse `protobuf:"bytes,9,opt,name=set_participants_response,json=setParticipantsResponse,proto3,oneof"`
}

type SimSupply_GetClockResponse struct {
	GetClockResponse *clock.GetClockResponse `protobuf:"bytes,10,opt,name=get_clock_response,json=getClockResponse,proto3,oneof"`
}

type SimSupply_SetClockResponse struct {
	SetClockResponse *clock.SetClockResponse `protobuf:"bytes,11,opt,name=set_clock_response,json=setClockResponse,proto3,oneof"`
}

type SimSupply_StartClockResponse struct {
	StartClockResponse *clock.StartClockResponse `protobuf:"bytes,12,opt,name=start_clock_response,json=startClockResponse,proto3,oneof"`
}

type SimSupply_StopClockResponse struct {
	StopClockResponse *clock.StopClockResponse `protobuf:"bytes,13,opt,name=stop_clock_response,json=stopClockResponse,proto3,oneof"`
}

type SimSupply_ForwardClockResponse struct {
	ForwardClockResponse *clock.ForwardClockResponse `protobuf:"bytes,14,opt,name=forward_clock_response,json=forwardClockResponse,proto3,oneof"`
}

type SimSupply_BackClockResponse struct {
	BackClockResponse *clock.BackClockResponse `protobuf:"bytes,15,opt,name=back_clock_response,json=backClockResponse,proto3,oneof"`
}

func (*SimSupply_GetSameAreaAgentsResponse) isSimSupply_Data() {}

func (*SimSupply_GetNeighborAreaAgentsResponse) isSimSupply_Data() {}

func (*SimSupply_SetAgentsResponse) isSimSupply_Data() {}

func (*SimSupply_GetAreaResponse) isSimSupply_Data() {}

func (*SimSupply_RegistParticipantResponse) isSimSupply_Data() {}

func (*SimSupply_SetParticipantsResponse) isSimSupply_Data() {}

func (*SimSupply_GetClockResponse) isSimSupply_Data() {}

func (*SimSupply_SetClockResponse) isSimSupply_Data() {}

func (*SimSupply_StartClockResponse) isSimSupply_Data() {}

func (*SimSupply_StopClockResponse) isSimSupply_Data() {}

func (*SimSupply_ForwardClockResponse) isSimSupply_Data() {}

func (*SimSupply_BackClockResponse) isSimSupply_Data() {}

func (m *SimSupply) GetData() isSimSupply_Data {
	if m != nil {
		return m.Data
	}
	return nil
}

func (m *SimSupply) GetGetSameAreaAgentsResponse() *agent.GetSameAreaAgentsResponse {
	if x, ok := m.GetData().(*SimSupply_GetSameAreaAgentsResponse); ok {
		return x.GetSameAreaAgentsResponse
	}
	return nil
}

func (m *SimSupply) GetGetNeighborAreaAgentsResponse() *agent.GetNeighborAreaAgentsResponse {
	if x, ok := m.GetData().(*SimSupply_GetNeighborAreaAgentsResponse); ok {
		return x.GetNeighborAreaAgentsResponse
	}
	return nil
}

func (m *SimSupply) GetSetAgentsResponse() *agent.SetAgentsResponse {
	if x, ok := m.GetData().(*SimSupply_SetAgentsResponse); ok {
		return x.SetAgentsResponse
	}
	return nil
}

func (m *SimSupply) GetGetAreaResponse() *area.GetAreaResponse {
	if x, ok := m.GetData().(*SimSupply_GetAreaResponse); ok {
		return x.GetAreaResponse
	}
	return nil
}

func (m *SimSupply) GetRegistParticipantResponse() *participant.RegistParticipantResponse {
	if x, ok := m.GetData().(*SimSupply_RegistParticipantResponse); ok {
		return x.RegistParticipantResponse
	}
	return nil
}

func (m *SimSupply) GetSetParticipantsResponse() *participant.SetParticipantsResponse {
	if x, ok := m.GetData().(*SimSupply_SetParticipantsResponse); ok {
		return x.SetParticipantsResponse
	}
	return nil
}

func (m *SimSupply) GetGetClockResponse() *clock.GetClockResponse {
	if x, ok := m.GetData().(*SimSupply_GetClockResponse); ok {
		return x.GetClockResponse
	}
	return nil
}

func (m *SimSupply) GetSetClockResponse() *clock.SetClockResponse {
	if x, ok := m.GetData().(*SimSupply_SetClockResponse); ok {
		return x.SetClockResponse
	}
	return nil
}

func (m *SimSupply) GetStartClockResponse() *clock.StartClockResponse {
	if x, ok := m.GetData().(*SimSupply_StartClockResponse); ok {
		return x.StartClockResponse
	}
	return nil
}

func (m *SimSupply) GetStopClockResponse() *clock.StopClockResponse {
	if x, ok := m.GetData().(*SimSupply_StopClockResponse); ok {
		return x.StopClockResponse
	}
	return nil
}

func (m *SimSupply) GetForwardClockResponse() *clock.ForwardClockResponse {
	if x, ok := m.GetData().(*SimSupply_ForwardClockResponse); ok {
		return x.ForwardClockResponse
	}
	return nil
}

func (m *SimSupply) GetBackClockResponse() *clock.BackClockResponse {
	if x, ok := m.GetData().(*SimSupply_BackClockResponse); ok {
		return x.BackClockResponse
	}
	return nil
}

// XXX_OneofWrappers is for the internal use of the proto package.
func (*SimSupply) XXX_OneofWrappers() []interface{} {
	return []interface{}{
		(*SimSupply_GetSameAreaAgentsResponse)(nil),
		(*SimSupply_GetNeighborAreaAgentsResponse)(nil),
		(*SimSupply_SetAgentsResponse)(nil),
		(*SimSupply_GetAreaResponse)(nil),
		(*SimSupply_RegistParticipantResponse)(nil),
		(*SimSupply_SetParticipantsResponse)(nil),
		(*SimSupply_GetClockResponse)(nil),
		(*SimSupply_SetClockResponse)(nil),
		(*SimSupply_StartClockResponse)(nil),
		(*SimSupply_StopClockResponse)(nil),
		(*SimSupply_ForwardClockResponse)(nil),
		(*SimSupply_BackClockResponse)(nil),
	}
}

func init() {
	proto.RegisterEnum("api.synerex.StatusType", StatusType_name, StatusType_value)
	proto.RegisterEnum("api.synerex.DemandType", DemandType_name, DemandType_value)
	proto.RegisterEnum("api.synerex.SupplyType", SupplyType_name, SupplyType_value)
	proto.RegisterType((*SimDemand)(nil), "api.synerex.SimDemand")
	proto.RegisterType((*SimSupply)(nil), "api.synerex.SimSupply")
}

func init() { proto.RegisterFile("simulation/synerex/synerex.proto", fileDescriptor_81d77a6394be03b4) }

var fileDescriptor_81d77a6394be03b4 = []byte{
	// 1091 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xac, 0x97, 0xdd, 0x52, 0xdb, 0x46,
	0x18, 0x86, 0xb1, 0x31, 0x06, 0x7f, 0x6e, 0xc0, 0x2c, 0x7f, 0xc6, 0x40, 0x02, 0xe9, 0x1f, 0xe5,
	0x00, 0x66, 0xd2, 0x83, 0xb6, 0x87, 0x02, 0x14, 0xe3, 0x21, 0xb5, 0x89, 0xa4, 0x24, 0x33, 0x99,
	0xe9, 0x68, 0xd6, 0xf6, 0x22, 0x34, 0xd8, 0x96, 0xaa, 0x5d, 0xb7, 0xa5, 0x57, 0xd2, 0x6b, 0xea,
	0x51, 0xaf, 0xa0, 0xd7, 0xd2, 0xd9, 0x1f, 0x4b, 0xeb, 0x95, 0x9c, 0xe9, 0x41, 0x4e, 0xe4, 0x9d,
	0xf7, 0xfb, 0xf6, 0xd9, 0xd5, 0xcb, 0xea, 0x15, 0x82, 0x63, 0x1a, 0x8e, 0xa7, 0x23, 0xcc, 0xc2,
	0x68, 0x72, 0x41, 0x9f, 0x26, 0x24, 0x21, 0x7f, 0xcc, 0x7e, 0xcf, 0xe3, 0x24, 0x62, 0x11, 0xaa,
	0xe3, 0x38, 0x3c, 0x57, 0x52, 0xeb, 0x50, 0x6b, 0xc7, 0x01, 0x99, 0x30, 0x79, 0x95, 0xad, 0xad,
	0x96, 0x5e, 0x4d, 0x08, 0x16, 0x17, 0x55, 0xd3, 0x67, 0x0e, 0x46, 0xd1, 0xe0, 0x51, 0x5e, 0x55,
	0xf5, 0x54, 0xab, 0xc6, 0x38, 0x61, 0xe1, 0x20, 0x8c, 0xf1, 0x84, 0xe9, 0x63, 0xd9, 0xf9, 0xf2,
	0xaf, 0x1a, 0xd4, 0xdc, 0x70, 0x7c, 0x4d, 0xc6, 0x78, 0x32, 0x44, 0x3f, 0x42, 0x7d, 0x28, 0x46,
	0x3e, 0x7b, 0x8a, 0x49, 0xb3, 0x74, 0x5c, 0x3a, 0x5d, 0x7f, 0xb5, 0x77, 0xae, 0x6d, 0xf9, 0x5c,
	0x76, 0x7a, 0x4f, 0x31, 0x71, 0x60, 0x98, 0x8e, 0xf9, 0x4c, 0xca, 0x30, 0x9b, 0x52, 0x39, 0xb3,
	0x5c, 0x30, 0xd3, 0x15, 0x75, 0x39, 0x93, 0xa6, 0x63, 0x84, 0xa0, 0x32, 0x26, 0x0c, 0x37, 0x97,
	0x8f, 0x4b, 0xa7, 0x35, 0x47, 0x8c, 0x11, 0x81, 0xc3, 0x80, 0x30, 0x9f, 0xe2, 0x31, 0xf1, 0xf9,
	0x4d, 0xfb, 0xc2, 0x16, 0xea, 0x27, 0xe4, 0xd7, 0x29, 0xa1, 0xac, 0x59, 0x39, 0x2e, 0x9d, 0xd6,
	0x5f, 0x7d, 0x29, 0xf0, 0xd2, 0xb1, 0x36, 0x61, 0x2e, 0x1e, 0x13, 0x2b, 0x21, 0xd8, 0x12, 0xbd,
	0x8e, 0x6c, 0xbd, 0x59, 0x72, 0x9a, 0xc1, 0x82, 0x1a, 0xba, 0x05, 0x44, 0x09, 0x33, 0xe1, 0x2b,
	0x02, 0x7e, 0xa0, 0xc1, 0x5d, 0xc2, 0x4c, 0x68, 0x83, 0x1a, 0x1a, 0xfa, 0x05, 0x9a, 0xbf, 0x85,
	0x74, 0x8a, 0x47, 0x21, 0x25, 0x26, 0xb2, 0x2a, 0x90, 0x27, 0x1a, 0xf2, 0xbd, 0x6c, 0xfd, 0x93,
	0x98, 0xe0, 0xdd, 0x14, 0x32, 0x8f, 0xbf, 0x86, 0x06, 0xb7, 0x44, 0xb8, 0x31, 0xc3, 0xae, 0x0a,
	0x6c, 0x53, 0x62, 0xf9, 0xd9, 0x68, 0x13, 0xc6, 0xef, 0x32, 0xa3, 0xad, 0x07, 0x73, 0x0a, 0x0a,
	0xa1, 0x95, 0x10, 0x1a, 0x52, 0xe6, 0x6b, 0x47, 0x21, 0xe5, 0xad, 0x09, 0xde, 0x77, 0x82, 0xa7,
	0x1f, 0x15, 0x87, 0x04, 0x21, 0x65, 0x77, 0x99, 0xa2, 0x99, 0x2b, 0x71, 0xf9, 0x1a, 0x1a, 0x40,
	0x93, 0x9b, 0xab, 0x71, 0x32, 0x3f, 0x6a, 0x62, 0xa1, 0x6f, 0x73, 0x0b, 0xb9, 0x44, 0x27, 0xe9,
	0xae, 0xd0, 0xc2, 0x0a, 0xba, 0x81, 0x4d, 0xee, 0x8a, 0x38, 0xfb, 0x29, 0x1d, 0x04, 0xbd, 0x25,
	0xe8, 0xf2, 0xa9, 0x68, 0x13, 0x76, 0xc5, 0x07, 0x19, 0x70, 0x23, 0x98, 0x97, 0x38, 0x89, 0xe6,
	0x48, 0xf5, 0x1c, 0xc9, 0xcd, 0x93, 0xa8, 0x41, 0xea, 0xc2, 0x16, 0x65, 0x38, 0x31, 0x59, 0x5f,
	0x08, 0xd6, 0xa1, 0xce, 0xe2, 0x5d, 0x06, 0x6d, 0x93, 0x9a, 0xa2, 0x38, 0xa5, 0x2c, 0x8a, 0x0d,
	0xdc, 0x33, 0xed, 0x94, 0xce, 0x70, 0x51, 0x6c, 0xd0, 0x1a, 0xd4, 0xd0, 0x90, 0x07, 0x3b, 0xf7,
	0x51, 0xf2, 0x3b, 0x4e, 0x86, 0x06, 0x6f, 0x5d, 0xf0, 0x9e, 0x6b, 0xbc, 0xd7, 0xb2, 0xcf, 0x40,
	0x6e, 0xdd, 0xe7, 0x65, 0xbe, 0xc5, 0x3e, 0x1e, 0x3c, 0x1a, 0xc8, 0x8d, 0xdc, 0x16, 0x2f, 0xf1,
	0xe0, 0xd1, 0xdc, 0x62, 0xdf, 0xd0, 0x2e, 0xab, 0x50, 0x19, 0x62, 0x86, 0x5f, 0xfe, 0x23, 0xa3,
	0xc9, 0x9d, 0xc6, 0xf1, 0xe8, 0x49, 0x04, 0x8c, 0x18, 0x2d, 0x8e, 0x26, 0xd9, 0xa9, 0x02, 0x26,
	0x1d, 0x7f, 0xe6, 0x68, 0x7a, 0x80, 0xa3, 0x05, 0xd1, 0x44, 0xe3, 0x68, 0x42, 0x89, 0xca, 0xa6,
	0xaf, 0x3e, 0x9d, 0x4d, 0xb2, 0xf7, 0x66, 0xc9, 0xd9, 0x0f, 0x16, 0x15, 0x11, 0x83, 0x13, 0xbe,
	0xd2, 0x84, 0x84, 0xc1, 0x43, 0x3f, 0x4a, 0x8a, 0x57, 0x93, 0x61, 0x75, 0x3a, 0xbf, 0x5a, 0x57,
	0x4d, 0x29, 0x5c, 0x91, 0x6f, 0x7f, 0x71, 0x83, 0x38, 0xbd, 0x7a, 0x26, 0xaa, 0x75, 0xaa, 0xda,
	0xe9, 0xcd, 0x85, 0x62, 0xca, 0xde, 0xa4, 0xa6, 0x88, 0xda, 0xf2, 0x09, 0x55, 0xb9, 0xa5, 0x68,
	0x32, 0xb8, 0xf6, 0x0b, 0x82, 0x2b, 0x45, 0x6d, 0x04, 0xf3, 0x12, 0x1a, 0xc1, 0x41, 0x22, 0x72,
	0xc8, 0x88, 0x2e, 0x85, 0x94, 0xd9, 0x75, 0xf6, 0x7f, 0xb2, 0x2b, 0x33, 0x3f, 0x59, 0x54, 0x44,
	0xf7, 0xb0, 0x5f, 0x90, 0x5e, 0x6a, 0xad, 0x9a, 0x66, 0xfa, 0x27, 0xe3, 0x2b, 0x5d, 0x69, 0x8f,
	0x16, 0x97, 0xf8, 0x93, 0xa3, 0x07, 0x98, 0x5a, 0x00, 0x72, 0x4f, 0x4e, 0x96, 0x60, 0x29, 0xb3,
	0x11, 0x18, 0xda, 0xec, 0x7d, 0x66, 0xc0, 0xea, 0xf9, 0xa4, 0x28, 0x80, 0x51, 0x13, 0xf6, 0x16,
	0xb6, 0xe7, 0x63, 0x4c, 0xe1, 0x64, 0x8e, 0x1d, 0x2d, 0xc8, 0xb1, 0x14, 0x88, 0x68, 0x4e, 0x95,
	0xc9, 0xa8, 0x25, 0x99, 0x22, 0x3e, 0x2b, 0x48, 0xc6, 0x34, 0xb6, 0xb4, 0xb3, 0x65, 0x8a, 0xe8,
	0x03, 0xec, 0x9a, 0x61, 0xa6, 0x90, 0x32, 0xcd, 0x5e, 0x2c, 0x4c, 0xb3, 0x94, 0xba, 0x7d, 0x5f,
	0xa0, 0xf3, 0x8d, 0xce, 0xe5, 0x99, 0xa2, 0x6e, 0xe4, 0x36, 0xaa, 0x05, 0x5a, 0xb6, 0xd1, 0xbe,
	0x29, 0xce, 0x22, 0xed, 0xec, 0x1b, 0x80, 0x2c, 0x6a, 0x50, 0x15, 0xca, 0xbd, 0xdb, 0xc6, 0x12,
	0xff, 0xed, 0xb6, 0x1b, 0x25, 0xb4, 0x06, 0x95, 0x6e, 0xaf, 0x6b, 0x37, 0xca, 0x67, 0x7f, 0x97,
	0x01, 0xb2, 0x7f, 0xb4, 0xd0, 0x31, 0x1c, 0xb6, 0x6d, 0xcf, 0x77, 0xad, 0x9f, 0x6d, 0xdf, 0x72,
	0x6c, 0xcb, 0xb7, 0xda, 0x76, 0xd7, 0x73, 0x7d, 0xc7, 0x7e, 0xfb, 0xce, 0x76, 0xbd, 0xc6, 0x12,
	0xda, 0x05, 0xe4, 0xda, 0x9e, 0xa9, 0x97, 0xd0, 0x21, 0x34, 0xdf, 0x77, 0xdc, 0x77, 0xd6, 0x9b,
	0xce, 0x47, 0xdb, 0xac, 0x96, 0xd1, 0x36, 0x34, 0x38, 0x57, 0x20, 0x67, 0xea, 0x32, 0x7a, 0x0e,
	0x2d, 0xc7, 0x76, 0x3b, 0xae, 0xe7, 0xdf, 0x59, 0x8e, 0xd7, 0xb9, 0xea, 0xdc, 0x59, 0x5d, 0x2f,
	0xad, 0x57, 0x38, 0x93, 0xaf, 0xa5, 0x15, 0x33, 0xe6, 0x0a, 0xda, 0x81, 0x4d, 0xce, 0xbc, 0x7a,
	0xd3, 0xbb, 0xba, 0x4d, 0xe5, 0x2a, 0x97, 0xdd, 0x9c, 0xbc, 0x8a, 0xf6, 0x60, 0xcb, 0xf5, 0x2c,
	0xc7, 0x2c, 0xac, 0x89, 0x1b, 0xf2, 0x7a, 0x77, 0x86, 0x5e, 0x43, 0xfb, 0xb0, 0xf3, 0xba, 0xe7,
	0x7c, 0xb0, 0x9c, 0x6b, 0xa3, 0x04, 0x7c, 0xca, 0xa5, 0x75, 0x75, 0x6b, 0xe8, 0xf5, 0xb3, 0x7f,
	0xcb, 0x00, 0xd9, 0xab, 0x01, 0x9d, 0xc0, 0xd1, 0x02, 0x33, 0xdd, 0xbb, 0x5e, 0xd7, 0xb5, 0x1b,
	0x4b, 0xe8, 0x6b, 0x38, 0xe1, 0x2d, 0x5d, 0xbb, 0xd3, 0xbe, 0xb9, 0xec, 0x39, 0xc5, 0x6d, 0x25,
	0xb1, 0x79, 0xdd, 0x74, 0x55, 0x28, 0xcf, 0x3c, 0x50, 0xbe, 0x2a, 0x79, 0x19, 0xbd, 0x80, 0x83,
	0x42, 0x63, 0x55, 0x43, 0x05, 0x1d, 0xc1, 0x7e, 0x81, 0xb3, 0xaa, 0xbc, 0xc2, 0x6f, 0x50, 0xb7,
	0x56, 0xe9, 0xd5, 0xd9, 0x1f, 0xdf, 0xd0, 0x57, 0x51, 0x13, 0xb6, 0xe7, 0xcd, 0x55, 0x95, 0x35,
	0x69, 0xbb, 0xe6, 0xae, 0x2a, 0xd4, 0x50, 0x0b, 0x76, 0x4d, 0x7b, 0x55, 0x0d, 0xf8, 0xa4, 0x39,
	0x7f, 0x55, 0xa1, 0x7e, 0xf9, 0xd3, 0xc7, 0x1f, 0x82, 0x90, 0x3d, 0x4c, 0xfb, 0xe7, 0x83, 0x68,
	0x6c, 0x7e, 0xf6, 0xf8, 0x78, 0x14, 0x3f, 0xe0, 0x0b, 0x1c, 0x87, 0x17, 0xf9, 0x6f, 0xa3, 0x7e,
	0x55, 0x7c, 0x85, 0x7c, 0xff, 0x5f, 0x00, 0x00, 0x00, 0xff, 0xff, 0xf3, 0xd6, 0x00, 0x63, 0x38,
	0x0d, 0x00, 0x00,
}

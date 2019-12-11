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

type SetClockRequest struct {
	Clock                *Clock   `protobuf:"bytes,1,opt,name=clock,proto3" json:"clock,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SetClockRequest) Reset()         { *m = SetClockRequest{} }
func (m *SetClockRequest) String() string { return proto.CompactTextString(m) }
func (*SetClockRequest) ProtoMessage()    {}
func (*SetClockRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_e96fed1976809896, []int{0}
}

func (m *SetClockRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SetClockRequest.Unmarshal(m, b)
}
func (m *SetClockRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SetClockRequest.Marshal(b, m, deterministic)
}
func (m *SetClockRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SetClockRequest.Merge(m, src)
}
func (m *SetClockRequest) XXX_Size() int {
	return xxx_messageInfo_SetClockRequest.Size(m)
}
func (m *SetClockRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_SetClockRequest.DiscardUnknown(m)
}

var xxx_messageInfo_SetClockRequest proto.InternalMessageInfo

func (m *SetClockRequest) GetClock() *Clock {
	if m != nil {
		return m.Clock
	}
	return nil
}

type SetClockResponse struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SetClockResponse) Reset()         { *m = SetClockResponse{} }
func (m *SetClockResponse) String() string { return proto.CompactTextString(m) }
func (*SetClockResponse) ProtoMessage()    {}
func (*SetClockResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_e96fed1976809896, []int{1}
}

func (m *SetClockResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SetClockResponse.Unmarshal(m, b)
}
func (m *SetClockResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SetClockResponse.Marshal(b, m, deterministic)
}
func (m *SetClockResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SetClockResponse.Merge(m, src)
}
func (m *SetClockResponse) XXX_Size() int {
	return xxx_messageInfo_SetClockResponse.Size(m)
}
func (m *SetClockResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_SetClockResponse.DiscardUnknown(m)
}

var xxx_messageInfo_SetClockResponse proto.InternalMessageInfo

type GetClockRequest struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetClockRequest) Reset()         { *m = GetClockRequest{} }
func (m *GetClockRequest) String() string { return proto.CompactTextString(m) }
func (*GetClockRequest) ProtoMessage()    {}
func (*GetClockRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_e96fed1976809896, []int{2}
}

func (m *GetClockRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetClockRequest.Unmarshal(m, b)
}
func (m *GetClockRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetClockRequest.Marshal(b, m, deterministic)
}
func (m *GetClockRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetClockRequest.Merge(m, src)
}
func (m *GetClockRequest) XXX_Size() int {
	return xxx_messageInfo_GetClockRequest.Size(m)
}
func (m *GetClockRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_GetClockRequest.DiscardUnknown(m)
}

var xxx_messageInfo_GetClockRequest proto.InternalMessageInfo

type GetClockResponse struct {
	Clock                *Clock   `protobuf:"bytes,1,opt,name=clock,proto3" json:"clock,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetClockResponse) Reset()         { *m = GetClockResponse{} }
func (m *GetClockResponse) String() string { return proto.CompactTextString(m) }
func (*GetClockResponse) ProtoMessage()    {}
func (*GetClockResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_e96fed1976809896, []int{3}
}

func (m *GetClockResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetClockResponse.Unmarshal(m, b)
}
func (m *GetClockResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetClockResponse.Marshal(b, m, deterministic)
}
func (m *GetClockResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetClockResponse.Merge(m, src)
}
func (m *GetClockResponse) XXX_Size() int {
	return xxx_messageInfo_GetClockResponse.Size(m)
}
func (m *GetClockResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_GetClockResponse.DiscardUnknown(m)
}

var xxx_messageInfo_GetClockResponse proto.InternalMessageInfo

func (m *GetClockResponse) GetClock() *Clock {
	if m != nil {
		return m.Clock
	}
	return nil
}

type ForwardClockRequest struct {
	StepNum              uint64   `protobuf:"varint,1,opt,name=step_num,json=stepNum,proto3" json:"step_num,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ForwardClockRequest) Reset()         { *m = ForwardClockRequest{} }
func (m *ForwardClockRequest) String() string { return proto.CompactTextString(m) }
func (*ForwardClockRequest) ProtoMessage()    {}
func (*ForwardClockRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_e96fed1976809896, []int{4}
}

func (m *ForwardClockRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ForwardClockRequest.Unmarshal(m, b)
}
func (m *ForwardClockRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ForwardClockRequest.Marshal(b, m, deterministic)
}
func (m *ForwardClockRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ForwardClockRequest.Merge(m, src)
}
func (m *ForwardClockRequest) XXX_Size() int {
	return xxx_messageInfo_ForwardClockRequest.Size(m)
}
func (m *ForwardClockRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_ForwardClockRequest.DiscardUnknown(m)
}

var xxx_messageInfo_ForwardClockRequest proto.InternalMessageInfo

func (m *ForwardClockRequest) GetStepNum() uint64 {
	if m != nil {
		return m.StepNum
	}
	return 0
}

type ForwardClockResponse struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ForwardClockResponse) Reset()         { *m = ForwardClockResponse{} }
func (m *ForwardClockResponse) String() string { return proto.CompactTextString(m) }
func (*ForwardClockResponse) ProtoMessage()    {}
func (*ForwardClockResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_e96fed1976809896, []int{5}
}

func (m *ForwardClockResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ForwardClockResponse.Unmarshal(m, b)
}
func (m *ForwardClockResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ForwardClockResponse.Marshal(b, m, deterministic)
}
func (m *ForwardClockResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ForwardClockResponse.Merge(m, src)
}
func (m *ForwardClockResponse) XXX_Size() int {
	return xxx_messageInfo_ForwardClockResponse.Size(m)
}
func (m *ForwardClockResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_ForwardClockResponse.DiscardUnknown(m)
}

var xxx_messageInfo_ForwardClockResponse proto.InternalMessageInfo

type BackClockRequest struct {
	StepNum              uint64   `protobuf:"varint,1,opt,name=step_num,json=stepNum,proto3" json:"step_num,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *BackClockRequest) Reset()         { *m = BackClockRequest{} }
func (m *BackClockRequest) String() string { return proto.CompactTextString(m) }
func (*BackClockRequest) ProtoMessage()    {}
func (*BackClockRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_e96fed1976809896, []int{6}
}

func (m *BackClockRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_BackClockRequest.Unmarshal(m, b)
}
func (m *BackClockRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_BackClockRequest.Marshal(b, m, deterministic)
}
func (m *BackClockRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_BackClockRequest.Merge(m, src)
}
func (m *BackClockRequest) XXX_Size() int {
	return xxx_messageInfo_BackClockRequest.Size(m)
}
func (m *BackClockRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_BackClockRequest.DiscardUnknown(m)
}

var xxx_messageInfo_BackClockRequest proto.InternalMessageInfo

func (m *BackClockRequest) GetStepNum() uint64 {
	if m != nil {
		return m.StepNum
	}
	return 0
}

type BackClockResponse struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *BackClockResponse) Reset()         { *m = BackClockResponse{} }
func (m *BackClockResponse) String() string { return proto.CompactTextString(m) }
func (*BackClockResponse) ProtoMessage()    {}
func (*BackClockResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_e96fed1976809896, []int{7}
}

func (m *BackClockResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_BackClockResponse.Unmarshal(m, b)
}
func (m *BackClockResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_BackClockResponse.Marshal(b, m, deterministic)
}
func (m *BackClockResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_BackClockResponse.Merge(m, src)
}
func (m *BackClockResponse) XXX_Size() int {
	return xxx_messageInfo_BackClockResponse.Size(m)
}
func (m *BackClockResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_BackClockResponse.DiscardUnknown(m)
}

var xxx_messageInfo_BackClockResponse proto.InternalMessageInfo

type StartClockRequest struct {
	StepNum              uint64   `protobuf:"varint,1,opt,name=step_num,json=stepNum,proto3" json:"step_num,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *StartClockRequest) Reset()         { *m = StartClockRequest{} }
func (m *StartClockRequest) String() string { return proto.CompactTextString(m) }
func (*StartClockRequest) ProtoMessage()    {}
func (*StartClockRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_e96fed1976809896, []int{8}
}

func (m *StartClockRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_StartClockRequest.Unmarshal(m, b)
}
func (m *StartClockRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_StartClockRequest.Marshal(b, m, deterministic)
}
func (m *StartClockRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_StartClockRequest.Merge(m, src)
}
func (m *StartClockRequest) XXX_Size() int {
	return xxx_messageInfo_StartClockRequest.Size(m)
}
func (m *StartClockRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_StartClockRequest.DiscardUnknown(m)
}

var xxx_messageInfo_StartClockRequest proto.InternalMessageInfo

func (m *StartClockRequest) GetStepNum() uint64 {
	if m != nil {
		return m.StepNum
	}
	return 0
}

type StartClockResponse struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *StartClockResponse) Reset()         { *m = StartClockResponse{} }
func (m *StartClockResponse) String() string { return proto.CompactTextString(m) }
func (*StartClockResponse) ProtoMessage()    {}
func (*StartClockResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_e96fed1976809896, []int{9}
}

func (m *StartClockResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_StartClockResponse.Unmarshal(m, b)
}
func (m *StartClockResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_StartClockResponse.Marshal(b, m, deterministic)
}
func (m *StartClockResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_StartClockResponse.Merge(m, src)
}
func (m *StartClockResponse) XXX_Size() int {
	return xxx_messageInfo_StartClockResponse.Size(m)
}
func (m *StartClockResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_StartClockResponse.DiscardUnknown(m)
}

var xxx_messageInfo_StartClockResponse proto.InternalMessageInfo

type StopClockRequest struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *StopClockRequest) Reset()         { *m = StopClockRequest{} }
func (m *StopClockRequest) String() string { return proto.CompactTextString(m) }
func (*StopClockRequest) ProtoMessage()    {}
func (*StopClockRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_e96fed1976809896, []int{10}
}

func (m *StopClockRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_StopClockRequest.Unmarshal(m, b)
}
func (m *StopClockRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_StopClockRequest.Marshal(b, m, deterministic)
}
func (m *StopClockRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_StopClockRequest.Merge(m, src)
}
func (m *StopClockRequest) XXX_Size() int {
	return xxx_messageInfo_StopClockRequest.Size(m)
}
func (m *StopClockRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_StopClockRequest.DiscardUnknown(m)
}

var xxx_messageInfo_StopClockRequest proto.InternalMessageInfo

type StopClockResponse struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *StopClockResponse) Reset()         { *m = StopClockResponse{} }
func (m *StopClockResponse) String() string { return proto.CompactTextString(m) }
func (*StopClockResponse) ProtoMessage()    {}
func (*StopClockResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_e96fed1976809896, []int{11}
}

func (m *StopClockResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_StopClockResponse.Unmarshal(m, b)
}
func (m *StopClockResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_StopClockResponse.Marshal(b, m, deterministic)
}
func (m *StopClockResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_StopClockResponse.Merge(m, src)
}
func (m *StopClockResponse) XXX_Size() int {
	return xxx_messageInfo_StopClockResponse.Size(m)
}
func (m *StopClockResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_StopClockResponse.DiscardUnknown(m)
}

var xxx_messageInfo_StopClockResponse proto.InternalMessageInfo

type Clock struct {
	GlobalTime           float64  `protobuf:"fixed64,1,opt,name=global_time,json=globalTime,proto3" json:"global_time,omitempty"`
	TimeStep             float64  `protobuf:"fixed64,2,opt,name=time_step,json=timeStep,proto3" json:"time_step,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Clock) Reset()         { *m = Clock{} }
func (m *Clock) String() string { return proto.CompactTextString(m) }
func (*Clock) ProtoMessage()    {}
func (*Clock) Descriptor() ([]byte, []int) {
	return fileDescriptor_e96fed1976809896, []int{12}
}

func (m *Clock) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Clock.Unmarshal(m, b)
}
func (m *Clock) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Clock.Marshal(b, m, deterministic)
}
func (m *Clock) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Clock.Merge(m, src)
}
func (m *Clock) XXX_Size() int {
	return xxx_messageInfo_Clock.Size(m)
}
func (m *Clock) XXX_DiscardUnknown() {
	xxx_messageInfo_Clock.DiscardUnknown(m)
}

var xxx_messageInfo_Clock proto.InternalMessageInfo

func (m *Clock) GetGlobalTime() float64 {
	if m != nil {
		return m.GlobalTime
	}
	return 0
}

func (m *Clock) GetTimeStep() float64 {
	if m != nil {
		return m.TimeStep
	}
	return 0
}

func init() {
	proto.RegisterType((*SetClockRequest)(nil), "api.clock.SetClockRequest")
	proto.RegisterType((*SetClockResponse)(nil), "api.clock.SetClockResponse")
	proto.RegisterType((*GetClockRequest)(nil), "api.clock.GetClockRequest")
	proto.RegisterType((*GetClockResponse)(nil), "api.clock.GetClockResponse")
	proto.RegisterType((*ForwardClockRequest)(nil), "api.clock.ForwardClockRequest")
	proto.RegisterType((*ForwardClockResponse)(nil), "api.clock.ForwardClockResponse")
	proto.RegisterType((*BackClockRequest)(nil), "api.clock.BackClockRequest")
	proto.RegisterType((*BackClockResponse)(nil), "api.clock.BackClockResponse")
	proto.RegisterType((*StartClockRequest)(nil), "api.clock.StartClockRequest")
	proto.RegisterType((*StartClockResponse)(nil), "api.clock.StartClockResponse")
	proto.RegisterType((*StopClockRequest)(nil), "api.clock.StopClockRequest")
	proto.RegisterType((*StopClockResponse)(nil), "api.clock.StopClockResponse")
	proto.RegisterType((*Clock)(nil), "api.clock.Clock")
}

func init() { proto.RegisterFile("simulation/clock/clock.proto", fileDescriptor_e96fed1976809896) }

var fileDescriptor_e96fed1976809896 = []byte{
	// 299 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x92, 0x4f, 0x4b, 0xc3, 0x40,
	0x10, 0xc5, 0xa9, 0x58, 0x6d, 0xa7, 0x87, 0xa6, 0x9b, 0x22, 0x15, 0x05, 0x25, 0x07, 0xf1, 0xe2,
	0x46, 0x14, 0x11, 0x3d, 0x56, 0xb4, 0x37, 0x0f, 0x89, 0x27, 0x2f, 0x61, 0x13, 0x97, 0x76, 0x69,
	0x36, 0xbb, 0x66, 0x27, 0xa8, 0xdf, 0x5e, 0x76, 0x13, 0x6b, 0xd3, 0x53, 0x2e, 0xf9, 0xf3, 0xde,
	0x9b, 0xdf, 0x0c, 0xc3, 0xc0, 0xa9, 0x11, 0xb2, 0xca, 0x19, 0x0a, 0x55, 0x84, 0x59, 0xae, 0xb2,
	0x75, 0xfd, 0xa4, 0xba, 0x54, 0xa8, 0xc8, 0x90, 0x69, 0x41, 0x9d, 0x10, 0x3c, 0xc0, 0x38, 0xe6,
	0xf8, 0x64, 0xbf, 0x23, 0xfe, 0x59, 0x71, 0x83, 0xe4, 0x02, 0xfa, 0xce, 0x9b, 0xf5, 0xce, 0x7b,
	0x97, 0xa3, 0x1b, 0x8f, 0x6e, 0xd2, 0xb4, 0xce, 0xd5, 0x76, 0x40, 0xc0, 0xfb, 0x2f, 0x35, 0x5a,
	0x15, 0x86, 0x07, 0x13, 0x18, 0x2f, 0xda, 0xb8, 0xe0, 0x11, 0xbc, 0xc5, 0x4e, 0xac, 0x73, 0x8b,
	0x6b, 0xf0, 0x5f, 0x54, 0xf9, 0xc5, 0xca, 0x8f, 0xd6, 0x84, 0xc7, 0x30, 0x30, 0xc8, 0x75, 0x52,
	0x54, 0xd2, 0x11, 0xf6, 0xa3, 0x43, 0xfb, 0xff, 0x5a, 0xc9, 0xe0, 0x08, 0xa6, 0xed, 0x8a, 0x66,
	0xb0, 0x2b, 0xf0, 0xe6, 0x2c, 0x5b, 0x77, 0xc5, 0xf8, 0x30, 0xd9, 0x8a, 0x37, 0x0c, 0x0a, 0x93,
	0x18, 0x59, 0x89, 0x5d, 0x21, 0x53, 0x20, 0xdb, 0xf9, 0x86, 0x62, 0xd7, 0x86, 0x4a, 0xb7, 0x76,
	0xe4, 0x5b, 0xf2, 0x46, 0x6b, 0x82, 0xcf, 0xd0, 0x77, 0x02, 0x39, 0x83, 0xd1, 0x32, 0x57, 0x29,
	0xcb, 0x13, 0x14, 0x92, 0xbb, 0x2e, 0xbd, 0x08, 0x6a, 0xe9, 0x4d, 0x48, 0x4e, 0x4e, 0x60, 0x68,
	0x9d, 0xc4, 0x36, 0x9e, 0xed, 0x39, 0x7b, 0x60, 0x85, 0x18, 0xb9, 0x9e, 0xdf, 0xbf, 0xdf, 0x2d,
	0x05, 0xae, 0xaa, 0x94, 0x66, 0x4a, 0x86, 0xe6, 0xa7, 0xe0, 0x25, 0xff, 0xfe, 0x7b, 0x27, 0x2c,
	0xd7, 0x2b, 0x16, 0x32, 0x2d, 0xc2, 0xdd, 0x8b, 0x49, 0x0f, 0xdc, 0xb1, 0xdc, 0xfe, 0x06, 0x00,
	0x00, 0xff, 0xff, 0xc3, 0xa1, 0x3e, 0x5e, 0x4c, 0x02, 0x00, 0x00,
}

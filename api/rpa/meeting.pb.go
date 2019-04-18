// Code generated by protoc-gen-go. DO NOT EDIT.
// source: rpa/meeting.proto

package rpa // import "github.com/synerex/synerex_alpha/api/rpa"

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import timestamp "github.com/golang/protobuf/ptypes/timestamp"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type MeetingService struct {
	ProviderId           uint64                 `protobuf:"fixed64,1,opt,name=provider_id,json=providerId,proto3" json:"provider_id,omitempty"`
	Rooms                []*MeetingService_Room `protobuf:"bytes,2,rep,name=rooms,proto3" json:"rooms,omitempty"`
	XXX_NoUnkeyedLiteral struct{}               `json:"-"`
	XXX_unrecognized     []byte                 `json:"-"`
	XXX_sizecache        int32                  `json:"-"`
}

func (m *MeetingService) Reset()         { *m = MeetingService{} }
func (m *MeetingService) String() string { return proto.CompactTextString(m) }
func (*MeetingService) ProtoMessage()    {}
func (*MeetingService) Descriptor() ([]byte, []int) {
	return fileDescriptor_meeting_82b606e469d747c9, []int{0}
}
func (m *MeetingService) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MeetingService.Unmarshal(m, b)
}
func (m *MeetingService) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MeetingService.Marshal(b, m, deterministic)
}
func (dst *MeetingService) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MeetingService.Merge(dst, src)
}
func (m *MeetingService) XXX_Size() int {
	return xxx_messageInfo_MeetingService.Size(m)
}
func (m *MeetingService) XXX_DiscardUnknown() {
	xxx_messageInfo_MeetingService.DiscardUnknown(m)
}

var xxx_messageInfo_MeetingService proto.InternalMessageInfo

func (m *MeetingService) GetProviderId() uint64 {
	if m != nil {
		return m.ProviderId
	}
	return 0
}

func (m *MeetingService) GetRooms() []*MeetingService_Room {
	if m != nil {
		return m.Rooms
	}
	return nil
}

type MeetingService_Room struct {
	RoomId               int32    `protobuf:"varint,1,opt,name=room_id,json=roomId,proto3" json:"room_id,omitempty"`
	Name                 string   `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Capacity             int32    `protobuf:"varint,3,opt,name=capacity,proto3" json:"capacity,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *MeetingService_Room) Reset()         { *m = MeetingService_Room{} }
func (m *MeetingService_Room) String() string { return proto.CompactTextString(m) }
func (*MeetingService_Room) ProtoMessage()    {}
func (*MeetingService_Room) Descriptor() ([]byte, []int) {
	return fileDescriptor_meeting_82b606e469d747c9, []int{0, 0}
}
func (m *MeetingService_Room) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MeetingService_Room.Unmarshal(m, b)
}
func (m *MeetingService_Room) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MeetingService_Room.Marshal(b, m, deterministic)
}
func (dst *MeetingService_Room) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MeetingService_Room.Merge(dst, src)
}
func (m *MeetingService_Room) XXX_Size() int {
	return xxx_messageInfo_MeetingService_Room.Size(m)
}
func (m *MeetingService_Room) XXX_DiscardUnknown() {
	xxx_messageInfo_MeetingService_Room.DiscardUnknown(m)
}

var xxx_messageInfo_MeetingService_Room proto.InternalMessageInfo

func (m *MeetingService_Room) GetRoomId() int32 {
	if m != nil {
		return m.RoomId
	}
	return 0
}

func (m *MeetingService_Room) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *MeetingService_Room) GetCapacity() int32 {
	if m != nil {
		return m.Capacity
	}
	return 0
}

type Booking struct {
	BookingId            uint64               `protobuf:"fixed64,1,opt,name=booking_id,json=bookingId,proto3" json:"booking_id,omitempty"`
	Date                 *timestamp.Timestamp `protobuf:"bytes,2,opt,name=date,proto3" json:"date,omitempty"`
	Start                *timestamp.Timestamp `protobuf:"bytes,3,opt,name=start,proto3" json:"start,omitempty"`
	End                  *timestamp.Timestamp `protobuf:"bytes,4,opt,name=end,proto3" json:"end,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *Booking) Reset()         { *m = Booking{} }
func (m *Booking) String() string { return proto.CompactTextString(m) }
func (*Booking) ProtoMessage()    {}
func (*Booking) Descriptor() ([]byte, []int) {
	return fileDescriptor_meeting_82b606e469d747c9, []int{1}
}
func (m *Booking) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Booking.Unmarshal(m, b)
}
func (m *Booking) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Booking.Marshal(b, m, deterministic)
}
func (dst *Booking) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Booking.Merge(dst, src)
}
func (m *Booking) XXX_Size() int {
	return xxx_messageInfo_Booking.Size(m)
}
func (m *Booking) XXX_DiscardUnknown() {
	xxx_messageInfo_Booking.DiscardUnknown(m)
}

var xxx_messageInfo_Booking proto.InternalMessageInfo

func (m *Booking) GetBookingId() uint64 {
	if m != nil {
		return m.BookingId
	}
	return 0
}

func (m *Booking) GetDate() *timestamp.Timestamp {
	if m != nil {
		return m.Date
	}
	return nil
}

func (m *Booking) GetStart() *timestamp.Timestamp {
	if m != nil {
		return m.Start
	}
	return nil
}

func (m *Booking) GetEnd() *timestamp.Timestamp {
	if m != nil {
		return m.End
	}
	return nil
}

func init() {
	proto.RegisterType((*MeetingService)(nil), "api.meeting.MeetingService")
	proto.RegisterType((*MeetingService_Room)(nil), "api.meeting.MeetingService.Room")
	proto.RegisterType((*Booking)(nil), "api.meeting.Booking")
}

func init() { proto.RegisterFile("rpa/meeting.proto", fileDescriptor_meeting_82b606e469d747c9) }

var fileDescriptor_meeting_82b606e469d747c9 = []byte{
	// 308 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x90, 0xbf, 0x6a, 0xf3, 0x30,
	0x14, 0xc5, 0x71, 0xe2, 0x24, 0x5f, 0xae, 0xe1, 0x83, 0x6a, 0xa9, 0x31, 0x94, 0x98, 0x4c, 0xa6,
	0x14, 0xb9, 0xa4, 0xd0, 0x07, 0xc8, 0xe6, 0xa1, 0x14, 0xdc, 0x4e, 0x5d, 0x82, 0x6c, 0xa9, 0x8e,
	0x68, 0xe4, 0x2b, 0x64, 0x25, 0x34, 0x6f, 0x56, 0xfa, 0x74, 0xc5, 0x52, 0xdc, 0x3f, 0x4b, 0x3b,
	0xe9, 0x9e, 0xa3, 0x9f, 0xce, 0x15, 0x07, 0xce, 0x8c, 0x66, 0xb9, 0x12, 0xc2, 0xca, 0xb6, 0xa1,
	0xda, 0xa0, 0x45, 0x12, 0x31, 0x2d, 0xe9, 0xc9, 0x4a, 0x16, 0x0d, 0x62, 0xb3, 0x13, 0xb9, 0xbb,
	0xaa, 0xf6, 0xcf, 0xb9, 0x95, 0x4a, 0x74, 0x96, 0x29, 0xed, 0xe9, 0xe5, 0x7b, 0x00, 0xff, 0xef,
	0x3c, 0xfc, 0x20, 0xcc, 0x41, 0xd6, 0x82, 0x2c, 0x20, 0xd2, 0x06, 0x0f, 0x92, 0x0b, 0xb3, 0x91,
	0x3c, 0x0e, 0xd2, 0x20, 0x9b, 0x96, 0x30, 0x58, 0x05, 0x27, 0xb7, 0x30, 0x31, 0x88, 0xaa, 0x8b,
	0x47, 0xe9, 0x38, 0x8b, 0x56, 0x29, 0xfd, 0xb6, 0x91, 0xfe, 0x0c, 0xa3, 0x25, 0xa2, 0x2a, 0x3d,
	0x9e, 0xdc, 0x43, 0xd8, 0x4b, 0x72, 0x0e, 0xb3, 0xde, 0x18, 0xc2, 0x27, 0xe5, 0xb4, 0x97, 0x05,
	0x27, 0x04, 0xc2, 0x96, 0x29, 0x11, 0x8f, 0xd2, 0x20, 0x9b, 0x97, 0x6e, 0x26, 0x09, 0xfc, 0xab,
	0x99, 0x66, 0xb5, 0xb4, 0xc7, 0x78, 0xec, 0xe8, 0x4f, 0xbd, 0x7c, 0x0b, 0x60, 0xb6, 0x46, 0x7c,
	0x91, 0x6d, 0x43, 0x2e, 0x00, 0x2a, 0x3f, 0x7e, 0x7d, 0x7a, 0x7e, 0x72, 0x0a, 0x4e, 0x28, 0x84,
	0x9c, 0x59, 0x1f, 0x1d, 0xad, 0x12, 0xea, 0x7b, 0xa1, 0x43, 0x2f, 0xf4, 0x71, 0xe8, 0xa5, 0x74,
	0x1c, 0xb9, 0x86, 0x49, 0x67, 0x99, 0xb1, 0x6e, 0xe7, 0xef, 0x0f, 0x3c, 0x48, 0xae, 0x60, 0x2c,
	0x5a, 0x1e, 0x87, 0x7f, 0xf2, 0x3d, 0xb6, 0xbe, 0x7c, 0xca, 0x1a, 0x69, 0xb7, 0xfb, 0x8a, 0xd6,
	0xa8, 0xf2, 0xee, 0xd8, 0x0a, 0x23, 0x5e, 0x87, 0x73, 0xc3, 0x76, 0x7a, 0xcb, 0x72, 0xa6, 0x65,
	0x6e, 0x34, 0xab, 0xa6, 0x2e, 0xe4, 0xe6, 0x23, 0x00, 0x00, 0xff, 0xff, 0xfb, 0x02, 0xba, 0x6f,
	0xed, 0x01, 0x00, 0x00,
}

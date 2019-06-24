// Code generated by protoc-gen-go. DO NOT EDIT.
// source: rpa/meeting.proto

package rpa // import "github.com/synerex/synerex_alpha/api/rpa"

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

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
	Cid                  string   `protobuf:"bytes,1,opt,name=cid,proto3" json:"cid,omitempty"`
	Status               string   `protobuf:"bytes,2,opt,name=status,proto3" json:"status,omitempty"`
	Year                 string   `protobuf:"bytes,3,opt,name=year,proto3" json:"year,omitempty"`
	Month                string   `protobuf:"bytes,4,opt,name=month,proto3" json:"month,omitempty"`
	Day                  string   `protobuf:"bytes,5,opt,name=day,proto3" json:"day,omitempty"`
	Week                 string   `protobuf:"bytes,6,opt,name=week,proto3" json:"week,omitempty"`
	Start                string   `protobuf:"bytes,7,opt,name=start,proto3" json:"start,omitempty"`
	End                  string   `protobuf:"bytes,8,opt,name=end,proto3" json:"end,omitempty"`
	People               string   `protobuf:"bytes,9,opt,name=people,proto3" json:"people,omitempty"`
	Title                string   `protobuf:"bytes,10,opt,name=title,proto3" json:"title,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *MeetingService) Reset()         { *m = MeetingService{} }
func (m *MeetingService) String() string { return proto.CompactTextString(m) }
func (*MeetingService) ProtoMessage()    {}
func (*MeetingService) Descriptor() ([]byte, []int) {
	return fileDescriptor_meeting_a072459768f6fb8f, []int{0}
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

func (m *MeetingService) GetCid() string {
	if m != nil {
		return m.Cid
	}
	return ""
}

func (m *MeetingService) GetStatus() string {
	if m != nil {
		return m.Status
	}
	return ""
}

func (m *MeetingService) GetYear() string {
	if m != nil {
		return m.Year
	}
	return ""
}

func (m *MeetingService) GetMonth() string {
	if m != nil {
		return m.Month
	}
	return ""
}

func (m *MeetingService) GetDay() string {
	if m != nil {
		return m.Day
	}
	return ""
}

func (m *MeetingService) GetWeek() string {
	if m != nil {
		return m.Week
	}
	return ""
}

func (m *MeetingService) GetStart() string {
	if m != nil {
		return m.Start
	}
	return ""
}

func (m *MeetingService) GetEnd() string {
	if m != nil {
		return m.End
	}
	return ""
}

func (m *MeetingService) GetPeople() string {
	if m != nil {
		return m.People
	}
	return ""
}

func (m *MeetingService) GetTitle() string {
	if m != nil {
		return m.Title
	}
	return ""
}

func init() {
	proto.RegisterType((*MeetingService)(nil), "api.meeting.MeetingService")
}

func init() { proto.RegisterFile("rpa/meeting.proto", fileDescriptor_meeting_a072459768f6fb8f) }

var fileDescriptor_meeting_a072459768f6fb8f = []byte{
	// 227 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x34, 0x90, 0xb1, 0x4a, 0xc4, 0x40,
	0x10, 0x86, 0x89, 0x77, 0x17, 0xbd, 0x15, 0x44, 0x17, 0x91, 0x29, 0xc5, 0xea, 0xb0, 0x48, 0x0a,
	0xdf, 0xc0, 0xde, 0x46, 0x3b, 0x1b, 0x99, 0x4b, 0x86, 0xcb, 0x62, 0xb2, 0x3b, 0x6c, 0xe6, 0xd4,
	0xbc, 0xad, 0x8f, 0x22, 0x3b, 0x73, 0x56, 0xfb, 0x7f, 0x3f, 0x7c, 0xc3, 0xcf, 0xba, 0x9b, 0xcc,
	0xd8, 0x4e, 0x44, 0x12, 0xe2, 0xa1, 0xe1, 0x9c, 0x24, 0xf9, 0x4b, 0xe4, 0xd0, 0x9c, 0xaa, 0x87,
	0xdf, 0xca, 0x5d, 0xbd, 0x58, 0x7e, 0xa3, 0xfc, 0x15, 0x3a, 0xf2, 0xd7, 0x6e, 0xd5, 0x85, 0x1e,
	0xaa, 0xfb, 0x6a, 0xb7, 0x7d, 0x2d, 0xd1, 0xdf, 0xb9, 0x7a, 0x16, 0x94, 0xe3, 0x0c, 0x67, 0x5a,
	0x9e, 0xc8, 0x7b, 0xb7, 0x5e, 0x08, 0x33, 0xac, 0xb4, 0xd5, 0xec, 0x6f, 0xdd, 0x66, 0x4a, 0x51,
	0x06, 0x58, 0x6b, 0x69, 0x50, 0x6e, 0xf6, 0xb8, 0xc0, 0xc6, 0x6e, 0xf6, 0xb8, 0x14, 0xf7, 0x9b,
	0xe8, 0x13, 0x6a, 0x73, 0x4b, 0x2e, 0xee, 0x2c, 0x98, 0x05, 0xce, 0xcd, 0x55, 0x28, 0x2e, 0xc5,
	0x1e, 0x2e, 0xcc, 0xa5, 0xa8, 0x7b, 0x98, 0x12, 0x8f, 0x04, 0x5b, 0xdb, 0x63, 0x54, 0x7c, 0x09,
	0x32, 0x12, 0x38, 0xf3, 0x15, 0x9e, 0x1f, 0xdf, 0x77, 0x87, 0x20, 0xc3, 0x71, 0xdf, 0x74, 0x69,
	0x6a, 0xe7, 0x25, 0x52, 0xa6, 0x9f, 0xff, 0xf7, 0x03, 0x47, 0x1e, 0xb0, 0x45, 0x0e, 0x6d, 0x66,
	0xdc, 0xd7, 0xfa, 0x45, 0x4f, 0x7f, 0x01, 0x00, 0x00, 0xff, 0xff, 0xee, 0xba, 0x8a, 0xb2, 0x37,
	0x01, 0x00, 0x00,
}

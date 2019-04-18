// Code generated by protoc-gen-go. DO NOT EDIT.
// source: rpa/user.proto

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

type User struct {
	UserId               uint64               `protobuf:"fixed64,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	Name                 string               `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	LastUpdated          *timestamp.Timestamp `protobuf:"bytes,3,opt,name=last_updated,json=lastUpdated,proto3" json:"last_updated,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *User) Reset()         { *m = User{} }
func (m *User) String() string { return proto.CompactTextString(m) }
func (*User) ProtoMessage()    {}
func (*User) Descriptor() ([]byte, []int) {
	return fileDescriptor_user_f4ff96a69f088c75, []int{0}
}
func (m *User) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_User.Unmarshal(m, b)
}
func (m *User) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_User.Marshal(b, m, deterministic)
}
func (dst *User) XXX_Merge(src proto.Message) {
	xxx_messageInfo_User.Merge(dst, src)
}
func (m *User) XXX_Size() int {
	return xxx_messageInfo_User.Size(m)
}
func (m *User) XXX_DiscardUnknown() {
	xxx_messageInfo_User.DiscardUnknown(m)
}

var xxx_messageInfo_User proto.InternalMessageInfo

func (m *User) GetUserId() uint64 {
	if m != nil {
		return m.UserId
	}
	return 0
}

func (m *User) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *User) GetLastUpdated() *timestamp.Timestamp {
	if m != nil {
		return m.LastUpdated
	}
	return nil
}

func init() {
	proto.RegisterType((*User)(nil), "api.user.User")
}

func init() { proto.RegisterFile("rpa/user.proto", fileDescriptor_user_f4ff96a69f088c75) }

var fileDescriptor_user_f4ff96a69f088c75 = []byte{
	// 198 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x34, 0x8d, 0xbf, 0x4b, 0xc6, 0x30,
	0x10, 0x86, 0x89, 0x7e, 0x54, 0xcd, 0x27, 0x0e, 0x59, 0x2c, 0x5d, 0x2c, 0x4e, 0xc1, 0x21, 0x01,
	0x9d, 0x5d, 0xdc, 0x5c, 0x8b, 0x5d, 0x5c, 0xca, 0xd5, 0x9c, 0x6d, 0xa0, 0x69, 0x8e, 0xfc, 0x00,
	0xfd, 0xef, 0x3f, 0xd2, 0xd2, 0xe9, 0xee, 0x3d, 0x9e, 0x7b, 0x1f, 0xfe, 0x10, 0x08, 0x74, 0x8e,
	0x18, 0x14, 0x05, 0x9f, 0xbc, 0xb8, 0x05, 0xb2, 0xaa, 0xe4, 0xe6, 0x69, 0xf2, 0x7e, 0x5a, 0x50,
	0x6f, 0xf7, 0x31, 0xff, 0xea, 0x64, 0x1d, 0xc6, 0x04, 0x8e, 0x76, 0xf4, 0x39, 0xf0, 0x53, 0x1f,
	0x31, 0x88, 0x47, 0x7e, 0x53, 0x1e, 0x06, 0x6b, 0x6a, 0xd6, 0x32, 0x59, 0x75, 0x55, 0x89, 0x9f,
	0x46, 0x08, 0x7e, 0x5a, 0xc1, 0x61, 0x7d, 0xd5, 0x32, 0x79, 0xd7, 0x6d, 0xbb, 0x78, 0xe7, 0xf7,
	0x0b, 0xc4, 0x34, 0x64, 0x32, 0x90, 0xd0, 0xd4, 0xd7, 0x2d, 0x93, 0xe7, 0xd7, 0x46, 0xed, 0x32,
	0x75, 0xc8, 0xd4, 0xd7, 0x21, 0xeb, 0xce, 0x85, 0xef, 0x77, 0xfc, 0xe3, 0xe5, 0x5b, 0x4e, 0x36,
	0xcd, 0x79, 0x54, 0x3f, 0xde, 0xe9, 0xf8, 0xbf, 0x62, 0xc0, 0xbf, 0x63, 0x0e, 0xb0, 0xd0, 0x0c,
	0x1a, 0xc8, 0xea, 0x40, 0x30, 0x56, 0x5b, 0xd9, 0xdb, 0x25, 0x00, 0x00, 0xff, 0xff, 0xe8, 0xce,
	0x9a, 0x93, 0xe3, 0x00, 0x00, 0x00,
}

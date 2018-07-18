// Code generated by protoc-gen-go. DO NOT EDIT.
// source: policytap.proto

/*
Package dnstap is a generated protocol buffer package.

It is generated from these files:
	policytap.proto

It has these top-level messages:
	DnstapAttribute
	Extra
*/
package dnstap

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

// PolicyAction defines available policy action type
type PolicyAction int32

const (
	PolicyAction_INVALID     PolicyAction = 0
	PolicyAction_DENY        PolicyAction = 1
	PolicyAction_PASSTHROUGH PolicyAction = 2
	PolicyAction_NXDOMAIN    PolicyAction = 3
	PolicyAction_REDIRECT    PolicyAction = 4
	PolicyAction_REFUSE      PolicyAction = 5
)

var PolicyAction_name = map[int32]string{
	0: "INVALID",
	1: "DENY",
	2: "PASSTHROUGH",
	3: "NXDOMAIN",
	4: "REDIRECT",
	5: "REFUSE",
}
var PolicyAction_value = map[string]int32{
	"INVALID":     0,
	"DENY":        1,
	"PASSTHROUGH": 2,
	"NXDOMAIN":    3,
	"REDIRECT":    4,
	"REFUSE":      5,
}

func (x PolicyAction) String() string {
	return proto.EnumName(PolicyAction_name, int32(x))
}
func (PolicyAction) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

type DnstapAttribute struct {
	Id    string `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	Value string `protobuf:"bytes,2,opt,name=value" json:"value,omitempty"`
}

func (m *DnstapAttribute) Reset()                    { *m = DnstapAttribute{} }
func (m *DnstapAttribute) String() string            { return proto.CompactTextString(m) }
func (*DnstapAttribute) ProtoMessage()               {}
func (*DnstapAttribute) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *DnstapAttribute) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *DnstapAttribute) GetValue() string {
	if m != nil {
		return m.Value
	}
	return ""
}

type Extra struct {
	Attrs []*DnstapAttribute `protobuf:"bytes,1,rep,name=attrs" json:"attrs,omitempty"`
}

func (m *Extra) Reset()                    { *m = Extra{} }
func (m *Extra) String() string            { return proto.CompactTextString(m) }
func (*Extra) ProtoMessage()               {}
func (*Extra) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *Extra) GetAttrs() []*DnstapAttribute {
	if m != nil {
		return m.Attrs
	}
	return nil
}

func init() {
	proto.RegisterType((*DnstapAttribute)(nil), "dnstap.DnstapAttribute")
	proto.RegisterType((*Extra)(nil), "dnstap.Extra")
	proto.RegisterEnum("dnstap.PolicyAction", PolicyAction_name, PolicyAction_value)
}

func init() { proto.RegisterFile("policytap.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 227 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x5c, 0x8f, 0xc1, 0x4b, 0xc3, 0x30,
	0x14, 0xc6, 0x6d, 0xb7, 0xd6, 0xf9, 0x3a, 0x6c, 0x78, 0x08, 0xf6, 0x38, 0x76, 0x1a, 0x82, 0x3d,
	0x28, 0xe8, 0x39, 0x98, 0xe8, 0x0a, 0xda, 0x8d, 0x74, 0x13, 0xbd, 0x08, 0xd9, 0xba, 0x43, 0x60,
	0x2c, 0x21, 0x7b, 0x15, 0xfd, 0xef, 0xc5, 0xf4, 0xb6, 0xdb, 0xfb, 0x3d, 0x3e, 0x7e, 0x7c, 0x1f,
	0xe4, 0xce, 0xee, 0xcd, 0xf6, 0x97, 0xb4, 0x2b, 0x9d, 0xb7, 0x64, 0x31, 0x6d, 0x0f, 0x47, 0xd2,
	0x6e, 0xfa, 0x08, 0xb9, 0x08, 0x17, 0x27, 0xf2, 0x66, 0xd3, 0xd1, 0x0e, 0x2f, 0x21, 0x36, 0x6d,
	0x11, 0x4d, 0xa2, 0xd9, 0x85, 0x8a, 0x4d, 0x8b, 0x57, 0x90, 0x7c, 0xeb, 0x7d, 0xb7, 0x2b, 0xe2,
	0xf0, 0xea, 0x61, 0xfa, 0x00, 0x89, 0xfc, 0x21, 0xaf, 0xf1, 0x16, 0x12, 0x4d, 0xe4, 0x8f, 0x45,
	0x34, 0x19, 0xcc, 0xb2, 0xbb, 0xeb, 0xb2, 0x37, 0x97, 0x27, 0x5a, 0xd5, 0xa7, 0x6e, 0xbe, 0x60,
	0xbc, 0x0c, 0x5d, 0xf8, 0x96, 0x8c, 0x3d, 0x60, 0x06, 0xe7, 0x55, 0xfd, 0xce, 0x5f, 0x2b, 0xc1,
	0xce, 0x70, 0x04, 0x43, 0x21, 0xeb, 0x4f, 0x16, 0x61, 0x0e, 0xd9, 0x92, 0x37, 0xcd, 0x6a, 0xae,
	0x16, 0xeb, 0x97, 0x39, 0x8b, 0x71, 0x0c, 0xa3, 0xfa, 0x43, 0x2c, 0xde, 0x78, 0x55, 0xb3, 0xc1,
	0x3f, 0x29, 0x29, 0x2a, 0x25, 0x9f, 0x56, 0x6c, 0x88, 0x00, 0xa9, 0x92, 0xcf, 0xeb, 0x46, 0xb2,
	0x64, 0x93, 0x86, 0x7d, 0xf7, 0x7f, 0x01, 0x00, 0x00, 0xff, 0xff, 0xf1, 0xf6, 0x4f, 0x5b, 0xf2,
	0x00, 0x00, 0x00,
}

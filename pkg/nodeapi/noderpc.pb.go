// Code generated by protoc-gen-go. DO NOT EDIT.
// source: noderpc.proto

package nodeapi

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
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

type AddPodAnnotationRequest struct {
	Namespace            string   `protobuf:"bytes,1,opt,name=namespace,proto3" json:"namespace,omitempty"`
	Name                 string   `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Key                  string   `protobuf:"bytes,3,opt,name=key,proto3" json:"key,omitempty"`
	Value                string   `protobuf:"bytes,4,opt,name=value,proto3" json:"value,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *AddPodAnnotationRequest) Reset()         { *m = AddPodAnnotationRequest{} }
func (m *AddPodAnnotationRequest) String() string { return proto.CompactTextString(m) }
func (*AddPodAnnotationRequest) ProtoMessage()    {}
func (*AddPodAnnotationRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_noderpc_23fc216ecb65138b, []int{0}
}
func (m *AddPodAnnotationRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AddPodAnnotationRequest.Unmarshal(m, b)
}
func (m *AddPodAnnotationRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AddPodAnnotationRequest.Marshal(b, m, deterministic)
}
func (dst *AddPodAnnotationRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AddPodAnnotationRequest.Merge(dst, src)
}
func (m *AddPodAnnotationRequest) XXX_Size() int {
	return xxx_messageInfo_AddPodAnnotationRequest.Size(m)
}
func (m *AddPodAnnotationRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_AddPodAnnotationRequest.DiscardUnknown(m)
}

var xxx_messageInfo_AddPodAnnotationRequest proto.InternalMessageInfo

func (m *AddPodAnnotationRequest) GetNamespace() string {
	if m != nil {
		return m.Namespace
	}
	return ""
}

func (m *AddPodAnnotationRequest) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *AddPodAnnotationRequest) GetKey() string {
	if m != nil {
		return m.Key
	}
	return ""
}

func (m *AddPodAnnotationRequest) GetValue() string {
	if m != nil {
		return m.Value
	}
	return ""
}

type AddPodAnnotationReply struct {
	Error                string   `protobuf:"bytes,1,opt,name=error,proto3" json:"error,omitempty"`
	Metav1StatusReason   string   `protobuf:"bytes,2,opt,name=metav1_status_reason,json=metav1StatusReason,proto3" json:"metav1_status_reason,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *AddPodAnnotationReply) Reset()         { *m = AddPodAnnotationReply{} }
func (m *AddPodAnnotationReply) String() string { return proto.CompactTextString(m) }
func (*AddPodAnnotationReply) ProtoMessage()    {}
func (*AddPodAnnotationReply) Descriptor() ([]byte, []int) {
	return fileDescriptor_noderpc_23fc216ecb65138b, []int{1}
}
func (m *AddPodAnnotationReply) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AddPodAnnotationReply.Unmarshal(m, b)
}
func (m *AddPodAnnotationReply) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AddPodAnnotationReply.Marshal(b, m, deterministic)
}
func (dst *AddPodAnnotationReply) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AddPodAnnotationReply.Merge(dst, src)
}
func (m *AddPodAnnotationReply) XXX_Size() int {
	return xxx_messageInfo_AddPodAnnotationReply.Size(m)
}
func (m *AddPodAnnotationReply) XXX_DiscardUnknown() {
	xxx_messageInfo_AddPodAnnotationReply.DiscardUnknown(m)
}

var xxx_messageInfo_AddPodAnnotationReply proto.InternalMessageInfo

func (m *AddPodAnnotationReply) GetError() string {
	if m != nil {
		return m.Error
	}
	return ""
}

func (m *AddPodAnnotationReply) GetMetav1StatusReason() string {
	if m != nil {
		return m.Metav1StatusReason
	}
	return ""
}

type DeletePodAnnotationRequest struct {
	Namespace            string   `protobuf:"bytes,1,opt,name=namespace,proto3" json:"namespace,omitempty"`
	Name                 string   `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Key                  string   `protobuf:"bytes,3,opt,name=key,proto3" json:"key,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *DeletePodAnnotationRequest) Reset()         { *m = DeletePodAnnotationRequest{} }
func (m *DeletePodAnnotationRequest) String() string { return proto.CompactTextString(m) }
func (*DeletePodAnnotationRequest) ProtoMessage()    {}
func (*DeletePodAnnotationRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_noderpc_23fc216ecb65138b, []int{2}
}
func (m *DeletePodAnnotationRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DeletePodAnnotationRequest.Unmarshal(m, b)
}
func (m *DeletePodAnnotationRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DeletePodAnnotationRequest.Marshal(b, m, deterministic)
}
func (dst *DeletePodAnnotationRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DeletePodAnnotationRequest.Merge(dst, src)
}
func (m *DeletePodAnnotationRequest) XXX_Size() int {
	return xxx_messageInfo_DeletePodAnnotationRequest.Size(m)
}
func (m *DeletePodAnnotationRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_DeletePodAnnotationRequest.DiscardUnknown(m)
}

var xxx_messageInfo_DeletePodAnnotationRequest proto.InternalMessageInfo

func (m *DeletePodAnnotationRequest) GetNamespace() string {
	if m != nil {
		return m.Namespace
	}
	return ""
}

func (m *DeletePodAnnotationRequest) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *DeletePodAnnotationRequest) GetKey() string {
	if m != nil {
		return m.Key
	}
	return ""
}

type DeletePodAnnotationReply struct {
	Error                string   `protobuf:"bytes,1,opt,name=error,proto3" json:"error,omitempty"`
	Metav1StatusReason   string   `protobuf:"bytes,2,opt,name=metav1_status_reason,json=metav1StatusReason,proto3" json:"metav1_status_reason,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *DeletePodAnnotationReply) Reset()         { *m = DeletePodAnnotationReply{} }
func (m *DeletePodAnnotationReply) String() string { return proto.CompactTextString(m) }
func (*DeletePodAnnotationReply) ProtoMessage()    {}
func (*DeletePodAnnotationReply) Descriptor() ([]byte, []int) {
	return fileDescriptor_noderpc_23fc216ecb65138b, []int{3}
}
func (m *DeletePodAnnotationReply) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DeletePodAnnotationReply.Unmarshal(m, b)
}
func (m *DeletePodAnnotationReply) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DeletePodAnnotationReply.Marshal(b, m, deterministic)
}
func (dst *DeletePodAnnotationReply) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DeletePodAnnotationReply.Merge(dst, src)
}
func (m *DeletePodAnnotationReply) XXX_Size() int {
	return xxx_messageInfo_DeletePodAnnotationReply.Size(m)
}
func (m *DeletePodAnnotationReply) XXX_DiscardUnknown() {
	xxx_messageInfo_DeletePodAnnotationReply.DiscardUnknown(m)
}

var xxx_messageInfo_DeletePodAnnotationReply proto.InternalMessageInfo

func (m *DeletePodAnnotationReply) GetError() string {
	if m != nil {
		return m.Error
	}
	return ""
}

func (m *DeletePodAnnotationReply) GetMetav1StatusReason() string {
	if m != nil {
		return m.Metav1StatusReason
	}
	return ""
}

func init() {
	proto.RegisterType((*AddPodAnnotationRequest)(nil), "nodeapi.AddPodAnnotationRequest")
	proto.RegisterType((*AddPodAnnotationReply)(nil), "nodeapi.AddPodAnnotationReply")
	proto.RegisterType((*DeletePodAnnotationRequest)(nil), "nodeapi.DeletePodAnnotationRequest")
	proto.RegisterType((*DeletePodAnnotationReply)(nil), "nodeapi.DeletePodAnnotationReply")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// MidoNetKubeNodeClient is the client API for MidoNetKubeNode service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type MidoNetKubeNodeClient interface {
	AddPodAnnotation(ctx context.Context, in *AddPodAnnotationRequest, opts ...grpc.CallOption) (*AddPodAnnotationReply, error)
	DeletePodAnnotation(ctx context.Context, in *DeletePodAnnotationRequest, opts ...grpc.CallOption) (*DeletePodAnnotationReply, error)
}

type midoNetKubeNodeClient struct {
	cc *grpc.ClientConn
}

func NewMidoNetKubeNodeClient(cc *grpc.ClientConn) MidoNetKubeNodeClient {
	return &midoNetKubeNodeClient{cc}
}

func (c *midoNetKubeNodeClient) AddPodAnnotation(ctx context.Context, in *AddPodAnnotationRequest, opts ...grpc.CallOption) (*AddPodAnnotationReply, error) {
	out := new(AddPodAnnotationReply)
	err := c.cc.Invoke(ctx, "/nodeapi.MidoNetKubeNode/AddPodAnnotation", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *midoNetKubeNodeClient) DeletePodAnnotation(ctx context.Context, in *DeletePodAnnotationRequest, opts ...grpc.CallOption) (*DeletePodAnnotationReply, error) {
	out := new(DeletePodAnnotationReply)
	err := c.cc.Invoke(ctx, "/nodeapi.MidoNetKubeNode/DeletePodAnnotation", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MidoNetKubeNodeServer is the server API for MidoNetKubeNode service.
type MidoNetKubeNodeServer interface {
	AddPodAnnotation(context.Context, *AddPodAnnotationRequest) (*AddPodAnnotationReply, error)
	DeletePodAnnotation(context.Context, *DeletePodAnnotationRequest) (*DeletePodAnnotationReply, error)
}

func RegisterMidoNetKubeNodeServer(s *grpc.Server, srv MidoNetKubeNodeServer) {
	s.RegisterService(&_MidoNetKubeNode_serviceDesc, srv)
}

func _MidoNetKubeNode_AddPodAnnotation_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddPodAnnotationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MidoNetKubeNodeServer).AddPodAnnotation(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/nodeapi.MidoNetKubeNode/AddPodAnnotation",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MidoNetKubeNodeServer).AddPodAnnotation(ctx, req.(*AddPodAnnotationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MidoNetKubeNode_DeletePodAnnotation_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeletePodAnnotationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MidoNetKubeNodeServer).DeletePodAnnotation(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/nodeapi.MidoNetKubeNode/DeletePodAnnotation",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MidoNetKubeNodeServer).DeletePodAnnotation(ctx, req.(*DeletePodAnnotationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _MidoNetKubeNode_serviceDesc = grpc.ServiceDesc{
	ServiceName: "nodeapi.MidoNetKubeNode",
	HandlerType: (*MidoNetKubeNodeServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "AddPodAnnotation",
			Handler:    _MidoNetKubeNode_AddPodAnnotation_Handler,
		},
		{
			MethodName: "DeletePodAnnotation",
			Handler:    _MidoNetKubeNode_DeletePodAnnotation_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "noderpc.proto",
}

func init() { proto.RegisterFile("noderpc.proto", fileDescriptor_noderpc_23fc216ecb65138b) }

var fileDescriptor_noderpc_23fc216ecb65138b = []byte{
	// 269 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xb4, 0x52, 0xcd, 0x4a, 0x03, 0x31,
	0x10, 0xb6, 0xb6, 0x2a, 0x1d, 0x10, 0xcb, 0x58, 0x31, 0x2c, 0x22, 0x35, 0x5e, 0x3c, 0x2d, 0xfe,
	0x3c, 0x41, 0xc1, 0x9b, 0x58, 0x64, 0x05, 0xaf, 0x6b, 0xb6, 0x99, 0xc3, 0xe2, 0x36, 0x13, 0x93,
	0x6c, 0xa1, 0xaf, 0xe8, 0x53, 0xc9, 0x66, 0x17, 0x05, 0x6d, 0x7b, 0xd2, 0xdb, 0x7c, 0x3f, 0xe1,
	0xe3, 0x9b, 0x09, 0x1c, 0x1a, 0xd6, 0xe4, 0xec, 0x3c, 0xb5, 0x8e, 0x03, 0xe3, 0x41, 0x03, 0x95,
	0x2d, 0xa5, 0x87, 0xd3, 0xa9, 0xd6, 0x4f, 0xac, 0xa7, 0xc6, 0x70, 0x50, 0xa1, 0x64, 0x93, 0xd1,
	0x7b, 0x4d, 0x3e, 0xe0, 0x19, 0x0c, 0x8d, 0x5a, 0x90, 0xb7, 0x6a, 0x4e, 0xa2, 0x37, 0xe9, 0x5d,
	0x0d, 0xb3, 0x6f, 0x02, 0x11, 0x06, 0x0d, 0x10, 0xbb, 0x51, 0x88, 0x33, 0x8e, 0xa0, 0xff, 0x46,
	0x2b, 0xd1, 0x8f, 0x54, 0x33, 0xe2, 0x18, 0xf6, 0x96, 0xaa, 0xaa, 0x49, 0x0c, 0x22, 0xd7, 0x02,
	0x99, 0xc3, 0xc9, 0xef, 0x50, 0x5b, 0x45, 0x3b, 0x39, 0xc7, 0xae, 0x8b, 0x6b, 0x01, 0x5e, 0xc3,
	0x78, 0x41, 0x41, 0x2d, 0x6f, 0x72, 0x1f, 0x54, 0xa8, 0x7d, 0xee, 0x48, 0x79, 0x36, 0x5d, 0x34,
	0xb6, 0xda, 0x73, 0x94, 0xb2, 0xa8, 0xc8, 0x57, 0x48, 0xee, 0xa9, 0xa2, 0x40, 0xff, 0x55, 0x4c,
	0x16, 0x20, 0xd6, 0x26, 0xfc, 0x61, 0x8b, 0xdb, 0x8f, 0x1e, 0x1c, 0x3d, 0x96, 0x9a, 0x67, 0x14,
	0x1e, 0xea, 0x82, 0x66, 0xac, 0x09, 0x5f, 0x60, 0xf4, 0x73, 0x75, 0x38, 0x49, 0xbb, 0x6b, 0xa6,
	0x1b, 0x4e, 0x99, 0x9c, 0x6f, 0x71, 0xd8, 0x6a, 0x25, 0x77, 0x30, 0x87, 0xe3, 0x35, 0x7d, 0xf0,
	0xf2, 0xeb, 0xe1, 0xe6, 0x7d, 0x26, 0x17, 0xdb, 0x4d, 0x31, 0xa0, 0xd8, 0x8f, 0x1f, 0xef, 0xee,
	0x33, 0x00, 0x00, 0xff, 0xff, 0x44, 0xa0, 0x0e, 0x78, 0x89, 0x02, 0x00, 0x00,
}

// Code generated by protoc-gen-go. DO NOT EDIT.
// source: cluster.proto

package proto

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
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
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type PingRequest struct {
	Payload              []byte   `protobuf:"bytes,1,opt,name=payload,proto3" json:"payload,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PingRequest) Reset()         { *m = PingRequest{} }
func (m *PingRequest) String() string { return proto.CompactTextString(m) }
func (*PingRequest) ProtoMessage()    {}
func (*PingRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_3cfb3b8ec240c376, []int{0}
}

func (m *PingRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PingRequest.Unmarshal(m, b)
}
func (m *PingRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PingRequest.Marshal(b, m, deterministic)
}
func (m *PingRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PingRequest.Merge(m, src)
}
func (m *PingRequest) XXX_Size() int {
	return xxx_messageInfo_PingRequest.Size(m)
}
func (m *PingRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_PingRequest.DiscardUnknown(m)
}

var xxx_messageInfo_PingRequest proto.InternalMessageInfo

func (m *PingRequest) GetPayload() []byte {
	if m != nil {
		return m.Payload
	}
	return nil
}

type PingResponse struct {
	Payload              []byte   `protobuf:"bytes,1,opt,name=payload,proto3" json:"payload,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PingResponse) Reset()         { *m = PingResponse{} }
func (m *PingResponse) String() string { return proto.CompactTextString(m) }
func (*PingResponse) ProtoMessage()    {}
func (*PingResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_3cfb3b8ec240c376, []int{1}
}

func (m *PingResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PingResponse.Unmarshal(m, b)
}
func (m *PingResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PingResponse.Marshal(b, m, deterministic)
}
func (m *PingResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PingResponse.Merge(m, src)
}
func (m *PingResponse) XXX_Size() int {
	return xxx_messageInfo_PingResponse.Size(m)
}
func (m *PingResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_PingResponse.DiscardUnknown(m)
}

var xxx_messageInfo_PingResponse proto.InternalMessageInfo

func (m *PingResponse) GetPayload() []byte {
	if m != nil {
		return m.Payload
	}
	return nil
}

type Node struct {
	ConnectionString     string   `protobuf:"bytes,1,opt,name=connectionString,proto3" json:"connectionString,omitempty"`
	Uuid                 string   `protobuf:"bytes,2,opt,name=uuid,proto3" json:"uuid,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Node) Reset()         { *m = Node{} }
func (m *Node) String() string { return proto.CompactTextString(m) }
func (*Node) ProtoMessage()    {}
func (*Node) Descriptor() ([]byte, []int) {
	return fileDescriptor_3cfb3b8ec240c376, []int{2}
}

func (m *Node) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Node.Unmarshal(m, b)
}
func (m *Node) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Node.Marshal(b, m, deterministic)
}
func (m *Node) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Node.Merge(m, src)
}
func (m *Node) XXX_Size() int {
	return xxx_messageInfo_Node.Size(m)
}
func (m *Node) XXX_DiscardUnknown() {
	xxx_messageInfo_Node.DiscardUnknown(m)
}

var xxx_messageInfo_Node proto.InternalMessageInfo

func (m *Node) GetConnectionString() string {
	if m != nil {
		return m.ConnectionString
	}
	return ""
}

func (m *Node) GetUuid() string {
	if m != nil {
		return m.Uuid
	}
	return ""
}

type HelloRequest struct {
	Iam                  *Node    `protobuf:"bytes,1,opt,name=iam,proto3" json:"iam,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *HelloRequest) Reset()         { *m = HelloRequest{} }
func (m *HelloRequest) String() string { return proto.CompactTextString(m) }
func (*HelloRequest) ProtoMessage()    {}
func (*HelloRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_3cfb3b8ec240c376, []int{3}
}

func (m *HelloRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_HelloRequest.Unmarshal(m, b)
}
func (m *HelloRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_HelloRequest.Marshal(b, m, deterministic)
}
func (m *HelloRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_HelloRequest.Merge(m, src)
}
func (m *HelloRequest) XXX_Size() int {
	return xxx_messageInfo_HelloRequest.Size(m)
}
func (m *HelloRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_HelloRequest.DiscardUnknown(m)
}

var xxx_messageInfo_HelloRequest proto.InternalMessageInfo

func (m *HelloRequest) GetIam() *Node {
	if m != nil {
		return m.Iam
	}
	return nil
}

type AliveNodesResponse struct {
	AliveNodes           []*Node  `protobuf:"bytes,1,rep,name=aliveNodes,proto3" json:"aliveNodes,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *AliveNodesResponse) Reset()         { *m = AliveNodesResponse{} }
func (m *AliveNodesResponse) String() string { return proto.CompactTextString(m) }
func (*AliveNodesResponse) ProtoMessage()    {}
func (*AliveNodesResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_3cfb3b8ec240c376, []int{4}
}

func (m *AliveNodesResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AliveNodesResponse.Unmarshal(m, b)
}
func (m *AliveNodesResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AliveNodesResponse.Marshal(b, m, deterministic)
}
func (m *AliveNodesResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AliveNodesResponse.Merge(m, src)
}
func (m *AliveNodesResponse) XXX_Size() int {
	return xxx_messageInfo_AliveNodesResponse.Size(m)
}
func (m *AliveNodesResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_AliveNodesResponse.DiscardUnknown(m)
}

var xxx_messageInfo_AliveNodesResponse proto.InternalMessageInfo

func (m *AliveNodesResponse) GetAliveNodes() []*Node {
	if m != nil {
		return m.AliveNodes
	}
	return nil
}

func init() {
	proto.RegisterType((*PingRequest)(nil), "proto.PingRequest")
	proto.RegisterType((*PingResponse)(nil), "proto.PingResponse")
	proto.RegisterType((*Node)(nil), "proto.Node")
	proto.RegisterType((*HelloRequest)(nil), "proto.HelloRequest")
	proto.RegisterType((*AliveNodesResponse)(nil), "proto.AliveNodesResponse")
}

func init() {
	proto.RegisterFile("cluster.proto", fileDescriptor_3cfb3b8ec240c376)
}

var fileDescriptor_3cfb3b8ec240c376 = []byte{
	// 431 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x91, 0x51, 0x6b, 0xd4, 0x40,
	0x14, 0x85, 0x59, 0xbb, 0xdb, 0x65, 0xef, 0xa6, 0x22, 0x23, 0xd6, 0x6c, 0xb4, 0x50, 0xf2, 0x62,
	0x51, 0xac, 0x50, 0x1f, 0x0a, 0x3e, 0x08, 0xc1, 0xea, 0x0a, 0x41, 0x29, 0x49, 0xfe, 0x40, 0x9a,
	0x5c, 0xca, 0xc0, 0x98, 0x89, 0x99, 0xc9, 0x60, 0x7e, 0xa0, 0xff, 0x4b, 0x32, 0x99, 0x49, 0xd2,
	0x6c, 0xec, 0xd3, 0x6e, 0xbe, 0x73, 0xee, 0xb9, 0x27, 0x37, 0x70, 0x92, 0xb1, 0x5a, 0x48, 0xac,
	0x2e, 0xcb, 0x8a, 0x4b, 0x4e, 0x56, 0xfa, 0xc7, 0xdb, 0x54, 0x65, 0xd6, 0x11, 0xff, 0x0d, 0x6c,
	0x6f, 0x69, 0x71, 0x1f, 0xe1, 0xef, 0x1a, 0x85, 0x24, 0x2e, 0xac, 0xcb, 0xb4, 0x61, 0x3c, 0xcd,
	0xdd, 0xc5, 0xf9, 0xe2, 0xc2, 0x89, 0xec, 0xa3, 0x7f, 0x01, 0x4e, 0x67, 0x14, 0x25, 0x2f, 0x04,
	0x3e, 0xe2, 0xfc, 0x06, 0xcb, 0x9f, 0x3c, 0x47, 0xf2, 0x16, 0x9e, 0x65, 0xbc, 0x28, 0x30, 0x93,
	0x94, 0x17, 0xb1, 0xac, 0x68, 0x71, 0xaf, 0xad, 0x9b, 0xe8, 0x80, 0x13, 0x02, 0xcb, 0xba, 0xa6,
	0xb9, 0xfb, 0x44, 0xeb, 0xfa, 0xbf, 0xff, 0x1e, 0x9c, 0xef, 0xc8, 0x18, 0xb7, 0xdd, 0xce, 0xe0,
	0x88, 0xa6, 0xbf, 0x74, 0xc4, 0xf6, 0x6a, 0xdb, 0xf5, 0xbf, 0x6c, 0x37, 0x45, 0x2d, 0xf7, 0x03,
	0x20, 0x01, 0xa3, 0x0a, 0x5b, 0x22, 0xfa, 0x9a, 0xef, 0x00, 0xd2, 0x9e, 0xba, 0x8b, 0xf3, 0xa3,
	0xe9, 0xec, 0x48, 0xbe, 0xfa, 0xbb, 0x82, 0xf5, 0x97, 0xee, 0x60, 0xe4, 0x03, 0x2c, 0x6f, 0x75,
	0x33, 0x63, 0x1e, 0x5d, 0xc9, 0x7b, 0xfe, 0x80, 0x99, 0x4d, 0xd7, 0xb0, 0xd2, 0x75, 0x89, 0x55,
	0xc7, 0xe5, 0xbd, 0x9d, 0x81, 0x33, 0x15, 0xaf, 0xe1, 0x64, 0x8f, 0x72, 0x10, 0x88, 0xed, 0xa7,
	0x38, 0xcd, 0x1f, 0x1b, 0xfc, 0x04, 0xeb, 0x50, 0x89, 0x38, 0x55, 0x48, 0x4e, 0x8d, 0xab, 0x7d,
	0x96, 0xbc, 0x42, 0xbb, 0xf6, 0xe5, 0x01, 0x37, 0xb3, 0x7b, 0x70, 0x42, 0x25, 0x42, 0x6c, 0xbe,
	0xfe, 0xa1, 0x42, 0x0a, 0xe2, 0x0d, 0xc6, 0x1e, 0xda, 0x90, 0x57, 0xb3, 0x9a, 0x09, 0xba, 0x81,
	0x6d, 0xa8, 0x44, 0x84, 0xb2, 0xa2, 0xa8, 0x90, 0xec, 0x06, 0xaf, 0x65, 0x36, 0xc6, 0x9b, 0x93,
	0x4c, 0xca, 0x67, 0xd8, 0x84, 0x4a, 0xdc, 0x20, 0x43, 0x89, 0x64, 0x54, 0xba, 0x23, 0x36, 0xc1,
	0x3d, 0x14, 0xcc, 0x7c, 0x00, 0x10, 0x2a, 0xb1, 0x47, 0x19, 0x62, 0x23, 0xc8, 0xc8, 0x17, 0x30,
	0xd6, 0xa2, 0xe9, 0x67, 0x18, 0x2b, 0xfd, 0x67, 0x38, 0x4e, 0x62, 0x7d, 0xcc, 0x17, 0xc6, 0x94,
	0xc4, 0x0f, 0x6e, 0x79, 0x3a, 0xc5, 0xc3, 0xee, 0x24, 0xee, 0x0f, 0xe0, 0xf6, 0xae, 0xe9, 0xfb,
	0xef, 0x66, 0x14, 0x13, 0xf1, 0x03, 0x9e, 0x26, 0x71, 0xa0, 0x52, 0xca, 0xd2, 0x3b, 0xca, 0xa8,
	0x6c, 0xc8, 0xeb, 0xde, 0x3c, 0xc6, 0x36, 0xea, 0xec, 0x3f, 0x6a, 0x17, 0x77, 0x77, 0xac, 0xd5,
	0x8f, 0xff, 0x02, 0x00, 0x00, 0xff, 0xff, 0x41, 0x6a, 0x49, 0xd5, 0xfe, 0x03, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// ClusterClient is the client API for Cluster service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type ClusterClient interface {
	Ping(ctx context.Context, in *PingRequest, opts ...grpc.CallOption) (*PingResponse, error)
	Hello(ctx context.Context, in *HelloRequest, opts ...grpc.CallOption) (*AliveNodesResponse, error)
	GetAliveNodes(ctx context.Context, in *Void, opts ...grpc.CallOption) (*AliveNodesResponse, error)
	KvsSave(ctx context.Context, in *KvsStoreRequest, opts ...grpc.CallOption) (*KvsStoreResponse, error)
	KvsKeyExists(ctx context.Context, in *KvsKeyExistsRequest, opts ...grpc.CallOption) (*KvsKeyExistsResponse, error)
	KvsRetrieve(ctx context.Context, in *KvsRetrieveRequest, opts ...grpc.CallOption) (*KvsRetrieveResponse, error)
	KvsDelete(ctx context.Context, in *KvsDeleteRequest, opts ...grpc.CallOption) (*KvsDeleteResponse, error)
	KvsGetKeys(ctx context.Context, in *KvsAllKeysRequest, opts ...grpc.CallOption) (*KvsAllKeysResponse, error)
	TSSave(ctx context.Context, in *TSStoreRequest, opts ...grpc.CallOption) (*TSStoreResponse, error)
	TSRetrieve(ctx context.Context, in *TSRetrieveRequest, opts ...grpc.CallOption) (*TSRetrieveResponse, error)
	TSAvailability(ctx context.Context, in *TSAvailabilityRequest, opts ...grpc.CallOption) (*TSAvailabilityResponse, error)
}

type clusterClient struct {
	cc grpc.ClientConnInterface
}

func NewClusterClient(cc grpc.ClientConnInterface) ClusterClient {
	return &clusterClient{cc}
}

func (c *clusterClient) Ping(ctx context.Context, in *PingRequest, opts ...grpc.CallOption) (*PingResponse, error) {
	out := new(PingResponse)
	err := c.cc.Invoke(ctx, "/proto.Cluster/Ping", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *clusterClient) Hello(ctx context.Context, in *HelloRequest, opts ...grpc.CallOption) (*AliveNodesResponse, error) {
	out := new(AliveNodesResponse)
	err := c.cc.Invoke(ctx, "/proto.Cluster/Hello", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *clusterClient) GetAliveNodes(ctx context.Context, in *Void, opts ...grpc.CallOption) (*AliveNodesResponse, error) {
	out := new(AliveNodesResponse)
	err := c.cc.Invoke(ctx, "/proto.Cluster/GetAliveNodes", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *clusterClient) KvsSave(ctx context.Context, in *KvsStoreRequest, opts ...grpc.CallOption) (*KvsStoreResponse, error) {
	out := new(KvsStoreResponse)
	err := c.cc.Invoke(ctx, "/proto.Cluster/KvsSave", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *clusterClient) KvsKeyExists(ctx context.Context, in *KvsKeyExistsRequest, opts ...grpc.CallOption) (*KvsKeyExistsResponse, error) {
	out := new(KvsKeyExistsResponse)
	err := c.cc.Invoke(ctx, "/proto.Cluster/KvsKeyExists", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *clusterClient) KvsRetrieve(ctx context.Context, in *KvsRetrieveRequest, opts ...grpc.CallOption) (*KvsRetrieveResponse, error) {
	out := new(KvsRetrieveResponse)
	err := c.cc.Invoke(ctx, "/proto.Cluster/KvsRetrieve", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *clusterClient) KvsDelete(ctx context.Context, in *KvsDeleteRequest, opts ...grpc.CallOption) (*KvsDeleteResponse, error) {
	out := new(KvsDeleteResponse)
	err := c.cc.Invoke(ctx, "/proto.Cluster/KvsDelete", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *clusterClient) KvsGetKeys(ctx context.Context, in *KvsAllKeysRequest, opts ...grpc.CallOption) (*KvsAllKeysResponse, error) {
	out := new(KvsAllKeysResponse)
	err := c.cc.Invoke(ctx, "/proto.Cluster/KvsGetKeys", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *clusterClient) TSSave(ctx context.Context, in *TSStoreRequest, opts ...grpc.CallOption) (*TSStoreResponse, error) {
	out := new(TSStoreResponse)
	err := c.cc.Invoke(ctx, "/proto.Cluster/TSSave", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *clusterClient) TSRetrieve(ctx context.Context, in *TSRetrieveRequest, opts ...grpc.CallOption) (*TSRetrieveResponse, error) {
	out := new(TSRetrieveResponse)
	err := c.cc.Invoke(ctx, "/proto.Cluster/TSRetrieve", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *clusterClient) TSAvailability(ctx context.Context, in *TSAvailabilityRequest, opts ...grpc.CallOption) (*TSAvailabilityResponse, error) {
	out := new(TSAvailabilityResponse)
	err := c.cc.Invoke(ctx, "/proto.Cluster/TSAvailability", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ClusterServer is the server API for Cluster service.
type ClusterServer interface {
	Ping(context.Context, *PingRequest) (*PingResponse, error)
	Hello(context.Context, *HelloRequest) (*AliveNodesResponse, error)
	GetAliveNodes(context.Context, *Void) (*AliveNodesResponse, error)
	KvsSave(context.Context, *KvsStoreRequest) (*KvsStoreResponse, error)
	KvsKeyExists(context.Context, *KvsKeyExistsRequest) (*KvsKeyExistsResponse, error)
	KvsRetrieve(context.Context, *KvsRetrieveRequest) (*KvsRetrieveResponse, error)
	KvsDelete(context.Context, *KvsDeleteRequest) (*KvsDeleteResponse, error)
	KvsGetKeys(context.Context, *KvsAllKeysRequest) (*KvsAllKeysResponse, error)
	TSSave(context.Context, *TSStoreRequest) (*TSStoreResponse, error)
	TSRetrieve(context.Context, *TSRetrieveRequest) (*TSRetrieveResponse, error)
	TSAvailability(context.Context, *TSAvailabilityRequest) (*TSAvailabilityResponse, error)
}

// UnimplementedClusterServer can be embedded to have forward compatible implementations.
type UnimplementedClusterServer struct {
}

func (*UnimplementedClusterServer) Ping(ctx context.Context, req *PingRequest) (*PingResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Ping not implemented")
}
func (*UnimplementedClusterServer) Hello(ctx context.Context, req *HelloRequest) (*AliveNodesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Hello not implemented")
}
func (*UnimplementedClusterServer) GetAliveNodes(ctx context.Context, req *Void) (*AliveNodesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAliveNodes not implemented")
}
func (*UnimplementedClusterServer) KvsSave(ctx context.Context, req *KvsStoreRequest) (*KvsStoreResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method KvsSave not implemented")
}
func (*UnimplementedClusterServer) KvsKeyExists(ctx context.Context, req *KvsKeyExistsRequest) (*KvsKeyExistsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method KvsKeyExists not implemented")
}
func (*UnimplementedClusterServer) KvsRetrieve(ctx context.Context, req *KvsRetrieveRequest) (*KvsRetrieveResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method KvsRetrieve not implemented")
}
func (*UnimplementedClusterServer) KvsDelete(ctx context.Context, req *KvsDeleteRequest) (*KvsDeleteResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method KvsDelete not implemented")
}
func (*UnimplementedClusterServer) KvsGetKeys(ctx context.Context, req *KvsAllKeysRequest) (*KvsAllKeysResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method KvsGetKeys not implemented")
}
func (*UnimplementedClusterServer) TSSave(ctx context.Context, req *TSStoreRequest) (*TSStoreResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method TSSave not implemented")
}
func (*UnimplementedClusterServer) TSRetrieve(ctx context.Context, req *TSRetrieveRequest) (*TSRetrieveResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method TSRetrieve not implemented")
}
func (*UnimplementedClusterServer) TSAvailability(ctx context.Context, req *TSAvailabilityRequest) (*TSAvailabilityResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method TSAvailability not implemented")
}

func RegisterClusterServer(s *grpc.Server, srv ClusterServer) {
	s.RegisterService(&_Cluster_serviceDesc, srv)
}

func _Cluster_Ping_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PingRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ClusterServer).Ping(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Cluster/Ping",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ClusterServer).Ping(ctx, req.(*PingRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Cluster_Hello_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(HelloRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ClusterServer).Hello(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Cluster/Hello",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ClusterServer).Hello(ctx, req.(*HelloRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Cluster_GetAliveNodes_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Void)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ClusterServer).GetAliveNodes(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Cluster/GetAliveNodes",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ClusterServer).GetAliveNodes(ctx, req.(*Void))
	}
	return interceptor(ctx, in, info, handler)
}

func _Cluster_KvsSave_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(KvsStoreRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ClusterServer).KvsSave(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Cluster/KvsSave",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ClusterServer).KvsSave(ctx, req.(*KvsStoreRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Cluster_KvsKeyExists_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(KvsKeyExistsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ClusterServer).KvsKeyExists(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Cluster/KvsKeyExists",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ClusterServer).KvsKeyExists(ctx, req.(*KvsKeyExistsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Cluster_KvsRetrieve_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(KvsRetrieveRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ClusterServer).KvsRetrieve(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Cluster/KvsRetrieve",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ClusterServer).KvsRetrieve(ctx, req.(*KvsRetrieveRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Cluster_KvsDelete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(KvsDeleteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ClusterServer).KvsDelete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Cluster/KvsDelete",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ClusterServer).KvsDelete(ctx, req.(*KvsDeleteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Cluster_KvsGetKeys_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(KvsAllKeysRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ClusterServer).KvsGetKeys(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Cluster/KvsGetKeys",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ClusterServer).KvsGetKeys(ctx, req.(*KvsAllKeysRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Cluster_TSSave_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TSStoreRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ClusterServer).TSSave(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Cluster/TSSave",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ClusterServer).TSSave(ctx, req.(*TSStoreRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Cluster_TSRetrieve_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TSRetrieveRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ClusterServer).TSRetrieve(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Cluster/TSRetrieve",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ClusterServer).TSRetrieve(ctx, req.(*TSRetrieveRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Cluster_TSAvailability_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TSAvailabilityRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ClusterServer).TSAvailability(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Cluster/TSAvailability",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ClusterServer).TSAvailability(ctx, req.(*TSAvailabilityRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Cluster_serviceDesc = grpc.ServiceDesc{
	ServiceName: "proto.Cluster",
	HandlerType: (*ClusterServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Ping",
			Handler:    _Cluster_Ping_Handler,
		},
		{
			MethodName: "Hello",
			Handler:    _Cluster_Hello_Handler,
		},
		{
			MethodName: "GetAliveNodes",
			Handler:    _Cluster_GetAliveNodes_Handler,
		},
		{
			MethodName: "KvsSave",
			Handler:    _Cluster_KvsSave_Handler,
		},
		{
			MethodName: "KvsKeyExists",
			Handler:    _Cluster_KvsKeyExists_Handler,
		},
		{
			MethodName: "KvsRetrieve",
			Handler:    _Cluster_KvsRetrieve_Handler,
		},
		{
			MethodName: "KvsDelete",
			Handler:    _Cluster_KvsDelete_Handler,
		},
		{
			MethodName: "KvsGetKeys",
			Handler:    _Cluster_KvsGetKeys_Handler,
		},
		{
			MethodName: "TSSave",
			Handler:    _Cluster_TSSave_Handler,
		},
		{
			MethodName: "TSRetrieve",
			Handler:    _Cluster_TSRetrieve_Handler,
		},
		{
			MethodName: "TSAvailability",
			Handler:    _Cluster_TSAvailability_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "cluster.proto",
}

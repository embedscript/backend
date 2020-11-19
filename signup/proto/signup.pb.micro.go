// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: proto/signup.proto

package signup

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

import (
	context "context"
	api "github.com/micro/micro/v3/service/api"
	client "github.com/micro/micro/v3/service/client"
	server "github.com/micro/micro/v3/service/server"
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

// Reference imports to suppress errors if they are not otherwise used.
var _ api.Endpoint
var _ context.Context
var _ client.Option
var _ server.Option

// Api Endpoints for Signup service

func NewSignupEndpoints() []*api.Endpoint {
	return []*api.Endpoint{}
}

// Client API for Signup service

type SignupService interface {
	Call(ctx context.Context, in *Request, opts ...client.CallOption) (*Response, error)
	Stream(ctx context.Context, in *StreamingRequest, opts ...client.CallOption) (Signup_StreamService, error)
	PingPong(ctx context.Context, opts ...client.CallOption) (Signup_PingPongService, error)
}

type signupService struct {
	c    client.Client
	name string
}

func NewSignupService(name string, c client.Client) SignupService {
	return &signupService{
		c:    c,
		name: name,
	}
}

func (c *signupService) Call(ctx context.Context, in *Request, opts ...client.CallOption) (*Response, error) {
	req := c.c.NewRequest(c.name, "Signup.Call", in)
	out := new(Response)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *signupService) Stream(ctx context.Context, in *StreamingRequest, opts ...client.CallOption) (Signup_StreamService, error) {
	req := c.c.NewRequest(c.name, "Signup.Stream", &StreamingRequest{})
	stream, err := c.c.Stream(ctx, req, opts...)
	if err != nil {
		return nil, err
	}
	if err := stream.Send(in); err != nil {
		return nil, err
	}
	return &signupServiceStream{stream}, nil
}

type Signup_StreamService interface {
	Context() context.Context
	SendMsg(interface{}) error
	RecvMsg(interface{}) error
	Close() error
	Recv() (*StreamingResponse, error)
}

type signupServiceStream struct {
	stream client.Stream
}

func (x *signupServiceStream) Close() error {
	return x.stream.Close()
}

func (x *signupServiceStream) Context() context.Context {
	return x.stream.Context()
}

func (x *signupServiceStream) SendMsg(m interface{}) error {
	return x.stream.Send(m)
}

func (x *signupServiceStream) RecvMsg(m interface{}) error {
	return x.stream.Recv(m)
}

func (x *signupServiceStream) Recv() (*StreamingResponse, error) {
	m := new(StreamingResponse)
	err := x.stream.Recv(m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (c *signupService) PingPong(ctx context.Context, opts ...client.CallOption) (Signup_PingPongService, error) {
	req := c.c.NewRequest(c.name, "Signup.PingPong", &Ping{})
	stream, err := c.c.Stream(ctx, req, opts...)
	if err != nil {
		return nil, err
	}
	return &signupServicePingPong{stream}, nil
}

type Signup_PingPongService interface {
	Context() context.Context
	SendMsg(interface{}) error
	RecvMsg(interface{}) error
	Close() error
	Send(*Ping) error
	Recv() (*Pong, error)
}

type signupServicePingPong struct {
	stream client.Stream
}

func (x *signupServicePingPong) Close() error {
	return x.stream.Close()
}

func (x *signupServicePingPong) Context() context.Context {
	return x.stream.Context()
}

func (x *signupServicePingPong) SendMsg(m interface{}) error {
	return x.stream.Send(m)
}

func (x *signupServicePingPong) RecvMsg(m interface{}) error {
	return x.stream.Recv(m)
}

func (x *signupServicePingPong) Send(m *Ping) error {
	return x.stream.Send(m)
}

func (x *signupServicePingPong) Recv() (*Pong, error) {
	m := new(Pong)
	err := x.stream.Recv(m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

// Server API for Signup service

type SignupHandler interface {
	Call(context.Context, *Request, *Response) error
	Stream(context.Context, *StreamingRequest, Signup_StreamStream) error
	PingPong(context.Context, Signup_PingPongStream) error
}

func RegisterSignupHandler(s server.Server, hdlr SignupHandler, opts ...server.HandlerOption) error {
	type signup interface {
		Call(ctx context.Context, in *Request, out *Response) error
		Stream(ctx context.Context, stream server.Stream) error
		PingPong(ctx context.Context, stream server.Stream) error
	}
	type Signup struct {
		signup
	}
	h := &signupHandler{hdlr}
	return s.Handle(s.NewHandler(&Signup{h}, opts...))
}

type signupHandler struct {
	SignupHandler
}

func (h *signupHandler) Call(ctx context.Context, in *Request, out *Response) error {
	return h.SignupHandler.Call(ctx, in, out)
}

func (h *signupHandler) Stream(ctx context.Context, stream server.Stream) error {
	m := new(StreamingRequest)
	if err := stream.Recv(m); err != nil {
		return err
	}
	return h.SignupHandler.Stream(ctx, m, &signupStreamStream{stream})
}

type Signup_StreamStream interface {
	Context() context.Context
	SendMsg(interface{}) error
	RecvMsg(interface{}) error
	Close() error
	Send(*StreamingResponse) error
}

type signupStreamStream struct {
	stream server.Stream
}

func (x *signupStreamStream) Close() error {
	return x.stream.Close()
}

func (x *signupStreamStream) Context() context.Context {
	return x.stream.Context()
}

func (x *signupStreamStream) SendMsg(m interface{}) error {
	return x.stream.Send(m)
}

func (x *signupStreamStream) RecvMsg(m interface{}) error {
	return x.stream.Recv(m)
}

func (x *signupStreamStream) Send(m *StreamingResponse) error {
	return x.stream.Send(m)
}

func (h *signupHandler) PingPong(ctx context.Context, stream server.Stream) error {
	return h.SignupHandler.PingPong(ctx, &signupPingPongStream{stream})
}

type Signup_PingPongStream interface {
	Context() context.Context
	SendMsg(interface{}) error
	RecvMsg(interface{}) error
	Close() error
	Send(*Pong) error
	Recv() (*Ping, error)
}

type signupPingPongStream struct {
	stream server.Stream
}

func (x *signupPingPongStream) Close() error {
	return x.stream.Close()
}

func (x *signupPingPongStream) Context() context.Context {
	return x.stream.Context()
}

func (x *signupPingPongStream) SendMsg(m interface{}) error {
	return x.stream.Send(m)
}

func (x *signupPingPongStream) RecvMsg(m interface{}) error {
	return x.stream.Recv(m)
}

func (x *signupPingPongStream) Send(m *Pong) error {
	return x.stream.Send(m)
}

func (x *signupPingPongStream) Recv() (*Ping, error) {
	m := new(Ping)
	if err := x.stream.Recv(m); err != nil {
		return nil, err
	}
	return m, nil
}
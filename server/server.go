package main

import (
	"context"
	"fmt"
	proto "github.com/robinbraemer/cacher/cacher"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"net"
)

type CacheServer struct {
	m map[string][]byte
}

func (s *CacheServer) All(ctx *proto.AllRequest, stream proto.Cache_AllServer) error {
	for k, v := range s.m {
		if err := stream.Send(&proto.AllReply{Entry: &proto.Entry{Key: k, Val: v}}); err != nil {
			return err
		}
	}
	return nil
}

func (s *CacheServer) Del(ctx context.Context, in *proto.DelRequest) (*proto.DelReply, error) {
	delete(s.m, in.Key)
	return &proto.DelReply{}, nil
}

func (s *CacheServer) Set(ctx context.Context, in *proto.SetRequest) (*proto.SetReply, error) {
	s.m[in.Entry.Key] = in.Entry.Val
	return &proto.SetReply{}, nil
}

func (s *CacheServer) Get(ctx context.Context, in *proto.GetRequest) (*proto.GetReply, error) {
	v, ok := s.m[in.Key]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "Key %s not found", in.Key)
	}
	return &proto.GetReply{Val: v}, nil
}

func main() {
	s := grpc.NewServer()
	proto.RegisterCacheServer(s, &CacheServer{m: make(map[string][]byte)})
	ln, err := net.Listen("tcp", ":50001")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Cacher server started")
	if err := s.Serve(ln); err != nil {
		log.Fatal(err)
	}
}

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
	"sync"
)

type CacheServer struct {
	m sync.Map
}

func (s *CacheServer) All(ctx *proto.AllRequest, stream proto.Cache_AllServer) error {
	var e error
	s.m.Range(func(k, v interface{}) bool {
		if err := stream.Send(&proto.AllReply{Entry: &proto.Entry{Key: k.(string), Val: v.([]byte)}}); err != nil {
			e = err
			return false
		}
		return true
	})
	return e
}

func (s *CacheServer) Del(ctx context.Context, in *proto.DelRequest) (*proto.DelReply, error) {
	s.m.Delete(in.Key)
	return &proto.DelReply{}, nil
}

func (s *CacheServer) Set(ctx context.Context, in *proto.SetRequest) (*proto.SetReply, error) {
	s.m.Store(in.Entry.Key, in.Entry.Val)
	return &proto.SetReply{}, nil
}

func (s *CacheServer) Get(ctx context.Context, in *proto.GetRequest) (*proto.GetReply, error) {
	v, ok := s.m.Load(in.Key)
	if !ok {
		return nil, status.Errorf(codes.NotFound, "Key %s not found", in.Key)
	}
	return &proto.GetReply{Val: v.([]byte)}, nil
}

// main starts the server.
// It does not take any inputs or return anything.
func main() {
	s := grpc.NewServer()
	proto.RegisterCacheServer(s, &CacheServer{})
	ln, err := net.Listen("tcp", ":50001")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Cacher server started")
	if err := s.Serve(ln); err != nil {
		log.Fatal(err)
	}
}

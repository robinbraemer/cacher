package main

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	proto "github.com/robinbraemer/cacher/cacher"
	"google.golang.org/grpc"
	"log"
)

func fillWithoutConnReuse(ctx context.Context, count int) {
	// fill cacher
	r := uuid.New().String()
	for i := 0; i < count; i++ {
		cc, err := grpc.Dial(":50001", grpc.WithInsecure())
		if err != nil {
			log.Fatal(err)
		}
		c := proto.NewCacheClient(cc)

		next := fmt.Sprintf("%s-%d", r, i)
		_, err = c.Set(ctx, &proto.SetRequest{Entry: &proto.Entry{Key: next, Val: []byte(next)}})
		if err != nil {
			log.Fatal(err)
		}

		_ = cc.Close()
	}
}
func fillWithConnReuse(ctx context.Context, c proto.CacheClient, count int) {
	// fill cacher
	r := uuid.New().String()
	for i := 0; i < count; i++ {
		next := fmt.Sprintf("%s-%d", r, i)
		_, err := c.Set(ctx, &proto.SetRequest{Entry: &proto.Entry{Key: next, Val: []byte(next)}})
		if err != nil {
			log.Fatal(err)
		}
	}
}

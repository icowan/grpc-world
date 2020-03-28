/**
 * @Time: 2020/3/28 17:05
 * @Author: solacowa@gmail.com
 * @File: grpc
 * @Software: GoLand
 */

package main

import (
	"context"
	"github.com/icowan/grpc-world/pkg/grpc/pb"
	"google.golang.org/grpc"
	"log"
	"time"
)

func main() {
	conn, err := grpc.Dial("127.0.0.1:8081", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer func() {
		_ = conn.Close()
	}()

	svc := pb.NewServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := svc.Put(ctx, &pb.GetRequest{
		Key: "hello",
		Val: "world",
	})
	if err != nil {
		log.Fatalf("could not put: %v", err)
	}

	r, err = svc.Get(ctx, &pb.GetRequest{
		Key: "hello",
	})
	if err != nil {
		log.Fatalf("could not get: %v", err)
	}
	log.Printf("data: %s", r.GetData())
}

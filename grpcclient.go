package main

import (
	"context"
	"fmt"
	pb "github.com/programmer74/gotsdb/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"os"
)

func main() {
	opts := []grpc.DialOption{
		grpc.WithInsecure(),
	}
	args := os.Args
	conn, err := grpc.Dial("127.0.0.1:5300", opts...)

	if err != nil {
		grpclog.Fatalf("fail to dial: %v", err)
	}

	defer conn.Close()

	client := pb.NewReverseClient(conn)
	request := &pb.Request{
		Message: args[1],
	}
	response, err := client.DoStuff(context.Background(), request)

	if err != nil {
		grpclog.Fatalf("fail to dial: %v", err)
	}

	fmt.Println(response.Message)
}
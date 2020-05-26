package main

import (
	"context"
	"github.com/abiosoft/ishell"
	pb "github.com/programmer74/gotsdb/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"os"
)

func main() {
	argsWithoutProg := os.Args[1:]
	if len(argsWithoutProg) < 1 {
		println("Usage: " + os.Args[0] + " <hostPort>")
		return
	}

	hostPort := argsWithoutProg[0]

	println("Connecting to " + hostPort)

	opts := []grpc.DialOption{
		grpc.WithInsecure(),
	}
	conn, err := grpc.Dial(hostPort, opts...)

	if err != nil {
		grpclog.Fatalf("fail to dial: %v", err)
	}

	defer conn.Close()

	client := pb.NewReverseClient(conn)

	// create new shell.
	// by default, new shell includes 'exit', 'help' and 'clear' commands.
	shell := ishell.New()

	// display welcome info.
	shell.Println("Ready to accept commands")

	// register a function for "greet" command.
	//shell.AddCmd(&ishell.Cmd{
	//	Name: "store",
	//	Help: "store value by key",
	//	Func: func(c *ishell.Context) {
	//		if (len(c.Args)) != 2 {
	//			c.Println("Usage: store <key> <value>")
	//		} else {
	//			StoreCmd(c, conn)
	//		}
	//	},
	//})
	//
	//shell.AddCmd(&ishell.Cmd{
	//	Name: "retrieve",
	//	Help: "retrieve value by key",
	//	Func: func(c *ishell.Context) {
	//		if (len(c.Args)) != 1 {
	//			c.Println("Usage: retrieve")
	//		} else {
	//			RequestCmd(c, conn)
	//		}
	//	},
	//})

	shell.AddCmd(&ishell.Cmd{
		Name: "test",
		Help: "test grpc",
		Func: func(c *ishell.Context) {
			if (len(c.Args)) != 1 {
				c.Println("Usage: test <string>")
			} else {
				c.Println(testGrpc(client, c.Args[0]))
			}
		},
	})

	shell.Run()
}

func testGrpc(client pb.ReverseClient, str string) string {
	request := &pb.Request{
		Message: str,
	}

	response, err := client.DoStuff(context.TODO(), request)

	if err != nil {
		grpclog.Fatalf("fail to dial: %v", err)
	}
	return response.Message
}

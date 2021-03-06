package main

import (
	"context"
	"github.com/abiosoft/ishell"
	pb "github.com/nikita-tomilov/gotsdb/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"os"
)

var msgId uint32 = 0

func main1() { //wow you cannot have two mains in one package, cool
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

	client := pb.NewGoTSDBClient(conn)

	shell := ishell.New()
	shell.Println("Ready to accept commands")

	shell.AddCmd(&ishell.Cmd{
		Name: "store",
		Help: "store value by key",
		Func: func(c *ishell.Context) {
			if (len(c.Args)) != 2 {
				c.Println("Usage: store <key> <value>")
			} else {
				c.Println(KvsSaveCmd(c, client))
			}
		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "retrieve",
		Help: "retrieve value by key",
		Func: func(c *ishell.Context) {
			if (len(c.Args)) != 1 {
				c.Println("Usage: retrieve <key> ")
			} else {
				c.Println(KvsRetrieveCmd(c, client))
			}
		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "check",
		Help: "check if key exists",
		Func: func(c *ishell.Context) {
			if (len(c.Args)) != 1 {
				c.Println("Usage: check <key> ")
			} else {
				c.Println(KvsKeyExistsCmd(c, client))
			}
		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "delete",
		Help: "delete key",
		Func: func(c *ishell.Context) {
			if (len(c.Args)) != 1 {
				c.Println("Usage: delete <key> ")
			} else {
				c.Println(KvsDeleteCmd(c, client))
			}
		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "list",
		Help: "list all keys",
		Func: func(c *ishell.Context) {
			c.Println(KvsListCmd(c, client))
		},
	})

	shell.Run()
}

func KvsSaveCmd(c *ishell.Context, client pb.GoTSDBClient) string {
	msgId += 1
	stringKey := c.Args[0]
	stringValue := c.Args[1]
	key := []byte(stringKey)
	value := []byte(stringValue)
	request := &pb.KvsStoreRequest{
		MsgId: msgId,
		Key:   key,
		Value: value,
	}
	response, err := client.KvsSave(context.TODO(), request)
	if err != nil {
		grpclog.Fatalf("fail to dial: %v", err)
		return "err " + err.Error()
	}
	if response.Ok {
		return "ok"
	}
	return "fail"
}

func KvsRetrieveCmd(c *ishell.Context, client pb.GoTSDBClient) string {
	msgId += 1
	stringKey := c.Args[0]
	key := []byte(stringKey)

	request := &pb.KvsRetrieveRequest{
		MsgId: msgId,
		Key:   key,
	}
	response, err := client.KvsRetrieve(context.TODO(), request)
	if err != nil {
		grpclog.Fatalf("fail to dial: %v", err)
		return "err " + err.Error()
	}
	return "ok; " + string(response.Value)
}

func KvsKeyExistsCmd(c *ishell.Context, client pb.GoTSDBClient) string {
	msgId += 1
	stringKey := c.Args[0]
	key := []byte(stringKey)

	request := &pb.KvsKeyExistsRequest{
		MsgId: msgId,
		Key:   key,
	}
	response, err := client.KvsKeyExists(context.TODO(), request)
	if err != nil {
		grpclog.Fatalf("fail to dial: %v", err)
		return "err " + err.Error()
	}
	if response.Exists {
		return "key exists"
	}
	return "key does not exist"
}

func KvsDeleteCmd(c *ishell.Context, client pb.GoTSDBClient) string {
	msgId += 1
	stringKey := c.Args[0]
	key := []byte(stringKey)

	request := &pb.KvsDeleteRequest{
		MsgId: msgId,
		Key:   key,
	}
	response, err := client.KvsDelete(context.TODO(), request)
	if err != nil {
		grpclog.Fatalf("fail to dial: %v", err)
		return "err " + err.Error()
	}
	if response.Ok {
		return "ok"
	}
	return "fail"
}

func KvsListCmd(c *ishell.Context, client pb.GoTSDBClient) string {
	msgId += 1

	request := &pb.KvsAllKeysRequest{
		MsgId: msgId,
	}
	response, err := client.KvsGetKeys(context.TODO(), request)
	if err != nil {
		grpclog.Fatalf("fail to dial: %v", err)
		return "err " + err.Error()
	}

	ans := ""
	for _, key := range response.Keys {
		ans = ans + string(key) + "; "
	}

	return "ok; " + ans
}

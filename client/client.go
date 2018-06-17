package main

import (
	"context"
	"log"
	"fmt"

	pb "github.com/alexuserid/grpc-chat/proto"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:8080", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("grpc.Dial: %v", err)
	}

	c := pb.NewChatClient(conn)

	var user string
	fmt.Scanln(&user)
	sid, err := c.Login(context.Background(), &pb.LoginRequest{Name: user})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(sid)

	users, err := c.LastUsers(context.Background(), &pb.Empty{})
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, v := range users.Users {
		fmt.Println(v)
	}
}

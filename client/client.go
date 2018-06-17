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
	fmt.Print("Enter user name: ")
	fmt.Scanln(&user)
	sid, err := c.Login(context.Background(), &pb.LoginRequest{Name: user})
	if err != nil {
		log.Fatalf("Login: %s", err)
	}
	fmt.Println(sid)

	users, err := c.ListUsers(context.Background(), &pb.Empty{})
	if err != nil {
		log.Fatalf("Getting user list: %s", err)
		return
	}
	fmt.Println("Users online:")
	for i, v := range users.Users {
		fmt.Printf("%d. %s\n", i+1, v)
	}
}

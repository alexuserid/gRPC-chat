package main

import (
	"context"
	"log"
	"fmt"
	"io"

	pb "github.com/alexuserid/grpc-chat/proto"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:8080", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("grpc.Dial: %v", err)
	}

	c := pb.NewChatClient(conn)

	// Login
	var user string
	fmt.Print("Enter user name: ")
	fmt.Scanln(&user)
	sid, err := c.Login(context.Background(), &pb.LoginRequest{Name: user})
	if err != nil {
		log.Fatalf("Login: %s", err)
	}
	fmt.Println(sid)

	// User list
	users, err := c.ListUsers(context.Background(), &pb.Empty{})
	if err != nil {
		log.Fatalf("Getting user list: %s", err)
		return
	}
	fmt.Println("Users online:")
	for i, v := range users.Users {
		fmt.Printf("%d. %s\n", i+1, v)
	}

	// Show messaes
	stream, err := c.Watch(context.Background(), &pb.Empty{})
	// Show old messages
	fmt.Println("Unread messages: ")
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Println("stream.Recv(): %v", err)
			return
		}
		if in.Name == "" {
			continue
		}
		fmt.Printf("%s: %v\n", in.Name, in.GetText())
	}
	// Show new messages
	// TODO: implment new messages receive

	// Send message
	var msgText string
	fmt.Print("Message: ")
	fmt.Scanln(&msgText)
	// TODO: implement logout message
	msg := &pb.SendMessageRequest{
		Sid: sid.Sid,
		Text: msgText,
	}
	_, err = c.SendMessage(context.Background(), msg)
	if err != nil {
		log.Fatal("SendMessage: %v", err)
	}
}

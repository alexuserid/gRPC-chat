package main

import (
	"container/ring"
	"context"
	"log"
	"fmt"
	"sync"
	"net"

	pb "github.com/alexuserid/grpc-chat/proto"
	"google.golang.org/grpc"
	"github.com/alexuserid/id"
)

const messageN = 10

type server struct{
	sidUser map[string]string
	usernames map[string]struct{}
	messageRing *ring.Ring
	mutex sync.RWMutex
}


func NewServer() (*server, error) {
	return &server{
		sidUser: make(map[string]string),
		usernames: make(map[string]struct{}),
		messageRing: ring.New(messageN),
	}, nil
}

func (s *server) Login(ctx context.Context, in *pb.LoginRequest) (*pb.LoginResponse, error) {
	sid, err := id.GetRandomHexString(32)
	if err != nil {
		log.Printf("id.GetRandomHex: %v", err)
	}

	exists := struct{}{}

	s.mutex.Lock()
	defer s.mutex.Unlock()
	if _, has := s.usernames[in.Name]; has {
		return nil, fmt.Errorf("User with this name already exists")
	}
	if _, has := s.sidUser[sid]; has {
		return nil, fmt.Errorf("Please, try again")
	}
	s.sidUser[sid] = in.Name
	s.usernames[in.Name] = exists

	return &pb.LoginResponse{Sid: sid}, nil
}

func (s *server) Logout(ctx context.Context, in *pb.LogoutRequest) (*pb.Empty, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	name := s.sidUser[in.Sid]
	delete(s.sidUser, in.Sid)
	delete(s.usernames, name)

	return nil, nil
}

func (s *server) LastUsers(ctx context.Context, in *pb.Empty) (*pb.ListUsersResponse, error) {
	var users []string
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	for k, _ := range s.usernames {
		users = append(users, k)
	}

	return &pb.ListUsersResponse{Users: users}, nil
}

func (s *server) SendMessage(ctx context.Context, in *pb.SendMessageRequest) (*pb.Empty, error) {
	s.messageRing.Value = in
	s.messageRing = s.messageRing.Next()
	return nil, nil
}

func (s *server) Watch(in *pb.Empty, stream pb.Chat_WatchServer) error {
	return fmt.Errorf("Unimplemented")
}

func main() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("net.Listen: %v", err)
	}

	gs := grpc.NewServer()
	s, err := NewServer()
	if err != nil {
		log.Fatalf("NewServer(): %v", err)
	}
	pb.RegisterChatServer(gs, s)
	gs.Serve(ln)
}

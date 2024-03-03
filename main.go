package main

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"
	pb "user.service"
)

const (
	port = ":50052"
)

var users []*pb.UserDetailsResponse

type userServiceServer struct {
	pb.UnimplementedUserServiceServer
}

func main() {
	initUsers()
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to start server on port: "+port, err)
	}

	s := grpc.NewServer()

	pb.RegisterUserServiceServer(s, &userServiceServer{})

	log.Println("Stated server at %v", lis.Addr())

}

func initUsers() {
	user1 := &pb.UserDetailsResponse{Name: "Ankit", Id: "qwerty"}
	user2 := &pb.UserDetailsResponse{Name: "Anand", Id: "zxcvb"}
	users = append(users, user1, user2)
}

func (s *userServiceServer) GetUsers(in *pb.NoParam, stream pb.userServiceGetUsersServer) error {
	log.Printf("Received %v", in)
	for _, user := range users {
		if err := stream.Send(user); err != nil {
			return err
		}
	}
	return nil
}

func (s *userServiceServer) GetUser(ctx context.Context,
	in *pb.UserRequest) (*pb.UserDetailsResponse, error) {
	log.Printf("Received %v", in)
	res := &pb.UserDetailsResponse{}

	for _, user := range users {
		if user.GetId() == in.GetId() {
			res = user
			break
		}
	}
	return res, nil

}

func (s *userServiceServer) CreateUser(ctx context.Context,
	in *pb.UserDetailsRequest) (*pb.UserDetailsResponse, error) {
	log.Printf("Create user for %v", in)
	user := &pb.UserDetailsResponse{Name: in.GetName(), Id: in.GetId()}
	users = append(users, user)
	return user, nil
}

func (s *userServiceServer) UpdateUser(ctx context.Context,
	in *pb.UserDetailsRequest) (*pb.UserDetailsResponse, error) {
	log.Printf("Update user for %v", in)
	resp := &pb.UserDetailsResponse{}

	for index, user := range users {
		if user.GetId() == in.GetId() {
			users = append(users[:index], users[index+1:]...)
			user.Name = in.GetName()
			resp = user
			users = append(users, user)
		}
	}

	return resp, nil
}

func (s *userServiceServer) DeleteUser(ctx context.Context,
	in *pb.UserRequest) (*pb.UserDetailsResponse, error) {
	log.Printf("Delete user for %v", in)
	resp := &pb.UserDetailsResponse{}

	for index, user := range users {
		if user.GetId() == in.GetId() {
			resp = user
			users = append(users[:index], users[index+1:]...)
		}
	}
	return resp, nil
}

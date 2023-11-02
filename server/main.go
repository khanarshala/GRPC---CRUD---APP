syntax = "proto3";

package user_service;

service UserService {
    rpc GetUserById (UserRequest) returns (User);
    rpc GetUsersByIds (UserIds) returns (UserList);
}

message UserRequest {
    int32 id = 1;
}

message UserIds {
    repeated int32 ids = 1;
}

message User {
    int32 id = 1;
    string fname = 2;
    string city = 3;
    uint64 phone = 4;
    float height = 5;
    bool married = 6;
}

package main

import (
    "context"
    "log"
    "net"

    "github.com/user_service"
    "google.golang.org/grpc"
)
type User struct {
    ID      int32
    Fname   string
    City    string
    Phone   uint64
    Height  float32
    Married bool
}
type UserServiceServer struct {
    users map[int32]User
}

func (s *UserServiceServer) GetUserById(ctx context.Context, req *user_service.UserRequest) (*user_service.User, error) {
    user, found := s.users[req.Id]
    if !found {
        return nil, grpc.Errorf(codes.NotFound, "User not found")
    }
    return &user_service.User{
        Id:      user.ID,
        Fname:   user.Fname,
        City:    user.City,
        Phone:   user.Phone,
        Height:  user.Height,
        Married: user.Married,
    }, nil
}
func (s *UserServiceServer) GetUsersByIds(ctx context.Context, req *user_service.UserIds) (*user_service.UserList, error) {
    var users []*user_service.User
    for _, id := range req.Ids {
        user, found := s.users[id]
        if found {
            users = append(users, &user_service.User{
                Id:      user.ID,
                Fname:   user.Fname,
                City:    user.City,
                Phone:   user.Phone,
                Height:  user.Height,
                Married: user.Married,
            })
        }
    }
    return &user_service.UserList{Users: users}, nil
}

func main() {
    lis, err := net.Listen("tcp", ":50051")
    if err != nil {
        log.Fatalf("Failed to listen: %v", err)
    }
    s := grpc.NewServer()
    userServer := &UserServiceServer{
        users: map[int32]User{
            1: {ID: 1, Fname: "Steve", City: "LA", Phone: 1234567890, Height: 5.8, Married: true},
            // Add more users as needed
        },
    }
    user_service.RegisterUserServiceServer(s, userServer)
    if err := s.Serve(lis); err != nil {
        log.Fatalf("Failed to serve: %v", err)
    }
}

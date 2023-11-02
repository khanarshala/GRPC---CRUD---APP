//1st I created a directory structure to initiaized the GO module
mkdir grpc-user-service
cd grpc-user-service
go mod init grpc-user-service
syntax = "proto3";

package users;

message User {
    int32 id = 1;
    string fname = 2;
    string city = 3;
    int64 phone = 4;
    float height = 5;
    bool married = 6;
}

service UserService {
    rpc GetUserById (UserRequest) returns (User);
    rpc GetUsersByIds (UserIdsRequest) returns (stream User);
}

message UserRequest {
    int32 id = 1;
}

message UserIdsRequest {
    repeated int32 ids = 1;
}
protoc --go_out=plugins=grpc:. user.proto
package main

import (
    "context"
    "log"
    "net"
    "google.golang.org/grpc"
    pb "yourpackage" // Import your generated code
)

type userServiceServer struct {
    users map[int32]*pb.User // Mocked database
}

func (s *userServiceServer) GetUserById(ctx context.Context, req *pb.UserRequest) (*pb.User, error) {
    user, found := s.users[req.Id]
    if !found {
        return nil, grpc.Errorf(codes.NotFound, "User not found")
    }
    return user, nil
}

func (s *userServiceServer) GetUsersByIds(req *pb.UserIdsRequest, stream pb.UserService_GetUsersByIdsServer) error {
    for _, id := range req.Ids {
        user, found := s.users[id]
        if found {
            if err := stream.Send(user); err != nil {
                return err
            }
        }
    }
    return nil
}

func main() {
    lis, err := net.Listen("tcp", ":50051")
    if err != nil {
        log.Fatalf("Failed to listen: %v", err)
    }
    server := grpc.NewServer()
    pb.RegisterUserServiceServer(server, &userServiceServer{
        users: map[int32]*pb.User{
            1: &pb.User{Id: 1, Fname: "Steve", City: "LA", Phone: 1234567890, Height: 5.8, Married: true},
            // Add more user entries here
        },
    })
    if err := server.Serve(lis); err != nil {
        log.Fatalf("Failed to serve: %v", err)
    }
}

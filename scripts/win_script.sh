syntax = "proto3";

package user;

message User {
  int32 id = 1;
  string fname = 2;
  string city = 3;
  int64 phone = 4;
  double height = 5;
  bool married = 6;
}

service UserService {
  rpc GetUserByID (GetUserRequest) returns (User);
  rpc GetUsersByIDs (GetUsersRequest) returns (repeated User);
}

message GetUserRequest {
  int32 id = 1;
}

message GetUsersRequest {
  repeated int32 ids = 1;
}
//Here, I generated the Go code from the "proto" file
protoc --go_out=plugins=grpc:. user.proto
}
}

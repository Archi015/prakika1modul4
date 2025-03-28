package main

import (
	"log"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

func main() {
	log.Println("gRPC and Protobuf installed successfully")
	_ = grpc.NewServer()     // Проверяем, что gRPC работает
	_ = proto.String("test") // Проверяем, что Protobuf работает
}

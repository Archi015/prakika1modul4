package main

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "proto/auth.proto"

	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewUserServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	resp, err := client.GetUser(ctx, &pb.UserRequest{Id: 1})
	if err != nil {
		log.Fatalf("Error calling GetUser: %v", err)
	}

	fmt.Printf("User: ID=%d, Name=%s, Email=%s\n", resp.Id, resp.Name, resp.Email)
}

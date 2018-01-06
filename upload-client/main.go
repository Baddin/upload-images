package main

import (
	"context"
	"log"

	pb "github.com/baddin/upload-images/upload-service"
	"google.golang.org/grpc"
)

const address = "localhost:50051"

func main() {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	c := pb.NewUploadImagesClient(conn)
	path := "anothertest.png"
	description := "this is a another test image"
	r, err := c.UpImage(context.Background(), &pb.UpRequest{Path: path, Description: description})
	if err != nil {
		panic(err)
	}
	log.Printf("Created: %v\n", r.Created)

	get, err := c.GetImage(context.Background(), &pb.GetRequest{Id: 6}) // just testing it :)
	if err != nil {
		panic(err)
	}
	log.Printf("imageUrl: %v\ndescription: %v\n", get.Url, get.Description)

}

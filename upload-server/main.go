package main

import (
	"net"

	pb "github.com/baddin/upload-images/upload-service"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"google.golang.org/grpc/reflection"
)

const (
	//	host        = "localhost"
	//port        = 5423
	user        = "postgres"
	dbname      = "image_store"
	servicePort = ":50051"
)

type server struct{}

func (s *server) UpImage(ctx context.Context, in *pb.UpRequest) (*pb.UpResponse, error) {
	psqlInfo := fmt.Sprintf("user=%s dbname=%s sslmode=disable", user, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}
	sqlStatment := "INSERT INTO images(url, description) VALUES ($1, $2)"
	_, err = db.Exec(sqlStatment, in.Path, in.Description)
	if err != nil {
		panic(err)
	}

	return &pb.UpResponse{Created: true}, nil

}

func (s *server) GetImage(ctx context.Context, in *pb.GetRequest) (*pb.GetResponse, error) {
	psqlInfo := fmt.Sprintf("user=%s dbname=%s sslmode=disable", user, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}
	sqlStatment := "SELECT * FROM images WHERE ID=$1"
	rows, err := db.Query(sqlStatment)
	if err != nil {
		panic(err)
	}
	var id int
	var url string
	var description string
	err = rows.Scan(&id, &url, &description)
	if err != nil {
		panic(err)
	}
	return &pb.GetResponse{Id: int32(id), Url: url, Description: description}, nil

}

func main() {
	lis, err := net.Listen("tcp", servicePort)
	if err != nil {
		panic(err)
	}
	s := grpc.NewServer()
	pb.RegisterUploadImagesServer(s, &server{})
	reflection.Register(s)
	s.Serve(lis)
}

package main

import (
	"context"
	"google.golang.org/grpc"
	pb "grpc-project/gen/proto"
	"grpc-project/postgre"
	"log"
	"net"
)

type testApiServer struct {
	pb.UnimplementedTestApiServer
	postgresClient postgre.Client
}

func (s *testApiServer) Echo(ctx context.Context, req *pb.ResponseRequest) (*pb.ResponseRequest, error) {
	return req, nil
}

func (s *testApiServer) GetProduct(ctx context.Context, req *pb.ProductId) (*pb.Product, error) {
	err, product := s.postgresClient.GetProductByID(int(req.GetId()))
	if err != nil {
		return nil, err
	}
	p := pb.Product{
		Id:          int32(product.ID),
		Title:       product.Title,
		Description: product.Description,
		Price:       product.Price,
	}

	return &p, nil
}

func (s *testApiServer) GetProducts(ctx context.Context, req *pb.Empty) (*pb.ProductList, error) {
	err, products := s.postgresClient.GetProducts()
	if err != nil {
		return nil, err
	}
	var productList []*pb.Product
	for _, product := range products {
		p := pb.Product{
			Id:          int32(product.ID),
			Title:       product.Title,
			Description: product.Description,
			Price:       product.Price,
		}
		productList = append(productList, &p)
	}

	return &pb.ProductList{Product: productList}, nil
}

func main() {
	postgresConnection, err := postgre.NewPostgresConnection()
	if err != nil {
		panic(err)
	}
	defer postgresConnection.Close()

	listen, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		log.Fatalln(err)
	}

	grpcServer := grpc.NewServer()

	pb.RegisterTestApiServer(grpcServer, &testApiServer{postgresClient: postgresConnection})

	err = grpcServer.Serve(listen)

	if err != nil {
		log.Println(err)
	}

}

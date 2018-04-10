//consignment_service/main.go
package main

import (
	pb "github.com/liumenggc/containermp/consignment-service/proto/consignment"
	"golang.org/x/net/context"
	"net"
	"log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	port = ":50051"
)

type IRepository interface {
	Create( *pb.Consignment) (*pb.Consignment, error)
	GetAll() []*pb.Consignment
}

//Repository 模拟一个数据库，此后会用一个真正的数据库代替它.
type Repository struct {
	consignments []*pb.Consignment
}

func (repo *Repository) Create(consignment *pb.Consignment) (*pb.Consignment, error){
	updated := append(repo.consignments, consignment)
	repo.consignments = updated
	return consignment, nil
}

func (repo *Repository) GetAll() []*pb.Consignment {
	return repo.consignments
}

type service struct {
	repo IRepository
}

//CreateConsignment  在proto 中，我们只给微服务定义一个方法
//CreateConsignment 方法，它只接受一个context以及proto定义的
//Consignment消息，这个Consignment由gRPC的服务器处理后提供给你的
func (s *service)CreateConsignment(ctx context.Context, req *pb.Consignment) (*pb.Response, error) {
	//保存consignment
	consignment, err :=s.repo.Create(req)
	if err != nil {
		return nil, err
	}

	return &pb.Response{Created: true, Consignment: consignment}, nil
}

func (s *service) GetConsignment(ctx context.Context, req *pb.GetRequest) (*pb.Response, error) {
	consignments := s.repo.GetAll()
	return &pb.Response{Consignments : consignments},nil
}

func main()  {

	repo := &Repository{}
		//设置gRPC服务器
		lis, err := net.Listen("tcp", port)
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
	s := grpc.NewServer()

	//在我们的gRPC服务器上注册微服务，这将会我们的代码各*.pb.go
	//的各种Interface对应起来
	pb.RegisterShippingServiceServer(s, &service{repo})

	//在gRPC服务器上注册reflection
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
//consignment-cli/cli.go
package main

import (
	pb "github.com/liumenggc/containermp/consignment-service/proto/consignment"
	"io/ioutil"
	"encoding/json"
	"google.golang.org/grpc"
	"log"
	"os"
	"golang.org/x/net/context"
)

const (
	address   = "localhost:50051"
	defaultFilename = "consignment.json"
)

func parseFile(file string) (*pb.Consignment, error)  {
	var consigment *pb.Consignment
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	json.Unmarshal(data, &consigment)
	 return consigment, err
}

func main() {
	//set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Did not connect: %v", err)
	}
	defer conn.Close()
	client :=pb.NewShippingServiceClient(conn)

	//Contact the server and print out its response.
	file := defaultFilename
	if len(os.Args) > 1 {
		file = os.Args[1]
	}

	consignment, err := parseFile(file)
		if err != nil {
			log.Fatalf("Counld not parse file: %v", err)
		}
	r, err := client.CreateConsignment(context.Background(),consignment)
	if err != nil {
		log.Fatalf("Could not greet: %v", err)
	}
	log.Printf("Create: %t", r.Created)

	getAll, err := client.GetConsignments(context.Background(), &pb.GetRequest{})
	if err != nil {
		log.Fatalf("Counld not list consignment:%v", err)
	}
	for _, v:=range getAll.Consignments{
		log.Println(v)
	}
}


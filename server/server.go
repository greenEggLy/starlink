package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	pb "starlink/pb"
	"starlink/sat"

	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 8080, "The server port")
)

type server struct {
	pb.UnimplementedSatSysServer
}

func newServer() *server {
	s := &server{}
	return s
}

func (s *server) GetSatellites(ctx context.Context, in *pb.SearchContext) (*pb.Satellite, error) {
	sat := sat.GetSatBySysNameAndName(in.SysName, in.SatName)
	return &sat, nil
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v\n", err)
	}
	server := grpc.NewServer()
	pb.RegisterSatSysServer(server, newServer())
	log.Printf("server listening at %v\n", lis.Addr())
	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v\n", err)
	}
}

package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	pb "starlink/pb"
	"starlink/sat"
	"starlink/ssys"

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

func (s *server) GetSatellites(ctx context.Context, in *pb.SearchContext) (*pb.Satellite, error) {
	sat := sat.GetSatBySysNameAndName(in.SysName, in.SatName)
	return &sat[0], nil
}

func (s *server) CmdGetSystem(ctx context.Context, in *pb.CmdRequest) (*pb.CmdResponse, error) {
	response := parseCmdline(in.Cmd)
	var msg pb.CmdResponse
	fmt.Println(response)
	switch response.(type) {
	case pb.Satellite_System:
		msg.Message = response.(*pb.Satellite_System).String()
	case []pb.Satellite_System:
		for _, sys := range response.([]pb.Satellite_System) {
			msg.Message += sys.String() + "\n"
		}
	case []pb.Satellite:
		for _, sat := range response.([]pb.Satellite) {
			msg.Message += sat.String() + "\n"
		}
	default:
		msg.Message = "Invalid command"
	}
	return &msg, nil
}

func (s *server) UpdateSystem(ctx context.Context, in *pb.UpdateContext) (*pb.UpdateResponse, error) {
	ssys.Update_System(in.SysName)
	return &pb.UpdateResponse{Message: "Update system successfully"}, nil
}

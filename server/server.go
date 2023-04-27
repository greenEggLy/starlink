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
	return &sat, nil
}

func (s *server) CmdGetSystem(ctx context.Context, in *pb.CmdRequest) (*pb.CmdResponse, error) {
	ret := parseCmdline(in.Cmd)
	var msg pb.CmdResponse
	if ret.RetType == 0 {
		msg.Message = ret.OneSat.String()
	} else if ret.RetType == 1 {
		for _, sat := range ret.Sats {
			msg.Message += sat.String() + "\n"
		}
	} else if ret.RetType == 2 {
		msg.Message = ret.OneSys.String()
	} else if ret.RetType == 3 {
		for _, sys := range ret.Syss {
			msg.Message += sys.String() + "\n"
		}
	} else {
		msg.Message = "Invalid command"
	}
	return &msg, nil
}

func (s *server) UpdateSystem(ctx context.Context, in *pb.UpdateContext) (*pb.UpdateResponse, error) {
	ssys.Update_System(in.SysName)
	return &pb.UpdateResponse{Message: "Update system successfully"}, nil
}

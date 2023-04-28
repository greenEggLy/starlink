package main

import (
	"context"
	"flag"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "starlink/pb"
)

var (
	serverAddr = flag.String("addr", "localhost:8080", "The server address in the format of host:port")
)

func main() {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	conn, err := grpc.Dial(*serverAddr, opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	client := pb.NewSatSysClient(conn)

	// Search for a satellite
	printSatellite(client, pb.SearchContext{SatName: "HJ-1A", SysName: "DMC"})

	printSys(client, pb.CmdRequest{Cmd: "get DMC"})
}

func printSatellite(client pb.SatSysClient, names pb.SearchContext) {
	sat, err := client.GetSatellites(context.Background(), &names)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	log.Printf("Satellite: %v\n", sat)
}

func printSys(clent pb.SatSysClient, cmdLine pb.CmdRequest) {
	response, err := clent.CmdGetSystem(context.Background(), &cmdLine)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	log.Printf("System: %v\n", response)
}

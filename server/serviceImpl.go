package main

import (
	"context"
	pb "starlink/pb"
	"starlink/ssys"
)

func (s *server) UpdateSystem(ctx context.Context, in *pb.UpdateContext) (*pb.UpdateResponse, error) {
	ssys.Update_System(in.SysName)
	return &pb.UpdateResponse{Message: "Update system successfully"}, nil
}

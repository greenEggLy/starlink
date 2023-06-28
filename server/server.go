package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"sync"
	"time"

	pb "starlink/pb"
	"starlink/utils"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

var (
	// port2 = flag.Int("port2", 8080, "The server port communicating with unity")
	port = flag.Int("port", 8081, "The server port")
)
var findNewTarget = make(chan bool, 10)

// a timeout timer as a channel
// var timeout = make(chan bool, 10)

type server struct {
	pb.UnimplementedSatComServer
	mu          sync.RWMutex
	findTarget  bool
	satNotes    *utils.ExpiredMap
	redisClient *utils.Redis
}

func newServer() *server {
	s := &server{
		satNotes:    utils.NewExpiredMap(),
		redisClient: utils.NewRedis(60),
	}
	return s
}

func createBasePos() pb.PositionInfo {
	pos := pb.PositionInfo{
		Timestamp: fmt.Sprint(time.Now().Unix()),
		Alt:       30,
		Lat:       30,
		Lng:       30,
	}
	return pos
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v\n", err)
	}
	server := grpc.NewServer()
	pb.RegisterSatComServer(server, newServer())
	log.Printf("server listening at %v\n", lis.Addr())
	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v\n", err)
	}
}

// server <-> Unity
func (s *server) ReceiveFromUnityTemplate(stream pb.SatCom_ReceiveFromUnityTemplateServer) error {
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		errStatus, _ := status.FromError(err)
		if errStatus.Code() == 1 {
			log.Printf("stream cancelled")
			return nil
		}
		if err != nil {
			log.Fatalf("failed to receive: %v\n", err)
		}
		// if unity find target, then save in redis
		findNewTarget <- in.FindTarget
		go func() {
			find := <-findNewTarget
			log.Printf("receive from unity, %v", in.TargetPosition)
			if find || !find && s.findTarget {
				// if someone finds target, then put it in redis
				if find {
					s.findTarget = true
				}
				s.mu.Lock()
				for _, v := range in.TargetPosition {
					s.redisClient.SetPosition(v)
				}
				s.mu.Unlock()
				sats := s.satNotes.GetAll()
				notes := s.redisClient.GetAllPos()
				if len(notes) == 0 {
					s.findTarget = false
				}
				msg := pb.Base2UnityInfo{
					FindTarget:     s.findTarget,
					TargetPosition: notes,
					TrackingSat:    sats,
				}
				stream.Send(&msg)
			} else {
				msg := pb.Base2UnityInfo{
					FindTarget:     false,
					TargetPosition: nil,
					TrackingSat:    nil,
				}
				stream.Send(&msg)
			}
		}()

	}
}

// server <-> Satellite
func (s *server) CommuWizSat(stream pb.SatCom_CommuWizSatServer) error {
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		errStatus, _ := status.FromError(err)
		if errStatus.Code() == 1 {
			log.Printf("stream cancelled")
			return nil
		}
		if err != nil {
			return err
		}
		// whenever receive a new message from channel find_new_target, send message to client
		log.Printf("receive from satellite, %v", in.TargetPosition)
		if s.findTarget {
			satInfo := &pb.SatelliteInfo{
				SatName:     in.SatName,
				SatPosition: in.SatPosition,
			}
			s.satNotes.Set(in.SatName, *satInfo, int64(60))
			findNewTarget <- true
		} else {
			s.satNotes.Delete(in.SatName)
			findNewTarget <- false
		}

		go func() {
			find := <-findNewTarget
			p := createBasePos()
			var msg pb.Base2SatInfo
			if find || s.findTarget {
				if find {
					s.findTarget = true
				}
				notes := s.redisClient.GetAllPos()
				if len(notes) == 0 {
					s.findTarget = false
				}
				msg = pb.Base2SatInfo{
					FindTarget:     s.findTarget,
					BasePosition:   &p,
					TargetPosition: notes,
				}
				err = stream.Send(&msg)
				if err != nil {
					log.Printf("[server]send message error")
					grpc.WithReturnConnectionError()
				}
			} else {
				msg = pb.Base2SatInfo{
					FindTarget:     false,
					BasePosition:   &p,
					TargetPosition: nil,
				}
				err = stream.Send(&msg)
				if err != nil {
					log.Printf("[server]send message error")
					grpc.WithReturnConnectionError()
				}
			}
		}()

	}
}

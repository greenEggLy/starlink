package main

import (
	"context"
	"encoding/binary"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"time"

	"starlink/pb"
	"starlink/utils"
	async "starlink/utils/async"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

var (
	// port2 = flag.Int("port2", 8080, "The server port communicating with unity")
	port = flag.Int("port", 50051, "The server port")
)
var findNewTarget = make(chan bool, 10)

// server holds whether the system is tracking a target
// and the information of satellites and targets
// expiredMap is thread-safe
type server struct {
	pb.UnimplementedSatComServer
	// mu          sync.RWMutex
	findTarget       bool
	trackingSatNotes *utils.ExpiredMap[string, []string]    // [target_name, [satellite_name]]
	tarNotes         *utils.ExpiredMap[string, string]      // [target_name, target_name]
	photoNotes       *utils.ExpiredMap[string, chan []byte] // [zone_info, channel]
	redisClient      *utils.Redis
}

func newServer() *server {
	s := &server{
		findTarget:       false,
		trackingSatNotes: utils.NewExpiredMap[string, []string](),
		tarNotes:         utils.NewExpiredMap[string, string](),
		photoNotes:       utils.NewExpiredMap[string, chan []byte](),
		redisClient:      utils.NewRedis(60),
	}
	return s
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

func (s *server) CommuWizUnity(request *pb.UnityRequest, stream pb.SatCom_CommuWizUnityServer) error {
	if !request.StatusOk {
		return nil
	}
	log.Printf("[unity] display\n")
	// a ticker sending unmessage every half second
	// a timer for 10 seconds
	ticker := time.NewTicker(500 * time.Millisecond)
	timer := time.NewTimer(10 * time.Second)
	defer ticker.Stop()
	defer timer.Stop()
	for {
		select {
		case <-ticker.C:
			msg := s.createBase2UnityMsg(s.findTarget)
			// log.Printf("msg: %v", msg)
			err := stream.Send(&msg)
			if err == io.EOF {
				return nil
			}
			errStatus, _ := status.FromError(err)
			if errStatus.Code() == 1 {
				log.Printf("request cancelled")
				return nil
			}
			if err != nil {
				log.Printf("failed to send: %v\n", err)
				return nil
			}
		case <-timer.C:
			return nil
		}
	}
}

// server <-> Unity [target]
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
		findNewTarget <- in.FindTarget || s.findTarget
		go func() {
			find := <-findNewTarget
			log.Printf("[unity] target")
			if in.FindTarget {
				// handle information from unity
				s.findTarget = true
				// save tracking target information
				targetNames := getAllTargetNames(in.TargetPosition)
				for _, v := range targetNames {
					s.tarNotes.Set(v, v, 60)
				}
				setOperations := async.Exec(func() bool {
					for _, v := range in.TargetPosition {
						s.redisClient.SetTarPos(v)
					}
					return true
				})
				_ = setOperations.Await()
			}

			msg := s.createBase2UnityMsg(find)
			err := stream.Send(&msg)

			if err != nil {
				log.Printf("[server]send to unity error, %v", err)
				grpc.WithReturnConnectionError()
			}
		}()

	}
}

// server <-> Unity [photo]
func (s *server) SendPhotos(request *pb.UnityPhotoRequest, stream pb.SatCom_SendPhotosServer) error {
	timeoutSec := 5

	zoneInfo := request.GetZone()
	// check if the request is duplicate
	find, _ := s.photoNotes.Get(zoneInfo.String())
	if find {
		return nil
	}
	// create a channel to receive photo
	photoChan := make(chan []byte, 1)
	s.photoNotes.Set(zoneInfo.String(), photoChan, int64(timeoutSec))
	// wait for the channel and send it to unity
	timer := time.NewTimer(time.Second * time.Duration(timeoutSec))
	select {
	case <-timer.C:
		// delete the element in photoNotes
		log.Printf("[unity] photo timeout")
		s.photoNotes.Delete(zoneInfo.String())
		err := stream.Send(&pb.BasePhotoResponse{
			Timestamp: getTimeStamp(),
			ImageData: nil,
		})
		return err
	case photo := <-photoChan:
		log.Printf("[unity] receive photo, %s", hex.EncodeToString(photo))
		s.photoNotes.Delete(zoneInfo.String())
		err := stream.Send(&pb.BasePhotoResponse{
			Timestamp: getTimeStamp(),
			ImageData: photo,
		})
		if err != nil {
			log.Printf("send photo error, %v", err)
			return err
		}
	}
	return nil
}

// server <-> Satellite [target]
func (s *server) CommuWizSat(stream pb.SatCom_CommuWizSatServer) error {
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			errStatus, _ := status.FromError(err)
			if errStatus.Code() == 1 {
				log.Printf("stream cancelled")
				return nil
			}
			return err
		}
		// whenever receive a new message from channel find_new_target, send message to client
		log.Printf("[satellite] target")

		s.findTarget = in.FindTarget || s.findTarget
		findNewTarget <- in.FindTarget || s.findTarget

		go func() {
			find := <-findNewTarget

			// save satellite information to redis
			// satInfo := in.SatInfo
			s.redisClient.SetSatPos(*in.SatInfo)

			// save info in satellite message
			if in.FindTarget {
				satInfo := in.SatInfo
				for _, targetInfo := range in.TargetInfo {
					// save predicted target information
					s.redisClient.SetTarPos(targetInfo)
					// save target information
					s.tarNotes.Set(targetInfo.TargetName, targetInfo.TargetName, int64(60))
					// save tracking satellite information
					exist, value := s.trackingSatNotes.Get(targetInfo.TargetName)
					if exist {
						for index, name := range value {
							if name == satInfo.SatName {
								break
							}
							if index == len(value)-1 {
								value = append(value, satInfo.SatName)
								s.trackingSatNotes.Set(targetInfo.TargetName, value, int64(60))
							}
						}
					} else {
						s.trackingSatNotes.Set(targetInfo.TargetName, []string{satInfo.SatName}, int64(60))
					}
				}
			}

			msg := s.createBase2SatMsg(find)
			err := stream.Send(&msg)

			if err != nil {
				log.Printf("[server]send message error")
				grpc.WithReturnConnectionError()
			}
		}()

	}
}

func (s *server) TakePhotos(ctx context.Context, request *pb.SatPhotoRequest) (*pb.BasePhotoReceiveResponse, error) {
	log.Printf("[satellite] photo")
	// get zone information
	zoneInfo := request.GetZone()
	if binary.Size(request.ImageData) <= 0 {
		return &pb.BasePhotoReceiveResponse{
			Timestamp:    getTimeStamp(),
			ReceivePhoto: false,
		}, nil
	}
	// check if other satellite has took the photo
	check, channel := s.photoNotes.GetAndDelete(zoneInfo.String())
	if check {
		channel <- request.ImageData
		return &pb.BasePhotoReceiveResponse{
			Timestamp:    getTimeStamp(),
			ReceivePhoto: true,
		}, nil
	}
	return &pb.BasePhotoReceiveResponse{
		Timestamp:    getTimeStamp(),
		ReceivePhoto: true,
	}, nil
}

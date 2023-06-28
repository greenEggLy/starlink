package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "starlink/pb"
)

var (
	serverAddr = flag.String("addr", "localhost:8081", "The server address in the format of host:port")

	randomGenerator = rand.New(rand.NewSource(time.Now().UnixNano()))

	satelliteMsgNum = 0
	unityMsgNum     = 0
)

func main() {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	conn, err := grpc.Dial(*serverAddr, opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	client := pb.NewSatComClient(conn)

	postAndReceive(client)
}

func postAndReceive(client pb.SatComClient) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	sat_stream, err := client.CommuWizSat(ctx)
	if err != nil {
		log.Fatalf("satellite-base flow failed: %v", err)
	}
	unity_stream, err := client.ReceiveFromUnityTemplate(ctx)
	if err != nil {
		log.Fatalf("unity-base flow failed: %v", err)
	}
	// two tickers
	satTicker := time.NewTicker(1 * time.Second)
	defer satTicker.Stop()
	unityTicker := time.NewTicker(1 * time.Second)
	defer unityTicker.Stop()
	// one timer
	timeoutTimer := time.NewTimer(5 * time.Second)

	// wait for both stream end
	waitc := make(chan struct{})
	waitc2 := make(chan struct{})
	// sat-base
	go func() error {
		for {
			in, err := sat_stream.Recv()
			if err == io.EOF {
				close(waitc)
				return nil
			}
			if err != nil {
				log.Fatalf("satellite-base flow failed: %v", err)
				return err
			}
			if in.FindTarget {
				// find new target
				// judge if the target is in the range
				// ...
				log.Printf("[sat]:target in horizon")
			} else {
				log.Printf("[sat]:no target in horizon\n")
			}
		}
	}()
	// unity-base
	go func() error {
		for {
			in, err := unity_stream.Recv()
			if err == io.EOF {
				close(waitc2)
				return nil
			}
			if err != nil {
				log.Fatalf("unity-base flow failed: %v", err)
				return err
			}
			if in.FindTarget {
				// find target
				// show warning on screen
				// ...
				names := getSatNames(in.TrackingSat)

				log.Printf("[unity]:find target, position num: %d, satellites: %v\n", len(in.TargetPosition), names)
			} else {
				log.Printf("[unity]:no target")
			}
		}
	}()

	// send message
	go func(satt, unit *time.Ticker, timer *time.Timer, client pb.SatComClient) {
		for {
			select {
			case <-satt.C:
				// send satellite position info to server
				var msg pb.Sat2BaseInfo
				if satelliteMsgNum < 0 {
					msg = pb.Sat2BaseInfo{
						SatName:        fmt.Sprintf("%s%d", "satellite-", satelliteMsgNum),
						SatPosition:    generateOnePos(0),
						FindTarget:     false,
						TargetPosition: nil,
					}
				} else {
					msg = pb.Sat2BaseInfo{
						SatName:        fmt.Sprintf("%s%d", "satellite-", satelliteMsgNum),
						SatPosition:    generateOnePos(0),
						FindTarget:     true,
						TargetPosition: generateRandomTargetPos(3),
					}
				}
				satelliteMsgNum++
				if err := sat_stream.Send(&msg); err != nil {
					log.Fatalf("satellite-base flow failed\n")
				}

			case <-unit.C:
				// send unity position info to server
				var msg pb.Unity2BaseInfoTemplate
				if unityMsgNum < 0 {
					msg = pb.Unity2BaseInfoTemplate{
						FindTarget:     false,
						TargetPosition: nil,
					}
				} else {
					msg = pb.Unity2BaseInfoTemplate{
						FindTarget:     true,
						TargetPosition: generateRandomTargetPos(3),
					}
				}
				unityMsgNum++
				if err := unity_stream.Send(&msg); err != nil {
					log.Fatalf("satellite-base flow failed\n")
				}
			case <-timeoutTimer.C:
				close(waitc)
				close(waitc2)
				unit.Stop()
				satt.Stop()
				sat_stream.CloseSend()
				unity_stream.CloseSend()
				return
			}
		}
	}(satTicker, unityTicker, timeoutTimer, client)

	<-waitc
	<-waitc2
}

func generateOnePos(i int64) *pb.PositionInfo {
	pos := pb.PositionInfo{
		Timestamp: fmt.Sprint(time.Now().Unix() + i),
		Alt:       randomGenerator.Float32(),
		Lng:       randomGenerator.Float32(),
		Lat:       randomGenerator.Float32(),
	}
	return &pos
}

func generateRandomTargetPos(num int) []*pb.PositionInfo {
	ret := []*pb.PositionInfo{}
	for i := 0; i < num; i++ {
		ret = append(ret, generateOnePos(int64(i)))
	}
	return ret
}

func getSatNames(list []*pb.SatelliteInfo) []string {
	ret := []string{}
	for _, v := range list {
		ret = append(ret, v.SatName)
	}
	return ret
}

// func sendSatInfoAndGetResponse(client pb.SatComClient) {
// 	pos := pb.PositionInfo{
// 		PositionX: 10,
// 		PositionY: 10,
// 		PositionZ: 10,
// 	}
// 	posNotes := []*pb.PositionInfo{
// 		&pos,
// 		&pos,
// 		&pos,
// 		&pos,
// 		&pos,
// 	}
// 	cli_pos := pb.PositionInfo{
// 		PositionX: 20,
// 		PositionY: 20,
// 		PositionZ: 20,
// 	}
// 	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
// 	defer cancel()

// 	stream, err := client.CommuWizSat(ctx)
// 	if err != nil {
// 		log.Fatalf("client.RouteChat failed: %v", err)
// 	}
// 	stream2, err := client.ReceiveFromUnityTemplate(ctx)
// 	if err != nil {
// 		log.Fatalf("client.RouteChat failed: %v", err)
// 	}
// 	waitc := make(chan struct{})
// 	waitc2 := make(chan struct{})
// 	// satwaitc := make(chan struct{})
// 	// unity act as a client and send position info to server
// 	go func() {
// 		for {
// 			in, err := stream2.Recv()
// 			if err == io.EOF {
// 				close(waitc2)
// 				return
// 			}
// 			if err != nil {
// 				log.Fatalf("client.RouteChat failed: %v", err)
// 			}
// 			if in.FindTarget {
// 				log.Printf("find target: target info: %v\n", in.TargetPosition)
// 			} else {
// 				log.Printf("no target found\n")
// 			}
// 		}
// 	}()
// 	go func() {
// 		for {
// 			in, err := stream.Recv()
// 			if err == io.EOF {
// 				// read done.
// 				close(waitc)
// 				return
// 			}
// 			if err != nil {
// 				log.Fatalf("client.RouteChat failed: %v", err)
// 			}
// 			log.Printf("get message from base")
// 			if in.FindTarget {
// 				log.Printf(", target at (%f, %f, %f)\n", in.BasePosition.PositionX, in.BasePosition.PositionY, in.BasePosition.PositionZ)
// 			}
// 		}
// 	}()
// 	for i := 0; i < 10; i++ {
// 		msg := pb.Unity2BaseInfoTemplate{
// 			Timestamp:      time.Now().UTC().Format(time.RFC3339),
// 			FindTarget:     true,
// 			TargetPosition: posNotes,
// 		}
// 		if err := stream2.Send(&msg); err != nil {
// 			log.Panic("send msg error!\n")
// 		}
// 		// put true to channel to indicate that the unity has sent the message
// 		// satwaitc <- struct{}{}
// 	}

// 	for i := 0; i < 10; i++ {
// 		// wait for satwaitc
// 		// <-satwaitc

// 		msg := pb.Sat2BaseInfo{
// 			SatName:        "haha",
// 			SatPosition:    &cli_pos,
// 			FindTarget:     true,
// 			TargetPosition: posNotes,
// 		}
// 		if err := stream.Send(&msg); err != nil {
// 			log.Panic("send msg error!\n")
// 		}
// 		log.Printf("send message to base")
// 	}
// 	stream.CloseSend()
// 	<-waitc
// 	<-waitc2
// }

// func printSatellite(client pb.SatSysClient, names pb.SearchContext) {
// 	sat, err := client.GetSatellites(context.Background(), &names)
// 	if err != nil {
// 		log.Fatalf("fail to dial: %v", err)
// 	}
// 	log.Printf("Satellite: %v\n", sat)
// }

// func printSys(clent pb.SatSysClient, cmdLine pb.CmdRequest) {
// 	response, err := clent.CmdGetSystem(context.Background(), &cmdLine)
// 	if err != nil {
// 		log.Fatalf("fail to dial: %v", err)
// 	}
// 	log.Printf("System: %v\n", response)
// }

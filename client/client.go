package main

import (
	"context"
	"flag"
	"io"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "starlink/pb"
	cli "starlink/utils/client"
)

var (
	serverAddr = flag.String("addr", "localhost:50051", "The server address in the format of host:port")
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
	// needPhotoRequestNum := 0
	// one timer
	timeoutTimer := time.NewTimer(30 * time.Second)

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
				// do something
				log.Printf("[sat]:target in horizon")
			} else {
				log.Printf("[sat]:no target in horizon\n")
			}

			// should take a photo and send
			if in.TakePhoto {
				image := make([]byte, 0)
				image = append(image, "photo"...)

				photoRequest := pb.SatPhotoRequest{
					Timestamp: cli.GetTimeStamp(),
					SatInfo:   cli.GenerateSatInfo(),
					Zone:      in.Zone[0],
					ImageData: image,
				}
				// send a photo request to server
				// in the rpc service TakePhotos
				for {
					checkMsg, err := client.TakePhotos(ctx, &photoRequest)
					if err != nil {
						log.Fatalf("satellite-base flow failed: %v", err)
					}
					if checkMsg.ReceivePhoto {
						break
					}
				}
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
				names := cli.GetSatNames(in.TrackingSat)

				log.Printf("[unity]:find target, position num: %d, satellite num: %d\n", len(in.TargetPosition), len(names))
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
				msg := cli.CreateSatRequest(3)
				if err := sat_stream.Send(msg); err != nil {
					log.Fatalf("satellite-base flow failed\n")
				}

			case <-unit.C:
				// send unity position info to server
				msg := cli.CreateUnityRequestTemplate(3)
				if err := unity_stream.Send(msg); err != nil {
					log.Fatalf("satellite-base flow failed\n")
				}
				request := pb.UnityPhotoRequest{
					Timestamp: cli.GetTimeStamp(),
					Zone:      cli.GenerateZoneInfo(),
				}
				go func() {
					photoIn, err := client.SendPhotos(ctx, &request)
					if err != nil {
						log.Fatalln("get photo error")
					}
					response, err := photoIn.Recv()
					if err != nil {
						log.Fatalln("get photo error")
					}
					photo := response.ImageData
					log.Printf("[unity] receive photo, %v", photo)
				}()

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

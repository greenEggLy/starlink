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
	cli "starlink/test/client/client"
)

var (
	serverAddr = flag.String("addr", "localhost:50051", "The server address in the format of host:port")
	timeout    = 30 * time.Second
	interval   = 1 * time.Second
)

func main() {
	for {
		var opts []grpc.DialOption
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
		conn, err := grpc.Dial(*serverAddr, opts...)
		if err != nil {
			log.Fatalf("fail to dial: %v", err)
		}
		client := pb.NewSatComClient(conn)
		postAndReceive(client)
		conn.Close()
	}
}

func postAndReceive(client pb.SatComClient) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	sat_stream, err := client.CommuWizSat(ctx)
	if err != nil {
		log.Printf("satellite-base flow failed: %v", err)
		return err
	}
	satTicker := time.NewTicker(interval)
	defer satTicker.Stop()
	timeoutTimer := time.NewTimer(timeout)
	waitc := make(chan struct{})
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

			if in.TakePhoto {
				// should take a photo and send
				go func(in *pb.Base2Sat) {
					photo := cli.LoadPhoto()
					image := make([]byte, 0)
					image = append(image, photo...)

					photoRequest := pb.SatPhotoRequest{
						Timestamp: cli.GetTimeStamp(),
						SatInfo:   cli.GenerateSatInfo(),
						Zone:      in.Zone[0],
						ImageData: image,
					}
					// send a photo request to server
					// in the rpc service TakePhotos
					for {
						log.Printf("[sat]:take a photo\n")
						checkMsg, err := client.TakePhotos(ctx, &photoRequest)
						log.Printf("[sat]:send a photo request to server\n")

						if err != nil {
							log.Fatalf("satellite-base flow failed: %v", err)
							return
						}
						if checkMsg.ReceivePhoto {
							return
						}
					}
				}(in)
			}
		}
	}()
	go func(satt *time.Ticker, timer *time.Timer, client pb.SatComClient) {
		for {
			select {
			case <-satt.C:
				// send satellite position info to server
				msg := cli.CreateSatRequest(3)
				if err := sat_stream.Send(msg); err != nil {
					log.Fatalf("satellite-base flow failed\n")
				}

			case <-timeoutTimer.C:
				close(waitc)
				satt.Stop()
				sat_stream.CloseSend()
				return
			}
		}
	}(satTicker, timeoutTimer, client)
	<-waitc
	return nil
}

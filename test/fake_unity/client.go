package main

import (
	"context"
	"flag"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "starlink/pb"
	cli "starlink/test/client/client"
)

var (
	serverAddr = flag.String("addr", "43.142.83.201:50051", "The server address in the format of host:port")
	timeout    = 60 * time.Second
	interval   = 2 * time.Second
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

func postAndReceive(client pb.SatComClient) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	// create a new Ticker of 10 secends, and a Timer for 60 seconds
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	timer := time.NewTimer(timeout)
	defer timer.Stop()
	waitC := make(chan struct{})
	go func() {
		select {
		case <-ticker.C:
			// send message to server and read the photo info
			zone := cli.GenerateZoneInfo()
			msg := pb.UnityPhotoRequest{
				Timestamp: cli.GetTimeStamp(),
				Zone:      zone,
			}
			go func() {
				log.Printf("Send photo message")
				photo, err := client.SendPhotos(ctx, &msg)
				log.Printf("Send photo message end")
				if err != nil {
					log.Fatalf("send photo failed: %v", err)
				}
				log.Printf("Photo: %v", photo)
			}()
		case <-timer.C:
			waitC <- struct{}{}
			log.Printf("Timeout")
			return
		}
	}()
	<-waitC
	return nil
}

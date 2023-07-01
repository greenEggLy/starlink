package main

import (
	"context"
	"flag"
	"io"
	"log"
	"math/rand"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "starlink/pb"
)

var (
	serverAddr      = flag.String("addr", "localhost:50051", "The server address in the format of host:port")
	randomGenerator = rand.New(rand.NewSource(time.Now().UnixNano()))
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
				msg := createSatRequest(3)
				if err := sat_stream.Send(msg); err != nil {
					log.Fatalf("satellite-base flow failed\n")
				}

			case <-unit.C:
				// send unity position info to server
				msg := createUnityRequestTemplate(3)
				if err := unity_stream.Send(msg); err != nil {
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

package main

import (
	"fmt"
	pb "starlink/pb"
	"starlink/utils"
)

func main() {
	zone := pb.ZoneInfo{
		UpperLeft: &pb.LLPosition{
			Lat: 30,
			Lng: 120,
		},
		BottomRight: &pb.LLPosition{
			Lat: 20,
			Lng: 110,
		},
	}
	str := zone.String()

	fmt.Printf("%v", utils.String2ZoneInfo(str))
}

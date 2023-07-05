package main

import (
	"fmt"
	pb "starlink/pb"
	"starlink/utils"
)

func main() {
	zone := pb.ZoneInfo{
		RequestIdentify: true,
		UpperLeft: &pb.LLPosition{
			Timestamp: "333",
			Lat:       30.11,
			Lng:       120,
		},
		BottomRight: &pb.LLPosition{
			Timestamp: "444",
			Lat:       20,
			Lng:       110.222,
		},
	}
	str := zone.String()

	fmt.Printf("%v\n", str)

	fmt.Printf("%v", utils.String2ZoneInfo(str))
}

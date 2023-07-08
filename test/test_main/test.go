package main

import (
	"encoding/base64"
	"fmt"
	pb "starlink/pb"
)

func main() {
	zone1 := &pb.ZoneInfo{
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
	zone2 := &pb.ZoneInfo{
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
	bytes1, err := zone1.XXX_Marshal(nil, false)
	if err != nil {
		fmt.Println(err)
	}
	bytes2, err := zone2.XXX_Marshal(nil, false)
	if err != nil {
		fmt.Println(err)
	}
	str1 := base64.StdEncoding.EncodeToString(bytes1)
	str2 := base64.StdEncoding.EncodeToString(bytes2)
	// judge whether bytes1 equals to bytes2
	if str1 == str2 {
		fmt.Println("equal")
	} else {
		fmt.Printf("%v\n", str2)
		fmt.Println("not equal")
	}

	b1, err := base64.StdEncoding.DecodeString(str1)
	if err != nil {
		fmt.Println(err)
	}
	b2, err := base64.StdEncoding.DecodeString(str2)
	if err != nil {
		fmt.Println(err)
	}
	zone3 := &pb.ZoneInfo{}
	zone4 := &pb.ZoneInfo{}
	err = zone3.XXX_Unmarshal(b1)
	if err != nil {
		fmt.Println(err)
	}
	err = zone4.XXX_Unmarshal(b2)
	if err != nil {
		fmt.Println(err)
	}
	if zone3.String() == zone1.String() {
		fmt.Println("equal")
	} else {
		fmt.Println("not equal")
	}
	// fmt.Printf("%v", zone.String())
	// fmt.Printf("%v", zone2.String())

}

package utils

import (
	"log"
	"regexp"
	pb "starlink/pb"
	"strconv"
)

func ZoneInfos2Strings(zoneInfos []*pb.ZoneInfo) []string {
	strings := make([]string, 0)
	for _, zoneInfo := range zoneInfos {
		strings = append(strings, zoneInfo.String())
	}
	return strings
}

func String2ZoneInfo(str string) *pb.ZoneInfo {
	shouldIdentify := false
	reBool := regexp.MustCompile(`request_identify:(true|false)`)
	matchBool := reBool.FindStringSubmatch(str)

	if matchBool[1] == "true" {
		shouldIdentify = true
	} else {
		shouldIdentify = false
	}

	re := regexp.MustCompile(`[-+]?\d+(\.\d+)?`)
	matches := re.FindAllString(str, -1)

	var numbers []float64

	if len(matches) != 6 {
		log.Printf("zoneInfo format wrong, %v", str)
		return nil
	}

	timestamp1 := matches[0]
	timestamp2 := matches[3]
	for index, match := range matches {
		if index == 0 || index == 3 {
			continue
		}
		num, err := strconv.ParseFloat(match, 64)
		if err == nil {
			numbers = append(numbers, num)
		}
	}
	ul_lat := numbers[0]
	ul_lng := numbers[1]
	br_lat := numbers[2]
	br_lng := numbers[3]
	zone := pb.ZoneInfo{
		RequestIdentify: shouldIdentify,
		UpperLeft: &pb.LLPosition{
			Timestamp: timestamp1,
			Lat:       float32(ul_lat),
			Lng:       float32(ul_lng),
		},
		BottomRight: &pb.LLPosition{
			Timestamp: timestamp2,
			Lat:       float32(br_lat),
			Lng:       float32(br_lng),
		},
	}
	return &zone
}

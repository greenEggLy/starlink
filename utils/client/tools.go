package client

import (
	"fmt"
	"math/rand"
	"starlink/pb"
	"strconv"
	"time"
)

var randomGenerator = rand.New(rand.NewSource(time.Now().UnixNano()))

func CreateSatRequest(predictSec int) *pb.SatRequest {
	satInfo := pb.SatelliteInfo{
		SatName:     "satellite-1",
		SatPosition: GenerateOneLLAPos(0),
	}

	msg := pb.SatRequest{
		SatInfo:    &satInfo,
		FindTarget: true,
		TargetInfo: GenerateRandomTargetPos(predictSec),
	}

	return &msg
}

func CreateUnityRequestTemplate(predictSec int) *pb.UnityRequestTemplate {
	msg := pb.UnityRequestTemplate{
		FindTarget:     true,
		TargetPosition: GenerateRandomTargetPos(predictSec),
	}
	return &msg
}

func GenerateOneLLAPos(timeStampSec int64) *pb.LLAPosition {
	pos := pb.LLAPosition{
		Timestamp: fmt.Sprint(time.Now().Unix() + timeStampSec),
		Alt:       randomGenerator.Float32(),
		Lng:       randomGenerator.Float32(),
		Lat:       randomGenerator.Float32(),
	}
	return &pos
}

func GenerateOneLLPos(timeStampSec int64) *pb.LLPosition {
	pos := pb.LLPosition{
		Timestamp: fmt.Sprint(time.Now().Unix() + timeStampSec),
		Lng:       randomGenerator.Float32(),
		Lat:       randomGenerator.Float32(),
	}
	return &pos
}

func GenerateSatInfo() *pb.SatelliteInfo {
	info := pb.SatelliteInfo{
		SatName:     "satellite-1",
		SatPosition: GenerateOneLLAPos(0),
	}
	return &info
}

func GenerateRandomTargetPos(num int) []*pb.TargetInfo {
	ret := []*pb.TargetInfo{}
	for i := 0; i < num; i++ {
		ele := pb.TargetInfo{
			TargetName:     "target1",
			TargetPosition: GenerateOneLLAPos(int64(i)),
		}
		ret = append(ret, &ele)
	}
	return ret
}

func GenerateZoneInfo() *pb.ZoneInfo {
	zone := pb.ZoneInfo{
		UpperLeft:   GenerateOneLLPos(0),
		BottomRight: GenerateOneLLPos(0),
	}
	return &zone
}

func GetSatNames(list []*pb.SatelliteInfo) []string {
	ret := []string{}
	for _, v := range list {
		ret = append(ret, v.SatName)
	}
	return ret
}

func GetTimeStamp() string {
	bytes := make([]byte, 0)
	bytes = append(bytes, strconv.FormatInt(time.Now().Unix(), 10)...)
	return string(bytes)
}

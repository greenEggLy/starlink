package main

import (
	"fmt"
	"starlink/pb"
	"strconv"
	"time"
)

func createSatRequest(predictSec int) *pb.SatRequest {
	satInfo := pb.SatelliteInfo{
		SatName:     "satellite-1",
		SatPosition: generateOneLLAPos(0),
	}

	msg := pb.SatRequest{
		SatInfo:    &satInfo,
		FindTarget: true,
		TargetInfo: generateRandomTargetPos(predictSec),
	}

	return &msg
}

func createUnityRequestTemplate(predictSec int) *pb.UnityRequestTemplate {
	msg := pb.UnityRequestTemplate{
		FindTarget:     true,
		TargetPosition: generateRandomTargetPos(predictSec),
	}
	return &msg
}

func generateOneLLAPos(timeStampSec int64) *pb.LLAPosition {
	pos := pb.LLAPosition{
		Timestamp: fmt.Sprint(time.Now().Unix() + timeStampSec),
		Alt:       randomGenerator.Float32(),
		Lng:       randomGenerator.Float32(),
		Lat:       randomGenerator.Float32(),
	}
	return &pos
}

func generateOneLLPos(timeStampSec int64) *pb.LLPosition {
	pos := pb.LLPosition{
		Timestamp: fmt.Sprint(time.Now().Unix() + timeStampSec),
		Lng:       randomGenerator.Float32(),
		Lat:       randomGenerator.Float32(),
	}
	return &pos
}

func generateSatInfo() *pb.SatelliteInfo {
	info := pb.SatelliteInfo{
		SatName:     "satellite-1",
		SatPosition: generateOneLLAPos(0),
	}
	return &info
}

func generateRandomTargetPos(num int) []*pb.TargetInfo {
	ret := []*pb.TargetInfo{}
	for i := 0; i < num; i++ {
		ele := pb.TargetInfo{
			TargetName:     "target1",
			TargetPosition: generateOneLLAPos(int64(i)),
		}
		ret = append(ret, &ele)
	}
	return ret
}

func generateZoneInfo() *pb.ZoneInfo {
	zone := pb.ZoneInfo{
		UpperLeft:   generateOneLLPos(0),
		BottomRight: generateOneLLPos(0),
	}
	return &zone
}

func getSatNames(list []*pb.SatelliteInfo) []string {
	ret := []string{}
	for _, v := range list {
		ret = append(ret, v.SatName)
	}
	return ret
}

func getTimeStamp() string {
	bytes := make([]byte, 0)
	bytes = append(bytes, strconv.FormatInt(time.Now().Unix(), 10)...)
	return string(bytes)
}

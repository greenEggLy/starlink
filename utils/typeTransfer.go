package utils

import (
	"encoding/base64"
	"log"
	pb "starlink/pb"
)

func ZoneInfo2String(zoneInfo *pb.ZoneInfo) string {
	bytes, err := zoneInfo.XXX_Marshal(nil, false)
	if err != nil {
		log.Printf("ZoneInfo2String error: %v", err)
		return ""
	}
	return base64.StdEncoding.EncodeToString(bytes)
}

func ZoneInfos2Strings(zoneInfos []*pb.ZoneInfo) []string {
	strings := make([]string, 0)
	for _, zoneInfo := range zoneInfos {
		strings = append(strings, ZoneInfo2String(zoneInfo))
	}
	return strings
}

func String2ZoneInfo(str string) *pb.ZoneInfo {
	bytes, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		log.Printf("String2ZoneInfo error: %v", err)
		return nil
	}
	zoneInfo := &pb.ZoneInfo{}
	err = zoneInfo.XXX_Unmarshal(bytes)
	if err != nil {
		log.Printf("String2ZoneInfo error: %v", err)
		return nil
	}
	return zoneInfo
}

func SatelliteInfo2String(satelliteInfo *pb.SatelliteInfo) string {
	bytes, err := satelliteInfo.XXX_Marshal(nil, false)
	if err != nil {
		log.Printf("SatelliteInfo2String error: %v", err)
		return ""
	}
	return base64.StdEncoding.EncodeToString(bytes)
}

func String2SatelliteInfo(str string) *pb.SatelliteInfo {
	bytes, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		log.Printf("String2SatelliteInfo error: %v", err)
		return nil
	}
	satelliteInfo := &pb.SatelliteInfo{}
	err = satelliteInfo.XXX_Unmarshal(bytes)
	if err != nil {
		log.Printf("String2SatelliteInfo error: %v", err)
		return nil
	}
	return satelliteInfo
}

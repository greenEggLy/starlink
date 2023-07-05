package main

import (
	"fmt"
	pb "starlink/pb"
	"starlink/utils"
	"starlink/utils/async"
	"strconv"
	"time"
)

// utils
func getAllTargetNames(targets []*pb.TargetInfo) []string {
	var names []string
	for _, v := range targets {
		names = append(names, v.TargetName)
	}
	names = utils.RemoveDuplicateElement(names)
	return names
}

func (s *server) createBase2UnityMsg(hasTracking bool) pb.Base2Unity {
	if !hasTracking {
		msg := pb.Base2Unity{
			FindTarget:       false,
			TargetPosition:   nil,
			TargetSatellites: nil,
		}
		return msg
	}
	// return information
	targets := s.tarNotes.GetAll()                // target names
	target_sats := s.trackingSatNotes.GetMapAll() // tracking satellites
	notes := s.redisClient.GetAllTarPos(targets)

	if len(notes) == 0 || len(target_sats) == 0 {
		s.findTarget = false
		msg := pb.Base2Unity{
			FindTarget:       false,
			TargetPosition:   nil,
			TargetSatellites: nil,
		}
		return msg
	}

	targetSatellitesMap := make(map[string]*pb.TrackingSatellites)
	for k, v := range target_sats {
		pos := s.redisClient.GetSelectedSatPos(v)
		targetSatellitesMap[k] = &pb.TrackingSatellites{
			Satellite: pos,
		}
	}

	msg := pb.Base2Unity{
		FindTarget:       true,
		TargetPosition:   notes,
		TargetSatellites: targetSatellitesMap,
	}
	return msg
}

func (s *server) createBase2SatMsg(hasTracking bool) pb.Base2Sat {
	basePosition := createBasePos()
	var takePhoto = false
	var zone []*pb.ZoneInfo

	if s.photoNotes.Size() > 0 {
		zones := s.photoNotes.GetAllKeys()
		// if zone type is pb.ZoneInfo, then append to msg
		// else then do nothing
		takePhoto = true
		for _, v := range zones {
			zone = append(zone, utils.String2ZoneInfo(v))
		}

	}

	if !hasTracking {
		msg := pb.Base2Sat{
			FindTarget:   false,
			BasePosition: &basePosition,
			TargetInfo:   nil,
			TakePhoto:    takePhoto,
			Zone:         zone,
		}
		return msg
	}

	targets := s.tarNotes.GetAll()
	positionNotes := async.Exec(func() []*pb.TargetInfo {
		return s.redisClient.GetAllTarPos(targets)
	})
	notes := positionNotes.Await()

	if len(notes) == 0 {
		s.findTarget = false
		msg := pb.Base2Sat{
			FindTarget:   false,
			BasePosition: &basePosition,
			TargetInfo:   nil,
			TakePhoto:    takePhoto,
			Zone:         zone,
		}
		return msg
	} else {
		msg := pb.Base2Sat{
			FindTarget:   true,
			BasePosition: &basePosition,
			TargetInfo:   notes,
			TakePhoto:    takePhoto,
			Zone:         zone,
		}
		return msg
	}
}

func createBasePos() pb.LLPosition {
	pos := pb.LLPosition{
		Timestamp: fmt.Sprint(time.Now().Unix()),
		Lat:       30,
		Lng:       30,
	}
	return pos
}

func getTimeStamp() string {
	bytes := make([]byte, 0)
	bytes = append(bytes, strconv.FormatInt(time.Now().Unix(), 10)...)
	return string(bytes)
}

package prefabs

import "starlink/pb"

type CmdRet struct {
	OneSat  pb.Satellite
	Sats    []pb.Satellite
	OneSys  pb.Satellite_System
	Syss    []pb.Satellite_System
	RetType int32
}

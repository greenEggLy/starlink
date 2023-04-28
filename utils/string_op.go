package utils

import (
	pb "starlink/pb"
	"strings"
)

func decode(tle [3]string) pb.Satellite {
	TLE1 := tle[1]
	TLE2 := tle[2]
	var result pb.Satellite

	result.Name = strings.TrimSpace(tle[0])
	result.NumSat = TLE1[2:7]
	result.Inter = TLE1[9:17]
	result.Year = TLE1[18:20]
	result.Day = TLE1[20:32]
	if TLE1[33] == '-' {
		result.FirstMotion = "-0"
	} else {
		result.FirstMotion = "0"
	}
	result.FirstMotion += TLE1[34:44]
	result.SecondMotion = TLE1[44:52]
	result.Drag = TLE1[53:61]
	result.Number = strings.TrimSpace(TLE1[64:68])
	result.Incl = strings.TrimSpace(TLE2[8:16])
	result.RA = strings.TrimSpace(TLE2[17:25])
	result.Eccentricity = "0." + TLE2[26:33]
	result.ArgPer = strings.TrimSpace(TLE2[34:42])
	result.Anomaly = strings.TrimSpace(TLE2[43:51])
	result.Motion = strings.TrimSpace(TLE2[52:63])
	result.Epoch = strings.TrimSpace(TLE2[63:68])

	return result
}

func trim_sys_name(sys string) string {
	sys = strings.TrimSpace(sys)
	sys = sys[1 : len(sys)-1]
	return sys
}

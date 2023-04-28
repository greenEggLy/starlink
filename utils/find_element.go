package utils

import (
	"starlink/pb"
	"starlink/prefabs"
)

func FindByName(val interface{}, args ...interface{}) (int, bool) {
	sys, ok := val.(prefabs.System)
	if ok {
		for i, item := range args {
			if item.(prefabs.System).NAME == sys.NAME {
				return i, true
			}
		}
		return -1, false
	}
	pbSys, ok := val.(pb.Satellite_System)
	pbSysSet, ok2 := args[0].([]pb.Satellite_System)
	if ok && ok2 {
		for i, item := range pbSysSet {
			if item.Name == pbSys.Name {
				return i, true
			}
		}
		return -1, false

	}
	return -1, false
}

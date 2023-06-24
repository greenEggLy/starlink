package main

import (
	"starlink/sat"
	"starlink/ssys"
	"strings"
)

func parseCmdline(cmd string) interface{} {
	args := strings.Fields(cmd)
	switch args[0] {
	case "update":
		ssys.Update_System(args[1])
		return nil
	case "getsys":
		if args[1] == "-a" {
			return ssys.GetAllSys()
		} else {
			return ssys.GetOneSys(args[1])
		}
	case "get":
		if len(args) == 2 {
			return sat.GetSatsBySysName(args[1])
		} else if len(args) == 3 {
			return sat.GetSatBySysNameAndName(args[1], args[2])
		}
	default:
	}
	return nil
}

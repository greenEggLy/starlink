package main

import (
	"starlink/prefabs"
	"starlink/sat"
	"starlink/ssys"
	"strings"
)

func parseCmdline(cmd string) prefabs.CmdRet {
	args := strings.Fields(cmd)
	var ret prefabs.CmdRet
	switch args[0] {
	case "update":
		ssys.Update_System(args[1])
		ret.RetType = -1
	case "getsys":
		if args[1] == "-a" {
			ret.Syss = ssys.GetAllSys()
			ret.RetType = 3
		} else {
			ret.OneSys = ssys.GetOneSys(args[1])
			ret.RetType = 2
		}

	case "get":
		ret.OneSat = sat.GetSatBySysNameAndName(args[1], args[2])
		ret.RetType = 0
	default:
	}
	return ret
}

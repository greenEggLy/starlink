package ssys

import (
	"starlink/pb"
	"starlink/utils"
)

func Update_System(sys_name string) {
	utils.Fetch_System_U(sys_name)
}

func GetAllSys() []pb.Satellite_System {
	return utils.GetAllSys_U()
}

func GetOneSys(sys_name string) pb.Satellite_System {
	return utils.GetOneSys_U(sys_name)
}

package sat

import (
	pb "starlink/pb"
	"starlink/utils"
)

func GetSatBySysNameAndName(sys_name, name string) []pb.Satellite {
	return utils.GetSatBySysNameAndName_U(sys_name, name)
}

func GetSatsBySysName(sys_name string) []pb.Satellite {
	system := utils.GetOneSys_U(sys_name)
	return utils.GetSatsBySysName_U(system)
}

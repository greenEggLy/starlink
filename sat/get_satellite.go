package sat

import (
	pb "starlink/pb"
	"starlink/utils"
)

func GetSatBySysNameAndName(sys_name, name string) pb.Satellite {
	return utils.GetSatBySysNameAndName_U(sys_name, name)
}

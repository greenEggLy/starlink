package utils

import (
	"database/sql"
	"fmt"
	"net/http"
	"starlink/globaldata"
	"starlink/prefabs"

	pb "starlink/pb"

	"github.com/gin-gonic/gin"
)

// api
func GetSatBySysIdAndName_HTTP_U(c *gin.Context) {
	sys_id_s := c.Param("sys_id")
	name := c.Param("name")
	var res pb.Satellite = get_sat_by_sysid_name(sys_id_s, name)
	var json_res = convert_pb_to_prefabs(res)
	if res.Name != "" {
		c.IndentedJSON(http.StatusOK, json_res)
		// c.IndentedJSON(http.StatusOK, res)
		return
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "satellite not found"})
}

func GetSatBySysNameAndName_U(sys_name, name string) pb.Satellite {
	fmt.Printf("GetSatBySysNameAndName_U: %s, %s\n", sys_name, name)
	sys_id_s := get_sysid_by_sysname(sys_name)
	// if not found, return an empty satellite
	if sys_id_s == "-1" {
		return pb.Satellite{}
	}
	// if found, query the database
	return get_sat_by_sysid_name(sys_id_s, name)
}

// private

func convert_pb_to_prefabs(sat pb.Satellite) prefabs.Satellite {
	var res prefabs.Satellite
	res.Name = sat.Name
	res.NumSat = sat.NumSat
	res.Inter = sat.Inter
	res.Year = sat.Year
	res.Day = sat.Day
	res.FirstMotion = sat.FirstMotion
	res.SecondMotion = sat.SecondMotion
	res.Drag = sat.Drag
	res.Number = sat.Number
	res.Incl = sat.Incl
	res.RA = sat.RA
	res.Eccentricity = sat.Eccentricity
	res.ArgPer = sat.ArgPer
	res.Anomaly = sat.Anomaly
	res.Motion = sat.Motion
	res.Epoch = sat.Epoch
	return res
}

func get_sysid_by_sysname(sys_name string) string {
	for _, sys := range globaldata.System_Info {
		if sys.NAME == sys_name {
			return sys.ID
		}
	}
	return "-1"
}

func get_sat_by_sysid_name(sys_id, name string) pb.Satellite {
	var res pb.Satellite
	db, err := sql.Open("mysql", "root:@tcp(localhost:3306)/starlink")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	query := "select name,num_sat,inter,year,day,first_motion,second_motion,drag,number,incl,r_a,eccentricity,arg_per,anomaly,motion,epoch from satellite where sys_id=? AND name=? ;"
	stmt, err := db.Prepare(query)
	if err != nil {
		panic(err.Error())
	}
	defer stmt.Close()
	stmt.QueryRow(sys_id, name).Scan(&res.Name, &res.NumSat, &res.Inter, &res.Year, &res.Day,
		&res.FirstMotion, &res.SecondMotion, &res.Drag, &res.Number,
		&res.Incl, &res.RA, &res.Eccentricity, &res.ArgPer, &res.Anomaly,
		&res.Motion, &res.Epoch)
	if res.Name != "" {
		return res
	}
	return pb.Satellite{}
}

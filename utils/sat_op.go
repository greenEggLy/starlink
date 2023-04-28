package utils

import (
	"database/sql"
	"fmt"
	"strconv"

	pb "starlink/pb"
)

func GetSatBySysNameAndName_U(sys_name, name string) []pb.Satellite {
	sys := GetOneSys_U(sys_name)
	// if not found, return an empty satellite
	if sys.Name == "" {
		return []pb.Satellite{}
	}
	// if found, query the database
	return get_sat_by_sysid_name(sys.GetId(), name)
}

// private
// func convert_pb_to_prefabs(sat pb.Satellite) prefabs.Satellite {
// 	var res prefabs.Satellite
// 	res.Name = sat.Name
// 	res.NumSat = sat.NumSat
// 	res.Inter = sat.Inter
// 	res.Year = sat.Year
// 	res.Day = sat.Day
// 	res.FirstMotion = sat.FirstMotion
// 	res.SecondMotion = sat.SecondMotion
// 	res.Drag = sat.Drag
// 	res.Number = sat.Number
// 	res.Incl = sat.Incl
// 	res.RA = sat.RA
// 	res.Eccentricity = sat.Eccentricity
// 	res.ArgPer = sat.ArgPer
// 	res.Anomaly = sat.Anomaly
// 	res.Motion = sat.Motion
// 	res.Epoch = sat.Epoch
// 	return res
// }

func get_sat_by_sysid_name(sys_id, name string) []pb.Satellite {
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
		var ress []pb.Satellite
		ress = append(ress, res)
		return ress
	}
	return []pb.Satellite{}
}

// get all satellites of a system
func GetSatsBySysName_U(sys pb.Satellite_System) []pb.Satellite {
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/starlink")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	sys_index, err := strconv.Atoi(sys.GetId())
	if err != nil {
		panic(err.Error())
	}
	var sats []pb.Satellite
	query := "SELECT * FROM satellite WHERE sys_id =?"
	symt, err := db.Prepare(query)
	if err != nil {
		panic(err.Error())
	}
	defer symt.Close()
	rows, err := symt.Query(sys_index)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(rows)
	for rows.Next() {
		var id string
		var sys_id string
		var sat pb.Satellite

		err = rows.Scan(&id, &sat.Name, &sat.NumSat, &sat.Inter,
			&sat.Year, &sat.Day, &sat.FirstMotion, &sat.SecondMotion,
			&sat.Drag, &sat.Number, &sat.Incl, &sat.RA, &sat.Eccentricity,
			&sat.ArgPer, &sat.Anomaly, &sat.Motion, &sat.Epoch, &sys_id)

		if err != nil {
			panic(err.Error())
		}
		sats = append(sats, sat)
		// sats = append(sats, sat)
		// fmt.Println(sat)
	}
	defer rows.Close()
	return sats
}

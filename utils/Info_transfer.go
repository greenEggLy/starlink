package utils

import (
	"database/sql"
	"starlink/prefabs"

	_ "github.com/go-sql-driver/mysql"
)

func Get_Satellites(tle []string, sys_index int64) []prefabs.Satellite {
	db, err := sql.Open("mysql", "root:@localhost(127.0.0.1:3306)/starlink")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	var sat_sys prefabs.System
	db.QueryRow("SELECT id, name from system where id = ?", sys_index).Scan(&sat_sys.ID, &sat_sys.NAME)
	if err != nil {
		panic(err.Error())
	}
	var sats []prefabs.Satellite
	rows, err := db.Query("get_sat_info WHERE sys_id = ?", sys_index)
	index := 0
	if err != nil {
		panic(err.Error())
	}
	for rows.Next() {
		rows.Scan(&sats[index].Name, &sats[index].NumSat, &sats[index].Inter,
			&sats[index].Year, &sats[index].Day, &sats[index].FirstMotion, &sats[index].SecondMotion,
			&sats[index].Drag, &sats[index].Number, &sats[index].Incl, &sats[index].RA, &sats[index].Eccentricity,
			&sats[index].ArgPer, &sats[index].Anomaly, &sats[index].Motion, &sats[index].Epoch)
		index += 1
	}
	defer rows.Close()
	return sats
}

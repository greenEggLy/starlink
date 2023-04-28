package utils

import (
	"database/sql"
	"starlink/pb"
)

func GetAllSys_U() []pb.Satellite_System {
	var res []pb.Satellite_System
	db, err := sql.Open("mysql", "root:@tcp(localhost:3306)/starlink")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	query := "SELECT * from systems"
	symt, err := db.Prepare(query)
	if err != nil {
		panic(err.Error())
	}
	defer symt.Close()
	rows, err := symt.Query()
	if err != nil {
		panic(err.Error())
	}
	defer rows.Close()
	for rows.Next() {
		var sat_sys pb.Satellite_System
		err = rows.Scan(&sat_sys.Id, &sat_sys.Name)
		if err != nil {
			panic(err.Error())
		}
		res = append(res, sat_sys)
	}
	return res
}

func GetOneSys_U(sys_name string) pb.Satellite_System {
	var res pb.Satellite_System
	db, err := sql.Open("mysql", "root:@tcp(localhost:3306)/starlink")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	query := "SELECT * from systems WHERE name =?"
	symt, err := db.Prepare(query)
	if err != nil {
		panic(err.Error())
	}
	defer symt.Close()

	symt.QueryRow(sys_name).Scan(&res.Id, &res.Name)
	return res
}

func InsertSys_U(sys_name string) {
	db, err := sql.Open("mysql", "root:@tcp(localhost:3306)/starlink")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	query := "INSERT INTO systems (name) VALUES (?);"
	symt, err := db.Prepare(query)
	if err != nil {
		panic(err.Error())
	}
	defer symt.Close()

	_, err = symt.Exec(sys_name)
	if err != nil {
		panic(err.Error())
	}
}

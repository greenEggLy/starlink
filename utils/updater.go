package utils

import (
	"bufio"
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"os"
	"starlink/globaldata"
	"starlink/prefabs"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)


	


func decode(tle [3]string) prefabs.Satellite{
	TLE1 := tle[1]
	TLE2 := tle[2]
	var result prefabs.Satellite

	result.NAME = strings.TrimSpace(tle[0])
	result.NumSat = TLE1[2:7]
	result.Inter = TLE1[9:17]
	result.Year = TLE1[18:20]
	result.Day = TLE1[20:32]
	if(TLE1[33] == '-'){
		result.FirstMotion = "-0"
	}else{
		result.FirstMotion = "0"
	}
	result.FirstMotion += TLE1[34:44]
	result.SecondMotion = TLE1[44:52]
	result.Drag = TLE1[53:61]
	result.Number = strings.TrimSpace(TLE1[64:68]) 
	result.Incl = strings.TrimSpace(TLE2[8:16])
	result.R_A = strings.TrimSpace(TLE2[17:25])
	result.Eccentricity = "0." + TLE2[26:33]
	result.ArgPer = strings.TrimSpace(TLE2[34:42])
	result.Anomaly = strings.TrimSpace(TLE2[43:51])
	result.Motion = strings.TrimSpace(TLE2[52:63])
	result.Epoch = strings.TrimSpace(TLE2[63:68])

	return result

}

// api
func Update_System(sys_name string) {
	path, _ :=os.Getwd()
	url := "http://celestrak.com/NORAD/elements/" + sys_name + ".txt"
	filepath := path + "/Data/" + sys_name + ".txt"

	// get raw resp from url
	resp, err := http.Get(url)
	if err != nil {
		 panic(err.Error())
	}
	defer resp.Body.Close()

	// open the file that will be written
	out, err := os.Open(filepath)
	if err != nil {
		panic(err.Error())
	}
	defer out.Close()
	
	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		panic(err.Error())
	}else{
		fmt.Print("load"+sys_name)
	}
}


// api
// or used after update data from internet
// update one system's info
func Update_data(sys_id int){
	lines := read_file(sys_id)
	sat_num := len(lines)/3
	fmt.Printf("sat_num: %d\n",sat_num)
	var tle [3]string
	// one_satellite stores tmp data, and use this to update database
	var one_satellite prefabs.Satellite 
	for i := 0; i < sat_num; i++{ // i is the satellite index
		tle[0] = lines[i*3]
		tle[1] = lines[i*3+1]
		tle[2] = lines[i*3+2]
		one_satellite = decode(tle)
		if len(globaldata.System_Info[sys_id-1].SysTLE) <= i{
			globaldata.System_Info[sys_id-1].SysTLE = append(globaldata.System_Info[sys_id].SysTLE, tle)
		}else{
			globaldata.System_Info[sys_id-1].SysTLE[i] = tle
		}
		// fmt.Printf("%s: %s, %s, %s\n",one_satellite.NAME, one_satellite.Anomaly, one_satellite.NumSat, one_satellite.Day)
		update_database(sys_id,one_satellite)
	}
}

// read file util
func read_file(sys_id int) []string{
	path, _ :=os.Getwd()
	fp,err := os.Open(path + "/Data/" + globaldata.System_Info[sys_id-1].NAME + ".txt")
	// fp,err := os.Open(path + "/Data/2012-044.txt")
	if err != nil {
        panic(err.Error())
    }
	defer fp.Close()
	buf := bufio.NewScanner(fp)
	var lines []string
	for {
		// exit if the file is over
		if !buf.Scan() {
			break
		}
		lines = append(lines, buf.Text()) // read every line
	}
	return lines
}

// use update_database in update_data
// write new information to database
func update_database(sys_id int, sat prefabs.Satellite){
	db, err := sql.Open("mysql","root:@tcp(localhost:3306)/starlink")
	if err != nil {
        panic(err.Error())
    }
	defer db.Close()
	insert_query := "insert into satellite(name, num_sat, inter, year, day, first_motion, second_motion, drag, number, incl, r_a, eccentricity, arg_per, anomaly, motion, epoch, sys_id)"
	value_query := " values(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
	on_duplicate := " on duplicate key"
	update_query := " update name = ?, num_sat = ?, inter = ?, year = ?, day = ?, first_motion = ?, second_motion = ?, drag = ?, number = ?, incl = ?, r_a = ?, eccentricity =?, arg_per = ?,anomaly = ?, motion = ?, epoch = ?;"

	query := insert_query + value_query + on_duplicate + update_query
	stmt,err :=db.Prepare(query)
	if err != nil {
        panic(err.Error())
    }
	defer stmt.Close()
	res,err:=stmt.Exec(sat.NAME, sat.NumSat, sat.Inter, sat.Year, sat.Day, sat.FirstMotion,sat.SecondMotion, sat.Drag,
		sat.Number,sat.Incl, sat.R_A, sat.Eccentricity, sat.ArgPer, sat.Anomaly, sat.Motion, sat.Epoch, sys_id,
		sat.NAME, sat.NumSat, sat.Inter, sat.Year, sat.Day, sat.FirstMotion,sat.SecondMotion, sat.Drag,
		sat.Number,sat.Incl, sat.R_A, sat.Eccentricity, sat.ArgPer, sat.Anomaly, sat.Motion, sat.Epoch)	
		
	if err != nil {
        panic(err.Error())
    }else{
		res.LastInsertId();
	}
	stmt.Close()
}

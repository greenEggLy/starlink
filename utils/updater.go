package utils

import (
	"bufio"
	"database/sql"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"

	pb "starlink/pb"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

// api
// fetch information from website and update the database
func Fetch_System_U(sys_name string) {
	// if the system is not exist on the website, then we cannot fetch it
	sys_from_web := get_sys_from_website()
	var value pb.Satellite_System
	value.Name = sys_name
	_, found := FindByName(value, sys_from_web)
	if !found {
		fmt.Printf("system %s not on the website\n", sys_name)
		return
	}

	// check if the system is exist in local database
	tar_sys := GetOneSys_U(sys_name)
	if tar_sys.Name == "" {
		// if the system is not exist, then insert
		InsertSys_U(sys_name)
		tar_sys = GetOneSys_U(sys_name)
	}

	// write to file and update database
	write_file(sys_name)
	update_data(tar_sys)
}

// read from file
// update database
func update_data(system pb.Satellite_System) {
	// read file
	lines := read_file(system.GetName())
	sat_num := len(lines) / 3
	fmt.Printf("sat_num: %d\n", sat_num)
	var tle [3]string
	var one_satellite pb.Satellite

	// decode and update
	for sat_index := 0; sat_index < sat_num; sat_index++ { // i is the satellite index
		tle[0] = lines[sat_index*3]
		tle[1] = lines[sat_index*3+1]
		tle[2] = lines[sat_index*3+2]
		one_satellite = decode(tle)
		sys_id, err := strconv.Atoi(system.GetId())
		if err != nil {
			log.Fatal(err)
		}
		update_database(sys_id, one_satellite)
	}
}

// get accessible system from website
func get_sys_from_website() []pb.Satellite_System {
	var res []pb.Satellite_System
	cmd := exec.Command("python3", "utils/soup/get_all_sys.py")
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
	out = out[1 : len(out)-2]
	for _, sys := range strings.Split(string(out), ",") {
		if sys != "" {
			sys = trim_sys_name(sys)
			res = append(res, pb.Satellite_System{Name: sys})
		}
	}
	return res
}

// read file util
func read_file(sys_name string) []string {
	path, _ := os.Getwd()
	fp, err := os.Open(path + "/Data/" + sys_name + ".txt")
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

// write to file
func write_file(sys_name string) {
	// fetch the system
	path, _ := os.Getwd()
	url := "http://celestrak.com/NORAD/elements/" + sys_name + ".txt"
	filepath := path + "/Data/" + sys_name + ".txt"

	// get raw resp from url
	resp, err := http.Get(url)
	if err != nil {
		panic(err.Error())
	}
	defer resp.Body.Close()

	// open the file that will be written
	out, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		panic(err.Error())
	}
	defer out.Close()

	// write the response body to the file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		panic(err.Error())
	}
}

// use update_database in update_data
// write new information to database
func update_database(sys_id int, sat pb.Satellite) {
	db, err := sql.Open("mysql", "root:@tcp(localhost:3306)/starlink")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	insert_query := "insert into satellite(name, num_sat, inter, year, day, first_motion, second_motion, drag, number, incl, r_a, eccentricity, arg_per, anomaly, motion, epoch, sys_id)"
	value_query := " values(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
	on_duplicate := " on duplicate key"
	update_query := " update name = ?, num_sat = ?, inter = ?, year = ?, day = ?, first_motion = ?, second_motion = ?, drag = ?, number = ?, incl = ?, r_a = ?, eccentricity =?, arg_per = ?,anomaly = ?, motion = ?, epoch = ?;"

	query := insert_query + value_query + on_duplicate + update_query
	stmt, err := db.Prepare(query)
	if err != nil {
		panic(err.Error())
	}
	defer stmt.Close()
	res, err := stmt.Exec(sat.Name, sat.NumSat, sat.Inter, sat.Year, sat.Day, sat.FirstMotion, sat.SecondMotion, sat.Drag,
		sat.Number, sat.Incl, sat.RA, sat.Eccentricity, sat.ArgPer, sat.Anomaly, sat.Motion, sat.Epoch, sys_id,
		sat.Name, sat.NumSat, sat.Inter, sat.Year, sat.Day, sat.FirstMotion, sat.SecondMotion, sat.Drag,
		sat.Number, sat.Incl, sat.RA, sat.Eccentricity, sat.ArgPer, sat.Anomaly, sat.Motion, sat.Epoch)

	if err != nil {
		panic(err.Error())
	} else {
		res.LastInsertId()
	}
	stmt.Close()
}

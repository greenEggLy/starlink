package getter

import (
	"database/sql"
	"fmt"
	"net/http"
	"starlink/prefabs"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetSatBySysIdAndName(c *gin.Context){
	sys_id_s := c.Param("sys_id")
	name := c.Param("name")
	sys_id,_:= strconv.ParseInt(sys_id_s,10,64)
	fmt.Printf("name: %s, sys_id: %d\n",name, sys_id)
	var res prefabs.Satellite
	db, err := sql.Open("mysql","root:@tcp(localhost:3306)/starlink")
	if err != nil {
        panic(err.Error())
    }
	defer db.Close()

	query := "select name,num_sat,inter,year,day,first_motion,second_motion,drag,number,incl,r_a,eccentricity,arg_per,anomaly,motion,epoch from satellite where sys_id=? AND name=? ;"
	stmt,err :=db.Prepare(query)
	if err != nil {
        panic(err.Error())
    }
	defer stmt.Close()
	stmt.QueryRow(sys_id_s, name).Scan(&res.NAME, &res.NumSat, &res.Inter, &res.Year, &res.Day, 
									   &res.FirstMotion, &res.SecondMotion, &res.Drag, &res.Number,
									   &res.Incl, &res.R_A, &res.Eccentricity, &res.ArgPer, &res.Anomaly,
									   &res.Motion, &res.Epoch)
	if(res.NAME != ""){
		c.IndentedJSON(http.StatusOK, res)
		return
	}
    c.IndentedJSON(http.StatusNotFound, gin.H{"message": "satellite not found"})
}
package main

// import "fmt"
import (
	// "fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"

	"starlink/getter"
	"starlink/prefabs"
)


var satellites = []prefabs.Satellite{
	{	 NAME: "0",NumSat: "0",Inter: "1",Year: "2021",Day: "122",FirstMotion: "0.2",SecondMotion: "0.3",Drag: "22",
		Number: "2",Incl: "2121",R_A: "222",Eccentricity: "21",ArgPer: "i9",Anomaly: "fjdsl",Motion: "left",Epoch: "2"},
}

	



func main() {
    router := gin.Default()
    router.GET("/starlink/get_test", get_test_satellite)
	router.GET("/starlink/get_satb_sidname/:sys_id/:name",getter.GetSatBySysIdAndName)

    router.Run("localhost:8080")

	// sys_id = 9 name = HAWK-A for test
	// fmt.Print(lines)
}


func get_test_satellite(c *gin.Context){
	c.IndentedJSON(http.StatusOK, satellites)
}

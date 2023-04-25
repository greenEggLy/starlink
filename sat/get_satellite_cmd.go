package sat

import (
	"fmt"

	"github.com/spf13/cobra"
)

var sys_name string
var sat_name string
var updateSys = &cobra.Command{
	Use:     "get",
	Aliases: []string{"g"},
	Short:   "get one satellite",
	Long:    `get one satellite by system name and satellite name`,
	Run: func(cmd *cobra.Command, args []string) {
		res := GetSatBySysNameAndName(sys_name, sat_name)
		fmt.Printf("%v\n", res)
	},
}

func init() {
	rootCmd.AddCommand(updateSys)
	updateSys.Flags().StringVarP(&sys_name, "sysname", "y", "", "satellite system name")
	updateSys.Flags().StringVarP(&sat_name, "satname", "s", "", "satellite name")
}

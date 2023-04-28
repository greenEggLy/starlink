package ssys

import (
	"fmt"

	"github.com/spf13/cobra"
)

var sys_name string
var updateSys = &cobra.Command{
	Use:     "update",
	Aliases: []string{"u"},
	Short:   "update satellite system",
	Long:    `update satellite system from the internet`,
	Run: func(cmd *cobra.Command, args []string) {
		if sys_name == "" {
			fmt.Println("please input satellite system name")
			return
		}
		Update_System(sys_name)
	},
}

func init() {
	rootCmd.AddCommand(updateSys)
	updateSys.Flags().StringVarP(&sys_name, "sysname", "y", "", "satellite system name")
}

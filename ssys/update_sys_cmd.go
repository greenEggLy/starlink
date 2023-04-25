package ssys

import (
	"fmt"

	"starlink/globaldata"

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
		find := false
		for _, sys := range globaldata.System_Info {
			if sys.NAME == sys_name {
				Update_System(sys_name)
				find = true
				break
			}
		}
		if !find {
			fmt.Printf("satellite system not found!\n")
		}
	},
}

func init() {
	rootCmd.AddCommand(updateSys)
	updateSys.Flags().StringVarP(&sys_name, "sysname", "y", "", "satellite system name")
}

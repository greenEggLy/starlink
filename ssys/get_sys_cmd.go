package ssys

import (
	"fmt"

	"github.com/spf13/cobra"
)

var get_all_sys bool
var get_sys = &cobra.Command{
	Use:   "getsys",
	Short: "get satellite system",
	Long:  `get satellite system from database`,
	Run: func(cmd *cobra.Command, args []string) {
		if get_all_sys && len(args) > 0 {
			fmt.Printf("wrong command! usage: getsys -a or getsys <system_name>\n")
		} else if get_all_sys {
			sys_set := GetAllSys()
			fmt.Printf("systems:\n%v\n", sys_set)
		} else {
			one_sys := GetOneSys(args[0])
			if one_sys.Name != "" {
				fmt.Printf("system:\n%v\n", one_sys)
			} else {
				fmt.Printf("no such system!\n")
			}

		}
	},
}

func init() {
	rootCmd.AddCommand(get_sys)
	get_sys.Flags().BoolVarP(&get_all_sys, "getall", "a", false, "whether to get all sys")
}

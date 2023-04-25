package ssys

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "satellite-system",
	Aliases: []string{"ssys"},
	Short:   "ssys",
	Long: `a satellite system
   
One can use stringer to modify, add or update satellite systems from the terminal`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("this is satellite system\n")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Whoops. There was an error while executing your CLI '%s'", err)
		os.Exit(1)
	}
}

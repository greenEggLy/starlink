package sat

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "satellite",
	Aliases: []string{"sat"},
	Short:   "operations to satellites",
	Long: `lalala I'm satellite!
   
One can use stringer to modify, add or update satellites from the terminal`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("this is satellite\n")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Whoops. There was an error while executing your CLI '%s'", err)
		os.Exit(1)
	}
}

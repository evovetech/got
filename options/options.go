package options

import "github.com/spf13/cobra"

var (
	Verbose bool
)

func AddTo(cmd *cobra.Command) {
	pflags := cmd.PersistentFlags()
	pflags.BoolVarP(&Verbose, "verbose", "v", false, "verbose")
}

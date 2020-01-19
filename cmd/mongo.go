package cmd

import (
	"github.com/spf13/cobra"

	"github.com/sv-tools/grade/pkg/driver/mongodb"
)

var mongoCmd = &cobra.Command{
	Use:   "mongo [flags] [file ...]",
	Short: "Store the benchmarks in MongoDB",
	Long:  `Driver to store the benchmarks in a Mongo Database`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return mongodb.Execute(&cfg)
	},
}

func init() {
	AddCommonFlags(mongoCmd, &cfg)
	AddDBFlags(mongoCmd, &cfg, "collection")

	RootCmd.AddCommand(mongoCmd)
}

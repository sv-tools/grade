package cmd

import (
	"github.com/sv-go-tools/grade/pkg/driver/json"

	"github.com/spf13/cobra"
)

var jsonCmd = &cobra.Command{
	Use:   "json [flags] [file ...]",
	Short: "Print the benchmarks in Json format",
	Long:  `Driver to print the benchmarks in the Json format.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return json.Execute(&cfg)
	},
}

func init() {
	AddCommonFlags(jsonCmd, &cfg)
	jsonCmd.PersistentFlags().IntVar(&cfg.JSONIndent, "indent", 2, "Json indent")

	RootCmd.AddCommand(jsonCmd)
}

package cmd

import (
	"github.com/spf13/cobra"

	"github.com/sv-go-tools/grade/pkg/driver/influxdb"
)

var influxCmd = &cobra.Command{
	Use:   "influx",
	Short: "Store the benchmarks in InfluxDB",
	Long:  `Driver to store the benchmarks in a Influx Database`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return influxdb.Execute(&cfg)
	},
}

func init() {
	AddCommonFlags(influxCmd, &cfg)
	AddDBFlags(influxCmd, &cfg, "measurement")
	influxCmd.PersistentFlags().BoolVarP(&cfg.Insecure, "insecure", "i", false, "Skip SSL verification if set")

	RootCmd.AddCommand(influxCmd)
}

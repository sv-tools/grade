package cmd

import (
	"github.com/sv-go-tools/grade/pkg/driver"
	"github.com/sv-go-tools/grade/pkg/driver/influxdb"

	"github.com/spf13/cobra"
)

var influxCmd = &cobra.Command{
	Use:   "influx",
	Short: "Store the benchmarks in InfluxDB",
	Long:  `Driver to store the benchmarsk in a InfluDB database`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return influxdb.Execute(&cfg)
	},
}

func init() {
	driver.AddFlags(influxCmd, &cfg, "measurement")

	RootCmd.AddCommand(influxCmd)
}

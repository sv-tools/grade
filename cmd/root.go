package cmd

import (
	"errors"
	"os"
	"time"

	"github.com/sv-go-tools/grade/internal/parse"

	"github.com/sv-go-tools/grade/pkg/driver"

	"github.com/spf13/cobra"
)

var (
	cfg      driver.Config
	unixTime int64
	version  string = "v0.0.0"
)

// RootCmd is a root command
var RootCmd = &cobra.Command{
	Use:   "grade",
	Short: "Grade uploads Go benchmark data into a database.",
	Long: `Grade ingests Go benchmark data into a database so that you can track performance over time.
Just pipe the output of go test into grade.
Complete example is available at https://github.com/sv-go-tools/grade
Prints the data in JSON if no driver selected.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return driver.Execute(&cfg)
	},
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeCharDevice) == 0 {
			benchmarks, err := parse.MultipleBenchmarks(os.Stdin)
			if err != nil {
				return err
			}
			cfg.Benchmarks = benchmarks
		} else {
			return errors.New("please pipe the output of go test into grade")
		}
		cfg.Timestamp = time.Unix(unixTime, 0)
		return nil
	},
	Version: version,
}

func init() {
	RootCmd.PersistentFlags().StringVar(&cfg.GoVersion, "goversion", "", "Go version used to run benchmarks")
	RootCmd.PersistentFlags().Int64Var(&unixTime, "timestamp", 0, "Unix epoch timestamp (in seconds) to apply when storing all benchmark results")
	RootCmd.PersistentFlags().StringVar(&cfg.Revision, "revision", "", "Revision of the repository used to generate benchmark results")
	RootCmd.PersistentFlags().StringVar(&cfg.HardwareID, "hardwareid", "", "User-specified string to represent the hardware on which the benchmarks were run")
	RootCmd.PersistentFlags().StringVar(&cfg.Branch, "branch", "", "Branch of the repository used to generate benchmark results. The flag is optional and can be omitted")

	_ = RootCmd.MarkPersistentFlagRequired("goversion")
	_ = RootCmd.MarkPersistentFlagRequired("timestamp")
	_ = RootCmd.MarkPersistentFlagRequired("revision")
	_ = RootCmd.MarkPersistentFlagRequired("hardwareid")
}

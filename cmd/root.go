package cmd

import (
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/spf13/cobra"

	"github.com/sv-go-tools/grade/internal/parse"
	"github.com/sv-go-tools/grade/pkg/driver"
)

var (
	cfg     driver.Config
	rawTime string
	version string = "v0.0.0"
)

// RootCmd is a root command
var RootCmd = &cobra.Command{
	Use:   "grade",
	Short: "Grade uploads Go benchmark data into a database.",
	Long: `Grade ingests Go benchmark data into a database so that you can track performance over time.
Just pipe the output of go test into grade.
Complete example is available at https://github.com/sv-go-tools/grade`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeCharDevice) == 0 {
			benchmarks, err := parse.MultipleBenchmarks(os.Stdin)
			if err != nil {
				return err
			}
			cfg.Records = driver.Records(&cfg, benchmarks)
		} else {
			return errors.New("please pipe the output of go test into grade")
		}
		seconds, err := strconv.Atoi(rawTime)
		if err == nil {
			cfg.Timestamp = time.Unix(int64(seconds), 0)
		} else {
			parsedTime, err := time.Parse(time.RFC3339, rawTime)
			if err != nil {
				return err
			}
			cfg.Timestamp = parsedTime
		}
		return nil
	},
	Version: version,
}

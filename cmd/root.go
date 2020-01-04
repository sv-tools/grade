package cmd

import (
	"errors"
	"io"
	"os"
	"strconv"
	"time"

	"github.com/spf13/cobra"

	"github.com/sv-go-tools/grade/pkg/driver"
)

var (
	cfg       driver.Config
	rawTime   string
	rawTags   []string
	rawOutput string
)

// RootCmd is a root command
var RootCmd = &cobra.Command{
	Use:   "grade",
	Short: "Grade uploads Go benchmark data into a database.",
	Long: `Grade ingests Go benchmark data into a database so that you can track performance over time.
Just pipe the output of go test into grade.
Complete example is available at https://github.com/sv-go-tools/grade`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if rawTime != "" {
			seconds, err := strconv.Atoi(rawTime)
			if err == nil {
				tm := time.Unix(int64(seconds), 0)
				cfg.Timestamp = &tm
			} else {
				parsedTime, err := time.Parse(time.RFC3339, rawTime)
				if err != nil {
					return err
				}
				cfg.Timestamp = &parsedTime
			}
		}
		tags, err := driver.ParseTags(rawTags)
		if err != nil {
			return err
		}
		cfg.Tags = tags
		stat, _ := os.Stdin.Stat()
		var reader io.Reader
		if (stat.Mode() & os.ModeCharDevice) == 0 {
			reader = os.Stdin
		} else if len(args) > 0 {
			readers := make([]io.Reader, len(args))
			for i, name := range args {
				f, err := os.Open(name)
				if err != nil {
					return err
				}
				defer f.Close()
				readers[i] = f
			}
			reader = io.MultiReader(readers...)
		} else {
			return errors.New("please pipe the output of go test into grade")
		}
		cfg.Records = driver.ParseRecords(&cfg, reader)
		if rawOutput != "" {
			writer, err := os.Create(rawOutput)
			if err != nil {
				return err
			}
			cmd.SetOut(writer)
		}
		cfg.Output = cmd.OutOrStdout()
		return nil
	},
}

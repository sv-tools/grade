package cmd

import (
	"github.com/spf13/cobra"

	"github.com/sv-tools/grade/pkg/driver"
)

// AddCommonFlags adds the common flags for a command
func AddCommonFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringArrayVarP(&rawTags, "tag", "t", nil, "Custom tag in key=value format")
	cmd.PersistentFlags().StringVar(&rawTime, "timestamp", "", "Unix epoch timestamp (in seconds) or RFC3339 to apply when storing all benchmark results")
	cmd.PersistentFlags().StringVarP(&rawOutput, "output", "o", "", "A filename to write an output")
}

// AddDBFlags adds the DB related flags for a command
func AddDBFlags(cmd *cobra.Command, cfg *driver.Config, collectionName string) {
	cmd.PersistentFlags().StringVar(&cfg.ConnectionURL, "connection-url", "", "Connection URL of Database instance to store benchmark results (set to empty string to print to stdout)")
	cmd.PersistentFlags().StringVar(&cfg.Database, "database", "benchmarks", "Name of database to store benchmark results")
	cmd.PersistentFlags().StringVar(&cfg.Collection, collectionName, "go", "Name of "+collectionName+" to store benchmark results")
}

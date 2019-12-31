package cmd

import (
	"github.com/spf13/cobra"

	"github.com/sv-go-tools/grade/pkg/driver"
)

// AddCommonFlags adds the common flags for a command
func AddCommonFlags(cmd *cobra.Command, cfg *driver.Config) {
	cmd.PersistentFlags().StringVar(&cfg.GoVersion, "goversion", "", "Go version used to run benchmarks")
	cmd.PersistentFlags().StringVar(&rawTime, "timestamp", "", "Unix epoch timestamp (in seconds) or RFC3339 to apply when storing all benchmark results")
	cmd.PersistentFlags().StringVar(&cfg.Revision, "revision", "", "Revision of the repository used to generate benchmark results")
	cmd.PersistentFlags().StringVar(&cfg.HardwareID, "hardwareid", "", "User-specified string to represent the hardware on which the benchmarks were run")
	cmd.PersistentFlags().StringVar(&cfg.Branch, "branch", "", "Branch of the repository used to generate benchmark results. The flag is optional and can be omitted")

	_ = cmd.MarkPersistentFlagRequired("goversion")
	_ = cmd.MarkPersistentFlagRequired("timestamp")
	_ = cmd.MarkPersistentFlagRequired("revision")
	_ = cmd.MarkPersistentFlagRequired("hardwareid")
}

// AddDBFlags adds the DB related flags for a command
func AddDBFlags(cmd *cobra.Command, cfg *driver.Config, collectionName string) {
	cmd.PersistentFlags().BoolVarP(&cfg.Insecure, "insecure", "i", false, "Skip SSL verification if set")
	cmd.PersistentFlags().StringVar(&cfg.ConnectionURL, "connection-url", "", "Connection URL of Database instance to store benchmark results (set to empty string to print to stdout)")
	cmd.PersistentFlags().StringVar(&cfg.Database, "database", "benchmarks", "Name of database to store benchmark results")
	cmd.PersistentFlags().StringVar(&cfg.Collection, collectionName, "", "Name of "+collectionName+" to store benchmark results")

	_ = cmd.MarkPersistentFlagRequired("connection-url")
	_ = cmd.MarkPersistentFlagRequired(collectionName)
}

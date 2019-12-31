package driver

import (
	"github.com/spf13/cobra"
)

// AddFlags adds default flags for a command
func AddFlags(cmd *cobra.Command, cfg *Config, collectionName string) {
	cmd.PersistentFlags().BoolVarP(&cfg.Insecure, "insecure", "i", false, "Skip SSL verification if set")
	cmd.PersistentFlags().StringVar(&cfg.ConnectionURL, "connection-url", "", "Connection URL of Database instance to store benchmark results (set to empty string to print to stdout)")
	cmd.PersistentFlags().StringVar(&cfg.Database, "database", "benchmarks", "Name of database to store benchmark results")
	cmd.PersistentFlags().StringVar(&cfg.Collection, collectionName, "", "Name of "+collectionName+" to store benchmark results")

	_ = cmd.MarkPersistentFlagRequired("connection-url")
	_ = cmd.MarkPersistentFlagRequired(collectionName)
}

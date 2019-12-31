package driver

import (
	"time"

	"github.com/sv-go-tools/grade/internal/parse"
)

// Config represents the settings to process benchmarks.
type Config struct {
	// ConnectionURL is a connection url
	ConnectionURL string `json:"-"`

	// Insecure is a flag to skip SSL verification if set.
	Insecure bool `json:"-"`

	// Database is the name of the database into which to store the processed benchmark results.
	Database string `json:"-"`

	// Collection is the name of the collection into which to store the processed benchmark results.
	Collection string `json:"-"`

	// GoVersion is the tag value to use to indicate which version of Go was used for the benchmarks that have run.
	GoVersion string

	// Timestamp is the time to use when recording all of the benchmark results,
	// and is typically the timestamp of the commit used for the benchmark.
	Timestamp time.Time

	// Revision is the tag value to use to indicate which revision of the repository was used for the benchmarks that have run.
	// Feel free to use a SHA, tag name, or whatever will be useful to you when querying.
	Revision string

	// HardwareID is a user-specified string to represent the hardware on which the benchmarks have run.
	HardwareID string

	// Branch is the tag value to use to indicate which branch of the repository was used for the benchmarks that have run.
	// The tag is optional and can be omitted.
	Branch string

	Benchmarks map[string][]*parse.Benchmark
}

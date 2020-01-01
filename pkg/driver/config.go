package driver

import (
	"time"

	"github.com/sv-go-tools/grade/internal/parse"
)

// Config represents the settings to process benchmarks.
// ConnectionURL is a connection url.
// Insecure is a flag to skip SSL verification if set.
// Database is the name of the database into which to store the processed benchmark results.
// Collection is the name of the collection into which to store the processed benchmark results.
// GoVersion is the tag value to use to indicate which version of Go was used for the benchmarks that have run.
// Timestamp is the time to use when recording all of the benchmark results, and is typically the timestamp of the commit used for the benchmark.
// Revision is the tag value to use to indicate which revision of the repository was used for the benchmarks that have run. Feel free to use a SHA, tag name, or whatever will be useful to you when querying.
// HardwareID is a user-specified string to represent the hardware on which the benchmarks have run.
// Branch is the tag value to use to indicate which branch of the repository was used for the benchmarks that have run. The tag is optional and can be omitted.
type Config struct {
	ConnectionURL string                        `json:"-" bson:"-"`
	Insecure      bool                          `json:"-" bson:"-"`
	Database      string                        `json:"-" bson:"-"`
	Collection    string                        `json:"-" bson:"-"`
	GoVersion     string                        `json:"goVersion" bson:"goVersion"`
	Timestamp     time.Time                     `json:"timestamp" bson:"timestamp"`
	Revision      string                        `json:"revision" bson:"revision"`
	HardwareID    string                        `json:"hardwareID" bson:"hardwareID"`
	Branch        string                        `json:"branch,omitempty" bson:"branch,omitempty"`
	Benchmarks    map[string][]*parse.Benchmark `json:"benchmarks" bson:"benchmarks"`
}

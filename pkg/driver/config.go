package driver

import (
	"io"
	"time"
)

// Config represents the settings to process benchmarks.
// ConnectionURL is a connection url.
// Insecure is a flag to skip SSL verification if set.
// Database is the name of the database into which to store the processed benchmark results.
// Collection is the name of the collection into which to store the processed benchmark results.
// GoVersion is the tag value to use to indicate which version of Go was used for the benchmarks that have run.
// Timestamp is the time to use when recording all the benchmark results, and is typically the timestamp of the commit used for the benchmark.
// Revision is the tag value to use to indicate which revision of the repository was used for the benchmarks that have run. Feel free to use an SHA, tag name, or whatever will be useful to you when querying.
// HardwareID is a user-specified string to represent the hardware on which the benchmarks have run.
// Branch is the tag value to use to indicate which branch of the repository was used for the benchmarks that have run. The tag is optional and can be omitted.
// Records is a list of parsed benchmarks
type Config struct {
	ConnectionURL string
	Insecure      bool
	Database      string
	Collection    string
	JSONIndent    int
	Timestamp     *time.Time
	Tags          Tags
	Records       Records
	Output        io.Writer
}

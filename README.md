# grade
[![Build Status](https://github.com/sv-go-tools/grade/workflows/Go/badge.svg)](https://github.com/sv-go-tools/grade/actions?query=branch%3Amaster+event%3Apush)
[![License](https://img.shields.io/github/license/sv-go-tools/grade.svg)](/LICENSE)
[![Release](https://img.shields.io/github/release/sv-go-tools/grade.svg)](https://github.com/sv-go-tools/grade.svg/releases/latest)


`grade` transforms Go benchmark data into various format so that you can uploads the data to a database and track the performance over time.

This is a fork of the [influxdata/grade](https://github.com/influxdata/grade). The main reason for the forking the project is the support of various drivers, such as `Json` and `MongoDB`.


## Installation

To download and install the `grade` executable into your `$GOPATH/bin`:

```sh
go get github.com/sv-go-tools/grade/cmd/grade
```

## Usage

Although you can pipe the output of `go test` directly into `grade`,
for now we recommend placing the output of `go test` in a file first so that if something goes wrong,
you don't have to wait again to run all the benchmarks.

For example, to run all the benchmarks in your current Go project:

```sh
go test -run=XXX -bench=. -benchmem ./... > bench.log
```

Then, assuming you are in the directory of your Go project and
git has checked out the same commit corresponding with the tests that have run,
this is the bare set of options to load the benchmark results into InfluxDB via `grade`:

### Json

```sh
grade json \
  --hardwareid="my dev machine" \
  --goversion="$(go version | cut -d' ' -f3-)" \
  --revision="$(git log -1 --format=%H)" \
  --timestamp="$(git log -1 --format=%ct)" \
  --branch="$(git rev-parse --abbrev-ref HEAD)" \
  < bench.log
```

Notes on this style of invocation:

* The hardware ID is a string that you specify to identify the hardware on which the benchmarks were run.
* The Go version subcommand will produce a string like `go1.6.2 darwin/amd64`, but you can use any string you'd like.
* The revision subcommand is the full SHA of the commit, but feel free to use a git tag name or any other string.
* The timestamp is a Unix epoch timestamp in seconds.
The above subcommand produces the Unix timestamp for the committer of the most recent commit.
This assumes that the commits whose benchmarks are being run, all are ascending in time;
git does not enforce that commits' timestamps are ascending, so if this assumption is broken,
your data may look strange when you visualize it.
* The branch subcommand is the name of the current branch. The `-branch` flag is optional.

### InfluxDB

```sh
grade influx \
  --connection-url="http://localhost:8086" \
  --database="grade_bencmarks" \
  --measurement="go" \
  --hardwareid="my dev machine" \
  --goversion="$(go version | cut -d' ' -f3-)" \
  --revision="$(git log -1 --format=%H)" \
  --timestamp="$(git log -1 --format=%ct)" \
  --branch="$(git rev-parse --abbrev-ref HEAD)" \
  < bench.log
```

The data from Go benchmarks tends to be very time-sparse (up to perhaps dozens of commits per day),
so we recommend creating your database with an infinite retention and a large shard duration.
Issue this command to your InfluxDB instance:

```sql
CREATE DATABASE benchmarks WITH DURATION INF SHARD DURATION 90d
```

The common flags are same with Json flags, but the `influx` has some additional:

* `--connection-url` is a connection url.
Basic auth credentials can be embedded in the URL if needed.
HTTPS is supported; supply `--insecure` if you need to skip SSL verification.
If you set it to an empty string `""`, `grade` will print line protocol to stdout.
* `--database` is a name of a Database, defaults to `benchmarks`.
* `--measurement` is not provided but defaults to `go`.

For each benchmark result from a run of `go test -bench`:

* Tags:
	* `goversion` is the same string as passed in to the `--goversion` flag.
	* `hwid` is the same string as passed in to the `--hardwareid` flag.
	* `name` is the name of the benchmark function, stripped of the `Benchmark` prefix.
	* `pkg` is the name of Go package containing the benchmark, e.g. `github.com/sv-go-tools/grade`.
	* `procs` is the number of CPUs used to run the benchmark. This is a tag because you are more likely to group by `procs` rather than chart them over time.
	* `branch` is the same string as passed in to the `--branch` flag.
	Since the `--branch` flag is optional and can be omited, the tag will be present only if the flag is set.
* Fields:
	* `alloced_bytes_per_op` is the allocated bytes per iteration of the benchmark.
	* `allocs_per_op` is how many allocations occurred per iteration of the benchmark.
	* `mb_per_s` is how many megabytes processed per second when running the benchmark.
	* `n` is the number of iterations in the benchmark.
	* `ns_per_op` is the number of wall nanoseconds taken per iteration of the benchmark.
	* `revision` is the git revision specified in the `--revision` flag.
	This was chosen to be a field so that the information is quickly available but not at the cost of a growing series cardinality per benchmark run.

## Sample

For a benchmark like this:

```
goos: darwin
goarch: amd64
pkg: github.com/sv-go-tools/grade
BenchmarkFib                     4185866               300 ns/op               0 B/op          0 allocs/op
BenchmarkFib-2                   4070166               287 ns/op               0 B/op          0 allocs/op
BenchmarkFibParallel             4145983               291 ns/op               0 B/op          0 allocs/op
BenchmarkFibParallel-2           8255342               147 ns/op               0 B/op          0 allocs/op
PASS
ok      github.com/sv-go-tools/grade    6.202s
```

Which is passed to `grade` like this:

### Json

```
grade json \
  --hardwareid="my dev machine" \
  --goversion="$(go version | cut -d' ' -f3-)" \
  --revision="$(git log -1 --format=%H)" \
  --timestamp="$(git log -1 --format=%ct)" \
  --branch="$(git rev-parse --abbrev-ref HEAD)" \
  < bench.log
```

You will see output like:
```json
{
  "GoVersion": "go1.13.5 darwin/amd64",
  "Timestamp": "2019-12-30T23:44:39-06:00",
  "Revision": "bf145c5423671fce4d252ce164cdc1e65f66de15",
  "HardwareID": "my dev machine",
  "Branch": "master",
  "Benchmarks": {
    "github.com/sv-go-tools/grade": [
      {
        "Name": "Fib",
        "Procs": 1,
        "N": 4185866,
        "NsPerOp": 300,
        "AllocedBytesPerOp": 0,
        "AllocsPerOp": 0,
        "MBPerS": 0,
        "Measured": 13
      },
      {
        "Name": "Fib",
        "Procs": 2,
        "N": 4070166,
        "NsPerOp": 287,
        "AllocedBytesPerOp": 0,
        "AllocsPerOp": 0,
        "MBPerS": 0,
        "Measured": 13
      },
      {
        "Name": "FibParallel",
        "Procs": 1,
        "N": 4145983,
        "NsPerOp": 291,
        "AllocedBytesPerOp": 0,
        "AllocsPerOp": 0,
        "MBPerS": 0,
        "Measured": 13
      },
      {
        "Name": "FibParallel",
        "Procs": 2,
        "N": 8255342,
        "NsPerOp": 147,
        "AllocedBytesPerOp": 0,
        "AllocsPerOp": 0,
        "MBPerS": 0,
        "Measured": 13
      }
    ]
  }
}
```

### InfluxDB

```sh
grade influx \
  --connection-url="" \
  --database="grade_bencmarks" \
  --measurement="go" \
  --hardwareid="my dev machine" \
  --goversion="$(go version | cut -d' ' -f3-)" \
  --revision="$(git log -1 --format=%H)" \
  --timestamp="$(git log -1 --format=%ct)" \
  --branch="$(git rev-parse --abbrev-ref HEAD)" \
  < bench.log
```

the `--connection-url` is empty to print the data

You will see output like:
```
go,branch=master,goversion=go1.13.5\ darwin/amd64,hwid=my\ dev\ machine,name=Fib,pkg=github.com/sv-go-tools/grade,procs=1 alloced_bytes_per_op=0i,allocs_per_op=0i,n=4185866i,ns_per_op=300,revision="bf145c5423671fce4d252ce164cdc1e65f66de15" 1577771079000000000
go,branch=master,goversion=go1.13.5\ darwin/amd64,hwid=my\ dev\ machine,name=Fib,pkg=github.com/sv-go-tools/grade,procs=2 alloced_bytes_per_op=0i,allocs_per_op=0i,n=4070166i,ns_per_op=287,revision="bf145c5423671fce4d252ce164cdc1e65f66de15" 1577771079000000000
go,branch=master,goversion=go1.13.5\ darwin/amd64,hwid=my\ dev\ machine,name=FibParallel,pkg=github.com/sv-go-tools/grade,procs=1 alloced_bytes_per_op=0i,allocs_per_op=0i,n=4145983i,ns_per_op=291,revision="bf145c5423671fce4d252ce164cdc1e65f66de15" 1577771079000000000
go,branch=master,goversion=go1.13.5\ darwin/amd64,hwid=my\ dev\ machine,name=FibParallel,pkg=github.com/sv-go-tools/grade,procs=2 alloced_bytes_per_op=0i,allocs_per_op=0i,n=8255342i,ns_per_op=147,revision="bf145c5423671fce4d252ce164cdc1e65f66de15" 1577771079000000000
```

# grade

[![Test](https://img.shields.io/github/workflow/status/sv-tools/grade/Benchmarks)](https://github.com/sv-tools/grade/actions?query=workflow%3A%22Benchmarks%22)
[![Version](https://img.shields.io/github/release/sv-tools/grade.svg)](https://github.com/sv-tools/grade/releases/latest)

`grade` transforms Go benchmark data into various format so that you can uploads the data to a database and track the performance over time.

This is a fork of the [influxdata/grade](https://github.com/influxdata/grade). The main reason for the forking the project is the support of various drivers, such as `Json` and `MongoDB`.


## Installation

### Using `brew`
```sh
brew install sv-tools/apps/grade
```

### Docker
```sh
docker pull docker.pkg.github.com/sv-tools/grade/grade:latest
```

### Go way
To install latest master:
```sh
go get github.com/sv-tools/grade
```

### Binary mode

Download a build for your OS from the [latest release](https://github.com/sv-tools/grade/releases/latest).

The checksums is signed by gpg key [7D76B375F08A7D93584B36D766538F03CDA385C7](https://keys.openpgp.org/search?q=sv.go.tools@gmail.com)

## Usage

Although you can pipe the output of `go test` directly into `grade` or pass the list of files as agruments,
for now we recommend placing the output of `go test` in a file first so that if something goes wrong,
you don't have to wait again to run all the benchmarks.

For example, to run all the benchmarks in your current Go project:

```sh
go test -run=XXX -bench= ... -benchmem ./... > bench.log
```

Then, assuming you are in the directory of your Go project and
git has checked out the same commit corresponding with the tests that have run,
this is the bare set of options to load the benchmark results via `grade`:

### Json

```sh
grade json \
  --tag "hardwareID=my dev machine" \
  --tag "goVersion=$(go version | cut -d' ' -f3-)" \
  --tag "revision=$(git log -1 --format=%H)" \
  --tag "branch=$(git rev-parse --abbrev-ref HEAD)" \
  --timestamp="$(git log -1 --format=%ct)" \
  bench.log
```

Notes on this style of invocation:

* The hardware ID is a string that you specify to identify the hardware on which the benchmarks were run.
* The Go version subcommand will produce a string like `go1.6.2 darwin/amd64`, but you can use any string you'd like.
* The revision subcommand is the full SHA of the commit, but feel free to use a git tag name or any other string.
* The branch subcommand is the name of the current branch. The `-branch` flag is optional.
The above arguments are the custom tags, you can define as many as you want tags with any names and values.

* The timestamp is a Unix epoch timestamp in seconds.
The above subcommand produces the Unix timestamp for the committer of the most recent commit.
This assumes that the commits whose benchmarks are being run, all are ascending in time;
git does not enforce that commits' timestamps are ascending, so if this assumption is broken,
your data may look strange when you visualize it.
The timestamp is an optional argument, not a tag. You can pass a time in Unix or RFC3339 format.
If the tag is not set, then will be used a current timestamp. In case of Json the local time is used,
for InfluxDB and MongoGo will be used a server's time.

#### Output:

```json
[
  {
    "allocedBytesPerOp": 0,
    "allocsPerOp": 0,
    "branch": "mongodb",
    "coverage": 0,
    "goArch": "amd64",
    "goOS": "darwin",
    "goVersion": "go1.13.5 darwin/amd64",
    "hardwareID": "my dev machine",
    "mbPerS": 0,
    "n": 578371,
    "name": "Fib",
    "nsPerOp": 2155,
    "package": "github.com/sv-tools/grade",
    "procs": 16,
    "revision": "e9fb4afc7bf1afcf6b98af1ccc3077f5c50ab8ec",
    "timestamp": "2020-01-02T14:24:18-06:00"
  },
  {
    "allocedBytesPerOp": 0,
    "allocsPerOp": 0,
    "branch": "mongodb",
    "coverage": 0,
    "goArch": "amd64",
    "goOS": "darwin",
    "goVersion": "go1.13.5 darwin/amd64",
    "hardwareID": "my dev machine",
    "mbPerS": 0,
    "n": 4771570,
    "name": "FibParallel",
    "nsPerOp": 261,
    "package": "github.com/sv-tools/grade",
    "procs": 16,
    "revision": "e9fb4afc7bf1afcf6b98af1ccc3077f5c50ab8ec",
    "timestamp": "2020-01-02T14:24:18-06:00"
  },
  {
    "allocedBytesPerOp": 0,
    "allocsPerOp": 0,
    "branch": "mongodb",
    "coverage": 0,
    "goArch": "amd64",
    "goOS": "darwin",
    "goVersion": "go1.13.5 darwin/amd64",
    "hardwareID": "my dev machine",
    "mbPerS": 0,
    "n": 587794,
    "name": "FibDriver",
    "nsPerOp": 2096,
    "package": "github.com/sv-tools/grade/pkg/driver",
    "procs": 16,
    "revision": "e9fb4afc7bf1afcf6b98af1ccc3077f5c50ab8ec",
    "timestamp": "2020-01-02T14:24:18-06:00"
  },
  {
    "allocedBytesPerOp": 0,
    "allocsPerOp": 0,
    "branch": "mongodb",
    "coverage": 0,
    "goArch": "amd64",
    "goOS": "darwin",
    "goVersion": "go1.13.5 darwin/amd64",
    "hardwareID": "my dev machine",
    "mbPerS": 0,
    "n": 4550692,
    "name": "FibParallelDriver",
    "nsPerOp": 271,
    "package": "github.com/sv-tools/grade/pkg/driver",
    "procs": 16,
    "revision": "e9fb4afc7bf1afcf6b98af1ccc3077f5c50ab8ec",
    "timestamp": "2020-01-02T14:24:18-06:00"
  }
]
```

### InfluxDB

```sh
grade influx \
  --connection-url="" \
  --database="grade_benchmarks" \
  --measurement="go" \
  --tag "hardwareID=my dev machine" \
  --tag "goVersion=$(go version | cut -d' ' -f3-)" \
  --tag "revision=$(git log -1 --format=%H)" \
  --tag "branch=$(git rev-parse --abbrev-ref HEAD)" \
  --timestamp="$(git log -1 --format=%ct)" \
  bench.log
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

#### Output

```ini
go,branch=mongodb,goArch=amd64,goOS=darwin,goVersion=go1.13.5\ darwin/amd64,hardwareID=my\ dev\ machine,name=Fib,package=github.com/sv-tools/grade,procs=16,revision=e9fb4afc7bf1afcf6b98af1ccc3077f5c50ab8ec allocedBytesPerOp=0i,allocsPerOp=0i,coverage=0,n=578371i,nsPerOp=2155 1577996658000000000
go,branch=mongodb,goArch=amd64,goOS=darwin,goVersion=go1.13.5\ darwin/amd64,hardwareID=my\ dev\ machine,name=FibParallel,package=github.com/sv-tools/grade,procs=16,revision=e9fb4afc7bf1afcf6b98af1ccc3077f5c50ab8ec allocedBytesPerOp=0i,allocsPerOp=0i,coverage=0,n=4771570i,nsPerOp=261 1577996658000000000
go,branch=mongodb,goArch=amd64,goOS=darwin,goVersion=go1.13.5\ darwin/amd64,hardwareID=my\ dev\ machine,name=FibDriver,package=github.com/sv-tools/grade/pkg/driver,procs=16,revision=e9fb4afc7bf1afcf6b98af1ccc3077f5c50ab8ec allocedBytesPerOp=0i,allocsPerOp=0i,coverage=0,n=587794i,nsPerOp=2096 1577996658000000000
go,branch=mongodb,goArch=amd64,goOS=darwin,goVersion=go1.13.5\ darwin/amd64,hardwareID=my\ dev\ machine,name=FibParallelDriver,package=github.com/sv-tools/grade/pkg/driver,procs=16,revision=e9fb4afc7bf1afcf6b98af1ccc3077f5c50ab8ec allocedBytesPerOp=0i,allocsPerOp=0i,coverage=0,n=4550692i,nsPerOp=271 1577996658000000000
```

For each benchmark result from a run of `go test -bench`:

* Tags:
  * All the tags passed to `grade` as `--tag`.
  * `package` is the name of Go package containing the benchmark, e.g. `github.com/sv-tools/grade`.
  * `name` is the name of the benchmark function, stripped of the `Benchmark` prefix, e.g. `Fib`.
  * `goArch` is the architecture of a your machine, e.g. `amd64`.
  * `goOS` is your operating system, e.g. `darwin`.
  * `procs` is the number of CPUs used to run the benchmark. This is a tag because you are more likely to group by `procs` rather than chart them over time.

* Fields:
  * `n` is the number of iterations in the benchmark.
  * `allocedBytesPerOp` is the allocated bytes per iteration of the benchmark.
  * `allocsPerOp` is how many allocations occurred per iteration of the benchmark.
  * `mbPerS` is how many megabytes processed per second when running the benchmark.
  * `nsPerOp` is the number of wall nanoseconds taken per iteration of the benchmark.
  * `coverage` is the percentage of code coverage for unit testing.

### MongoDB

```sh
grade mongo \
  --connection-url="" \
  --database="grade_benchmarks" \
  --collection="go" \
  --tag "hardwareID=my dev machine" \
  --tag "goVersion=$(go version | cut -d' ' -f3-)" \
  --tag "revision=$(git log -1 --format=%H)" \
  --tag "branch=$(git rev-parse --abbrev-ref HEAD)" \
  --timestamp="$(git log -1 --format=%ct)" \
  bench.log
```

The differnce is that this requires `--collection` tag instead of `--measurement`

#### Output

```json
[{"allocedBytesPerOp":0,"allocsPerOp":0,"branch":"mongodb","coverage":0,"goArch":"amd64","goOS":"darwin","goVersion":"go1.13.5 darwin/amd64","hardwareID":"my dev machine","mbPerS":0,"n":578371,"name":"Fib","nsPerOp":2155,"package":"github.com/sv-tools/grade","procs":"16","revision":"e9fb4afc7bf1afcf6b98af1ccc3077f5c50ab8ec","timestamp":"2020-01-02T14:24:18-06:00"},{"allocedBytesPerOp":0,"allocsPerOp":0,"branch":"mongodb","coverage":0,"goArch":"amd64","goOS":"darwin","goVersion":"go1.13.5 darwin/amd64","hardwareID":"my dev machine","mbPerS":0,"n":4771570,"name":"FibParallel","nsPerOp":261,"package":"github.com/sv-tools/grade","procs":"16","revision":"e9fb4afc7bf1afcf6b98af1ccc3077f5c50ab8ec","timestamp":"2020-01-02T14:24:18-06:00"},{"allocedBytesPerOp":0,"allocsPerOp":0,"branch":"mongodb","coverage":0,"goArch":"amd64","goOS":"darwin","goVersion":"go1.13.5 darwin/amd64","hardwareID":"my dev machine","mbPerS":0,"n":587794,"name":"FibDriver","nsPerOp":2096,"package":"github.com/sv-tools/grade/pkg/driver","procs":"16","revision":"e9fb4afc7bf1afcf6b98af1ccc3077f5c50ab8ec","timestamp":"2020-01-02T14:24:18-06:00"},{"allocedBytesPerOp":0,"allocsPerOp":0,"branch":"mongodb","coverage":0,"goArch":"amd64","goOS":"darwin","goVersion":"go1.13.5 darwin/amd64","hardwareID":"my dev machine","mbPerS":0,"n":4550692,"name":"FibParallelDriver","nsPerOp":271,"package":"github.com/sv-tools/grade/pkg/driver","procs":"16","revision":"e9fb4afc7bf1afcf6b98af1ccc3077f5c50ab8ec","timestamp":"2020-01-02T14:24:18-06:00"}]
```

Or with a non empty `--connection-url`:

```sh
./grade mongo --connection-url="mongodb://admin:secret@localhost:27017" bench.log
```

the output contains Object IDs:

```js
ObjectID("5e0ffd8a9033c7ff19f23c32")
ObjectID("5e0ffd8a9033c7ff19f23c33")
ObjectID("5e0ffd8a9033c7ff19f23c34")
ObjectID("5e0ffd8a9033c7ff19f23c35")
```

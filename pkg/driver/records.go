package driver

import (
	"bufio"
	"io"
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/tools/benchmark/parse"
)

type Record map[string]interface{}
type Records []Record

var reProcsSuffix = regexp.MustCompile(`(?:-(\d+))$`)

func makeRecord(cfg *Config, goOS, goArch, packageName string, benchmark *parse.Benchmark) Record {
	rec := make(Record)
	for k, v := range cfg.Tags {
		rec[k] = v
	}
	rec["goOS"] = goOS
	rec["goArch"] = goArch
	rec["package"] = packageName
	rec["n"] = benchmark.N
	rec["nsPerOp"] = benchmark.NsPerOp
	rec["allocedBytesPerOp"] = benchmark.AllocedBytesPerOp
	rec["allocsPerOp"] = benchmark.AllocsPerOp
	rec["mbPerS"] = benchmark.MBPerS
	rec["measured"] = benchmark.Measured

	name := strings.TrimPrefix(benchmark.Name, "Benchmark")
	name = reProcsSuffix.ReplaceAllLiteralString(name, "")
	rec["name"] = name
	if match := reProcsSuffix.FindStringSubmatch(benchmark.Name); match == nil {
		rec["procs"] = 1
	} else {
		n, err := strconv.Atoi(match[1])
		if err != nil {
			rec["procs"] = 1
		} else {
			rec["procs"] = n
		}
	}

	if cfg.Timestamp != nil {
		rec["timestamp"] = *cfg.Timestamp
	}
	return rec
}

func ParseRecords(cfg *Config, r io.Reader) Records {
	var (
		goOS, goArch, packageName string
		res                       Records
	)

	scan := bufio.NewScanner(r)
	for scan.Scan() {
		line := scan.Text()
		fields := strings.Fields(line)
		if len(fields) > 1 {
			switch fields[0] {
			case "pkg:":
				packageName = fields[1]
				continue
			case "goos:":
				goOS = fields[1]
				continue
			case "goarch:":
				goArch = fields[1]
				continue
			case "coverage:":
				res = addCoverage(res, packageName, fields[1])
				continue
			}
		}

		if b, err := parse.ParseLine(line); err == nil {
			res = append(res, makeRecord(cfg, goOS, goArch, packageName, b))
		}
	}
	return res
}

func addCoverage(records Records, packageName, coverage string) Records {
	coverage = strings.TrimSuffix(coverage, "%")
	c, err := strconv.ParseFloat(coverage, 64)
	if err != nil {
		c = 0
	}
	for i, rec := range records {
		if rec["package"] == packageName {
			rec["coverage"] = c
			records[i] = rec
		}
	}
	return records
}

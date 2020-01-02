package driver

import (
	"io"
	"time"

	"golang.org/x/tools/benchmark/parse"
)

type Record struct {
	GoVersion         string    `json:"goVersion" bson:"goVersion"`
	Revision          string    `json:"revision" bson:"revision"`
	HardwareID        string    `json:"hardwareID" bson:"hardwareID"`
	Branch            string    `json:"branch,omitempty" bson:"branch,omitempty"`
	Timestamp         time.Time `json:"timestamp" bson:"timestamp"`
	Name              string    `json:"name" bson:"name"`
	Procs             int       `json:"procs" bson:"procs"`
	N                 int       `json:"n" bson:"n"`
	NsPerOp           float64   `json:"nsPerOp" bson:"nsPerOp"`
	AllocedBytesPerOp uint64    `json:"allocedBytesPerOp" bson:"allocedBytesPerOp"`
	AllocsPerOp       uint64    `json:"allocsPerOp" bson:"allocsPerOp"`
	MBPerS            float64   `json:"mbPerS" bson:"mbPerS"`
	Measured          int       `json:"-" bson:"-"`
}

func makeRecords(cfg *Config, data map[string][]*parse.Benchmark) []*Record {
	var records []*Record
	for _, benchmarks := range data {
		for _, benchmark := range benchmarks {
			records = append(records, &Record{
				GoVersion:         cfg.GoVersion,
				Timestamp:         cfg.Timestamp,
				Revision:          cfg.Revision,
				HardwareID:        cfg.HardwareID,
				Branch:            cfg.Branch,
				Name:              benchmark.Name,
				N:                 benchmark.N,
				NsPerOp:           benchmark.NsPerOp,
				AllocedBytesPerOp: benchmark.AllocedBytesPerOp,
				AllocsPerOp:       benchmark.AllocsPerOp,
				MBPerS:            benchmark.MBPerS,
				Measured:          benchmark.Measured,
			})
		}
	}
	return records
}

func Parse(cfg *Config, r io.Reader) ([]*Record, error) {
	benchmarks, err := parse.ParseSet(r)
	if err != nil {
		return nil, err
	}
	return makeRecords(cfg, benchmarks), nil
}

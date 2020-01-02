package driver

import (
	"time"

	"github.com/sv-go-tools/grade/internal/parse"
)

type Record struct {
	GoVersion         string    `json:"goVersion" bson:"goVersion"`
	Timestamp         time.Time `json:"timestamp" bson:"timestamp"`
	Revision          string    `json:"revision" bson:"revision"`
	HardwareID        string    `json:"hardwareID" bson:"hardwareID"`
	Branch            string    `json:"branch,omitempty" bson:"branch,omitempty"`
	Package           string    `json:"package" bson:"package"`
	Name              string    `json:"name" bson:"name"`
	Procs             int       `json:"procs" bson:"procs"`
	N                 int       `json:"n" bson:"n"`
	NsPerOp           float64   `json:"nsPerOp" bson:"nsPerOp"`
	AllocedBytesPerOp uint64    `json:"allocedBytesPerOp" bson:"allocedBytesPerOp"`
	AllocsPerOp       uint64    `json:"allocsPerOp" bson:"allocsPerOp"`
	MBPerS            float64   `json:"mbPerS" bson:"mbPerS"`
	Measured          int       `json:"-" bson:"-"`
}

func Records(cfg *Config, data map[string][]*parse.Benchmark) []*Record {
	var records []*Record
	for packageName, benchmarks := range data {
		for _, benchmark := range benchmarks {
			records = append(records, &Record{
				GoVersion:         cfg.GoVersion,
				Timestamp:         cfg.Timestamp,
				Revision:          cfg.Revision,
				HardwareID:        cfg.HardwareID,
				Branch:            cfg.Branch,
				Package:           packageName,
				Name:              benchmark.Name,
				Procs:             benchmark.Procs,
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

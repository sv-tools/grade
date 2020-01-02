package influxdb

import (
	"fmt"
	"net/url"

	client "github.com/influxdata/influxdb1-client/v2"
	"golang.org/x/tools/benchmark/parse"

	"github.com/sv-go-tools/grade/pkg/driver"
)

func Execute(cfg *driver.Config) error {
	points, err := Points(cfg)
	if err != nil {
		return err
	}

	if cfg.ConnectionURL == "" {
		// Dry run requested.
		for _, p := range points.Points() {
			fmt.Println(p.String())
		}
		return nil
	}
	cl, err := buildClient(cfg)
	if err != nil {
		return err
	}
	defer cl.Close()

	if err := cl.Write(points); err != nil {
		return err
	}
	return nil
}

func buildClient(cfg *driver.Config) (client.Client, error) {
	u, err := url.Parse(cfg.ConnectionURL)
	if err != nil {
		return nil, err
	}

	c := client.HTTPConfig{
		Addr:               cfg.ConnectionURL,
		UserAgent:          "influxdata.Grade",
		InsecureSkipVerify: cfg.Insecure,
	}

	if u.User != nil {
		c.Username = u.User.Username()
		c.Password, _ = u.User.Password()
	}

	return client.NewHTTPClient(c)
}

// Points parses the benchmark output from r and creates a batch of points using cfg.
func Points(cfg *driver.Config) (client.BatchPoints, error) {
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Precision: "s",
		Database:  cfg.Database,
	})
	if err != nil {
		return nil, err
	}

	for _, rec := range cfg.Records {
		tags := map[string]string{
			"goversion": rec.GoVersion,
			"hwid":      rec.HardwareID,
			"name":      rec.Name,
		}
		if cfg.Branch != "" {
			tags["branch"] = rec.Branch
		}
		p, err := client.NewPoint(
			cfg.Collection,
			tags,
			makeFields(rec),
			cfg.Timestamp,
		)
		if err != nil {
			return nil, err
		}

		bp.AddPoint(p)
	}

	return bp, nil
}

func makeFields(rec *driver.Record) map[string]interface{} {
	f := make(map[string]interface{}, 6)

	f["revision"] = rec.Revision
	f["n"] = rec.N

	if (rec.Measured & parse.NsPerOp) != 0 {
		f["ns_per_op"] = rec.NsPerOp
	}
	if (rec.Measured & parse.MBPerS) != 0 {
		f["mb_per_s"] = rec.MBPerS
	}
	if (rec.Measured & parse.AllocedBytesPerOp) != 0 {
		f["alloced_bytes_per_op"] = int64(rec.AllocedBytesPerOp)
	}
	if (rec.Measured & parse.AllocsPerOp) != 0 {
		f["allocs_per_op"] = int64(rec.AllocsPerOp)
	}

	return f
}

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
		fields := makeFields(rec)
		tags := makeTags(rec, fields)
		var (
			err error
			p   *client.Point
		)
		if cfg.Timestamp == nil {
			p, err = client.NewPoint(
				cfg.Collection,
				tags,
				fields,
			)
		} else {
			p, err = client.NewPoint(
				cfg.Collection,
				tags,
				fields,
				*cfg.Timestamp,
			)
		}
		if err != nil {
			return nil, err
		}

		bp.AddPoint(p)
	}

	return bp, nil
}

func makeFields(rec driver.Record) map[string]interface{} {
	fields := make(map[string]interface{})

	if c, ok := rec["coverage"]; ok {
		fields["coverage"] = c
	}
	fields["n"] = rec["n"]

	if (rec["measured"].(int) & parse.NsPerOp) != 0 {
		fields["nsPerOp"] = rec["nsPerOp"]
	}
	if (rec["measured"].(int) & parse.MBPerS) != 0 {
		fields["mbPerS"] = rec["mbPerS"]
	}
	if (rec["measured"].(int) & parse.AllocedBytesPerOp) != 0 {
		fields["allocedBytesPerOp"] = int64(rec["allocedBytesPerOp"].(uint64))
	}
	if (rec["measured"].(int) & parse.AllocsPerOp) != 0 {
		fields["allocsPerOp"] = int64(rec["allocsPerOp"].(uint64))
	}
	return fields
}

func makeTags(rec driver.Record, fields map[string]interface{}) map[string]string {
	tags := make(map[string]string)
	for k, v := range rec {
		if _, ok := fields[k]; !ok {
			if v, ok := v.(string); ok {
				tags[k] = v
			}
		}
	}
	return tags
}

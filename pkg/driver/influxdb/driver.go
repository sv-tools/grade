package influxdb

import (
	"fmt"
	"net/url"
	"strconv"

	client "github.com/influxdata/influxdb1-client/v2"

	"github.com/sv-go-tools/grade/internal/parse"
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

	for pkg, bs := range cfg.Benchmarks {
		for _, b := range bs {
			tags := map[string]string{
				"goversion": cfg.GoVersion,
				"hwid":      cfg.HardwareID,
				"pkg":       pkg,
				"procs":     strconv.Itoa(b.Procs),
				"name":      b.Name,
			}
			if cfg.Branch != "" {
				tags["branch"] = cfg.Branch
			}
			p, err := client.NewPoint(
				cfg.Collection,
				tags,
				makeFields(b, cfg),
				cfg.Timestamp,
			)
			if err != nil {
				return nil, err
			}

			bp.AddPoint(p)
		}
	}

	return bp, nil
}

func makeFields(b *parse.Benchmark, cfg *driver.Config) map[string]interface{} {
	f := make(map[string]interface{}, 6)

	f["revision"] = cfg.Revision
	f["n"] = b.N

	if (b.Measured & parse.NsPerOp) != 0 {
		f["ns_per_op"] = b.NsPerOp
	}
	if (b.Measured & parse.MBPerS) != 0 {
		f["mb_per_s"] = b.MBPerS
	}
	if (b.Measured & parse.AllocedBytesPerOp) != 0 {
		f["alloced_bytes_per_op"] = int64(b.AllocedBytesPerOp)
	}
	if (b.Measured & parse.AllocsPerOp) != 0 {
		f["allocs_per_op"] = int64(b.AllocsPerOp)
	}

	return f
}

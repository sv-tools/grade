package json

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/sv-go-tools/grade/pkg/driver"
)

func Execute(cfg *driver.Config) error {
	var (
		data []byte
		err  error
	)
	if cfg.JSONIndent == 0 {
		data, err = json.Marshal(makeRecords(cfg.Records))
	} else {
		data, err = json.MarshalIndent(makeRecords(cfg.Records), "", strings.Repeat(" ", cfg.JSONIndent))
	}
	if err != nil {
		return err
	}
	fmt.Fprintln(cfg.Output, string(data))
	return nil
}

func makeRecords(records driver.Records) driver.Records {
	for _, rec := range records {
		delete(rec, "measured")
		if _, ok := rec["timestamp"]; !ok {
			rec["timestamp"] = time.Now()
		}
	}
	return records
}

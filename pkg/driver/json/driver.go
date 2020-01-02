package json

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/sv-go-tools/grade/pkg/driver"
)

func Execute(cfg *driver.Config) error {
	indent := strings.Repeat(" ", cfg.JSONIndent)
	data, err := json.MarshalIndent(makeRecords(cfg.Records), "", indent)
	if err != nil {
		return err
	}
	fmt.Println(string(data))
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

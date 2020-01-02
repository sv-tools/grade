package json

import (
	"encoding/json"
	"fmt"

	"github.com/sv-go-tools/grade/pkg/driver"
)

func Execute(cfg *driver.Config) error {
	data, err := json.MarshalIndent(cfg.Records, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(data))
	return nil
}

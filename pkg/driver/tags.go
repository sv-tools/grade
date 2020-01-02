package driver

import (
	"fmt"
	"strings"
)

type Tags map[string]string

func ParseTags(tags []string) (Tags, error) {
	res := make(Tags)
	for _, tag := range tags {
		fields := strings.SplitN(tag, "=", 2)
		if len(fields) != 2 {
			return nil, fmt.Errorf("incorrect tag: %s", tag)
		}
		res[strings.TrimSpace(fields[0])] = strings.TrimSpace(fields[1])
	}
	return res, nil
}

package templates

import (
	"encoding/json"
	"fmt"
	"strings"
)

func _unmarshalString[T any](value string) (buf T, err error) {
	var data []byte

	if strings.Contains(value, ",") {
		data, err = json.Marshal(strings.Split(value, ","))
		if err != nil {
			return buf, fmt.Errorf("list marshal: %w", err)
		}
	} else {
		data = []byte(value)
	}

	return _unmarshalBytes[T](data)
}

func _unmarshalBytes[T any](value []byte) (buf T, err error) {
	if err := json.Unmarshal(value, &buf); err != nil {
		return buf, fmt.Errorf("unmarshal: %w", err)
	}

	return buf, nil
}

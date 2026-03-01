package templates

import (
	"encoding/json"
	"fmt"
)

func _unmarshalString[T any](value string) (buf T, err error) {
	return _unmarshalBytes[T]([]byte(value))
}

func _unmarshalBytes[T any](value []byte) (buf T, err error) {
	if err = json.Unmarshal(value, &buf); err != nil {
		return buf, fmt.Errorf("unmarshal: %w", err)
	}

	return buf, nil
}

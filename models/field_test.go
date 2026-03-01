package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestField_IsError(t *testing.T) {
	tt := []struct {
		name    string
		field   Field
		isError bool
	}{
		{
			name:    "Field isn't an error",
			field:   Field{TypeExpr: "notError"},
			isError: false,
		},
		{
			name:    "Field is an error",
			field:   Field{TypeExpr: "error"},
			isError: true,
		},
		{
			name:    "Empty type expr",
			field:   Field{TypeExpr: ""},
			isError: false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.isError, tc.field.IsError())
		})
	}
}

func TestField_NameCapitalized(t *testing.T) {
	tt := []struct {
		name    string
		field   Field
		nameCap string
	}{
		{
			name:    "Name not capitalized yet",
			field:   Field{Name: "name"},
			nameCap: "Name",
		},
		{
			name:    "Name already capitalized",
			field:   Field{Name: "Name"},
			nameCap: "Name",
		},
		{
			name:    "Name uppercased",
			field:   Field{Name: "NAME"},
			nameCap: "NAME",
		},
		{
			name:    "Name empty",
			field:   Field{Name: ""},
			nameCap: "",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.nameCap, tc.field.NameCapitalized())
		})
	}
}

func TestField_JSONTag(t *testing.T) {
	tt := []struct {
		name    string
		field   Field
		jsonTag string
	}{
		{
			name:    "Name not capitalized yet",
			field:   Field{Name: "name"},
			jsonTag: "`json:\"name\"`",
		},
		{
			name:    "Name already capitalized",
			field:   Field{Name: "Name"},
			jsonTag: "`json:\"Name\"`",
		},
		{
			name:    "Name uppercased",
			field:   Field{Name: "NAME"},
			jsonTag: "`json:\"NAME\"`",
		},
		{
			name:    "Name empty",
			field:   Field{Name: ""},
			jsonTag: "`json:\"\"`",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.jsonTag, tc.field.JSONTag())
		})
	}
}

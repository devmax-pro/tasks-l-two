package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCut(t *testing.T) {
	tests := []struct {
		name   string
		row    string
		config *Config
		want   string
	}{
		{
			name:   "Simple input with single field",
			row:    "field1\tfield2\tfield3",
			config: &Config{fields: "1", delimiter: "\t", separated: false},
			want:   "field1",
		},
		{
			name:   "Simple input with multiple fields",
			row:    "field1\tfield2\tfield3",
			config: &Config{fields: "1,3", delimiter: "\t", separated: false},
			want:   "field1\tfield3",
		},
		{
			name:   "Custom delimiter",
			row:    "field1,field2,field3",
			config: &Config{fields: "2", delimiter: ",", separated: false},
			want:   "field2",
		},
		{
			name:   "Separated flag with delimiter",
			row:    "field1\tfield2\tfield3",
			config: &Config{fields: "2", delimiter: "\t", separated: true},
			want:   "field2",
		},
		{
			name:   "Separated flag without delimiter",
			row:    "field1 field2 field3",
			config: &Config{fields: "2", delimiter: "\t", separated: true},
			want:   "",
		},
		{
			name:   "Field out of range",
			row:    "field1\tfield2\tfield3",
			config: &Config{fields: "4", delimiter: "\t", separated: false},
			want:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := cut(tt.row, tt.config)
			assert.Equal(t, tt.want, got)
			if err != nil {
				t.Errorf("Expected not error, got %v", err)
			}
		})
	}
}

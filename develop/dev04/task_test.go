package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAnagram(t *testing.T) {
	testCases := []struct {
		desc string
		data []string
		want map[string][]string
	}{
		{
			desc: "normal",
			data: []string{
				"пятак", "пятка", "тяпка", "листок", "слиток", "столик", "столик", "листок",
			},
			want: map[string][]string{
				"листок": {"слиток", "столик"},
				"пятак":  {"пятка", "тяпка"},
			},
		},
		{
			desc: "none",
			data: []string{
				"4hffiof", "erhjbhopjas", "sffdfaasdfk",
			},
			want: map[string][]string{},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			got := anagram(tC.data)
			assert.Equal(t, tC.want, got)
		})
	}
}

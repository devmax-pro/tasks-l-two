package main

import "testing"

func TestUnpack(t *testing.T) {
	testCases := []struct {
		desc      string
		input     string
		want      string
		haveError bool
	}{
		{
			desc:  "normal",
			input: "a4df5cvs1",
			want:  "aaaadfffffcvs",
		},
		{
			desc:  "no numbers, only letters",
			input: "abcd",
			want:  "abcd",
		},
		{
			desc:      "only numbers",
			input:     "45",
			want:      "",
			haveError: true,
		},
		{
			desc:  "empty string",
			input: "",
			want:  "",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			got, err := unpackString(tC.input)
			if err == nil && tC.haveError {
				t.Error("Expected error, but err is nil")
				return
			}
			if err != nil && !tC.haveError {
				t.Errorf("Error: %v", err)
				return
			}
			if got != tC.want {
				t.Errorf("got: %s, want: %s", got, tC.want)
			}
		})
	}
}

func TestUnpackEscape(t *testing.T) {
	data := map[string]string{
		"qwe\\4\\5": "qwe45",
		"qwe\\45":   "qwe44444",
		"qwe\\\\5":  "qwe\\\\\\\\\\",
	}

	for s, e := range data {
		r, err := unpackString(s)
		if err != nil {
			t.Errorf("bad unpack for %s: got error %v", s, err)
		}
		if r != e {
			t.Errorf("bad unpack for %s: got %v expected %v", s, r, e)
		}
	}
}

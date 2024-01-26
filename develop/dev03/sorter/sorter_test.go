package sorter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSort(t *testing.T) {
	testCases := []struct {
		desc  string
		flags inputFlags
		data  []string
		want  []string
	}{
		{
			desc:  "normal",
			flags: inputFlags{},
			data: []string{
				"Standing on one's head at job interviews forms a lasting impression.",
				"I want to play game.",
				"The chic gangster liked to start the day with a pink scarf.",
				"He kept telling himself that one day it would all somehow make sense.",
				"Cats are good pets, for they are clean and are not noisy.",
				"I am my aunt's sister's daughter.",
			},
			want: []string{
				"Cats are good pets, for they are clean and are not noisy.",
				"He kept telling himself that one day it would all somehow make sense.",
				"I am my aunt's sister's daughter.",
				"I want to play game.",
				"Standing on one's head at job interviews forms a lasting impression.",
				"The chic gangster liked to start the day with a pink scarf.",
			},
		},
		{
			desc:  "not numeric order",
			flags: inputFlags{},
			data: []string{
				"1",
				"5",
				"13",
				"23",
				"11",
				"21",
				"31",
			},
			want: []string{
				"1",
				"11",
				"13",
				"21",
				"23",
				"31",
				"5",
			},
		},
		{
			desc: "reverse order",
			flags: inputFlags{
				reverse: true,
			},
			data: []string{
				"1",
				"5",
				"13",
				"23",
				"11",
				"21",
				"31",
			},
			want: []string{
				"5",
				"31",
				"23",
				"21",
				"13",
				"11",
				"1",
			},
		},
		{
			desc: "delete duplicate",
			flags: inputFlags{
				unique: true,
			},
			data: []string{
				"1",
				"1",
				"5",
				"13",
				"23",
				"11",
				"11",
				"21",
				"31",
				"31",
				"31",
			},
			want: []string{
				"1",
				"11",
				"13",
				"21",
				"23",
				"31",
				"5",
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			got := sortLines(tC.data, tC.flags)
			assert.Equal(t, tC.want, got)
		})
	}
}

func TestSortColumns(t *testing.T) {
	testCases := []struct {
		desc  string
		flags inputFlags
		data  []string
		want  []string
	}{
		{
			desc: "by 2nd column",
			flags: inputFlags{
				column: 2,
			},
			data: []string{
				"I want to play game.",
				"Standing on one's head at job interviews forms a lasting impression.",
				"The chic gangster liked to start the day with a pink scarf.",
				"He kept telling himself that one day it would all somehow make sense.",
				"Cats are good pets, for they are clean and are not noisy.",
				"I am my aunt's sister's daughter.",
			},
			want: []string{
				"I am my aunt's sister's daughter.",
				"Cats are good pets, for they are clean and are not noisy.",
				"The chic gangster liked to start the day with a pink scarf.",
				"He kept telling himself that one day it would all somehow make sense.",
				"Standing on one's head at job interviews forms a lasting impression.",
				"I want to play game.",
			},
		},
		{
			desc: "by column out of range",
			flags: inputFlags{
				column: 200,
			},
			data: []string{
				"Standing on one's head at job interviews forms a lasting impression.",
				"The chic gangster liked to start the day with a pink scarf.",
				"He kept telling himself that one day it would all somehow make sense.",
				"Cats are good pets, for they are clean and are not noisy.",
				"I am my aunt's sister's daughter.",
			},
			want: []string{
				"Cats are good pets, for they are clean and are not noisy.",
				"He kept telling himself that one day it would all somehow make sense.",
				"I am my aunt's sister's daughter.",
				"Standing on one's head at job interviews forms a lasting impression.",
				"The chic gangster liked to start the day with a pink scarf.",
			},
		},
		{
			desc: "numbers numeric order",
			flags: inputFlags{
				numeric: true,
			},
			data: []string{
				"5",
				"23",
				"1",
				"21",
				"31",
				"13",
				"11",
			},
			want: []string{
				"1",
				"5",
				"11",
				"13",
				"21",
				"23",
				"31",
			},
		},
		{
			desc: "by 2nd column, in reverse",
			flags: inputFlags{
				column:  2,
				reverse: true,
			},
			data: []string{
				"Standing on one's head at job interviews forms a lasting impression.",
				"The chic gangster liked to start the day with a pink scarf.",
				"He kept telling himself that one day it would all somehow make sense.",
				"Cats are good pets, for they are clean and are not noisy.",
				"I am my aunt's sister's daughter.",
			},
			want: []string{
				"Standing on one's head at job interviews forms a lasting impression.",
				"He kept telling himself that one day it would all somehow make sense.",
				"The chic gangster liked to start the day with a pink scarf.",
				"Cats are good pets, for they are clean and are not noisy.",
				"I am my aunt's sister's daughter.",
			},
		},
		{
			desc: "delete duplicate",
			flags: inputFlags{
				column: 1,
				unique: true,
			},
			data: []string{
				"1",
				"1",
				"5",
				"13",
				"23",
				"11",
				"11",
				"21",
				"31",
				"31",
				"31",
			},
			want: []string{
				"1",
				"11",
				"13",
				"21",
				"23",
				"31",
				"5",
			},
		},
		{
			desc: "numeric sort, but column starts with letter",
			flags: inputFlags{
				column:  1,
				numeric: true,
			},
			data: []string{
				"d1",
				"ad5",
				"asbv13",
				"sfg23",
				"fa11",
				"gh21",
				"31",
			},
			want: []string{
				"31",
				"ad5",
				"asbv13",
				"d1",
				"fa11",
				"gh21",
				"sfg23",
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			got := sortLines(tC.data, tC.flags)
			assert.Equal(t, tC.want, got)
		})
	}
}

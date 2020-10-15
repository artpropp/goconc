package main

import (
	"reflect"
	"testing"
)

func TestParseArgs(t *testing.T) {
	tests := []struct {
		name  string
		input []string
		want  []cmdArgs
	}{
		{
			"two commands",
			[]string{"abc", "def", "::", "abcd", "efgh"},
			[]cmdArgs{
				{"abc", []string{"def"}},
				{"abcd", []string{"efgh"}},
			},
		},
		{
			"single command",
			[]string{"abc", "def", "ghi"},
			[]cmdArgs{
				{"abc", []string{"def", "ghi"}},
			},
		},
		{
			"no command",
			[]string{},
			nil,
		},
		{
			"two commands with short notation",
			[]string{"abc", "def", ":", "ghi"},
			[]cmdArgs{
				{"abc", []string{"def"}},
				{"abc", []string{"ghi"}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseArgs(tt.input)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseArgs() = %v, want %v", got, tt.want)
			}
		})
	}

}

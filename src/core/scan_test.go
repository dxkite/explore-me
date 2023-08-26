package core

import (
	"reflect"
	"regexp"
	"testing"
)

func TestParseTag(t *testing.T) {
	reg, _ := regexp.Compile("\\[(.+?)\\]")
	tests := []struct {
		name string
		want []string
	}{
		{"[CTF][Web][Pwn] Ctf.md", []string{"CTF", "Web", "Pwn"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseTag(tt.name, reg)

			if err != nil {
				t.Error(err)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseTag() = %v, want %v", got, tt.want)
			}
		})
	}
}

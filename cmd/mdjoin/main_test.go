package main

import (
	"path/filepath"
	"testing"
)

func Test_checkFile(t *testing.T) {
	type args struct {
		fpath string
		skip  string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			"skip",
			args{
				filepath.Join("ab", "cd", "skip.md"),
				"abc.md, skip.md, def.txt",
			},
			false,
		},
		{
			"don't skip",
			args{
				filepath.Join("ab", "cd", "file.md"),
				"abc.md, skip.md, def.txt",
			},
			true,
		},
		{
			"not md file",
			args{
				filepath.Join("ab", "cd", "file.txt"),
				"abc.md, skip.md, def.txt",
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := checkFile(tt.args.fpath, tt.args.skip); got != tt.want {
				t.Errorf("checkFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

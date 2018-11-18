package md

import (
	"io"
	"io/ioutil"
	"reflect"
	"strings"
	"testing"
)

func TestRemoveHeader(t *testing.T) {
	type args struct {
		r         io.Reader
		headerEnd string
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{
			"with header",
			args{strings.NewReader("abc\ndef\n---\nddd\neee\n"), "---"},
			[]byte("ddd\neee\n"),
		},
		{
			"with header 2",
			args{strings.NewReader("abc\ndef\n---  \nddd\neee\n"), "---"},
			[]byte("ddd\neee\n"),
		},
		{
			"without header",
			args{strings.NewReader("abc\ndef\nddd\neee\n"), "---"},
			[]byte("abc\ndef\nddd\neee\n"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotReader := RemoveHeader(tt.args.r, tt.args.headerEnd)
			got, _ := ioutil.ReadAll(gotReader)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RemoveHeader() = %v, want %v", string(got), string(tt.want))
			}
		})
	}
}

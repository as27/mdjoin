package md

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"strings"
)

// RemoveHeader ready the io.Reader until headerEnd and returns a io.Reader
// which contains just the body
func RemoveHeader(r io.Reader, headerEnd string) io.Reader {
	// prepare if there is no header
	buf := &bytes.Buffer{}
	tee := io.TeeReader(r, buf)
	br := bufio.NewReader(tee)
	for {
		s, err := br.ReadString('\n')
		if err != nil {
			// when there is no more data and no header found the
			// buffer of the teeReader is returned
			if err == io.EOF {
				return buf
			}
			log.Println("error reading string", err)
		}
		if strings.HasPrefix(s, headerEnd) {
			break
		}
	}
	return br
}

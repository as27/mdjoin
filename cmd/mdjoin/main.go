package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/as27/mdjoin/pkg/md"
)

var (
	flagOutFile   = flag.String("out", "document.md", "name of the output file")
	flagSkipFiles = flag.String("skip", "", "comma separated list of filenames to skip")
)

func main() {
	flag.Parse()
	root := flag.Arg(0)
	if root == "" {
		root = "./"
	}
	f, err := os.OpenFile(*flagOutFile, os.O_CREATE, 0777)
	if err != nil {
		fmt.Println("cannot open output file", err)
		os.Exit(1)
	}
	defer f.Close()
	walk(root, f, *flagSkipFiles)
}

func walk(root string, w io.Writer, skips string) {
	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() || !checkFile(path, skips) {
			return err
		}
		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()
		body := md.RemoveHeader(f, "---")

		fileinfo := fmt.Sprintf("\n*File: %s*\n---\n", info.Name())
		_, err = w.Write([]byte(fileinfo))
		if err != nil {
			log.Println("error writing file info", err)
		}

		_, err = io.Copy(w, body)
		if err != nil {
			log.Println("error when copy file", err)
		}

		return err
	})
}

func checkFile(fpath string, skips string) bool {
	if strings.ToLower(filepath.Ext(fpath)) != ".md" {
		return false
	}
	skipFile := strings.Contains(skips, filepath.Base(fpath))
	return !skipFile
}

package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"

	"github.com/as27/mdjoin/pkg/md"
)

var (
	flagOutFile    = flag.String("out", "document.md", "name of the output file")
	flagSkipFiles  = flag.String("skip", "", "comma separated list of filenames to skip")
	flagConfigFile = flag.String("conf", "", "use a yaml config file")
)

func config() {
	viper.SetDefault("skip", *flagSkipFiles)
	viper.SetDefault("out", *flagOutFile)
	if *flagConfigFile != "" {
		viper.SetConfigFile(*flagConfigFile)
		viper.SetConfigType("yaml")
		err := viper.ReadInConfig()
		if err != nil {
			log.Println("error reading config", err)
		}
	}
}

func main() {
	flag.Parse()
	root := flag.Arg(0)
	if root == "" {
		root = "./"
	}
	config()
	outFile := viper.GetString("out")
	skips := viper.GetString("skip")
	f, err := os.OpenFile(outFile, os.O_CREATE, 0777)
	if err != nil {
		fmt.Println("cannot open output file", err)
		os.Exit(1)
	}
	defer f.Close()
	skips = fmt.Sprintf("%s,%s", skips, outFile)
	walk(root, f, skips)
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

		_, err = io.Copy(w, body)
		if err != nil {
			log.Println("error when copy file", err)
		}

		fileinfo := fmt.Sprintf("\n_File: %s_\n---\n", info.Name())
		_, err = w.Write([]byte(fileinfo))
		if err != nil {
			log.Println("error writing file info", err)
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

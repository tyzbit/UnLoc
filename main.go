package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
	"unicode/utf8"
)

func CollectFiles(path, extensions string) (filePaths []string) {
	if !strings.HasSuffix("/", path) {
		path = path + "/"
	}
	dir, err := os.ReadDir(path)
	if err != nil {
		fmt.Println("error reading dir ", path, ": ", err)
		return
	}
	for _, f := range dir {
		if strings.HasPrefix(f.Name(), ".") {
			continue
		}
		i, err := f.Info()
		if err != nil {
			fmt.Println("error getting info for file ", f.Name(), ": ", err)
		}
		if i.IsDir() {
			fmt.Println("found dir ", i.Name())
			files := CollectFiles(path+i.Name()+"/", extensions)
			if err != nil {
				fmt.Println("error collecting files from dir ", i.Name(), ": ", err)
			}
			filePaths = append(filePaths, files...)
		} else {
			hasSuffix := false
			for _, e := range strings.Split(extensions, ",") {
				if strings.HasSuffix(f.Name(), e) {
					hasSuffix = true
				}
			}
			if !hasSuffix {
				continue
			}
			fmt.Println("found file", i.Name())
			filePaths = append(filePaths, path+i.Name())
		}
	}
	return filePaths
}

func CalculateWords(content []byte) int {
	return len(strings.Fields(string(content)))
}

func main() {
	path := flag.String("path", "", "Path to traverse")
	extensions := flag.String("extensions", ".go", "Comma-separated list of extensions to include")
	wordsPerMinute := flag.Int("wpm", 60, "Words per minute")
	flag.Parse()
	if *path == "" {
		fmt.Println("Calculate how long it would take to write your app if you " +
			"wrote it all in one go with no mistakes.")
		flag.PrintDefaults()
		os.Exit(0)
	}
	files := CollectFiles(*path, *extensions)
	var contents [][]byte
	for _, file := range files {
		r, err := os.ReadFile(file)
		if err != nil {
			fmt.Println("error reading file ", file, ": ", err)
			os.Exit(1)
		}
		contents = append(contents, r)
	}
	var wordCount float64
	for _, c := range contents {
		if utf8.ValidString(string(c)) {
			wordCount = wordCount + float64(CalculateWords(c))
		}
	}
	e := wordCount / float64(*wordsPerMinute)
	fmt.Println("estimate: ", time.Duration(e*float64(time.Minute)).String())
}

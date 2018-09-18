package main

import (
	"flag"
	"fmt"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	inputFile := ""
	flag.StringVar(&inputFile, "f", "", "output file")
	flag.Parse()

	time.Now().Sub(time.Now()).String()

	data, err := parseFile(inputFile)
	if err != nil {
		fmt.Printf("parse file error:%s", err)
		return
	}

	data.debugPrint()

	outputPath := getOutputPath(inputFile)
	render(outputPath, data)
}

func getOutputPath(inputPath string) string {
	if !strings.HasSuffix(inputPath, ".go") {
		return ""
	}
	dir, file := filepath.Split(inputPath[:len(inputPath)-3])
	return filepath.Join(dir, fmt.Sprintf("%s_goaop.go", file))
}

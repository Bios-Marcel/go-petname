package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	path := os.Args[1]

	fmt.Printf(
		`package %s

var (
	Adjectives = []string{%s}
	Adverbs = []string{%s}
	Names = []string{%s}
)

// vim: nowrap`,
		filepath.Base(path),
		toCommaSeparatedStringLiterals(filepath.Join(path, "adjectives.txt")),
		toCommaSeparatedStringLiterals(filepath.Join(path, "adverbs.txt")),
		toCommaSeparatedStringLiterals(filepath.Join(path, "names.txt")))
}

func toCommaSeparatedStringLiterals(path string) string {
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	scan := bufio.NewScanner(f)
	var list strings.Builder
	for scan.Scan() {
		list.WriteRune('"')
		list.WriteString(scan.Text())
		list.WriteString(`", `)
	}

	return list.String()
}

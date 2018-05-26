package main

import (
	"flag"
	"io/ioutil"
	"regexp"
)

var (
	syntax *regexp.Regexp
	file   string
	delay  int
)

func init() {
	syntax = regexp.MustCompile(`[^><+-.,[\]]`)

	flag.StringVar(&file, "file", "", "The file to load, including extension.")
	flag.IntVar(&delay, "delay", 1, "The time in ms to delay each operation.")

	flag.Parse()
}

func main() {
	file, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}

	bf := BrainfuckConstructor(string(file))
	bf.Execute()
}

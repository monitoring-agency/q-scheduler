package main

import (
	"github.com/hellflame/argparse"
)

func main() {
	parser := argparse.NewParser("q-scheduler", "", &argparse.ParserConfig{})

	parser.Parse(nil)

}

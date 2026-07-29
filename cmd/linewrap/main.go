package main

import (
	"golang.org/x/tools/go/analysis/singlechecker"

	"lesiw.io/linewrap"
)

func main() { singlechecker.Main(linewrap.Analyzer) }

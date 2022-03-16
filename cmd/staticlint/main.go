package main

import (
	"golang.org/x/tools/go/analysis/multichecker"
	"shortener/pkg/linter"
)

func main() {
	multichecker.Main(
		linter.DefaultAnalysers()...,
	)
}

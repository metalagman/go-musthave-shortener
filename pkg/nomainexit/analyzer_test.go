package nomainexit

import (
	"golang.org/x/tools/go/analysis/multichecker"
)

func Example() {
	multichecker.Main(Analyzer)
}

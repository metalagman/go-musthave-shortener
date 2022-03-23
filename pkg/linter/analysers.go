package linter

import (
	errname "github.com/Antonboom/errname/pkg/analyzer"
	errcheck "github.com/kisielk/errcheck/errcheck"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/asmdecl"
	"golang.org/x/tools/go/analysis/passes/assign"
	"golang.org/x/tools/go/analysis/passes/atomic"
	"golang.org/x/tools/go/analysis/passes/atomicalign"
	"golang.org/x/tools/go/analysis/passes/bools"
	"golang.org/x/tools/go/analysis/passes/buildssa"
	"golang.org/x/tools/go/analysis/passes/buildtag"
	"golang.org/x/tools/go/analysis/passes/cgocall"
	"golang.org/x/tools/go/analysis/passes/composite"
	"golang.org/x/tools/go/analysis/passes/copylock"
	"golang.org/x/tools/go/analysis/passes/ctrlflow"
	"golang.org/x/tools/go/analysis/passes/deepequalerrors"
	"golang.org/x/tools/go/analysis/passes/errorsas"
	"golang.org/x/tools/go/analysis/passes/findcall"
	"golang.org/x/tools/go/analysis/passes/framepointer"
	"golang.org/x/tools/go/analysis/passes/httpresponse"
	"golang.org/x/tools/go/analysis/passes/ifaceassert"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/analysis/passes/loopclosure"
	"golang.org/x/tools/go/analysis/passes/lostcancel"
	"golang.org/x/tools/go/analysis/passes/nilfunc"
	"golang.org/x/tools/go/analysis/passes/nilness"
	"golang.org/x/tools/go/analysis/passes/pkgfact"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/reflectvaluecompare"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/shift"
	"golang.org/x/tools/go/analysis/passes/sigchanyzer"
	"golang.org/x/tools/go/analysis/passes/sortslice"
	"golang.org/x/tools/go/analysis/passes/stdmethods"
	"golang.org/x/tools/go/analysis/passes/stringintconv"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"golang.org/x/tools/go/analysis/passes/testinggoroutine"
	"golang.org/x/tools/go/analysis/passes/tests"
	"golang.org/x/tools/go/analysis/passes/unmarshal"
	"golang.org/x/tools/go/analysis/passes/unreachable"
	"golang.org/x/tools/go/analysis/passes/unsafeptr"
	"golang.org/x/tools/go/analysis/passes/unusedresult"
	"golang.org/x/tools/go/analysis/passes/unusedwrite"
	"golang.org/x/tools/go/analysis/passes/usesgenerics"
	"honnef.co/go/tools/quickfix"
	"honnef.co/go/tools/simple"
	"honnef.co/go/tools/staticcheck"
	"honnef.co/go/tools/stylecheck"
	"shortener/pkg/nomainexit"
)

// All available analyzers collection
func All() []*analysis.Analyzer {
	var list []*analysis.Analyzer
	list = append(list, DefaultAnalysers()...)
	list = append(list, StaticCheckSA()...)
	//list = append(list, StaticCheckST()...)
	list = append(list, StaticCheckS()...)
	list = append(list, StaticCheckQF()...)
	list = append(list, CustomAnalysers()...)
	list = append(list, nomainexit.Analyzer)
	return list
}

// CustomAnalysers collection
func CustomAnalysers() []*analysis.Analyzer {
	return []*analysis.Analyzer{
		errname.New(),
		errcheck.Analyzer,
	}
}

// DefaultAnalysers collection
func DefaultAnalysers() []*analysis.Analyzer {
	return []*analysis.Analyzer{
		// The traditional vet suite:
		asmdecl.Analyzer,
		assign.Analyzer,
		atomic.Analyzer,
		bools.Analyzer,
		buildssa.Analyzer,
		buildtag.Analyzer,
		cgocall.Analyzer,
		composite.Analyzer,
		copylock.Analyzer,
		ctrlflow.Analyzer,
		errorsas.Analyzer,
		findcall.Analyzer,
		framepointer.Analyzer,
		httpresponse.Analyzer,
		ifaceassert.Analyzer,
		inspect.Analyzer,
		loopclosure.Analyzer,
		lostcancel.Analyzer,
		nilfunc.Analyzer,
		pkgfact.Analyzer,
		printf.Analyzer,
		reflectvaluecompare.Analyzer,
		shift.Analyzer,
		sigchanyzer.Analyzer,
		stdmethods.Analyzer,
		stringintconv.Analyzer,
		structtag.Analyzer,
		tests.Analyzer,
		unmarshal.Analyzer,
		unreachable.Analyzer,
		unsafeptr.Analyzer,
		unusedresult.Analyzer,
		usesgenerics.Analyzer,

		// Non-vet analyzers:
		atomicalign.Analyzer,
		deepequalerrors.Analyzer,
		//fieldalignment.Analyzer,
		nilness.Analyzer,
		shadow.Analyzer,
		sortslice.Analyzer,
		testinggoroutine.Analyzer,
		unusedwrite.Analyzer,
	}
}

// StaticCheckS is a collection of staticcheck.io S analyzers
func StaticCheckS() []*analysis.Analyzer {
	var res []*analysis.Analyzer
	for _, v := range simple.Analyzers {
		res = append(res, v.Analyzer)
	}
	return res
}

// StaticCheckSA is a collection of staticcheck.io SA analyzers
func StaticCheckSA() []*analysis.Analyzer {
	var res []*analysis.Analyzer
	for _, v := range staticcheck.Analyzers {
		res = append(res, v.Analyzer)
	}
	return res
}

// StaticCheckST is a collection of staticcheck.io ST analyzers
func StaticCheckST() []*analysis.Analyzer {
	var res []*analysis.Analyzer
	for _, v := range stylecheck.Analyzers {
		res = append(res, v.Analyzer)
	}
	return res
}

// StaticCheckQF is a collection of staticcheck.io QF analyzers
func StaticCheckQF() []*analysis.Analyzer {
	var res []*analysis.Analyzer
	for _, v := range quickfix.Analyzers {
		res = append(res, v.Analyzer)
	}
	return res
}

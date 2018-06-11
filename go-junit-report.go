package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/jstemmer/go-junit-report/formatter"
	"github.com/jstemmer/go-junit-report/parser"
)

var (
	noXMLHeader   bool
	packageName   string
	goVersionFlag string
	setExitCode   bool
	outputFile    string
	verbose       bool
)

func init() {
	flag.BoolVar(&noXMLHeader, "no-xml-header", false, "do not print xml header")
	flag.StringVar(&packageName, "package-name", "", "specify a package name (compiled test have no package name in output)")
	flag.StringVar(&goVersionFlag, "go-version", "", "specify the value to use for the go.version property in the generated XML")
	flag.BoolVar(&setExitCode, "set-exit-code", false, "set exit code to 1 if tests failed")
	flag.StringVar(&outputFile, "output-file", "", "specify the output file to print the test report to, if specified then stdout will print the original output")
	flag.BoolVar(&verbose, "verbose", false, "print verbose test output, equivalent to test -v")
}

func main() {
	flag.Parse()

	outputLevel :=  parser.NoOutput
	if(outputFile != "") {
		if(verbose) {
			outputLevel = parser.FullOutput
		} else {
			outputLevel = parser.BasicOutput
		}
	}

	// Read input
	report, err := parser.Parse(os.Stdin, packageName, outputLevel)
	if err != nil {
		fmt.Printf("Error reading input: %s\n", err)
		os.Exit(1)
	}

	writer := os.Stdout

	if(outputFile != "") {
		f, err := os.Create(outputFile)
		if err != nil {
			fmt.Printf("Error opening output file: %s\n", err)
			os.Exit(1)
		}
		writer = f
	
		defer f.Close()
	}

	// Write xml
	err = formatter.JUnitReportXML(report, noXMLHeader, goVersionFlag, writer)
	if err != nil {
		fmt.Printf("Error writing XML: %s\n", err)
		os.Exit(1)
	}

	if setExitCode && report.Failures() > 0 {
		os.Exit(1)
	}
}

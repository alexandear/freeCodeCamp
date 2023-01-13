package gotest

import (
	"strings"
)

// status represents single test status.
type status string

const (
	Passed  status = "passed"
	Skipped status = "skipped"
	Failed  status = "failed"
	Unknown status = "unknown"
)

var knownStatuses = map[status]struct{}{
	Passed:  {},
	Skipped: {},
	Failed:  {},
}

// TestResult represents single test result (status and output).
type TestResult struct {
	Status status
	Output string
}

func (tr *TestResult) IndentedOutput() string {
	return strings.ReplaceAll(tr.Output, "\n", "\n\t")
}

// TestResults represents results of a multiple tests.
//
// They are returned by runners.
type TestResults struct {
	// Test results by full test name.
	TestResults map[string]TestResult
}

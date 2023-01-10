package gotest

import (
	"time"
)

type action string

//nolint:godot // false positive for unexported identifiers
const (
	// actionRun means the test has started running.
	actionRun action = "run"
	// actionPause means the test has been paused.
	actionPause action = "pause"
	// actionCont means the test has continued running.
	actionCont action = "cont"
	// actionPass means the test passed.
	actionPass action = "pass"
	// actionBench means the benchmark printed log output but did not fail.
	actionBench action = "bench"
	// actionFail means the test or benchmark failed.
	actionFail action = "fail"
	// actionOutput means the test printed output.
	actionOutput action = "output"
	// actionSkip means the test was skipped or the package contained no tests.
	actionSkip action = "skip"
)

type testEvent struct {
	Time           time.Time `json:"Time"`
	Action         action    `json:"Action"`
	Package        string    `json:"Package"`
	Test           string    `json:"Test"`
	Output         string    `json:"Output"`
	ElapsedSeconds float64   `json:"Elapsed"`
}

func (te testEvent) Elapsed() time.Duration {
	return time.Duration(te.ElapsedSeconds * float64(time.Second))
}

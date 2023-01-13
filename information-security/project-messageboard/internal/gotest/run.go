package gotest

import (
	"context"
	"encoding/json"
	"io"
	"os"
	"os/exec"
)

func Run(ctx context.Context, dir string, args []string, verbose bool) (*TestResults, error) {
	args = append([]string{"test", "-v", "-json", "-count=1"}, args...)
	cmd := exec.CommandContext(ctx, "go", args...)
	cmd.Dir = dir
	cmd.Stderr = os.Stderr
	p, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	var r io.Reader = p
	if verbose {
		r = io.TeeReader(p, os.Stdout)
	}

	if err = cmd.Start(); err != nil {
		return nil, err
	}

	d := json.NewDecoder(r)
	d.DisallowUnknownFields()

	res := &TestResults{
		TestResults: make(map[string]TestResult),
	}

	for {
		var event testEvent
		if err = d.Decode(&event); err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		// skip package failures
		if event.Test == "" {
			continue
		}

		testName := event.Test
		result := res.TestResults[testName]
		result.Output += event.Output
		result.Assertions = []Assertion{MakeEqualAssertion(), MakePropertyAssertion()} // TODO: replace with real assertions
		switch event.Action {
		case actionPass:
			result.Status = Passed
		case actionFail:
			result.Status = Failed
		case actionSkip:
			result.Status = Skipped
		case actionBench, actionCont, actionOutput, actionPause, actionRun:
			fallthrough
		default:
			result.Status = Unknown
		}
		res.TestResults[testName] = result
	}

	if err = cmd.Wait(); err != nil {
		if _, ok := err.(*exec.ExitError); ok {
			err = nil
		}
	}

	return res, err
}

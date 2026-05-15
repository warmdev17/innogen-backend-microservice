package judge

import (
	"math"
	"strings"

	"innogen-backend/run_service/internal/dto"
)

// EvaluateResult holds the outcome of evaluating a single Piston execution.
type EvaluateResult struct {
	Status       string
	ActualOutput string
	RuntimeMs    *int
	ErrorMessage *string
}

// trimOutput removes leading and trailing whitespace from a string.
func trimOutput(s string) string {
	return strings.TrimSpace(s)
}

// Evaluate compares the Piston execution result against the expected output.
//
// Parameters:
//   - expectedOutput: the expected output from the test case
//   - compileStderr: non-empty if compilation failed (passed when compile code != 0)
//   - runStdout: stdout from the Piston run stage
//   - runStderr: stderr from the Piston run stage
//   - runCode: exit code from the Piston run stage (0 = success)
//   - runSignal: signal that killed the process, if any (nil if no signal)
//   - runTime: execution time in seconds (0 if not measured)
func Evaluate(expectedOutput, compileStderr, runStdout, runStderr string, runCode int, runSignal *string, runTime float64) EvaluateResult {
	// Compilation error
	if compileStderr != "" {
		msg := trimOutput(compileStderr)
		if msg == "" {
			msg = "Compilation error"
		}
		return EvaluateResult{
			Status:       dto.StatusCompilationError,
			ErrorMessage: &msg,
		}
	}

	// Killed by signal
	if runSignal != nil && *runSignal != "" {
		msg := "Killed by signal: " + *runSignal
		actual := trimOutput(runStdout)
		return EvaluateResult{
			Status:       dto.StatusRuntimeError,
			ActualOutput: actual,
			ErrorMessage: &msg,
		}
	}

	// Runtime error (non-zero exit code)
	if runCode != 0 {
		msg := trimOutput(runStderr)
		if msg == "" {
			msg = "Runtime error"
		}
		actual := trimOutput(runStdout)
		return EvaluateResult{
			Status:       dto.StatusRuntimeError,
			ActualOutput: actual,
			ErrorMessage: &msg,
		}
	}

	// Successful execution — compare outputs
	actual := trimOutput(runStdout)
	expected := trimOutput(expectedOutput)
	rtMs := runtimeToMs(runTime)

	status := dto.StatusWrongAnswer
	if actual == expected {
		status = dto.StatusAccepted
	}

	return EvaluateResult{
		Status:       status,
		ActualOutput: actual,
		RuntimeMs:    rtMs,
	}
}

// runtimeToMs converts a duration in seconds to milliseconds integer.
// If seconds is 0 or negative, returns nil.
func runtimeToMs(seconds float64) *int {
	if seconds <= 0 {
		return nil
	}
	ms := int(math.Round(seconds * 1000))
	return &ms
}

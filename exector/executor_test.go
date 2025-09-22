package exector

import (
	"context"
	"slices"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestRun(t *testing.T) {
	tests := []struct {
		name                                   string
		commands                               []string
		options                                []Option
		wantResultsWithoutTimesSortedByCommand []*Result
		within1Sec                             bool
		wantErr                                bool
	}{
		{
			name:     "single echo command",
			commands: []string{"echo hello"},
			options:  nil,
			wantResultsWithoutTimesSortedByCommand: []*Result{
				{
					Command:  "echo hello",
					Stdout:   []byte("hello\n"),
					Stderr:   nil,
					Combined: []byte("hello\n"),
					ExitCode: 0,
				},
			},
			within1Sec: true,
			wantErr:    false,
		},
		{
			name:     "multiple echo commands",
			commands: []string{"echo first", "echo second"},
			options:  nil,
			wantResultsWithoutTimesSortedByCommand: []*Result{
				{
					Command:  "echo first",
					Stdout:   []byte("first\n"),
					Stderr:   nil,
					Combined: []byte("first\n"),
					ExitCode: 0,
				},
				{
					Command:  "echo second",
					Stdout:   []byte("second\n"),
					Stderr:   nil,
					Combined: []byte("second\n"),
					ExitCode: 0,
				},
			},
			within1Sec: true,
			wantErr:    false,
		},
		{
			name:     "command with non-zero exit code",
			commands: []string{"false"},
			options:  nil,
			wantResultsWithoutTimesSortedByCommand: []*Result{
				{
					Command:  "false",
					Stdout:   nil,
					Stderr:   nil,
					Combined: nil,
					ExitCode: 1,
				},
			},
			within1Sec: true,
			wantErr:    false,
		},
		{
			name:     "with shell option",
			commands: []string{"echo shell-test"},
			options:  []Option{Shell("sh")},
			wantResultsWithoutTimesSortedByCommand: []*Result{
				{
					Command:  "echo shell-test",
					Stdout:   []byte("shell-test\n"),
					Stderr:   nil,
					Combined: []byte("shell-test\n"),
					ExitCode: 0,
				},
			},
			within1Sec: true,
			wantErr:    false,
		},
		{
			name:     "command with stderr output",
			commands: []string{"echo error message >&2"},
			options:  nil,
			wantResultsWithoutTimesSortedByCommand: []*Result{
				{
					Command:  "echo error message >&2",
					Stdout:   nil,
					Stderr:   []byte("error message\n"),
					Combined: []byte("error message\n"),
					ExitCode: 0,
				},
			},
			within1Sec: true,
			wantErr:    false,
		},
		{
			name:     "with fail fast option",
			commands: []string{"false", "sleep 10"},
			options:  []Option{FailFast(true)},
			wantResultsWithoutTimesSortedByCommand: []*Result{
				{
					Command:  "false",
					Stdout:   nil,
					Stderr:   nil,
					Combined: nil,
					ExitCode: 1,
				},
				{
					Command:  "sleep 10",
					Stdout:   nil,
					Stderr:   nil,
					Combined: nil,
					ExitCode: -1,
				},
			},
			within1Sec: true,
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e, err := New(tt.commands, tt.options...)
			if err != nil {
				t.Fatalf("failed to create executor: %v", err)
			}

			ctx := context.Background()
			start := time.Now()
			resultCh := make(chan *Result)
			errCh := make(chan error)

			go e.Run(ctx, resultCh, errCh)
			var results []*Result
		L:
			for r := range resultCh {
				results = append(results, r)
				select {
				case errr := <-errCh:
					err = errr
					break L
				default:
				}
			}
			if err == nil {
				err = <-errCh
			}
			elapsed := time.Since(start)

			if (err != nil) != tt.wantErr {
				t.Errorf("unexpected error status: got %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.within1Sec && elapsed > time.Second {
				t.Errorf("execution took too long: %v", elapsed)
				return
			}

			if len(results) != len(tt.wantResultsWithoutTimesSortedByCommand) {
				t.Errorf("unexpected results count: got %d, want %d", len(results), len(tt.wantResultsWithoutTimesSortedByCommand))
				return
			}

			slices.SortFunc(results, func(a, b *Result) int {
				if a.Command < b.Command {
					return -1
				}
				if a.Command > b.Command {
					return 1
				}
				return 0
			})

			opts := cmpopts.IgnoreFields(Result{}, "StartTime", "EndTime")
			if diff := cmp.Diff(tt.wantResultsWithoutTimesSortedByCommand, results, opts); diff != "" {
				t.Errorf("results mismatch (-want +got):\n%s", diff)
			}

			// Verify time fields separately
			for _, result := range results {
				if result.StartTime.IsZero() {
					t.Errorf("command %q has zero start time", result.Command)
				}
				if result.EndTime.IsZero() {
					t.Errorf("command %q has zero end time", result.Command)
				}
				if result.EndTime.Before(result.StartTime) {
					t.Errorf("command %q end time %v is before start time %v", result.Command, result.EndTime, result.StartTime)
				}
			}
		})
	}
}

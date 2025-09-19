package exector

import (
	"bytes"
	"context"
	"io"
	"os"
	"runtime"
	"time"

	"github.com/k1LoW/errors"
	"github.com/k1LoW/exec"
	"golang.org/x/sync/errgroup"
)

const DefaultShell = "bash"

type Executor struct {
	commands []string
	stdout   io.Writer
	stderr   io.Writer
	shell    string
	failFast bool
}

type Result struct {
	Command   string
	Stdout    []byte
	Stderr    []byte
	Combined  []byte
	ExitCode  int
	StartTime time.Time
	EndTime   time.Time
}

func New(commands []string, opts ...Option) (*Executor, error) {
	e := &Executor{
		commands: commands,
		stdout:   os.Stdout,
		stderr:   os.Stderr,
		shell:    DefaultShell,
	}
	for _, opt := range opts {
		if err := opt(e); err != nil {
			return nil, err
		}
	}
	return e, nil
}

func (e *Executor) Run(ctx context.Context) ([]*Result, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	sh, err := exec.LookPath(e.shell)
	if err != nil {
		return nil, err
	}
	var eg errgroup.Group
	var results []*Result
	eg.SetLimit(runtime.GOMAXPROCS(0) - 1)
	for _, c := range e.commands {
		eg.Go(func() error {
			cmd := exec.CommandContext(ctx, sh, "-c", c)
			stdout := &bytes.Buffer{}
			stderr := &bytes.Buffer{}
			combined := &bytes.Buffer{}
			cmd.Stdout = io.MultiWriter(stdout, combined)
			cmd.Stderr = io.MultiWriter(stderr, combined)
			exitCode := -1
			var start time.Time
			defer func() {
				result := &Result{
					Command:   c,
					Stdout:    stdout.Bytes(),
					Stderr:    stderr.Bytes(),
					Combined:  combined.Bytes(),
					ExitCode:  exitCode,
					StartTime: start,
					EndTime:   time.Now(),
				}
				results = append(results, result)
			}()
			start = time.Now()
			if err := cmd.Run(); err != nil {
				var exitErr *exec.ExitError
				if errors.As(err, &exitErr) {
					if e.failFast {
						cancel()
					}
					exitCode = cmd.ProcessState.ExitCode()
					return nil
				}
				return err
			}
			exitCode = cmd.ProcessState.ExitCode()
			return nil
		})
	}
	if err := eg.Wait(); err != nil {
		return nil, err
	}
	return results, nil
}

package exector

import (
	"bytes"
	"context"
	"io"
	"os"
	"time"

	"github.com/k1LoW/errors"
	"github.com/k1LoW/exec"
	"golang.org/x/sync/errgroup"
)

const DefaultShell = "bash"

type Executor struct {
	commands             []string
	stdout               io.Writer
	stderr               io.Writer
	shell                string
	failFast             bool
	maxRetriesPerCommand int
}

type Result struct {
	Index     int
	Command   string
	Stdout    []byte
	Stderr    []byte
	Combined  []byte
	ExitCode  int
	StartTime time.Time
	EndTime   time.Time
	Retries   int
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

func (e *Executor) Run(ctx context.Context, ch chan<- *Result, errCh chan<- error) {
	defer close(ch)
	defer close(errCh)
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	sh, err := exec.LookPath(e.shell)
	if err != nil {
		if errCh != nil {
			errCh <- err
		}
		return
	}
	var eg errgroup.Group
	for i, c := range e.commands {
		eg.Go(func() error {
			retries := 0
			finished := false
			lastExitCode := -1
			for {
				err := func() error {
					cmd := exec.CommandContext(ctx, sh, "-c", c)
					stdout := &bytes.Buffer{}
					stderr := &bytes.Buffer{}
					combined := &bytes.Buffer{}
					cmd.Stdout = io.MultiWriter(stdout, combined)
					cmd.Stderr = io.MultiWriter(stderr, combined)
					exitCode := -1
					start := time.Now()
					defer func() {
						result := &Result{
							Index:     i,
							Command:   c,
							Stdout:    stdout.Bytes(),
							Stderr:    stderr.Bytes(),
							Combined:  combined.Bytes(),
							ExitCode:  exitCode,
							StartTime: start,
							EndTime:   time.Now(),
							Retries:   retries,
						}
						if ch != nil {
							ch <- result
						}
						lastExitCode = exitCode
						retries++
					}()
					if err := cmd.Run(); err != nil {
						var exitErr *exec.ExitError
						if !errors.As(err, &exitErr) {
							return err
						}
						select {
						case <-ctx.Done():
							finished = true
							return nil
						default:
						}
					}
					exitCode = cmd.ProcessState.ExitCode()
					if exitCode == 0 {
						finished = true
						return nil
					}
					if retries >= e.maxRetriesPerCommand {
						finished = true
						return nil
					}
					return nil
				}()
				if err != nil {
					return err
				}
				if finished {
					break
				}
			}
			if e.failFast {
				if lastExitCode != 0 {
					cancel()
				}
			}
			return nil
		})
	}
	if err := eg.Wait(); err != nil {
		if errCh != nil {
			errCh <- err
		}
	}
}

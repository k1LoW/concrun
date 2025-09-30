package exector

import "github.com/k1LoW/errors"

type Option func(*Executor) error

func FailFast(failFast bool) Option {
	return func(e *Executor) error {
		e.failFast = failFast
		return nil
	}
}

func Shell(shell string) Option {
	return func(e *Executor) error {
		e.shell = shell
		return nil
	}
}

func MaxRetriesPerCommand(n int) Option {
	return func(e *Executor) error {
		if n < 0 {
			return errors.New("max retries per command must be 0 or greater")
		}
		e.maxRetriesPerCommand = n
		return nil
	}
}

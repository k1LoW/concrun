package exector

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

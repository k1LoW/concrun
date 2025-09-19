/*
Copyright © 2025 Ken'ichiro Oyama <k1lowxb@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/k1LoW/concrun/exector"
	"github.com/k1LoW/concrun/version"
	"github.com/spf13/cobra"
)

var (
	commands []string
	shell    string
	failFast bool
)

const commandPointer = "▶"

var rootCmd = &cobra.Command{
	Use:          "concrun",
	Short:        "Run commands concurrently",
	Long:         `Run commands concurrently.`,
	Version:      version.Version,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		opts := []exector.Option{
			exector.Shell(shell),
			exector.FailFast(failFast),
		}
		e, err := exector.New(commands, opts...)
		if err != nil {
			return err
		}
		results, err := e.Run(cmd.Context())
		if err != nil {
			return err
		}
		exitCode := 0
		killed := false
		for _, r := range results {
			fmt.Println("--------------------------------------------------")
			_, _ = fmt.Fprintf(os.Stdout, "%s %s\n", commandPointer, r.Command)
			_, _ = os.Stdout.Write(r.Combined)
			if r.ExitCode == -1 {
				killed = true
				fmt.Fprintln(os.Stderr, "(command was terminated by signal)")
			}
			if r.ExitCode > exitCode {
				exitCode = r.ExitCode
			}
			d := r.EndTime.Sub(r.StartTime)
			fmt.Printf("---- [ exit code: %d, excution time: %s ]\n", r.ExitCode, d)
		}
		if exitCode != 0 {
			os.Exit(exitCode)
		}
		if killed {
			os.Exit(137) // 128 + SIGKILL(9)
		}
		return nil
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringArrayVarP(&commands, "command", "c", []string{}, "command to run")
	rootCmd.Flags().StringVarP(&shell, "shell", "s", exector.DefaultShell, "shell to use")
	rootCmd.Flags().BoolVarP(&failFast, "fail-fast", "", false, "exit on first error")
}

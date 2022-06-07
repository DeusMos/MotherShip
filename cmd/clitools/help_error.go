package clitools

import (
	"log"

	"github.com/spf13/cobra"
)

// HelpError ...
type HelpError struct {
	msg string
}

// NewHelpError ...
func NewHelpError(msg string) *HelpError {
	return &HelpError{
		msg: msg,
	}
}

// Error ...
func (he *HelpError) Error() string {
	return he.msg
}

// HelpWrap is used to wrap a CLI command and to not include help documentation
// when the issue isn't with user syntax.
func HelpWrap(fun func(*cobra.Command, []string) error) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) (err error) {
		if err = fun(cmd, args); nil != err {
			if _, ok := err.(*HelpError); ok {
				return
			}
			log.Fatalln(err)
		}
		return
	}
}

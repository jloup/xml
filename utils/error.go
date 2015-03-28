package utils

import (
	"fmt"

	"github.com/JLoup/errors"
	"github.com/JLoup/flag"
)

var (
	ErrorFlagCounter flag.Counter = 0
	XMLTokenError                 = flag.InitFlag(&ErrorFlagCounter, "XMLTokenError")
	XMLSyntaxError                = flag.InitFlag(&ErrorFlagCounter, "XMLSyntaxError")
	XMLError                      = flag.Join("XMLError", XMLSyntaxError, XMLTokenError)
	IOError                       = flag.InitFlag(&ErrorFlagCounter, "IOError")
)

type ParserError interface {
	errors.ErrorFlagged
}

type Error struct {
	flag flag.Flag
	msg  string
}

func NewError(f flag.Flag, msg string) ParserError {
	return Error{flag: f, msg: msg}
}

func (s Error) Error() string {
	return fmt.Sprintf("[%s] %s", s.FlagString(), s.msg)
}

func (s Error) ErrorWithCode(f flag.Flag) errors.ErrorFlagged {
	if flag.Intersect(f, s.flag) {
		return s
	}
	return nil
}

func (s Error) Flag() flag.Flag {
	return s.flag
}

func (s Error) FlagString() string {
	return s.flag.String()
}

func (s Error) Msg() string {
	return s.msg
}

type delegatedError struct {
	tokenName      string
	delegatedError ParserError
}

func (d *delegatedError) Error() string {
	return fmt.Sprintf("in '%s':\n%s", d.tokenName, d.delegatedError.Error())
}

func (d *delegatedError) ErrorWithCode(f flag.Flag) errors.ErrorFlagged {
	return d.delegatedError.ErrorWithCode(f)
}

func (d *delegatedError) Flag() flag.Flag {
	return d.delegatedError.Flag()
}

func (d *delegatedError) FlagString() string {
	return d.delegatedError.FlagString()
}

func (d *delegatedError) Msg() string {
	return fmt.Sprintf("in '%s': %s", d.tokenName, d.delegatedError.Error())
}

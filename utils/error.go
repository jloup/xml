package utils

import (
	"fmt"

	"github.com/jloup/utils"
)

var (
	ErrorFlagCounter utils.Counter = 0
	XMLTokenError                  = utils.InitFlag(&ErrorFlagCounter, "XMLTokenError")
	XMLSyntaxError                 = utils.InitFlag(&ErrorFlagCounter, "XMLSyntaxError")
	XMLError                       = utils.Join("XMLError", XMLSyntaxError, XMLTokenError)
	IOError                        = utils.InitFlag(&ErrorFlagCounter, "IOError")
)

type ParserError interface {
	utils.ErrorFlagged
}

type Error struct {
	flag utils.Flag
	msg  string
}

func NewError(f utils.Flag, msg string) ParserError {
	return Error{flag: f, msg: msg}
}

func (s Error) Error() string {
	return fmt.Sprintf("[%s] %s", s.FlagString(), s.msg)
}

func (s Error) ErrorWithCode(f utils.Flag) utils.ErrorFlagged {
	if utils.Intersect(f, s.flag) {
		return s
	}
	return nil
}

func (s Error) Flag() utils.Flag {
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

func (d *delegatedError) ErrorWithCode(f utils.Flag) utils.ErrorFlagged {
	return d.delegatedError.ErrorWithCode(f)
}

func (d *delegatedError) Flag() utils.Flag {
	return d.delegatedError.Flag()
}

func (d *delegatedError) FlagString() string {
	return d.delegatedError.FlagString()
}

func (d *delegatedError) Msg() string {
	return fmt.Sprintf("in '%s': %s", d.tokenName, d.delegatedError.Error())
}

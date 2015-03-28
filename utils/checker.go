package utils

import (
	"github.com/JLoup/errors"
	"github.com/JLoup/flag"
)

const (
	EnableAllError  uint = ((1 << 32) - 1)
	DisableAllError uint = 0

	AllError string = "_"
)

type FlagChecker interface {
	CheckFlag(element string, mask ParserError) bool
	ErrorWithCode(element string, mask ParserError) errors.ErrorFlagged
}

type eUserError struct {
	ElementName string
	flag        flag.Flag
}

type ErrorChecker struct {
	elementErrors    []eUserError
	defaultErrorFlag flag.Flag
}

func NewErrorChecker(defaultError uint) ErrorChecker {
	u := ErrorChecker{}
	if defaultError == EnableAllError {
		u.defaultErrorFlag = flag.NewBigInt(64)
	} else {
		u.defaultErrorFlag = flag.Flag{}
	}
	return u
}

func (u *ErrorChecker) findElement(name string) int {
	for i, el := range u.elementErrors {
		if el.ElementName == name {
			return i
		}
	}
	return -1
}

func (u *ErrorChecker) createNewElement(name string) int {
	if i := u.findElement(name); i != -1 {
		return i
	}

	u.elementErrors = append(u.elementErrors, eUserError{ElementName: name, flag: flag.Flag{}})
	return u.findElement(name)
}

func (u *ErrorChecker) EnableErrorChecking(element string, flags ...flag.Flag) {
	if element == AllError {
		flags1 := append(flags, u.defaultErrorFlag)
		u.defaultErrorFlag = flag.Join("", flags1...)

		for i, _ := range u.elementErrors {
			flags2 := append(flags, u.elementErrors[i].flag)
			u.elementErrors[i].flag = flag.Join("", flags2...)
		}
		return
	}

	i := u.findElement(element)
	if i == -1 {
		i = u.createNewElement(element)

	}
	flags = append(flags, u.elementErrors[i].flag)
	u.elementErrors[i].flag = flag.Join("", flags...)
}

func (u *ErrorChecker) DisableErrorChecking(element string, flags ...flag.Flag) {
	if element == AllError {
		u.defaultErrorFlag = flag.Exclude(u.defaultErrorFlag, flags...)

		for i, _ := range u.elementErrors {
			u.elementErrors[i].flag = flag.Exclude(u.elementErrors[i].flag, flags...)
		}
		return
	}

	i := u.findElement(element)
	if i == -1 {
		i = u.createNewElement(element)
		u.elementErrors[i].flag = u.defaultErrorFlag
	}

	u.elementErrors[i].flag = flag.Exclude(u.elementErrors[i].flag, flags...)
}

func (u *ErrorChecker) Flag(element string) flag.Flag {
	if i := u.findElement(element); i == -1 {
		return u.defaultErrorFlag
	} else {

		return u.elementErrors[i].flag
	}
}

func (u *ErrorChecker) CheckFlag(element string, mask ParserError) bool {

	if index := u.findElement(element); index == -1 {
		return flag.Intersect(u.defaultErrorFlag, mask.Flag())
	} else {
		return flag.Intersect(u.elementErrors[index].flag, mask.Flag())
	}
	return false
}

func (u *ErrorChecker) ErrorWithCode(element string, mask ParserError) errors.ErrorFlagged {

	if index := u.findElement(element); index == -1 {
		return mask.ErrorWithCode(u.defaultErrorFlag)
	} else {
		return mask.ErrorWithCode(u.elementErrors[index].flag)
	}
	return nil
}

package utils

import (
	"github.com/jloup/utils"
)

const (
	EnableAllError  uint = ((1 << 32) - 1)
	DisableAllError uint = 0

	AllError string = "_"
)

type FlagChecker interface {
	CheckFlag(element string, mask ParserError) bool
	ErrorWithCode(element string, mask ParserError) utils.ErrorFlagged
}

type eUserError struct {
	ElementName string
	flag        utils.Flag
}

type ErrorChecker struct {
	elementErrors    []eUserError
	defaultErrorFlag utils.Flag
}

func NewErrorChecker(defaultError uint) ErrorChecker {
	u := ErrorChecker{}
	if defaultError == EnableAllError {
		u.defaultErrorFlag = utils.NewBigInt(64)
	} else {
		u.defaultErrorFlag = utils.Flag{}
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

	u.elementErrors = append(u.elementErrors, eUserError{ElementName: name, flag: utils.Flag{}})
	return u.findElement(name)
}

func (u *ErrorChecker) EnableErrorChecking(element string, flags ...utils.Flag) {
	if element == AllError {
		flags1 := append(flags, u.defaultErrorFlag)
		u.defaultErrorFlag = utils.Join("", flags1...)

		for i, _ := range u.elementErrors {
			flags2 := append(flags, u.elementErrors[i].flag)
			u.elementErrors[i].flag = utils.Join("", flags2...)
		}
		return
	}

	i := u.findElement(element)
	if i == -1 {
		i = u.createNewElement(element)

	}
	flags = append(flags, u.elementErrors[i].flag)
	u.elementErrors[i].flag = utils.Join("", flags...)
}

func (u *ErrorChecker) DisableErrorChecking(element string, flags ...utils.Flag) {
	if element == AllError {
		u.defaultErrorFlag = utils.Exclude(u.defaultErrorFlag, flags...)

		for i, _ := range u.elementErrors {
			u.elementErrors[i].flag = utils.Exclude(u.elementErrors[i].flag, flags...)
		}
		return
	}

	i := u.findElement(element)
	if i == -1 {
		i = u.createNewElement(element)
		u.elementErrors[i].flag = u.defaultErrorFlag
	}

	u.elementErrors[i].flag = utils.Exclude(u.elementErrors[i].flag, flags...)
}

func (u *ErrorChecker) Flag(element string) utils.Flag {
	if i := u.findElement(element); i == -1 {
		return u.defaultErrorFlag
	} else {

		return u.elementErrors[i].flag
	}
}

func (u *ErrorChecker) CheckFlag(element string, mask ParserError) bool {

	if index := u.findElement(element); index == -1 {
		return utils.Intersect(u.defaultErrorFlag, mask.Flag())
	} else {
		return utils.Intersect(u.elementErrors[index].flag, mask.Flag())
	}
	return false
}

func (u *ErrorChecker) ErrorWithCode(element string, mask ParserError) utils.ErrorFlagged {

	if index := u.findElement(element); index == -1 {
		return mask.ErrorWithCode(u.defaultErrorFlag)
	} else {
		return mask.ErrorWithCode(u.elementErrors[index].flag)
	}
	return nil
}

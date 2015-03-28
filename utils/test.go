package utils

import (
	"fmt"
	"strings"
)

type VisitorValidator func(actual Visitor, expected Visitor) error
type VisitorConstructor func() Visitor

type TestVisitor struct {
	XML                string
	ExpectedError      ParserError
	ExpectedVisitor    Visitor
	VisitorConstructor VisitorConstructor
	Validator          VisitorValidator
	CustomError        *ErrorChecker
}

func (t *TestVisitor) CheckTestCaseErrorsEnabled() error {
	if t.CustomError == nil {
		u := NewErrorChecker(EnableAllError)
		t.CustomError = &u
	}

	v := t.VisitorConstructor()

	err := Walk(strings.NewReader(t.XML), v, t.CustomError)

	if err != nil && t.ExpectedError != nil {
		if !err.Flag().Cmp(t.ExpectedError.Flag()) {
			return fmt.Errorf("[Error flag returned is not valid] '%v' vs '%v' (expected)", err.FlagString(), t.ExpectedError.FlagString())
		}
		return nil
	}

	if err != nil {
		return fmt.Errorf("[Unexpected error] %v", err)
	}

	if t.ExpectedError != nil {
		return fmt.Errorf("[No error detected] expecting '%v'", t.ExpectedError.FlagString())
	}

	if err := t.Validator(v, t.ExpectedVisitor); err != nil {
		return fmt.Errorf("[XML wrong data] %s", err.Error())
	}

	return nil
}

func (t *TestVisitor) CheckTestCaseErrorsDisabled() error {
	errorCheck := NewErrorChecker(DisableAllError)
	v := t.VisitorConstructor()

	err := Walk(strings.NewReader(t.XML), v, &errorCheck)

	if err != nil {
		return fmt.Errorf("[Disabled Errors - Unexpected error] %v", err)
	}

	if err := t.Validator(v, t.ExpectedVisitor); err != nil {
		return fmt.Errorf("[Disabled Errors - XML wrong data] %s", err.Error())
	}

	return nil
}

func (t *TestVisitor) CheckTestCase() error {
	if err := t.CheckTestCaseErrorsEnabled(); err != nil {
		return err
	}

	if t.ExpectedError != nil {
		if err := t.CheckTestCaseErrorsDisabled(); err != nil {
			return err
		}
	}
	return nil
}

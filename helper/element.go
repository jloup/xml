package helper

type ElementValidator func(name, s string) ParserError

type Element struct {
	Name      string
	Value     string
	Validator ElementValidator

	Occurence *Occurence
}

func NewElement(name, value string, validator ElementValidator) Element {
	return Element{Name: name, Value: value, Validator: validator, Occurence: nil}
}

func (e *Element) SetOccurence(occ *Occurence) {
	e.Occurence = occ
}

func (e *Element) IncOccurence() {
	if e.Occurence != nil {
		e.Occurence.Inc()
	}
}

func (e Element) Validate() ParserError {

	if e.Occurence == nil || (e.Occurence != nil && e.Occurence.NbOccurences > 0) {
		if err := e.Validator(e.Name, e.Value); err != nil {
			return err
		}
	}

	if e.Occurence != nil {
		if err := e.Occurence.Validate(); err != nil {
			return err
		}
	}
	return nil
}

func (e Element) String() string {
	return e.Value
}

func (e *Element) Reset() {
	if e.Occurence != nil {
		e.Occurence.Reset()
	}
}

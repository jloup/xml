package helper

type OccurenceValidator func(o *Occurence) ParserError

type Occurence struct {
	Name         string
	NbOccurences int
	Validator    OccurenceValidator
}

func NewOccurence(name string, validator OccurenceValidator) *Occurence {
	return &Occurence{Name: name, NbOccurences: 0, Validator: validator}
}

func (o *Occurence) Inc() {
	o.NbOccurences += 1
}

func (o *Occurence) Reset() {
	o.NbOccurences = 0
}

func (o *Occurence) Validate() ParserError {
	return o.Validator(o)
}

type OccurenceCollection struct {
	Occurences []*Occurence
}

func NewOccurenceCollection(occs ...*Occurence) OccurenceCollection {
	o := OccurenceCollection{Occurences: occs}

	return o
}

func (o *OccurenceCollection) AddOccurence(occ *Occurence) {
	o.Occurences = append(o.Occurences, occ)
}

func (o *OccurenceCollection) Inc(name string) {
	for i, occ := range o.Occurences {
		if occ.Name == name {
			o.Occurences[i].Inc()
			return
		}
	}
}

func (o *OccurenceCollection) Count(name string) int {
	for _, occ := range o.Occurences {
		if occ.Name == name {
			return occ.NbOccurences
		}
	}
	return 0
}

func (o *OccurenceCollection) Validate(name string) ParserError {
	for i, occ := range o.Occurences {
		if occ.Name == name {
			return o.Occurences[i].Validate()
		}
	}
	return nil
}

func (o *OccurenceCollection) Reset() {
	for i, _ := range o.Occurences {
		o.Occurences[i].Reset()
	}
}

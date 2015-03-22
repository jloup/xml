package extension

import (
	"encoding/xml"

	"github.com/JLoup/errors"
	"github.com/JLoup/xml/helper"
)

type VisitorExtension struct {
	name       string
	Manager    Manager
	Repository Repository
	Store      Store
}

func InitExtension(name string, manager Manager) VisitorExtension {
	v := VisitorExtension{name: name, Manager: manager}
	v.Repository = manager.GetRepo(name)

	v.Store.Occ = v.Repository.Occ
	v.Store.Occ.Reset()

	return v
}

func (v *VisitorExtension) ProcessElement(el helper.StartElement, parent helper.Visitor) (helper.Visitor, helper.ParserError) {
	if constructor := v.Repository.GetElement(el.Name); constructor == nil {
		return nil, nil

	} else {
		ext := constructor()
		ext.SetParent(parent)

		nextV, err := ext.ProcessStartElement(el)

		v.Store.Add(ext.Name(), ext)
		return nextV, err
	}

}

func (v *VisitorExtension) ProcessAttr(attr xml.Attr, parent helper.Visitor) {
	if constructor := v.Repository.GetAttr(attr.Name); constructor != nil {
		ext := constructor()
		ext.SetParent(parent)
		ext.Set(attr.Value)
		v.Store.Add(ext.Name(), ext)
	}
}

func (v *VisitorExtension) Validate(errorAgg *errors.ErrorAggregator) {
	helper.ValidateOccurenceCollection(v.name, errorAgg, v.Store.Occ)
	v.Store.Validate(errorAgg)
}

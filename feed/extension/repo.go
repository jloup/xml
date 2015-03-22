package extension

import (
	"encoding/xml"
	"errors"
	"fmt"

	"github.com/JLoup/flag"

	"github.com/JLoup/xml/helper"
)

func xmlNameToString(name xml.Name) string {
	return fmt.Sprintf("%s:%s", name.Space, name.Local)
}

type eElement struct {
	name        xml.Name
	constructor ElementConstructor
}

type eAttr struct {
	name        xml.Name
	constructor AttrConstructor
}

// base element for Manager
type Repository struct {
	name     string
	Occ      helper.OccurenceCollection
	elements []eElement
	attrs    []eAttr
}

func (r *Repository) findAttr(name xml.Name) int {
	for i, attr := range r.attrs {
		if attr.name == name {
			return i
		}
	}
	return -1
}

func (r *Repository) AddAttr(name xml.Name, constructor AttrConstructor, attrDuplicatedFlag flag.Flag) error {

	if index := r.findAttr(name); index != -1 {
		return errors.New("an extension of this name already exists")
	}

	r.attrs = append(r.attrs, eAttr{name: name, constructor: constructor})
	r.Occ.AddOccurence(helper.NewOccurence(xmlNameToString(name), helper.UniqueValidator(attrDuplicatedFlag)))

	return nil
}

func (r *Repository) GetAttr(name xml.Name) AttrConstructor {
	if index := r.findAttr(name); index == -1 {
		return nil
	} else {
		return r.attrs[index].constructor
	}
}

func (r *Repository) findElement(name xml.Name) int {
	for i, el := range r.elements {
		if el.name == name {
			return i
		}
	}
	return -1
}

func (r *Repository) AddElement(name xml.Name, constructor ElementConstructor, occValidator helper.OccurenceValidator) error {

	if index := r.findElement(name); index != -1 {
		return errors.New("an extension of this name already exists")
	}

	r.elements = append(r.elements, eElement{name: name, constructor: constructor})
	r.Occ.AddOccurence(helper.NewOccurence(xmlNameToString(name), occValidator))

	return nil
}

func (r *Repository) GetElement(name xml.Name) ElementConstructor {
	if index := r.findElement(name); index == -1 {
		return nil
	} else {
		return r.elements[index].constructor
	}
}

/*
*
*
*
*
*
 */

type Manager struct {
	tags []Repository
}

func (m *Manager) findAndCreate(name string) int {
	length := -1
	for i, tag := range m.tags {
		if tag.name == name {
			return i
		}
		length = i
	}
	m.tags = append(m.tags, Repository{name: name})
	return length + 1
}

func (m *Manager) find(name string) int {
	for i, tag := range m.tags {
		if tag.name == name {
			return i
		}
	}
	return -1
}

func (m *Manager) AddAttrExtension(tagName string, name xml.Name, constructor AttrConstructor, attrDuplicatedFlag flag.Flag) error {
	index := m.findAndCreate(tagName)
	return m.tags[index].AddAttr(name, constructor, attrDuplicatedFlag)
}

func (m *Manager) AddElementExtension(tagName string, name xml.Name, constructor ElementConstructor, occValidator helper.OccurenceValidator) error {
	index := m.findAndCreate(tagName)
	return m.tags[index].AddElement(name, constructor, occValidator)
}

func (m *Manager) GetRepo(tagName string) Repository {
	if index := m.find(tagName); index == -1 {
		return Repository{}
	} else {
		return m.tags[index]
	}
}

package extension

import (
	"encoding/xml"

	"github.com/jloup/errors"
	"github.com/jloup/xml/utils"
)

type storeInterface interface {
	String() string
	Validate() utils.ParserError
}

type eStore struct {
	name       xml.Name
	extensions []storeInterface
}

type Store struct {
	stores []eStore
	Occ    utils.OccurenceCollection
}

func (s *Store) find(name xml.Name) int {
	for i, store := range s.stores {
		if store.name == name {
			return i
		}
	}

	return -1
}

func (s *Store) findAndCreate(name xml.Name) int {
	length := -1
	for i, store := range s.stores {
		if store.name == name {
			return i
		}
		length = i
	}

	s.stores = append(s.stores, eStore{name: name})
	return length + 1
}

func (s *Store) Add(name xml.Name, el storeInterface) {
	index := s.findAndCreate(name)

	s.stores[index].extensions = append(s.stores[index].extensions, el)

	s.Occ.Inc(xmlNameToString(name))
}

func (s *Store) Get(name xml.Name) (string, bool) {
	if i := s.find(name); i != -1 {
		return s.stores[i].extensions[0].String(), true
	}

	return "", false
}

func (s *Store) GetItf(name xml.Name) (interface{}, bool) {
	if i := s.find(name); i != -1 {
		return s.stores[i].extensions[0], true
	}

	return "", false
}

func (s *Store) GetCollection(name xml.Name) ([]storeInterface, bool) {
	if i := s.find(name); i != -1 {
		return s.stores[i].extensions, true
	}

	return nil, false

}

func (s *Store) Validate(errorAgg *errors.ErrorAggregator) {
	for _, store := range s.stores {
		for _, ext := range store.extensions {
			if _, ok := ext.(Attr); ok {
				if err := ext.Validate(); err != nil {
					errorAgg.NewError(err)
				}
			}
		}
	}
}

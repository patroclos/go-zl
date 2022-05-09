package zettel

import "fmt"

type Zetteler interface {
	Zettel(id string) (Z, error)
}

type Putter interface {
	Put(Z) error
}

type Resolver interface {
	// Find a non-empty []Zettel matching the given query
	Resolve(query string) ([]Z, error)
}

type Iterator interface {
	Next() bool
	Zet() Z
}

type Storage interface {
	Zetteler
	Putter
	Resolver
	Iter() Iterator
	Remove(Z) error
}

type ZettelerIter interface {
	Zetteler
	Iter() Iterator
}

func All(st ZettelerIter) []Z {
	iter := st.Iter()
	var zets []Z = nil
	if iter.Next() {
		zets = make([]Z, 1, 512)
		zets[0] = iter.Zet()
	}
	for iter.Next() {
		zets = append(zets, iter.Zet())
	}
	return zets
}

type slice struct {
	zets  []Z
	idmap map[string]Z
}

func (s *slice) Zettel(id string) (Z, error) {
	z, ok := s.idmap[id]
	if !ok {
		return nil, fmt.Errorf("id not found")
	}
	return z, nil
}

func (s *slice) Iter() Iterator {
	return &iter{s, -1}
}

type iter struct {
	s   *slice
	pos int
}

func (i *iter) Next() bool {
	i.pos++
	return len(i.s.zets) > i.pos
}

func (i *iter) Zet() Z {
	return i.s.zets[i.pos]
}

func Slice(in []Z) ZettelerIter {
	idmap := make(map[string]Z)
	for _, z := range in {
		idmap[z.Id()] = z
	}

	return &slice{in, idmap}
}

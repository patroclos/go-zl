package zettel

type Zetteler interface {
	Zettel(id string) (Zettel, error)
}

type Putter interface {
	Put(Zettel) error
}

type Resolver interface {
	// Find a non-empty []Zettel matching the given query
	Resolve(query string) ([]Zettel, error)
}

type Iterator interface {
	Next() bool
	Zet() Zettel
}

type Storage interface {
	Zetteler
	Putter
	Resolver
	Iter() Iterator
	Remove(Zettel) error
}

type ZettelerIter interface {
	Zetteler
	Iter() Iterator
}

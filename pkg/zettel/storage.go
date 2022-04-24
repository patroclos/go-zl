package zettel

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

package zettel

type Zetteler interface {
	Zettel(id string) (Zettel, error)
}

type Putter interface {
	Put(Zettel) error
}

type Storage interface {
	Zetteler
	Putter
}

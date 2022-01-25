package zettel

import (
	"math/rand"
	"time"
)

const (
	idCharset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	idLen     = 12
)

var rng *rand.Rand

func generateId(rng *rand.Rand) string {
	buf := [idLen]byte{}
	for i := 0; i < idLen; i++ {
		buf[i] = idCharset[rng.Intn(len(idCharset))]
	}
	return string(buf[:])
}

func plainGenerateId() string {
	if rng == nil {
		rng = rand.New(rand.NewSource(time.Now().UnixNano()))
	}
	return generateId(rng)
}

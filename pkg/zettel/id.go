package zettel

import (
	"fmt"
	"math/rand"
	"time"
)

const (
	idCharset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	idLen     = 4
)

var rng *rand.Rand

func MakeIdAt(rng *rand.Rand, t time.Time) string {
	ts := t.Format("060102")
	buf := [idLen]byte{}
	for i := 0; i < idLen; i++ {
		buf[i] = idCharset[rng.Intn(len(idCharset))]
	}
	return fmt.Sprintf("%s-%s", ts, buf[:])
}

func MakeId() string {
	if rng == nil {
		rng = rand.New(rand.NewSource(time.Now().UnixNano()))
	}
	return MakeIdAt(rng, time.Now())
}

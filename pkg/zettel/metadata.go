package zettel

import "time"

type Labels map[string]string

type MetaInfo struct {
	Labels Labels
	Link   *LinkInfo `yaml:"link"`
	CreationTimestamp time.Time
}

type LinkInfo struct {
	A   string `yaml:"from"`
	B   string `yaml:"to"`
	Ctx []string `yaml:"context"`
}

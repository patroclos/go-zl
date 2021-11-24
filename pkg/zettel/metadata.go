package zettel

import "time"

type MetaInfo struct {
	Labels map[string]string
	Link   *LinkInfo `yaml:"link"`
	CreationTimestamp time.Time
}

type LinkInfo struct {
	A   string `yaml:"from"`
	B   string `yaml:"to"`
	Ctx []string `yaml:"context"`
}

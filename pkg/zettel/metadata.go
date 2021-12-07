package zettel

import "time"

type Labels map[string]string

type MetaInfo struct {
	Labels     Labels    `yaml:"labels"`
	Link       *LinkInfo `yaml:"link,omitempty"`
	CreateTime time.Time `yaml:"creationTimestamp"`
}

type LinkInfo struct {
	A   string   `yaml:"from"`
	B   string   `yaml:"to"`
	Ctx []string `yaml:"context"`
}

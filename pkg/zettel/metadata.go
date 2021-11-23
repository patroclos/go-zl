package zettel

type MetaInfo struct {
	Labels map[string]string
	Link   *LinkInfo
}

type LinkInfo struct {
	A   string `yaml:"from"`
	B   string `yaml:"to"`
	Ctx []string `yaml:"context"`
}

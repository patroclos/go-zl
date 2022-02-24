package config

import (
	"io"
	"io/ioutil"
	"os"

	"github.com/kirsle/configdir"
	"gopkg.in/yaml.v2"
	"jensch.works/zl/pkg/zettel"
)

type Config struct {
	Profiles      []Profile
	ActiveProfile string
}

type Profile struct {
	Name      string
	Directory string
	Tolerate  []string
	Labels    []zettel.Labelspec
}

var def *Config

func Default() (*Config, error) {
	if def != nil {
		return def, nil
	}
	f, err := os.Open(configdir.LocalConfig("zl"))
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return Read(f)
	// def, err := Read(
}

func SetDefault(c *Config) error {
	f, err := os.Open(configdir.LocalConfig("zl"))
	if err != nil {
		return err
	}
	defer f.Close()
	return Write(c, f)
}

func Read(r io.Reader) (*Config, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	c := new(Config)
	if err := yaml.Unmarshal(b, c); err != nil {
		return nil, err
	}
	return c, nil
}

func Write(c *Config, w io.Writer) error {
	bytes, err := yaml.Marshal(c)
	if err != nil {
		return err
	}
	_, err = w.Write(bytes)
	return err
}

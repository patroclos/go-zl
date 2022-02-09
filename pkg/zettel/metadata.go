package zettel

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"time"

	"gopkg.in/yaml.v2"
)

type Labels map[string]string

var (
	ErrorFormat      = errors.New("format error")
	ErrorInvalidLink = errors.New("invalid link")
)

type MetaInfo struct {
	Labels     Labels    `yaml:"labels,omitempty"`
	Link       *LinkInfo `yaml:"link,omitempty"`
	CreateTime time.Time `yaml:"creationTimestamp,omitempty"`
}

func (i *MetaInfo) copy(from MetaInfo) {
	i.CreateTime = from.CreateTime
	i.Labels = make(map[string]string)
	for k, v := range from.Labels {
		i.Labels[k] = v
	}

	if from.Link != nil {
		i.Link = new(LinkInfo)
		i.Link.A, i.Link.B, i.Link.Ctx = from.Link.A, from.Link.B, from.Link.Ctx
	}
}

func (i *MetaInfo) Equal(o *MetaInfo) bool {
	if len(i.Labels) != len(o.Labels) {
		return false
	}

	for k, v := range i.Labels {
		if v2, ok := o.Labels[k]; !ok || v != v2 {
			return false
		}
	}

	if i.Link == nil && o.Link != nil {
		return false
	}
	if i.Link != nil {
		if o.Link == nil {
			return false
		}

		if i.Link.A != o.Link.A {
			return false
		}
		if i.Link.B != o.Link.B {
			return false
		}
		if len(i.Link.Ctx) != len(o.Link.Ctx) {
			return false
		}
		for idx := range i.Link.Ctx {
			if i.Link.Ctx[idx] != o.Link.Ctx[idx] {
				return false
			}
		}
	}

	return i.CreateTime == o.CreateTime
}

func (i *MetaInfo) NewTemp() (*os.File, error) {
	tmp, err := os.CreateTemp("", "zlmedit*.yaml")
	if err != nil {
		return nil, err
	}

	enc := yaml.NewEncoder(tmp)
	defer enc.Close()
	if err := enc.Encode(i); err != nil {
		os.Remove(tmp.Name())
		return nil, err
	}
	return tmp, nil
}

type LinkInfo struct {
	A   string   `yaml:"from"`
	B   string   `yaml:"to"`
	Ctx []string `yaml:"context"`
}

type metaDto struct {
	Labels     Labels                 `yaml:"labels"`
	Link       map[string]interface{} `yaml:"link,omitempty"`
	CreateTime time.Time              `yaml:"creationTimestamp"`
}

// Must either return (non-nil, nil) or (nil, non-nil)
func ParseMeta(r io.Reader) (*MetaInfo, error) {
	in, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	dto := new(metaDto)
	err = yaml.Unmarshal(in, dto)
	if err != nil {
		return nil, err
	}

	info := &MetaInfo{
		Labels:     dto.Labels,
		CreateTime: dto.CreateTime,
	}

	if err := _readLink(dto, info); err != nil {
		return nil, err
	}

	return info, nil
}

func _readLink(dto *metaDto, info *MetaInfo) error {
	lnk := &LinkInfo{}
	if dto.Link == nil {
		return nil
	}
	from, okF := dto.Link["from"]
	to, okT := dto.Link["to"]
	if !okF {
		return fmt.Errorf(`%w: dto.Link["from"] not found`, ErrorFormat)
	}
	if !okT {
		return fmt.Errorf(`%w: dto.Link["to"] not found`, ErrorFormat)
	}

	if str, ok := from.(string); ok {
		lnk.A = str
	}
	if m, ok := from.(map[interface{}]interface{}); ok {
		zet, ok := m["zet"].(string)
		if !ok {
			return fmt.Errorf(`%w: link["from"]["zet"] not a string`, ErrorFormat)
		}
		lnk.A = zet
	}

	if str, ok := to.(string); ok {
		lnk.B = str
	}

	if m, ok := to.(map[interface{}]interface{}); ok {
		zet, ok := m["zet"].(string)
		if !ok {
			return fmt.Errorf(`%w: link["to"]["zet"] not a string: %#v`, ErrorFormat, m["zet"])
		}
		lnk.B = zet
	}

	if ctxa, ok := dto.Link["context"].([]interface{}); ok {
		for _, c := range ctxa {
			if str, ok := c.(string); ok {
				lnk.Ctx = append(lnk.Ctx, str)
			}
		}
	}

	info.Link = lnk

	return nil
}

func validateLink(l *LinkInfo) error {
	if l == nil {
		return fmt.Errorf("(%w: nil)", ErrorInvalidLink)
	}

	if l.A == "" {
		return fmt.Errorf("(%w: A empty)", ErrorInvalidLink)
	}

	if l.B == "" {
		return fmt.Errorf("(%w: B empty)", ErrorInvalidLink)
	}

	return nil
}

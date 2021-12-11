package zettel

import (
	"errors"
	"fmt"
	"time"

	"gopkg.in/yaml.v2"
)

type Labels map[string]string

var (
	ErrorFormat      = errors.New("format error")
	ErrorInvalidLink = errors.New("invalid link")
)

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

// TODO: pros and cons of sugar representations? same key vs pluralization
//	as far as im concerned pluralization is stupid, bc we need the "listen up,
//	fully detailed" version to fit anywhere the shortform goes to encode k well
type metaDto struct {
	Labels     Labels                 `yaml:"labels"`
	Link       map[string]interface{} `yaml:"link,omitempty"`
	CreateTime time.Time              `yaml:"crecreationTimestampationTimestamp"`
}

func ParseMeta(in []byte) (*MetaInfo, error) {
	dto := new(metaDto)
	err := yaml.Unmarshal(in, dto)
	if err != nil {
		return nil, err
	}

	info := &MetaInfo{
		Labels:     dto.Labels,
		CreateTime: dto.CreateTime,
	}

	lnk := &LinkInfo{}

	if err := _readLink(dto, lnk); err != nil {
		return nil, err
	} else {
		info.Link = lnk
	}

	return info, nil
}

func _readLink(dto *metaDto, lnk *LinkInfo) error {
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

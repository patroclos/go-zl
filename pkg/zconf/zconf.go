package zconf

import (
	"os"
	"strings"
)

type Cfg struct {
	Elems *ElemSpec
}

type ElemSpec struct {
	Code *ElemCodeSpec
}

type ElemCodeSpec struct {
	Filters map[string]Filter
}

type Filter struct {
	Dir string
	Cmd string
}

func FromEnv() (*Cfg, error) {
	cfg := &Cfg{
		Elems: &ElemSpec{
			Code: &ElemCodeSpec{
				Filters: map[string]Filter{},
			},
		},
	}

	for _, env := range os.Environ() {
		kv := strings.SplitN(env, "=", 2)
		if strings.HasPrefix(kv[0], "ZLSRV_FILTER__") {
			param := strings.TrimPrefix(kv[0], "ZLSRV_FILTER__")
			cfg.Elems.Code.Filters[param] = Filter{Cmd: kv[1]}
		}
		switch kv[0] {
		case "ZLSRV_FILTER__<lang>":
		}
	}
	return cfg, nil
}

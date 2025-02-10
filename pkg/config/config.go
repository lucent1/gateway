package config

import (
	"github.com/doublecloud/transfer/pkg/abstract/model"
	"gopkg.in/yaml.v3"

	"github.com/doublecloud/transfer/pkg/abstract"
)

type Gateway struct {
	Type   abstract.ProviderType
	Params any
}

func (g *Gateway) Endpoint() (model.Source, error) {
	return model.NewSource(g.Type, g.ParamRaw())
}

func (g *Gateway) ParamRaw() string {
	switch p := g.Params.(type) {
	case []byte:
		return string(p)
	case string:
		return p
	default:
		data, _ := yaml.Marshal(p)
		return string(data)
	}
}

func FromYaml(raw []byte) (*Gateway, error) {
	var gw Gateway
	err := yaml.Unmarshal(raw, &gw)
	if err != nil {
		return nil, nil
	}
	return &gw, nil
}

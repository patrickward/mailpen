package mailpen

import (
	"context"
)

type Module struct {
	config   *Config
	mailpen  *Mailpen
	provider Provider
}

func NewModule(provider Provider, config *Config) *Module {
	return &Module{
		config:   config,
		provider: provider,
	}
}

func (m *Module) ID() string {
	return "hop.mail"
}

func (m *Module) Init() error {
	mp, err := New(m.provider, m.config)
	if err != nil {
		return err
	}
	m.mailpen = mp
	return nil
}

func (m *Module) Start(_ context.Context) error {
	return nil
}

func (m *Module) Stop(_ context.Context) error {
	return nil
}

func (m *Module) Mailpen() *Mailpen {
	return m.mailpen
}

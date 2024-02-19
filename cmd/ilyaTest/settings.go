package main

type (
	Settings struct {
		Database DatabaseSettings `yaml:"database"`
	}

	DatabaseSettings struct {
		Arguments string `yaml:"arguments"`
		Type      string `yaml:"type"`
	}
)

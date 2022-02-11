package hhot

import "flag"

type Config interface {
	BasePath() string
}

type config struct {
	basePath string
}

func NewConfig() Config {
	c := config{}

	flag.StringVar(&c.basePath, "base-path", "/", "URL base path")

	flag.Parse()

	return &c
}

func (c *config) BasePath() string {
	return c.basePath
}

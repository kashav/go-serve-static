package serve_static

import "errors"

const (
	Debug bool = true
	Port  int  = 8080
)

type Config struct {
	ID                string
	Repo              string
	Build             string
	Serve             string
	SparseCheckoutDir string `yaml:"sparse_directory"`
}

func (c *Config) Check() error {
	if c.ID == "" {
		return errors.New("expected non-empty `id` field")
	}
	if c.Repo == "" {
		return errors.New("expected non-empty `repo` field")
	}
	if c.Build == "" {
		return errors.New("expected non-empty `build` command")
	}
	if c.Serve == "" {
		return errors.New("expected non-empty `serve` directory")
	}
	return nil
}

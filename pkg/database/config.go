package database

import (
	"net/url"
)

type Config struct {
	ConnType string `env:"DB_CONN_TYPE" json:",omitempty"`
	ConnName string `env:"DB_CONN_NAME" json:",omitempty"`
	Name     string `env:"DB_NAME" json:",omitempty"`
	User     string `env:"DB_USER" json:",omitempty"`
	Host     string `env:"DB_HOST, default=localhost" json:",omitempty"`
	Port     string `env:"DB_PORT" json:",omitempty"`
	Password string `env:"DB_PASSWORD" json:"-"` // ignored by zap's JSON formatter
}

func (c *Config) DatabaseConfig() *Config {
	return c
}

func (c *Config) ConnectionURL() string {
	if c == nil {
		return ""
	}

	host := c.Host
	if v := c.Port; v != "" {
		host = host + ":" + v
	}

	u := &url.URL{
		Scheme: c.ConnType,
		Host:   host,
		Path:   c.Name,
	}

	if c.User != "" || c.Password != "" {
		u.User = url.UserPassword(c.User, c.Password)
	}

	return u.String()
}

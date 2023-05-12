package repository

import (
	"fmt"
	"strings"
)

// Protocols.
const (
	HTTPSProtocol = "https"
	SSHProtocol   = "ssh"
)

type Config struct {
	Protocol            string        `json:"protocol"`
	Host                string        `json:"host"`
	RepoURLPath         string        `json:"repo_url_path"`
	Remote              string        `json:"remote"`
	Branch              string        `json:"branch"`
	RelativeContentPath string        `json:"relative_content_path"`
	SSHAuth             SSHAuthConfig `json:"ssh_auth"`
}

type SSHAuthConfig struct {
	User           string `json:"user"`
	PrivateKeyPath string `json:"private_key_path"`
}

func (c *Config) URL() string {
	userPart := ""
	if c.Protocol == SSHProtocol && c.SSHAuth.User != "" {
		userPart = c.SSHAuth.User + "@"
	}

	return fmt.Sprintf(
		"%s://%s%s/%s",
		c.Protocol,
		userPart,
		c.Host,
		strings.TrimPrefix(c.RepoURLPath, "/"),
	)
}

func (c *Config) Validate() error {
	if c.Protocol != HTTPSProtocol && c.Protocol != SSHProtocol {
		return ErrInvalidProtocol
	}
	return nil
}

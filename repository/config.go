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
	RepoPath            string        `json:"repo_path"`
	Remote              string        `json:"remote"`
	Branch              string        `json:"branch"`
	RelativeContentPath string        `json:"relative_content_path"`
	SSHAuth             SSHAuthConfig `json:"ssh_auth"`
}

type SSHAuthConfig struct {
	// SSH username.
	User string `json:"user"`

	// Absolute path to the private key file for SSH authentication of the User.
	PrivateKeyPath string `json:"private_key_path"`

	// Password protecting the private key.
	PrivateKeyPassword string `json:"private_key_password"`
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
		strings.TrimPrefix(c.RepoPath, "/"),
	)
}

func (c *Config) Validate() error {
	if c.Protocol != HTTPSProtocol && c.Protocol != SSHProtocol {
		return ErrInvalidProtocol
	}

	if c.Protocol == SSHProtocol && (c.SSHAuth.User == "" || c.SSHAuth.PrivateKeyPath == "") {
		return ErrAuthParamMissing
	}

	return nil
}

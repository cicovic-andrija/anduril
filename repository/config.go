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
	// SSH username.
	User string `json:"user"`

	// Absolute path to the private key file for SSH authentication of the User.
	// Note: Currently, only non-encrypted PEM keys are supported, so ensure that the
	// key has read-only level of access in the repository.
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

	if c.Protocol == SSHProtocol && (c.SSHAuth.User == "" || c.SSHAuth.PrivateKeyPath == "") {
		return ErrAuthParamMissing
	}

	return nil
}

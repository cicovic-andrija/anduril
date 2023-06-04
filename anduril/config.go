package anduril

import (
	"fmt"
	"time"

	"github.com/cicovic-andrija/anduril/repository"
	"github.com/cicovic-andrija/libgo/https"
)

type Config struct {
	HTTPS      https.Config      `json:"https"`
	Repository repository.Config `json:"repository"`
	Settings   Settings          `json:"settings"`
}

type Settings struct {
	PublishPrivateArticles    bool          `json:"publish_private_articles"`
	PublishPersonalArticles   bool          `json:"publish_personal_articles"`
	RepositorySyncPeriod      string        `json:"repository_sync_period"`
	RepositorySyncPeriodDur   time.Duration `json:"-"`
	StaleFileCleanupPeriod    string        `json:"stale_file_cleanup_period"`
	StaleFileCleanupPeriodDur time.Duration `json:"-"`
}

func (s *Settings) Validate() error {
	dur, err := time.ParseDuration(s.RepositorySyncPeriod)
	if err != nil {
		return fmt.Errorf("repository sync period: %v", err)
	}
	s.RepositorySyncPeriodDur = dur

	dur, err = time.ParseDuration(s.StaleFileCleanupPeriod)
	if err != nil {
		return fmt.Errorf("stale file cleanup period: %v", err)
	}
	s.StaleFileCleanupPeriodDur = dur

	return nil
}

package storage

import (
	"errors"
	"log"
	"strings"

	"github.com/aceworksdev/go-stripe-eboekhouden/internal"
	"github.com/aceworksdev/go-stripe-eboekhouden/internal/config"
	"github.com/aceworksdev/go-stripe-eboekhouden/internal/server/domain/storage"
	"github.com/aceworksdev/go-stripe-eboekhouden/internal/storage/mysql"
)

func New(s *storage.Storage, dbConfig *config.AllDB, sa *internal.ServicesAvailable) error {
	if err := NewDatabase(dbConfig, s, sa); err != nil {
		return err
	}
	return nil
}

func NewDatabase(c *config.AllDB, s *storage.Storage, sa *internal.ServicesAvailable) (err error) {
	log.Printf("Connecting to %s", c.Current)
	switch strings.ToLower(c.Current) {
	case "postgres":
		return errors.New("not implemented")
	case "mysql":
		if err := mysql.New(c.MySQL, s, sa); err != nil {
			return err
		}
	default:
		return internal.ModeUnknown("storage", c.Current, "mysql", "postgres")
	}
	log.Printf("Connected to %s", c.Current)
	return nil
}

package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql" // The driver used for MySQL
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql" // The driver used for MySQL

	"github.com/aceworksdev/go-stripe-eboekhouden/internal"
	"github.com/aceworksdev/go-stripe-eboekhouden/internal/config"
	"github.com/aceworksdev/go-stripe-eboekhouden/internal/server/domain/storage"
	"github.com/aceworksdev/go-stripe-eboekhouden/internal/storage/mysql/migrations"
	queries "github.com/aceworksdev/go-stripe-eboekhouden/internal/storage/mysql/queries/generated"
)

var registerDriverOnce sync.Once

const (
	IntervalDBCheck = 10 * time.Minute
)

func New(c *config.StorageMySQL, s *storage.Storage, sa *internal.ServicesAvailable) (err error) {
	dbs := &DBs{}
	errs := make(chan error, 10)
	NewConnection(c, dbs, sa, errs)
	querier, err := queries.Prepare(context.TODO(), dbs.Master)
	if err != nil {
		return err
	}

	if s.Customer, err = NewCustomer(dbs, querier); err != nil {
		return err
	}

	return err
}

func NewConnection(config *config.StorageMySQL, dbs *DBs, servicesAvailable *internal.ServicesAvailable, errs chan error) {
	go func() {
		l := log.New(os.Stdout, "[DATABASE][ERROR]", 0)
		for {
			err := <-errs
			l.Printf("%v", err)
		}
	}()

	mySQLInfo := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?clientFoundRows=true", config.User, config.Password, config.Hostname, config.Port, config.Database)

	if db, err := ConnAndPingMySQL(mySQLInfo, config); err != nil {
		errs <- err
	} else {
		dbs.Master = db
	}

	go func() {
		//nolint, must be infinite loop
		for {
			select {
			case <-time.After(IntervalDBCheck):
				ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(time.Second))
				defer cancel()

				if servicesAvailable.DB {
					if _, err := dbs.Master.ExecContext(ctx, "SELECT 1 + 1;"); err != nil {
						log.Println(err)
						dbs.Master.Close()
					}
				}
				if !servicesAvailable.DB {
					if db, err := ConnAndPingMySQL(mySQLInfo, config); err != nil {
						errs <- err
					} else {
						dbs.Master = db
					}
				}
			}
		}
	}()
}

type DBs struct {
	Master *sql.DB
}

func ConnAndPingMySQL(sqlConn string, c *config.StorageMySQL) (*sql.DB, error) {
	db, err := sql.Open("mysql", sqlConn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}
	// Set settings fro DB
	db.SetConnMaxLifetime(c.ConnMaxLifetime)
	db.SetMaxOpenConns(c.MaxOpenConns)
	db.SetMaxIdleConns(c.MaxIdleConns)
	return db, nil
}

func DatabaseConnectionCheck(connections ...*sql.DB) error {
	for _, conn := range connections {
		if conn == nil {
			return errors.New("connection_nil")
		}
		if _, err := conn.Exec("SELECT 1 + 1;"); err != nil {
			return errors.New("connection_dead")
		}
	}
	return nil
}

func PerformMigrations(sqlConn string) error {
	migrator, err := setupMigrator(sqlConn)
	if err != nil {
		return err
	}

	if err := migrator.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("migrations are up-to-date")
			return nil
		}

		return err
	}

	return nil
}

func MigrationStatus(dsn string) (version uint, dirty bool, err error) {
	migrator, err := setupMigrator(dsn)
	if err != nil {
		return 0, false, err
	}

	return migrator.Version()
}

func PerformMigrateDown(sqlConn string) error {
	migrator, err := setupMigrator(sqlConn)
	if err != nil {
		return err
	}

	if err := migrator.Down(); err != nil {
		return err
	}

	return nil

}

func setupMigrator(sqlConn string) (*migrate.Migrate, error) {
	const driverName = "mysql"
	registerDriverOnce.Do(func() {
		migrations.RegisterDriver(driverName)
	})
	return migrate.New(fmt.Sprintf("%s://", driverName), sqlConn)
}

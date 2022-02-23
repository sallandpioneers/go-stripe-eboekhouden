package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/aceworksdev/go-stripe-eboekhouden/internal/storage/mysql"
	"github.com/spf13/cobra"
)

//nolint:gochecknoglobals // this is the recommended way to use cobra
var migrateOpts struct {
	mysqlConn string
}

//nolint:gochecknoinits // this is the recommended way to use cobra
func init() {
	migrateUpCommand.Flags().StringVarP(&migrateOpts.mysqlConn, "mysql-conn", "", "", "MySQL Connection URL")
	migrateStatusCommand.Flags().StringVarP(&migrateOpts.mysqlConn, "mysql-conn", "", "", "MySQL Connection URL")

	if err := migrateUpCommand.MarkFlagRequired("mysql-conn"); err != nil {
		log.Fatal(err)
	}

	if err := migrateStatusCommand.MarkFlagRequired("mysql-conn"); err != nil {
		log.Fatal(err)
	}

	migrateCommand.AddCommand(migrateUpCommand)
	migrateCommand.AddCommand(migrateDownCommand)
	migrateCommand.AddCommand(migrateStatusCommand)
}

//nolint:gochecknoglobals // this is the recommended way to use cobra
var migrateCommand = &cobra.Command{
	Use:   "migrate",
	Short: "Run the migration tool",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetFlags(log.Lshortfile | log.Ltime)
	},
}

//nolint:gochecknoglobals // this is the recommended way to use cobra
var migrateUpCommand = &cobra.Command{
	Use:   "up",
	Short: "Up the migrations",
	Run: func(cmd *cobra.Command, args []string) {
		err := mysql.PerformMigrations(migrateOpts.mysqlConn)
		if err != nil {
			log.Fatal(err)
		}

		os.Exit(0)
	},
}

//nolint:gochecknoglobals // this is the recommended way to use cobra
var migrateDownCommand = &cobra.Command{
	Use:   "down",
	Short: "Removal all migrations",
	Run: func(cmd *cobra.Command, args []string) {
		err := mysql.PerformMigrateDown(migrateOpts.mysqlConn)
		if err != nil {
			log.Fatal(err)
		}

		os.Exit(0)
	},
}

//nolint:gochecknoglobals // this is the recommended way to use cobra
var migrateStatusCommand = &cobra.Command{
	Use:   "status",
	Short: "Show migration status",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetFlags(log.Lshortfile | log.Ltime)

		version, dirty, err := mysql.MigrationStatus(migrateOpts.mysqlConn)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("dirty=%v version=%d\n", dirty, version)
	},
}

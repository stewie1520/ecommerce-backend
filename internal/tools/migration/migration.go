package migration

import (
	"github.com/pocketbase/dbx"
	pb_core "github.com/pocketbase/pocketbase/core"
	pb_migrations "github.com/pocketbase/pocketbase/migrations"
	pb_migration_log "github.com/pocketbase/pocketbase/migrations/logs"
	pb_migrate "github.com/pocketbase/pocketbase/tools/migrate"
)

type migrationsConnection struct {
	DB             *dbx.DB
	MigrationsList pb_migrate.MigrationsList
}

func migrationsConnectionsMap(app pb_core.App) map[string]migrationsConnection {
	return map[string]migrationsConnection{
		"db": {
			DB:             app.DB(),
			MigrationsList: pb_migrations.AppMigrations,
		},
		"logs": {
			DB:             app.LogsDB(),
			MigrationsList: pb_migration_log.LogsMigrations,
		},
	}
}

func RunMigrations(app pb_core.App) error {
	connections := migrationsConnectionsMap(app)

	for _, c := range connections {
		runner, err := pb_migrate.NewRunner(c.DB, c.MigrationsList)
		if err != nil {
			return err
		}

		if _, err := runner.Up(); err != nil {
			return err
		}
	}

	return nil
}

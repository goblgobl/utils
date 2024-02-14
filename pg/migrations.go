package pg

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"src.goblgobl.com/utils/log"
)

type Migrate func(tx pgx.Tx) error

type Migration struct {
	SQL     string
	Migrate Migrate
	Version uint16
}

func MigrateAll(db DB, appName string, migrations []Migration) error {
	latestVersion, err := GetCurrentMigrationVersion(db, appName)
	if err != nil {
		return err
	}

	log.Info("migration_check_start").String("app", appName).String("storage", "postgres").Int("installed_version", latestVersion).Log()
	for _, migration := range migrations {
		version := int(migration.Version)
		if version <= latestVersion {
			continue
		}

		err := db.Transaction(func(tx pgx.Tx) error {
			if sql := migration.SQL; sql != "" {
				if _, err := tx.Exec(context.Background(), sql); err != nil {
					return fmt.Errorf("Failed to run pg migration #%d - %w", version, err)
				}
			} else if err := migration.Migrate(tx); err != nil {
				return fmt.Errorf("Failed to run pg migration #%d - %w", version, err)
			}

			_, err = tx.Exec(context.Background(), `insert into migrations (app, version) values ($1, $2)`, appName, version)

			if err != nil {
				return fmt.Errorf("pg insert into migrations - %w", err)
			}
			return nil
		})

		if err != nil {
			log.Error("migration_fail").Int("version", version).Err(err).Log()
			return err
		}
		log.Info("migration_applied").Int("version", version).Log()
	}

	log.Info("migration_check_end").Log()

	return nil
}

func GetCurrentMigrationVersion(db DB, appName string) (int, error) {
	exists, err := db.TableExists("migrations")
	if err != nil {
		return 0, err
	}

	if !exists {
		_, err := db.Exec(context.Background(), `
			create table migrations (
				app text not null,
				version integer not null,
				created timestamptz not null default now(),
				primary key(app, version)
			)
		`)
		if err != nil {
			return 0, fmt.Errorf("pg create migrations - %w", err)
		}
		return 0, nil
	}

	value, err := Scalar[*int](db, `
		select max(version)
		from migrations
		where app = $1
	`, appName)

	if err != nil {
		return 0, fmt.Errorf("pg max migration - %w", err)
	}
	if value == nil {
		return 0, nil
	}
	return *value, nil
}

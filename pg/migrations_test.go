package pg

import (
	"context"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5"
	"src.goblgobl.com/tests/assert"
)

/*
TODO:
These tests mutate the migrations table, which
can cause some issue during local development.
We should just change the code so that the migration can receive a tx
which would allow us to run these tests within a tx and rollback

For now, we just try to copy the data and then restore it
*/

func Test_MigrateAll_NormalRun(t *testing.T) {
	realMigrations := testGetRealMigrations()
	defer func() {
		testRestoreRealMigrations(realMigrations)
	}()

	bg := context.Background()
	db.Exec(bg, "drop table if exists test_migrations")
	db.Exec(bg, "drop table if exists migrations")
	migrateTest := func(appName string) {
		err := MigrateAll(db, appName, []Migration{
			Migration{Version: 1, Migrate: MigrateOne},
			Migration{Version: 2, Migrate: MigrateTwo},
		})
		assert.Nil(t, err)

		value, err := Scalar[int](db, "select * from test_migrations")
		assert.Nil(t, err)
		assert.Equal(t, value, 9001)

		var version int
		rows, _ := db.Query(bg, "select version from migrations order by version")
		defer rows.Close()

		rows.Next()
		rows.Scan(&version)
		assert.Equal(t, version, 1)

		rows.Next()
		rows.Scan(&version)
		assert.Equal(t, version, 2)

		assert.False(t, rows.Next())

		current, err := GetCurrentMigrationVersion(db, appName)
		assert.Nil(t, err)
		assert.Equal(t, current, 2)
	}

	migrateTest("app1")
	migrateTest("app1") // this should be a noop
}

func Test_MigrateAll_Error(t *testing.T) {
	realMigrations := testGetRealMigrations()
	defer func() {
		testRestoreRealMigrations(realMigrations)
	}()

	bg := context.Background()
	db.Exec(bg, "drop table if exists test_migrations")
	db.Exec(bg, "drop table if exists migrations")
	migrateTest := func(appName string) {
		err := MigrateAll(db, appName, []Migration{
			Migration{Version: 1, Migrate: MigrateOne},
			Migration{Version: 2, Migrate: MigrateTwo},
			Migration{Version: 3, Migrate: MigrateErr},
		})
		assert.StringContains(t, err.Error(), "Failed to run pg migration #3")

		value, err := Scalar[int](db, "select * from test_migrations")
		assert.Nil(t, err)
		assert.Equal(t, value, 9001)

		var version int
		rows, _ := db.Query(bg, "select version from migrations order by version")
		defer rows.Close()

		rows.Next()
		rows.Scan(&version)
		assert.Equal(t, version, 1)

		rows.Next()
		rows.Scan(&version)
		assert.Equal(t, version, 2)

		assert.False(t, rows.Next())
	}

	migrateTest("app1")
	migrateTest("app1") // this should be a noop
}

func MigrateOne(tx pgx.Tx) error {
	_, err := tx.Exec(context.Background(), "create table test_migrations (id integer not null)")
	return err
}

func MigrateTwo(tx pgx.Tx) error {
	_, err := tx.Exec(context.Background(), "insert into test_migrations(id) values (9001)")
	return err
}

func MigrateErr(tx pgx.Tx) error {
	_, err := tx.Exec(context.Background(), "fail")
	return err
}

func testGetRealMigrations() map[string]int {
	rows, err := db.Query(context.Background(), "select app, version from migrations")
	if err != nil {
		if strings.Contains(err.Error(), `relation "migrations" does not exist`) {
			return nil
		}
		panic(err)
	}
	defer rows.Close()

	migrations := make(map[string]int)
	for rows.Next() {
		var app string
		var version int
		rows.Scan(&app, &version)
		migrations[app] = version
	}
	return migrations
}

func testRestoreRealMigrations(migrations map[string]int) {
	for app, version := range migrations {
		_, err := db.Exec(context.Background(), `
			insert into migrations (app, version)
			values ($1, $2)
			on conflict do nothing
		`, app, version)
		if err != nil {
			panic(err)
		}
	}
}

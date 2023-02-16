package pg

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/jackc/pgx/v5"
	"src.goblgobl.com/tests"
	"src.goblgobl.com/tests/assert"
)

var db DB

func init() {
	var err error
	db, err = New(tests.PG("gobl_utils_test"))
	if err != nil {
		panic(err)
	}
}

func Test_New_Invalid(t *testing.T) {
	_, err := New("nope")
	assert.Equal(t, err.Error(), "code: 3003 - cannot parse `nope`: failed to parse as DSN (invalid dsn)")
	assert.False(t, db.IsNotFound(err))
}

func Test_Scalar(t *testing.T) {
	value, err := Scalar[int](db, "select 566 from doesnotexist")
	assert.Equal(t, err.Error(), `ERROR: relation "doesnotexist" does not exist (SQLSTATE 42P01)`)
	assert.Equal(t, value, 0)

	value, err = Scalar[int](db, "select 566")
	assert.Nil(t, err)
	assert.Equal(t, value, 566)

	value, err = Scalar[int](db, "select 566 where $1", false)
	assert.True(t, db.IsNotFound(err))
	assert.True(t, errors.Is(err, pgx.ErrNoRows))
	assert.Equal(t, value, 0)

	str, err := Scalar[string](db, "select 'hello'")
	assert.Nil(t, err)
	assert.Equal(t, str, "hello")
}

func Test_DB_TableExist(t *testing.T) {
	db.MustExec("drop table if exists test_migrations")
	exists, err := db.TableExists("test_migrations")
	assert.Nil(t, err)
	assert.False(t, exists)

	db.MustExec("create table if not exists test_migrations (id int)")
	exists, err = db.TableExists("test_migrations")
	assert.Nil(t, err)
	assert.True(t, exists)
}

func Test_DB_Transaction_Rollback(t *testing.T) {
	bg := context.Background()
	db.MustExec("drop table if exists test_migrations")

	forcedErr := errors.New("forced error")
	err := db.Transaction(func(tx pgx.Tx) error {
		_, err := tx.Exec(bg, "create table if not exists test_migrations (id int)")
		assert.Nil(t, err)
		return forcedErr
	})
	assert.True(t, errors.Is(err, forcedErr))
	exists, err := db.TableExists("test_migrations")
	assert.Nil(t, err)
	assert.False(t, exists)
}

func Test_DB_Transaction_Commit(t *testing.T) {
	err := db.Transaction(func(tx pgx.Tx) error {
		_, err := tx.Exec(context.Background(), "create table if not exists test_migrations (id int)")
		assert.Nil(t, err)
		return nil
	})
	assert.Nil(t, err)
	exists, _ := db.TableExists("test_migrations")
	db.MustExec("drop table test_migrations")
	assert.True(t, exists)
}

func Test_DB_Placeholder(t *testing.T) {
	for i := 0; i < 50; i++ {
		assert.Equal(t, db.Placeholder(i), fmt.Sprintf("$%d", i+1))
	}
}

func Test_DB_RowToMap(t *testing.T) {
	m, err := db.RowToMap("select 1 where false")
	assert.True(t, errors.Is(err, ErrNoRows))
	assert.Equal(t, len(m), 0)

	m, err = db.RowToMap("select 1 as a, 'b' as b, '1200783b-3463-4a98-a527-fc61b6ac32f2'::uuid as uuid")
	assert.Nil(t, err)
	assert.Equal(t, len(m), 3)
	assert.Equal(t, m.Int("a"), 1)
	assert.Equal(t, m.String("b"), "b")
	assert.Equal(t, m.String("uuid"), "1200783b-3463-4a98-a527-fc61b6ac32f2")
}

func Test_DB_RowsToMap(t *testing.T) {
	m, err := db.RowsToMap("select 1 where false")
	assert.Nil(t, err)
	assert.Equal(t, len(m), 0)

	m, err = db.RowsToMap(`
			select 1 as a, 'b' as b, '1200783b-3463-4a98-a527-fc61b6ac32f2'::uuid as uuid
			union all
			select 2 as a, 'bee' as b, '0541242E-DE39-426B-8E71-BF12D00035FF'::uuid as uuid
	`)
	assert.Nil(t, err)
	assert.Equal(t, len(m), 2)

	row := m[0]
	assert.Equal(t, row.Int("a"), 1)
	assert.Equal(t, row.String("b"), "b")
	assert.Equal(t, row.String("uuid"), "1200783b-3463-4a98-a527-fc61b6ac32f2")

	row = m[1]
	assert.Equal(t, row.Int("a"), 2)
	assert.Equal(t, row.String("b"), "bee")
	assert.Equal(t, row.String("uuid"), "0541242e-de39-426b-8e71-bf12d00035ff")
}

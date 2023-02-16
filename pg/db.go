package pg

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"src.goblgobl.com/utils"
	"src.goblgobl.com/utils/log"
	"src.goblgobl.com/utils/typed"
	"src.goblgobl.com/utils/uuid"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrNoRows         = pgx.ErrNoRows
	ErrMoreThanOneRow = errors.New("Result contained more than 1 row")
)

type Row = pgx.Row

type DB struct {
	*pgxpool.Pool
}

func New(url string) (DB, error) {
	pool, err := pgxpool.New(context.Background(), url)
	if err != nil {
		return DB{}, log.Err(utils.ERR_PG_INIT, err).String("url", url)
	}
	return DB{pool}, nil
}

func Scalar[T any](db DB, sql string, args ...any) (T, error) {
	row := db.Pool.QueryRow(context.Background(), sql, args...)

	var value T
	err := row.Scan(&value)
	return value, err
}

func (db DB) TableExists(tableName string) (bool, error) {
	sql := `
		select exists (
			select 1 from pg_tables
			where schemaname = 'public' and tablename = $1
		)
	`
	exists, err := Scalar[bool](db, sql, tableName)
	if err != nil {
		if err == pgx.ErrNoRows {
			return false, nil
		}
		return false, fmt.Errorf("pg table exists - %w", err)
	}
	return exists, err
}

func (db DB) Transaction(fn func(tx pgx.Tx) error) error {
	bg := context.Background()
	tx, err := db.Pool.Begin(bg)
	if err != nil {
		return err
	}

	defer tx.Rollback(bg)
	if err := fn(tx); err != nil {
		return err
	}

	return tx.Commit(bg)
}

// Exists for our test factory which are designed to work with
// different databases
func (_ DB) Placeholder(i int) string {
	switch i {
	case 0:
		return "$1"
	case 1:
		return "$2"
	case 2:
		return "$3"
	case 3:
		return "$4"
	case 4:
		return "$5"
	case 5:
		return "$6"
	case 6:
		return "$7"
	case 7:
		return "$8"
	case 8:
		return "$9"
	case 9:
		return "$10"
	case 10:
		return "$11"
	case 11:
		return "$12"
	case 12:
		return "$13"
	case 13:
		return "$14"
	case 14:
		return "$15"
	case 15:
		return "$16"
	case 16:
		return "$17"
	case 17:
		return "$18"
	case 18:
		return "$19"
	case 19:
		return "$20"
	default:
		return "$" + strconv.Itoa(i+1)
	}
}

// Exists for our test factory which are designed to work with different databases
// For PG, the driver automatically handles JSON encoding.
func (_ DB) JSON(value any) (any, error) {
	return value, nil

}

func (db DB) MustExec(sql string, args ...any) {
	if _, err := db.Exec(context.Background(), sql, args...); err != nil {
		panic(err)
	}
}

func (db DB) RowToMap(sql string, args ...any) (typed.Typed, error) {
	rows, err := db.Query(context.Background(), sql, args...)
	if err != nil {
		return typed.Typed{}, err
	}

	slice, err := pgx.CollectRows(rows, pgx.RowToMap)
	if len(slice) == 0 {
		return nil, ErrNoRows
	}

	if len(slice) != 1 {
		return typed.Typed{}, ErrMoreThanOneRow
	}
	return typed.Typed(rowToMapTransform(slice[0])), nil
}

func (db DB) RowsToMap(sql string, args ...any) ([]typed.Typed, error) {
	rows, err := db.Query(context.Background(), sql, args...)
	if err != nil {
		return nil, err
	}

	slice, err := pgx.CollectRows(rows, pgx.RowToMap)
	if len(slice) == 0 {
		return nil, nil
	}

	t := make([]typed.Typed, len(slice))
	for i, slice := range slice {
		t[i] = typed.Typed(rowToMapTransform(slice))
	}
	return t, nil
}

func (db DB) IsNotFound(err error) bool {
	return err == ErrNoRows
}

// for now, just fix uuids
func rowToMapTransform(row map[string]any) map[string]any {
	for key, value := range row {
		if b, ok := value.([16]byte); ok && len(b) == 16 {
			uuid, err := uuid.FromBytes(b[:])
			if err != nil {
				panic(err)
			}
			row[key] = uuid
		}
	}
	return row
}

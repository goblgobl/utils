package sqlite

import (
	"strconv"

	"src.goblgobl.com/sqlite"
	"src.goblgobl.com/utils"
	"src.goblgobl.com/utils/json"
	"src.goblgobl.com/utils/log"
	"src.goblgobl.com/utils/typed"
)

var (
	ErrNoRows = sqlite.ErrNoRows
)

type Scanner sqlite.Scanner

type Conn struct {
	sqlite.Conn
}

func New(filePath string, create bool) (Conn, error) {
	conn, err := sqlite.Open(filePath, create)
	if err != nil {
		return Conn{}, log.Err(utils.ERR_SQLITE_INIT, err).String("path", filePath)
	}
	return Conn{conn}, nil
}

// Exists for our test factory which are designed to work with different databases
func (_ Conn) Placeholder(i int) string {
	switch i {
	case 0:
		return "?1"
	case 1:
		return "?2"
	case 2:
		return "?3"
	case 3:
		return "?4"
	case 4:
		return "?5"
	case 5:
		return "?6"
	case 6:
		return "?7"
	case 7:
		return "?8"
	case 8:
		return "?9"
	case 9:
		return "?10"
	case 10:
		return "?11"
	case 11:
		return "?12"
	case 12:
		return "?13"
	case 13:
		return "?14"
	case 14:
		return "?15"
	case 15:
		return "?16"
	case 16:
		return "?17"
	case 17:
		return "?18"
	case 18:
		return "?19"
	case 19:
		return "?20"
	default:
		return "?" + strconv.Itoa(i+1)
	}
}

// Exists for our test factory which are designed to work with different databases
// And, for SQLite, we can't just pass any object in to have it converted to JSON
func (_ Conn) JSON(value any) (any, error) {
	if value == nil {
		return nil, nil
	}
	return json.Marshal(value)
}

func Scalar[T any](conn Conn, sql string, args ...any) (T, error) {
	row := conn.Conn.Row(sql, args...)

	var value T
	err := row.Scan(&value)
	if err != nil {
		return value, err
	}
	return value, err
}

func (c Conn) RowToMap(sql string, args ...any) (typed.Typed, error) {
	m, err := c.Row(sql, args...).Map()
	if err != nil {
		return nil, err
	}
	return typed.Typed(m), err
}

func (c Conn) RowsToMap(sql string, args ...any) ([]typed.Typed, error) {
	rows := c.Rows(sql, args...)
	defer rows.Close()

	t := make([]typed.Typed, 0, 10)
	for rows.Next() {
		m, err := rows.Stmt.Map()
		if err != nil {
			return nil, err
		}
		t = append(t, typed.Typed(m))
	}

	if err := rows.Error(); err != nil {
		return nil, err
	}

	return t, nil
}

func (c Conn) IsNotFound(err error) bool {
	return err == ErrNoRows
}

func (c Conn) TableExists(tableName string) (bool, error) {
	sql := `
		select exists (
			select 1 from sqlite_master
			where type = 'table' and name = ?1
		)
	`
	exists, err := Scalar[bool](c, sql, tableName)
	if err == ErrNoRows {
		return false, nil
	}
	return exists, err
}

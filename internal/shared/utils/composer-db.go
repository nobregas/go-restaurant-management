package utils

import (
	"database/sql"
	"fmt"
	"reflect"
)

type ComposerDB struct {
	db      *sql.DB
	result  sql.Result
	table   string
	lastID  int
	columns []string
}

func (c *ComposerDB) Compose(model interface{}) interface{} {
	query := fmt.Sprintf("SELECT * FROM %s WHERE id = ?", c.table)

	row := c.db.QueryRow(query, c.lastID)

	modelValue := reflect.ValueOf(model)
	if modelValue.Kind() == reflect.Ptr {
		modelValue = modelValue.Elem()
	}

	numFields := modelValue.NumField()
	scanArgs := make([]interface{}, numFields)
	for i := 0; i < numFields; i++ {
		scanArgs[i] = modelValue.Field(i).Addr().Interface()
	}

	if err := row.Scan(scanArgs...); err != nil {
		return err
	}

	return modelValue.Interface()
}

func NewComposerDB(db *sql.DB, result sql.Result, table string) *ComposerDB {
	lastID, err := result.LastInsertId()
	if err != nil {
		panic(err)
	}

	return &ComposerDB{
		db:     db,
		result: result,
		table:  table,
		lastID: int(lastID),
	}
}

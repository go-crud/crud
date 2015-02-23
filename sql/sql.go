package sql

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/plimble/errs"
	"reflect"
	"strings"
)

type CRUD struct {
	db        *sqlx.DB
	dbName    string
	tableName string
}

func New(db *sqlx.DB, dbName, tableName string) *CRUD {
	return &CRUD{db, dbName, tableName}
}

func (c *CRUD) Create(v interface{}) error {
	colNames, vals := readDBTag(v)
	queryString := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", c.tableName, colNames, stringQuestionMark(len(vals)))
	_, err := c.db.Exec(queryString, vals...)
	return err
}

func (c *CRUD) Delete(id interface{}) error {
	queryString := fmt.Sprintf("DELETE FROM %s WHERE id=?", c.tableName)
	result, err := c.db.Exec(queryString, id)
	if i, _ := result.RowsAffected(); i == 0 {
		return errs.NewNotFound("not found")
	}
	return err
}

func (c *CRUD) Update(id interface{}, v map[string]interface{}) error {
	keys, vals := readUpdateKey(v)
	vals = append(vals, id)
	queryString := fmt.Sprintf("UPDATE %s SET %s WHERE id=?", c.tableName, keys)
	result, err := c.db.Exec(queryString, vals...)

	if i, _ := result.RowsAffected(); i == 0 {
		return errs.NewNotFound("not found")
	}

	return errs.Sql(err)
}

func (c *CRUD) Upsert(id, v interface{}) error {
	found, err := c.Exist(id)

	if found == true {
		newv := convertUpsertData(v)
		err = c.Update(id, newv)
	} else {
		err = c.Create(v)
	}
	return err
}

func (c *CRUD) Exist(id interface{}) (bool, error) {
	var count int
	err := db.Get(&count, "SELECT count(*) FROM "+c.tableName+" WHERE id=?", id)

	if count == 0 {
		return false, errs.NewNotFound("not found")
	}

	return true, err
}

func stringQuestionMark(count int) string {
	return "?" + strings.Repeat(",?", count-1)
}

func readUpdateKey(v map[string]interface{}) (string, []interface{}) {
	count := len(v)
	cols := make([]string, count)
	vals := make([]interface{}, count)
	i := 0
	for key, val := range v {
		cols[i] = key + "=?"
		vals[i] = val
		i++
	}
	return strings.Join(cols, ","), vals
}

func convertUpsertData(v interface{}) map[string]interface{} {
	s := reflect.ValueOf(v).Elem()
	typeOfT := s.Type()
	count := s.NumField()

	vals := make(map[string]interface{}, count)

	for i := 0; i < count; i++ {
		f := s.Field(i)
		vals[typeOfT.Field(i).Tag.Get("db")] = f.Interface()
	}

	return vals
}

func readDBTag(v interface{}) (string, []interface{}) {
	s := reflect.ValueOf(v).Elem()
	typeOfT := s.Type()
	count := s.NumField()

	cols := make([]string, count)
	vals := make([]interface{}, count)

	for i := 0; i < count; i++ {
		f := s.Field(i)
		vals[i] = f.Interface()
		cols[i] = typeOfT.Field(i).Tag.Get("db")
	}

	return strings.Join(cols, ","), vals
}

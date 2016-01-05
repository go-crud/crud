
package sql

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	errs "github.com/go-crud/errors2"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type testData struct {
	ID     int       `db:"id"`
	Name   string    `db:"name"`
	Time   time.Time `db:"time"`
	Height int       `db:"height"`
	Weight float64   `db:"weight"`
}

var db *sqlx.DB

func init() {
	time.Local = time.UTC
	var err error
	db, err = sqlx.Connect("sqlite3", ":memory:")
	if err != nil {
		panic(err)
	}
	schema := `CREATE TABLE IF NOT EXISTS testTable (
                id int(10) PRIMARY KEY,
                name varchar(50),
                time timestamp,
                height int(10),
                weight float(10,2)
            );`
	if _, err := db.Exec(schema); err != nil {
		panic(err)
	}
}

func getSetup() *CRUD {
	return NewCRUD(db, "testDB", "testTable")
}

func TestReadDBTag(t *testing.T) {
	now := time.Now().UTC()
	v := &testData{ID: 11, Name: "Xier", Time: now, Height: 1500, Weight: 25.55}
	cols, vals := readDBTag(v)
	assert.Equal(t, "id,name,time,height,weight", cols)
	assert.Equal(t, []interface{}{11, "Xier", now, 1500, 25.55}, vals)
}

func TestReadUpdateKey(t *testing.T) {
	x := map[string]interface{}{
		"foo": []string{"a", "b"},
		"bar": "foo",
		"baz": 10.4,
	}
	queryString, vals := readUpdateKey(x)

	assert.Contains(t, queryString, "foo=?")
	assert.Contains(t, queryString, "bar=?")
	assert.Contains(t, queryString, "baz=?")

	assert.Contains(t, vals, []string{"a", "b"})
	assert.Contains(t, vals, "foo")
	assert.Contains(t, vals, 10.4)
}

func TestUpsertCreate(t *testing.T) {
	crud := getSetup()
	now := time.Now()

	v := testData{ID: 25, Name: "Xier", Time: now, Height: 1500, Weight: 25.55}
	err := crud.Upsert("25", &v)
	assert.NoError(t, err)

	var tempData testData
	err = crud.db.Get(&tempData, "SELECT * FROM testTable WHERE id=25")
	assert.NoError(t, err)
	assert.Equal(t, v, tempData)
}

func TestUpsertUpdate(t *testing.T) {
	crud := getSetup()
	now := time.Now()

	v := testData{ID: 26, Name: "Xier", Time: now, Height: 1500, Weight: 25.55}
	err := crud.Create(&v)
	assert.NoError(t, err)

	v2 := testData{ID: 26, Name: "Zier", Time: now, Height: 200, Weight: 33.32}
	err = crud.Upsert("26", &v2)
	assert.NoError(t, err)

	var tempData testData
	err = crud.db.Get(&tempData, "SELECT * FROM testTable WHERE id=26")
	assert.NoError(t, err)
	assert.Equal(t, v2, tempData)
}

func TestconvertUpsertData(t *testing.T) {
	now := time.Now().UTC()

	checkVals := map[string]interface{}{
		"id":     24,
		"name":   "Xier",
		"time":   now,
		"height": 1500,
		"weight": 25.55,
	}

	v := &testData{ID: 24, Name: "Xier", Time: now, Height: 1500, Weight: 25.55}
	vals := convertUpsertData(v)
	assert.Equal(t, checkVals, vals)
}

func TestExist(t *testing.T) {
	crud := getSetup()
	now := time.Now()

	v := testData{ID: 18, Name: "Xier", Time: now, Height: 1500, Weight: 25.55}
	err := crud.Create(&v)
	assert.NoError(t, err)

	found, err := crud.Exist("18")
	assert.NoError(t, err)
	assert.True(t, found)
}

func TestExistNotFound(t *testing.T) {
	crud := getSetup()
	_, err := crud.Exist("19")
	assert.Error(t, err)
	assert.True(t, errs.IsNotFound(err))
}

func TestUpdate(t *testing.T) {
	crud := getSetup()
	now := time.Now()

	v := testData{ID: 13, Name: "Xier", Time: now, Height: 1500, Weight: 25.55}
	err := crud.Create(&v)
	assert.NoError(t, err)

	dataForUpdate := map[string]interface{}{
		"name":   "Wichit",
		"height": 900,
		"weight": 44.44,
	}
	err = crud.Update(13, dataForUpdate)
	assert.NoError(t, err)

	v2 := testData{ID: 13, Name: "Wichit", Time: now, Height: 900, Weight: 44.44}
	var tempData testData
	err = crud.db.Get(&tempData, "SELECT * FROM testTable WHERE id=13")
	assert.NoError(t, err)
	assert.Equal(t, v2, tempData)
}

func TestUpdateNotFound(t *testing.T) {
	crud := getSetup()

	dataForUpdate := map[string]interface{}{
		"name":   "Wichit",
		"height": 900,
		"weight": 44.44,
	}
	err := crud.Update(14, dataForUpdate)
	assert.Error(t, err)
	assert.True(t, errs.IsNotFound(err))
}

func TestCreate(t *testing.T) {
	crud := getSetup()
	now := time.Now()
	v := testData{ID: 11, Name: "Xier", Time: now, Height: 1500, Weight: 25.55}

	err := crud.Create(&v)
	assert.NoError(t, err)

	var tempData testData
	err = crud.db.Get(&tempData, "SELECT * FROM testTable WHERE id=11")
	assert.NoError(t, err)
	assert.Equal(t, v, tempData)
}

func TestDelete(t *testing.T) {
	crud := getSetup()
	now := time.Now()
	v := testData{ID: 12, Name: "Xier", Time: now, Height: 1500, Weight: 25.55}

	err := crud.Create(&v)
	assert.NoError(t, err)

	err = crud.Delete(v.ID)
	assert.NoError(t, err)

	var tempData testData

	err = crud.db.Get(&tempData, "SELECT * FROM testTable WHERE id=12")
	err = errs.Sql(err)
	assert.Error(t, err)
	assert.True(t, errs.IsNotFound(err))
}

func TestDeleteNotFound(t *testing.T) {
	crud := getSetup()
	err := crud.Delete("16")
	assert.Error(t, err)
	assert.True(t, errs.IsNotFound(err))
}

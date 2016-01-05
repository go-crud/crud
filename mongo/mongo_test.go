package mongo

import (
//    "fmt"
    "encoding/json"
    "github.com/go-crud/errors2"
    
	"github.com/stretchr/testify/assert"
	"gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
	"testing"
)

type testData struct {
	Id   bson.ObjectId `bson:"_id,omitempty"`
	Name string `bson:"name"`
    Age int `bson:"age"`
}

func testCRUD_JSON(t *testing.T) {
	assert := assert.New(t)

	session, dberr := mgo.Dial("192.168.1.178:27017")
	assert.NoError(dberr)
    defer session.Close()

	db := "test"
	c := "users"
	crud := NewCRUD(session, db, c)
    users:=session.DB(db).C(c)
	users.RemoveAll(nil)

    var jsonBlob = []byte(`{"Id":"","Name": "Alex"}`)
    var data map[string]interface{}
    err := json.Unmarshal(jsonBlob, &data)  
    assert.NoError(err)
	err = crud.Insert(data)
    assert.NoError(err)
   // fmt.Printf("%v\n",data);
	assert.NoError(err)
}
func TestMap(t *testing.T) {
    args:=make(map[string]interface{})
    args["Age"]=55
    v,ok:=args["Age"]
    assert := assert.New(t)
    assert.Equal(v.(int), 55)
    assert.True(ok)
}
func TestCRUD(t *testing.T) {
	assert := assert.New(t)

	session, dberr := mgo.Dial("192.168.1.178:27017")
	assert.NoError(dberr)
    defer session.Close()

	db := "test"
	c := "users"
	crud := NewCRUD(session, db, c)
    
    users:=session.DB(db).C(c)

	users.RemoveAll(nil)

	data1 := &testData{Name: "Tom",Age:44}
	err := crud.Insert(data1)
    assert.NoError(err)

	var data2 *testData
	//get none exist data
	err = users.FindId("2").One(&data2)
	assert.True(errors2.IsNotFound(errors2.Mgo(err)))
	assert.Nil(data2)

	//get exist data
	err = users.FindId(data1.Id).One(&data2)
	assert.NoError(err)
	assert.Equal(data2, data1)

	//check none exist data
	exist, err := crud.Exist("2")
	assert.NoError(err)
	assert.False(exist)

	//check exist with data
	exist, err = crud.Exist(data1.Id)
	assert.NoError(err)
	assert.True(exist)

	//update exist data
    args:=make(map[string]interface{})
    args["age"]=35
	err = crud.Update(data1.Id, args)
	assert.NoError(err)
	var data3 *testData
	err = users.FindId(data1.Id).One(&data3)
	assert.NoError(err)
	assert.NotEqual(data3, data2)

    err = crud.UpdateAll(data1.Id, data2)
	assert.NoError(err)
    err = users.FindId(data1.Id).One(&data3)
	assert.NoError(err)
    assert.Equal(data3, data2)
    
	//update none exist data
	err = crud.UpdateAll("2", data2)
	assert.True(errors2.IsNotFound(err))

	//delete exist data
	err = crud.Delete(data1.Id)
	assert.NoError(err)
	var data4 *testData
	err = users.FindId(data1.Id).One(&data4)
	assert.True(errors2.IsNotFound(errors2.Mgo(err)))
	assert.Nil(data4)

	//delete none exist data
	err = crud.Delete("2")
	assert.True(errors2.IsNotFound(err))
}

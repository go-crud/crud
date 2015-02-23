package mongo

import (
	"github.com/plimble/utils/errors2"

	"github.com/stretchr/testify/assert"
	"gopkg.in/mgo.v2"
	"testing"
)

type testData struct {
	ID   string `bson:"_id"`
	Name string `bson:"name"`
}

func TestCRUD(t *testing.T) {
	assert := assert.New(t)

	session, err := mgo.Dial("192.168.59.103:27017")
	assert.NoError(err)

	db := "mongocrud"
	c := "test"
	crud := New(session, db, c)

	session.DB(db).C(c).RemoveAll(nil)

	//add data
	data1 := &testData{"1", "name1"}
	err = crud.Insert(data1)
	assert.NoError(err)

	//add duplicate data
	err = crud.Insert(data1)
	assert.Error(err)

	var data2 *testData

	//get none exist data
	err = session.DB(db).C(c).FindId("2").One(&data2)
	assert.True(errors2.IsNotFound(errors2.Mgo(err)))
	assert.Nil(data2)

	//get exist data
	err = session.DB(db).C(c).FindId("1").One(&data2)
	assert.NoError(err)
	assert.Equal(data2, data1)

	//check none exist data
	exist, err := crud.Exist("2")
	assert.NoError(err)
	assert.False(exist)

	//check exist with data
	exist, err = crud.Exist("1")
	assert.NoError(err)
	assert.True(exist)

	//update exist data
	data2.Name = "name_updated"
	err = crud.UpdateAll("1", data2)
	assert.NoError(err)
	var data3 *testData
	err = session.DB(db).C(c).FindId("1").One(&data3)
	assert.NoError(err)
	assert.Equal(data3, data2)

	//update none exist data
	err = crud.UpdateAll("2", data2)
	assert.True(errors2.IsNotFound(err))

	//delete exist data
	err = crud.Delete("1")
	assert.NoError(err)
	var data4 *testData
	err = session.DB(db).C(c).FindId("1").One(&data4)
	assert.True(errors2.IsNotFound(errors2.Mgo(err)))
	assert.Nil(data4)

	//delete none exist data
	err = crud.Delete("2")
	assert.True(errors2.IsNotFound(err))
}

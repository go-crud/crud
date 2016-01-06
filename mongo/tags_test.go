package mongo

import (
    "encoding/json"
 	"github.com/stretchr/testify/assert"
	"gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
	"testing"
)



func testTags(t *testing.T) {
	assert := assert.New(t)

	session, dberr := mgo.Dial("192.168.1.178:27017")
	assert.NoError(dberr)
    defer session.Close()

	db := "test"
	c := "tags"
	crud := NewCRUD(session, db, c)
    tags:=session.DB(db).C(c)
	tags.RemoveAll(nil)

    var jsonBlob = []byte(`{"_id":"","name":"My Tags","children":[
                                {"name":"size","children":[
                                   {"name": "small"},
                                   {"name": "big"}
                                ]},   
                                {"name":"distance","children":[
                                    {"name":"far"},
                                    {"name":"near"}
                                ]}  
                           ]}`)
    var data map[string]interface{}
    err := json.Unmarshal(jsonBlob, &data)  
    assert.NoError(err)
	err = crud.Insert(data)
    assert.NoError(err)
   // fmt.Printf("%v\n",data);
	assert.NoError(err)
    args:=make(map[string]interface{})
    args["name"]="medium"
    tags.Update(bson.M{"children.name": "size"},bson.M{ "$addToSet": bson.M{"children.$.children":args} })
    
}

package mongo

import (
    "reflect"
	"github.com/go-crud/errors2"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type CRUD struct {
	session *mgo.Session
	db      string
	c       string
}

func NewCRUD(session *mgo.Session, db, c string) *CRUD {
	return &CRUD{session, db, c}
}

func (crud *CRUD) Insert(v interface{}) (error) {
	session := crud.session.Copy()
	defer session.Close()
    value:=reflect.ValueOf(v)
   
    id:=bson.NewObjectId()// bid,_:=id.MarshalJSON()
    if value.Type().Kind() == reflect.Map {
        m := v.(map[string]interface{})
        if len(m["_id"].(string))<1 {
            m["_id"]=id
        }
    }else{
      s := value.Elem()
      f := s.FieldByName("Id")
      oid := f.Interface()
      if len(oid.(bson.ObjectId))<1{
          f.SetString(string(id))
      }
    }
  
    err:=errors2.Mgo(session.DB(crud.db).C(crud.c).Insert(v))
	return err
}

func (crud *CRUD) Delete(id interface{}) error {
	session := crud.session.Copy()
	defer session.Close()
    // bson.ObjectIdHex(id)
	return errors2.Mgo(session.DB(crud.db).C(crud.c).RemoveId(id))
}

func (crud *CRUD) Upsert(id, v interface{}) error {
	session := crud.session.Copy()
	defer session.Close()
	_, err := session.DB(crud.db).C(crud.c).UpsertId(id, v)
	return errors2.Mgo(err)
}

func (crud *CRUD) Update(id interface{}, v map[string]interface{}) error {
	session := crud.session.Copy()
	defer session.Close()
	return errors2.Mgo(session.DB(crud.db).C(crud.c).UpdateId(id, bson.M{"$set": v}))
}

func (crud *CRUD) UpdateAll(id interface{}, v interface{}) error {
	session := crud.session.Copy()
	defer session.Close()
	return errors2.Mgo(session.DB(crud.db).C(crud.c).UpdateId(id, v))
}

func (crud *CRUD) Exist(id interface{}) (bool, error) {
	session := crud.session.Copy()
	defer session.Close()
	count, err := session.DB(crud.db).C(crud.c).FindId(id).Count()
	if count == 0 {
		return false, err
	}
	return true, err
}

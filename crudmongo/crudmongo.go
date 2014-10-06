package crudmongo

import (
	"gopkg.in/mgo.v2"
)

type CRUD struct {
	session *mgo.Session
	db      string
	c       string
}

func New(session *mgo.Session, db, c string) *CRUD {
	return &CRUD{session, db, c}
}

func (crud *CRUD) Add(v interface{}) error {
	session := crud.session.Copy()
	defer session.Close()

	return session.DB(crud.db).C(crud.c).Insert(v)
}

func (crud *CRUD) IsNotFound(err error) bool {
	if err != mgo.ErrNotFound {
		return false
	}

	return true
}

func (crud *CRUD) Get(id, v interface{}) error {
	session := crud.session.Copy()
	defer session.Close()

	err := session.DB(crud.db).C(crud.c).FindId(id).One(v)

	return err
}

func (crud *CRUD) Delete(id interface{}) error {
	session := crud.session.Copy()
	defer session.Close()

	return session.DB(crud.db).C(crud.c).RemoveId(id)
}

func (crud *CRUD) Update(id, v interface{}) error {
	session := crud.session.Copy()
	defer session.Close()

	return session.DB(crud.db).C(crud.c).UpdateId(id, v)
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

package crud

type CRUD interface {
	Create(v interface{}) error
	Delete(id interface{}) error
	Update(id, v interface{}) error
	Upsert(id, v interface{}) error
	Exist(id interface{}) (bool, error)
}

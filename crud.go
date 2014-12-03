package crud

type CRUD interface {
	Create(v interface{}) error
	Delete(id interface{}) error
	Update(id interface{}, v map[string]interface{}) error
	Upsert(id, v interface{}) error
	Exist(id interface{}) (bool, error)
	GetByID(id interface{}, v interface{}) error
}

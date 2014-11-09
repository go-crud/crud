package crud

type CRUD interface {
	Create(v interface{}) error
	Get(id, v interface{}) error
	Delete(id interface{}) error
	Update(id, v interface{}) error
	IsNotFound(err error) bool
	Exist(id interface{}) (bool, error)
}

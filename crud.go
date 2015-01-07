package crud

type CRUD interface {
	Insert(v interface{}) error
	Delete(id string) error
	Update(id string, v map[string]interface{}) error
	UpdateAll(id string, v interface{}) error
	Upsert(id string, v interface{}) error
	Exist(id string) (bool, error)
	Get(id string, v interface{}) error
}

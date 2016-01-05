package crud

type CRUD interface {
	Insert(v interface{})  error
	Delete(id interface{}) error
	Update(id interface{}, v map[string]interface{}) error
	UpdateAll(id interface{}, v interface{}) error
	Upsert(id interface{}, v interface{}) error
	Exist(id interface{}) (bool, error)
}

type TreeOp interface {
	Init(v map[string]interface{}) error
	UpdateNode(filter map[string]interface{},  v map[string]interface{}) error
}
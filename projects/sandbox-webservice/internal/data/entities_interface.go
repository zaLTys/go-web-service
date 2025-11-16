package data

type EntityRepository interface {
	Insert(entity *Entity) error
	Get(id int64) (*Entity, error)
	Update(entity *Entity) error
	Delete(id int64) error
	GetAll() ([]*Entity, error)
}

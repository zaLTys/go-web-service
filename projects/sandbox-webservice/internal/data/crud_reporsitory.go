package data

type CRUDRepository[T any] interface {
	Insert(*T) error
	Get(id int64) (*T, error)
	Update(*T) error
	Delete(id int64) error
	GetAll() ([]*T, error)
}

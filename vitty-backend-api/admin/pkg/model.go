package pkg

// type ModelInterface interface {
// 	GetAll() ([]ModelInterface, error)
// 	Get(id string) (ModelInterface, error)
// 	GetPKField() string
// 	GetItemName() string
// 	Create(model interface{}) error
// 	Update(model interface{}) error
// 	Delete(model interface{}) error
// }

type ModelInterface interface {
	GetAll() ([]ModelInterface, error)
	Get(id string) (ModelInterface, error)
	GetPKField() string
	GetItemName() string
	Create(model ModelInterface) error
	Update(model ModelInterface) error
	Delete(model ModelInterface) error
}

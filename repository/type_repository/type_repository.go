package type_repository

import "erp-service/model/entity"

type TypeRepository interface {
	AddType(*entity.TypeTransaction) (*string, error)
	GetAllType() ([]entity.TypeTransaction, error)
	GetDetailType(id string) (*entity.TypeTransaction, error)
	UpdateType(id string, types *entity.TypeTransaction) (*string, error)
	DeleteType(id string) error
}

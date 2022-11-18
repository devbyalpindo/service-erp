package type_repository

import (
	"erp-service/helper"
	"erp-service/model/entity"

	"gorm.io/gorm"
)

type TypeRepositoryImpl struct {
	DB *gorm.DB
}

func NewTypeRepository(DB *gorm.DB) TypeRepository {
	return &TypeRepositoryImpl{DB}
}

func (repository *TypeRepositoryImpl) AddType(types *entity.TypeTransaction) (*string, error) {
	if err := repository.DB.Create(&types).Error; err != nil {
		return nil, err
	}

	return &types.TypeID, nil
}

func (repository *TypeRepositoryImpl) GetAllType() ([]entity.TypeTransaction, error) {
	types := []entity.TypeTransaction{}
	err := repository.DB.Model(&entity.TypeTransaction{}).Scan(&types).Error
	helper.PanicIfError(err)
	if len(types) <= 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return types, nil
}

func (repository *TypeRepositoryImpl) GetDetailType(id string) (*entity.TypeTransaction, error) {
	types := entity.TypeTransaction{}
	result := repository.DB.Where("type_id = ?", id).Find(&types)

	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return &types, nil
}

func (repository *TypeRepositoryImpl) UpdateType(id string, types *entity.TypeTransaction) (*string, error) {
	result := repository.DB.Model(&types).Where("type_id = ?", id).Updates(entity.TypeTransaction{TypeTransaction: types.TypeTransaction, Description: types.Description})

	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return &id, nil
}

func (repository *TypeRepositoryImpl) DeleteType(id string) error {
	result := repository.DB.Where("type_id = ?", id).Delete(&entity.TypeTransaction{})

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

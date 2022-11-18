package seeder

import (
	"erp-service/model/entity"
	"testing"

	"github.com/google/uuid"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestSeedType(t *testing.T) {

	dsn := "root:@tcp(localhost:3306)/erp_salju?charset=utf8mb4&parseTime=True&loc=Local"

	DB, err := gorm.Open(mysql.Open(dsn))

	if err != nil {
		t.Fatal("failed to connect database")
	}

	types := []string{"DEPOSIT", "WITHDRAW", "BONUS"}
	typeList := []entity.TypeTransaction{}

	for _, s := range types {
		value := entity.TypeTransaction{
			TypeID:          uuid.New().String(),
			TypeTransaction: s,
			Description:     "",
		}

		typeList = append(typeList, value)
	}

	errDB := DB.Create(&typeList).Error
	if errDB != nil {
		t.Fatal("error batch insert to mysql", errDB)
	}
}

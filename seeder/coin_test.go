package seeder

import (
	"erp-service/model/entity"
	"testing"
	"time"

	"github.com/google/uuid"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestSeedCoin(t *testing.T) {

	dsn := "root:@tcp(localhost:3306)/erp_salju?charset=utf8mb4&parseTime=True&loc=Local"

	DB, err := gorm.Open(mysql.Open(dsn))

	if err != nil {
		t.Fatal("failed to connect database")
	}

	createID := uuid.New().String()
	role := []entity.Coin{
		{
			CoinID:    createID,
			CoinName:  "SALJU 88",
			Balance:   0,
			Note:      "Coin Games Salju88",
			CreatedAt: time.Now().Format("2006-01-02 15:04:05"),
			UpdatedAt: time.Now().Format("2006-01-02 15:04:05"),
		},
	}

	errDB := DB.Create(&role).Error
	if errDB != nil {
		t.Fatal("error batch insert to mysql", errDB)
	}
}

package seeder

import (
	"erp-service/model/entity"
	"testing"

	"github.com/google/uuid"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestSeedRole(t *testing.T) {

	dsn := "root:@tcp(localhost:3306)/erp_salju?charset=utf8mb4&parseTime=True&loc=Local"

	DB, err := gorm.Open(mysql.Open(dsn))

	if err != nil {
		t.Fatal("failed to connect database")
	}

	listedRole := []string{"ADMIN", "DEPOSITOR", "WIDTHDRAWER"}
	listRole := []entity.Role{}

	for _, s := range listedRole {
		role := entity.Role{
			RoleID:   uuid.New().String(),
			RoleName: s,
		}
		listRole = append(listRole, role)
	}

	errDB := DB.Create(&listRole).Error
	if errDB != nil {
		t.Fatal("error batch insert to mysql", errDB)
	}
}

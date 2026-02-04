package migration

import (
	"github.com/minisource/template_go/config"
	"github.com/minisource/template_go/domain/model"
	"github.com/minisource/go-common/db/gorm"
	"github.com/minisource/go-common/logging"
	"gorm.io/gorm"
)

const countStarExp = "count(*)"

var logger = logging.NewLogger(&config.GetConfig().Logger)

func Up1() {
	database := gormdb.GetDb()

	createTables(database)
	// createCountry(database)
}

func createTables(database *gorm.DB) {
	tables := []interface{}{}

	// User
	tables = addNewTable(database, model.User{}, tables)

	err := database.Migrator().CreateTable(tables...)
	if err != nil {
		logger.Error(logging.Postgres, logging.Migration, err.Error(), nil)
	}
	logger.Info(logging.Postgres, logging.Migration, "tables created", nil)
}

func addNewTable(database *gorm.DB, model interface{}, tables []interface{}) []interface{} {
	if !database.Migrator().HasTable(model) {
		tables = append(tables, model)
	}
	return tables
}

// func createCountry(database *gorm.DB) {
// 	count := 0
// 	database.
// 		Model(&models.Country{}).
// 		Select(countStarExp).
// 		Find(&count)
// 	if count == 0 {
// 		database.Create(&models.Country{Name: "Iran", Cities: []models.City{
// 			{Name: "Tehran"},
// 			{Name: "Isfahan"},
// 			{Name: "Shiraz"},
// 			{Name: "Chalus"},
// 			{Name: "Ahwaz"},
// 		}})
// 	}
// }

func Down1() {
	// nothing
}

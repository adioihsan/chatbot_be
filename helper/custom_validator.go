package helper

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

func ExistsValidator(db *gorm.DB) validator.Func {
	return func(fl validator.FieldLevel) bool {
		params := strings.Split(fl.Param(), "=")
		if len(params) != 2 {
			return false
		}

		tableName := params[0]
		columnName := params[1]
		fieldValue := fl.Field().Interface()

		var exists bool
		query := fmt.Sprintf("%s = ?", columnName)
		result := db.Table(tableName).Select("count(*) > 0").Where(query, fieldValue).Find(&exists)
		return result.Error == nil && exists
	}
}

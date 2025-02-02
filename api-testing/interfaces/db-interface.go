package interfaces

import (
	"gorm.io/gorm"
)

type DBInterface interface {
	First(dest interface{}, conds ...interface{}) *gorm.DB
	Create(value interface{}) *gorm.DB
}

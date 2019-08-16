package sql

import "github.com/jinzhu/gorm"

func GetDb() (db *gorm.DB, err error) {
	db, err = gorm.Open(config.SqlDialect, config.SqlConnectionString)
	return
}

//import (
//	"github.com/jinzhu/gorm"
//	_ "github.com/jinzhu/gorm/dialects/mysql"
//)
//
//type Scope struct {
//	gorm.Scope
//}
//
//// TableName return table name
//func (scope *Scope) TableName() string {
//	if scope.Search != nil && len(scope.Search.tableName) > 0 {
//		return scope.Search.tableName
//	}
//
//	if tabler, ok := scope.Value.(tabler); ok {
//		return tabler.TableName()
//	}
//
//	if tabler, ok := scope.Value.(dbTabler); ok {
//		return tabler.TableName(scope.db)
//	}
//
//	return scope.GetModelStruct().TableName(scope.db.Model(scope.Value))
//}
//
//type DB struct {
//	gorm.DB
//}
//
//
//// NewScope create a scope for current operation
//func (s *DB) NewScope(value interface{}) *gorm.Scope {
//	scope := s.DB.NewScope(value)
//
//	return scope
//}

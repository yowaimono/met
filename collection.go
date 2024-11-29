package met

import (
	"reflect"
	"strings"

	"go.mongodb.org/mongo-driver/mongo"
)

// GetCollection 获取集合
func (d *DB) GetCollection(model interface{}) *mongo.Collection {
	collectionName := d.CollectionName(model)
	return d.Cli.Database(d.Cfg.DbName).Collection(collectionName)
}

// CollectionName 获取模型的集合名称
func (d *DB) CollectionName(model interface{}) string {
	value := reflect.ValueOf(model)
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}
	method := value.MethodByName("CollectionName")
	if method.IsValid() {
		results := method.Call(nil)
		if len(results) > 0 {
			return results[0].String()
		}
	}
	return strings.ToLower(value.Type().Name())
}

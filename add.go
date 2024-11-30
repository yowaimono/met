package mongoq

import (
	"context"
	"fmt"
	"reflect"
)

func (d *DB) Add(data interface{}) error {
	collection := d.GetCollection(data)

	// Hooks

	// implement hooks BeforeAdd
	// call it
	if v, ok := data.(interface{ BeforeAdd() }); ok {
		v.BeforeAdd()
	}

	result, err := collection.InsertOne(context.Background(), data)
	if err != nil {
		return err
	}

	// 使用反射回写 _id
	v := reflect.ValueOf(data).Elem()
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		bsonTag := field.Tag.Get("bson")

		// 检查 bson 标签是否为 "_id"
		if bsonTag == "_id,omitempty" {
			idField := v.Field(i)
			if idField.IsValid() && idField.CanSet() {
				idField.Set(reflect.ValueOf(result.InsertedID))
			}
			break
		}
	}

	// Hooks

	// implement hooks AfterAdd
	// call it
	if v, ok := data.(interface{ AfterAdd() }); ok {
		v.AfterAdd()
	}

	return nil
}

// InsertMany 插入多个文档
func (d *DB) AddMany(data interface{}) error {
	collection := d.GetCollection(data)

	// 使用反射获取切片
	v := reflect.ValueOf(data)
	if v.Kind() != reflect.Slice {
		return fmt.Errorf("data must be a slice")
	}

	// 创建文档切片
	var docs []interface{}
	for i := 0; i < v.Len(); i++ {
		docs = append(docs, v.Index(i).Interface())
	}

	// Hooks
	// implement hooks BeforeAdd
	// call it
	for _, doc := range docs {
		if v, ok := doc.(interface{ BeforeAdd() }); ok {
			v.BeforeAdd()
		}
	}

	result, err := collection.InsertMany(context.Background(), docs)
	if err != nil {
		return err
	}

	// 使用反射回写 _id
	for i, insertedID := range result.InsertedIDs {
		elem := v.Index(i)
		t := elem.Type()

		for j := 0; j < t.NumField(); j++ {
			field := t.Field(j)
			bsonTag := field.Tag.Get("bson")

			// 检查 bson 标签是否为 "_id"
			if bsonTag == "_id,omitempty" {
				idField := elem.Field(j)
				if idField.IsValid() && idField.CanSet() {
					idField.Set(reflect.ValueOf(insertedID))
				}
				break
			}
		}
	}

	// Hooks
	// implement hooks AfterAdd
	// call it
	for _, doc := range docs {
		if v, ok := doc.(interface{ AfterAdd() }); ok {
			v.AfterAdd()
		}
	}

	return nil
}

package met

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (d *DB) Filter(query ...interface{}) *DB {
	d.filter = bson.M{}
	for i := 0; i < len(query); i += 2 {
		key := query[i].(string)
		value := query[i+1]
		d.filter[key] = value
	}
	return d
}

func (d *DB) UpdateMany(data interface{}, query ...interface{}) error {
	collection := d.GetCollection(data)

	// Hooks
	// implement hooks BeforeUpdate
	// call it
	if v, ok := data.(interface{ BeforeUpdate() }); ok {
		v.BeforeUpdate()
	}

	// 检查 bson 标签是否为 "_id"

	for i := 0; i < len(query); i += 2 {
		key := query[i].(string)
		value := query[i+1]
		d.update[key] = value
	}

	result, err := collection.UpdateMany(context.Background(), d.filter, bson.M{"$set": d.update})
	if err != nil {
		return err
	}

	if v, ok := data.(interface{ AfterUpdate() }); ok {
		v.AfterUpdate()
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("no document matched")
	}
	return nil

	//	return nil
}

func (d *DB) UpdateOne(data interface{}, query ...interface{}) error {
	collection := d.GetCollection(data)

	// Hooks
	// implement hooks BeforeUpdate
	// call it
	if v, ok := data.(interface{ BeforeUpdate() }); ok {
		v.BeforeUpdate()
	}

	// 检查 bson 标签是否为 "_id"

	for i := 0; i < len(query); i += 2 {
		key := query[i].(string)
		value := query[i+1]
		d.update[key] = value
	}

	result, err := collection.UpdateOne(context.Background(), d.filter, bson.M{"$set": d.update})
	if err != nil {
		return err
	}

	if v, ok := data.(interface{ AfterUpdate() }); ok {
		v.AfterUpdate()
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("no document matched")
	}
	return nil

	//	return nil
}

// func (d *DB) Add(data interface{}) error {
// 	collection := d.GetCollection(data)

// 	// Hooks

// 	// implement hooks BeforeAdd
// 	// call it
// 	if v, ok := data.(interface{ BeforeAdd() }); ok {
// 		v.BeforeAdd()
// 	}

// 	result, err := collection.InsertOne(context.Background(), data)
// 	if err != nil {
// 		return err
// 	}

// 	// 使用反射回写 _id
// 	v := reflect.ValueOf(data).Elem()
// 	t := v.Type()

// 	for i := 0; i < t.NumField(); i++ {
// 		field := t.Field(i)
// 		bsonTag := field.Tag.Get("bson")

// 		// 检查 bson 标签是否为 "_id"
// 		if bsonTag == "_id,omitempty" {
// 			idField := v.Field(i)
// 			if idField.IsValid() && idField.CanSet() {
// 				idField.Set(reflect.ValueOf(result.InsertedID))
// 			}
// 			break
// 		}
// 	}

// 	// Hooks

// 	// implement hooks AfterAdd
// 	// call it
// 	if v, ok := data.(interface{ AfterAdd() }); ok {
// 		v.AfterAdd()
// 	}

// 	return nil
// }

func (d *DB) Count(data interface{}, filter ...interface{}) (int64, error) {
	collection := d.GetCollection(data)

	for i := 0; i < len(filter); i += 2 {
		d.filter[filter[i].(string)] = filter[i+1]

	}
	count, err := collection.CountDocuments(context.Background(), d.filter)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (d *DB) OrderBy(data interface{}, order ...interface{}) error {
	// 使用反射获取切片的类型
	dt := reflect.TypeOf(data).Elem()

	// 打印切片的类型
	fmt.Printf("Slice type: %v\n", dt)

	// 获取切片的元素类型
	elemType := dt.Elem()

	// 创建新实例
	newInstance := reflect.New(elemType).Elem()

	// 尝试获取 CollectionName 方法
	method := newInstance.Addr().MethodByName("CollectionName")
	var collectionName string
	if method.IsValid() {
		results := method.Call(nil)
		collectionName = results[0].String()
		fmt.Printf("Collection name from method: %s\n", collectionName)
	} else {
		collectionName = strings.ToLower(newInstance.Type().Name())
		fmt.Printf("Collection name from type: %s\n", collectionName)
	}

	// 获取集合
	collection := d.Cli.Database(d.Cfg.DbName).Collection(collectionName)

	// 构建过滤器
	//d.filter = bson.M{}
	// if len(filter) > 0 {
	// 	for i := 0; i < len(filter); i += 2 {
	// 		key := filter[i].(string)
	// 		value := filter[i+1]
	// 		d.filter[key] = value
	// 	}
	// }
	sort := bson.M{}
	if len(order) > 0 {
		// 构建排序条件

		for i := 0; i < len(order); i += 2 {
			key := order[i].(string)
			value := order[i+1]
			sort[key] = value
		}
	}
	// 设置排序条件

	// 打印过滤器
	fmt.Printf("Filter: %+v\n", d.filter)

	option := options.Find()
	option.SetSort(sort)

	// 执行查询操作
	cursor, err := collection.Find(context.Background(), d.filter, option)
	if err != nil {
		fmt.Printf("Find error: %v\n", err)
		return err
	}
	defer cursor.Close(context.Background())

	// 检查是否有查询结果
	if !cursor.Next(context.Background()) {
		fmt.Println("No documents found")
		return nil
	}

	// 使用反射创建切片
	v := reflect.ValueOf(data).Elem()
	slice := reflect.MakeSlice(dt, 0, 0)

	// 解码查询结果到切片
	for {
		elem := reflect.New(elemType).Elem()
		if err := cursor.Decode(elem.Addr().Interface()); err != nil {
			fmt.Printf("Decode error: %v\n", err)
			return err
		}
		// 打印解码的对象
		//fmt.Printf("Decoded object: %+v\n", elem.Interface())
		slice = reflect.Append(slice, elem)

		// 检查是否有更多结果
		if !cursor.Next(context.Background()) {
			break
		}
	}

	// 将切片赋值给 data
	v.Set(slice)

	// 打印最终的切片
	//fmt.Printf("Final slice: %+v\n", v.Interface())

	return nil
}

func (d *DB) Paginate(data interface{}, page int, pageSize int, order ...interface{}) error {
	// 使用反射获取切片的类型
	dt := reflect.TypeOf(data).Elem()

	// 打印切片的类型
	fmt.Printf("Slice type: %v\n", dt)

	// 获取切片的元素类型
	elemType := dt.Elem()

	// 创建新实例
	newInstance := reflect.New(elemType).Elem()

	// 尝试获取 CollectionName 方法
	method := newInstance.Addr().MethodByName("CollectionName")
	var collectionName string
	if method.IsValid() {
		results := method.Call(nil)
		collectionName = results[0].String()
		fmt.Printf("Collection name from method: %s\n", collectionName)
	} else {
		collectionName = strings.ToLower(newInstance.Type().Name())
		fmt.Printf("Collection name from type: %s\n", collectionName)
	}

	// 获取集合
	collection := d.Cli.Database(d.Cfg.DbName).Collection(collectionName)

	// 构建排序条件
	sort := bson.M{}
	if len(order) > 0 {
		for i := 0; i < len(order); i += 2 {
			key := order[i].(string)
			value := order[i+1]
			sort[key] = value
		}
	}

	// 设置排序条件
	option := options.Find()
	option.SetSort(sort)

	// 设置分页条件
	option.SetSkip(int64((page - 1) * pageSize))
	option.SetLimit(int64(pageSize))

	// 打印过滤器
	fmt.Printf("Filter: %+v\n", d.filter)

	// 执行查询操作
	cursor, err := collection.Find(context.Background(), d.filter, option)
	if err != nil {
		fmt.Printf("Find error: %v\n", err)
		return err
	}
	defer cursor.Close(context.Background())

	// 检查是否有查询结果
	if !cursor.Next(context.Background()) {
		fmt.Println("No documents found")
		return nil
	}

	// 使用反射创建切片
	v := reflect.ValueOf(data).Elem()
	slice := reflect.MakeSlice(dt, 0, 0)

	// 解码查询结果到切片
	for {
		elem := reflect.New(elemType).Elem()
		if err := cursor.Decode(elem.Addr().Interface()); err != nil {
			fmt.Printf("Decode error: %v\n", err)
			return err
		}
		// 打印解码的对象
		//fmt.Printf("Decoded object: %+v\n", elem.Interface())
		slice = reflect.Append(slice, elem)

		// 检查是否有更多结果
		if !cursor.Next(context.Background()) {
			break
		}
	}

	// 将切片赋值给 data
	v.Set(slice)

	// 打印最终的切片
	//fmt.Printf("Final slice: %+v\n", v.Interface())

	return nil
}

package met

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Limit 设置查询结果的数量限制
func (d *DB) Limit(limit int64) *DB {
	d.limit = limit
	return d
}

func (d *DB) FindOne(data interface{}, filter ...interface{}) error {
	collecttion := d.GetCollection(data)

	for i := 0; i < len(filter); i += 2 {
		d.filter[filter[i].(string)] = filter[i+1]

	}

	result := collecttion.FindOne(context.Background(), d.filter)

	if result.Err() != nil {
		return result.Err()
	}

	result.Decode(data)

	return nil
}

func (d *DB) Find(data interface{}, filter ...interface{}) error {
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
	d.filter = bson.M{}
	if len(filter) > 0 {
		for i := 0; i < len(filter); i += 2 {
			key := filter[i].(string)
			value := filter[i+1]
			d.filter[key] = value
		}
	}

	// 设置查询选项
	findOptions := options.Find()
	if d.limit > 0 {
		findOptions.SetLimit(d.limit)
	}

	// 打印过滤器
	fmt.Printf("Filter: %+v\n", d.filter)

	// 执行查询操作
	cursor, err := collection.Find(context.Background(), d.filter, findOptions)
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

	d.limit = 0

	// 打印最终的切片
	//fmt.Printf("Final slice: %+v\n", v.Interface())

	return nil
}

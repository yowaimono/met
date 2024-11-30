package met

import (
	"go.mongodb.org/mongo-driver/bson"
)

// And 添加 $and 条件
func (d *DB) And(conditions ...bson.M) *DB {
	if d.filter == nil {
		d.filter = bson.M{}
	}
	if _, ok := d.filter["$and"]; !ok {
		d.filter["$and"] = []bson.M{}
	}
	d.filter["$and"] = append(d.filter["$and"].([]bson.M), conditions...)
	return d
}

// Or 添加 $or 条件
func (d *DB) Or(conditions ...bson.M) *DB {
	if d.filter == nil {
		d.filter = bson.M{}
	}
	if _, ok := d.filter["$or"]; !ok {
		d.filter["$or"] = []bson.M{}
	}
	d.filter["$or"] = append(d.filter["$or"].([]bson.M), conditions...)
	return d
}

// In 添加 $in 条件
func (d *DB) In(field string, values ...interface{}) *DB {
	if d.filter == nil {
		d.filter = bson.M{}
	}
	if _, ok := d.filter[field]; !ok {
		d.filter[field] = bson.M{}
	}
	d.filter[field].(bson.M)["$in"] = values
	return d
}

// Gt 添加 $gt 条件
func (d *DB) Gt(field string, value interface{}) *DB {
	if d.filter == nil {
		d.filter = bson.M{}
	}
	if _, ok := d.filter[field]; !ok {
		d.filter[field] = bson.M{}
	}
	d.filter[field].(bson.M)["$gt"] = value
	return d
}

// Lt 添加 $lt 条件
func (d *DB) Lt(field string, value interface{}) *DB {
	if d.filter == nil {
		d.filter = bson.M{}
	}
	if _, ok := d.filter[field]; !ok {
		d.filter[field] = bson.M{}
	}
	d.filter[field].(bson.M)["$lt"] = value
	return d
}

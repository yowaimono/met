package met

import "context"

func (d *DB) DeleteOne(data interface{}, filter ...interface{}) error {
	collecttion := d.GetCollection(data)

	for i := 0; i < len(filter); i += 2 {
		d.filter[filter[i].(string)] = filter[i+1]

	}
	_, err := collecttion.DeleteOne(context.Background(), d.filter)
	return err
}

func (d *DB) DeleteMany(data interface{}, filter ...interface{}) error {
	collecttion := d.GetCollection(data)
	for i := 0; i < len(filter); i += 2 {
		d.filter[filter[i].(string)] = filter[i+1]

	}
	_, err := collecttion.DeleteMany(context.Background(), d.filter)
	return err
}

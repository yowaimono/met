[中文文档](https://github.com/yowaimono/met/tree/main/docs/README_zh.md)

# Met - MongoDB Enhanced Toolkit

Met is a Go package that provides a simple and efficient way to interact with MongoDB. It offers a set of methods to perform CRUD operations, manage database connections, and handle common MongoDB operations with ease.

## Features

- **CRUD Operations**: Perform Create, Read, Update, and Delete operations on MongoDB collections.
- **Hook Methods**: Define hooks for `BeforeAdd`, `BeforeUpdate`, `AfterAdd`, and other lifecycle events.
- **Flexible Filtering**: Easily filter and query data using a simple API.
- **Automatic Timestamps**: Automatically manage `CreatedAt` and `UpdatedAt` timestamps.
- **Custom Collection Names**: Define custom collection names for your models.

## Installation

To install Met, use the following command:

```bash
go get github.com/yowaimono/met
```

## Usage

### Initialize the Database Connection

First, initialize the database connection using the `Init` function:

```go
package main

import (
	"fmt"
	"github.com/yowaimono/met"
)

func main() {
	db, err := met.Init("mongodb://admin:xxxx@localhost:27017", &met.Config{
		DbName: "test",
	})

	if err != nil {
		panic(err)
	}

	defer db.Close()
}
```

### Define a Model

Define your model by embedding the `met.Model` struct and implementing the necessary hooks:

```go
type User struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	Name         string             `bson:"name"`
	met.Model    `bson:",inline"`
}

func (u *User) BeforeAdd() {
	fmt.Println("before add")
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
}

func (u *User) CollectionName() string {
	return "user"
}

func (u *User) BeforeUpdate() {
	fmt.Println("before update")
	u.UpdatedAt = time.Now()
}

func (u *User) AfterAdd() {
	fmt.Println("after add")
}
```

### Perform CRUD Operations

#### Add a Document

```go
db.Add(&User{
	Name: "jack",
})
```

#### Update Documents

```go
db.Filter("name", "jack").UpdateMany(&User{}, "name", "mark")
```

#### Delete Documents

```go
db.DeleteMany(&User{}, "name", "mark")
db.DeleteOne(&User{}, "name", "jack")
```

#### Find Documents

```go
db.FindOne(&User{}, "name", "mark")
db.Find(&User{}, "name", "mark")
```

## Contributing

Contributions are welcome! Please feel free to submit a pull request or open an issue if you find any bugs or have suggestions for improvements.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Thanks to the MongoDB team for providing a robust database solution.
- Inspired by various Go packages and frameworks that simplify database interactions.


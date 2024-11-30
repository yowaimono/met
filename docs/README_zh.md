# Met - MongoDB 增强工具包

Met 是一个 Go 包，提供了一种简单高效的方式来与 MongoDB 进行交互。它提供了一组方法来执行 CRUD 操作、管理数据库连接，并轻松处理常见的 MongoDB 操作。

## 功能特性

- **CRUD 操作**：执行创建、读取、更新和删除操作。
- **钩子方法**：定义 `BeforeAdd`、`BeforeUpdate`、`AfterAdd` 等生命周期事件的钩子。
- **灵活过滤**：使用简单的 API 轻松过滤和查询数据。
- **自动时间戳**：自动管理 `CreatedAt` 和 `UpdatedAt` 时间戳。
- **自定义集合名称**：为你的模型定义自定义集合名称。

## 安装

使用以下命令安装 Met：

```bash
go get github.com/yowaimono/met
```

## 使用示例

### 初始化数据库连接

首先，使用 `Init` 函数初始化数据库连接：

```go
package main

import (
	"fmt"
	"github.com/yowaimono/met"
)

func main() {
	db, err := met.Init("mongodb://admin:xxxxxxx@localhost:27017", &met.Config{
		DbName: "test",
	})

	if err != nil {
		panic(err)
	}

	defer db.Close()
}
```

### 定义模型

通过嵌入 `met.Model` 结构体并实现必要的钩子来定义你的模型：

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

### 执行 CRUD 操作

#### 添加文档

```go
db.Add(&User{
	Name: "jack",
})
```

#### 更新文档

```go
db.Filter("name", "jack").UpdateMany(&User{}, "name", "mark")
```

#### 删除文档

```go
db.DeleteMany(&User{}, "name", "mark")
db.DeleteOne(&User{}, "name", "jack")
```

#### 查找文档

```go
db.FindOne(&User{}, "name", "mark")
db.Find(&User{}, "name", "mark")
```

## 贡献

欢迎贡献！如果你发现任何错误或有改进建议，请随时提交拉取请求或开启问题。

## 许可证

本项目采用 MIT 许可证。详情请参阅 [LICENSE](LICENSE) 文件。

## 致谢

- 感谢 MongoDB 团队提供强大的数据库解决方案。
- 受各种 Go 包和框架的启发，简化了数据库交互。


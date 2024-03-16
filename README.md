# README

Gen是由字节跳动无恒实验室与GORM作者联合研发的一个基于GORM的安全ORM框架，主要通过代码生成方式实现GORM代码封装。

Gen框架在GORM框架的基础上提供了以下能力：

- 基于原始SQL语句生成可重用的CRUD API
- 生成不使用interface{}的100%安全的DAO API
- 依据数据库生成遵循GORM约定的结构体Model
- 支持GORM的所有特性

简单来说，使用Gen框架后我们无需手动定义结构体Model，同时Gen框架也能帮我们生成类型安全的CRUD代码。

```bash
docker run -d -p 3306:3306 --privileged=true -e MYSQL_ROOT_PASSWORD=123456 --name mysql mysql:5.7 --character-set-server=utf8mb4 --collation-server=utf8mb4_general_ci
```


```sql
create databse db2;

CREATE TABLE book
(
    `id`     bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
    `title`  varchar(128) NOT NULL COMMENT '书籍名称',
    `author` varchar(128) NOT NULL COMMENT '作者',
    `price`  int NOT NULL DEFAULT '0' COMMENT '价格',
    `publish_date` datetime COMMENT '出版日期',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='书籍表';
```

```bash
go run generate.go
```

```
2024/03/16 12:41:19 find 1 table from db: [book]
2024/03/16 12:41:19 got 5 columns from table <book>
2024/03/16 12:41:19 Start generating code.
2024/03/16 12:41:19 generate model file(table <book> -> {model.Book}): /mnt/c/Users/ruben/db/dal/model/book.gen.go
2024/03/16 12:41:20 generate query file: /mnt/c/Users/ruben/db/dal/query/book.gen.go
2024/03/16 12:41:20 generate query file: /mnt/c/Users/ruben/db/dal/query/gen.go
2024/03/16 12:41:20 Generate code done.
```
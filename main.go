package main

import (
	"context"
	"fmt"
	"time"

	"gen_demo/dal"
	"gen_demo/dal/model"
	"gen_demo/dal/query"
)

// MySQLDSN MySQL data source name
const MySQLDSN = "root:123456@tcp(127.0.0.1:3306)/db2?charset=utf8mb4&parseTime=True"

func init() {
	dal.DB = dal.ConnectDB(MySQLDSN).Debug()
}

func main() {
	// 设置默认DB对象
	query.SetDefault(dal.DB)

	// 创建
	b1 := model.Book{
		Title:       "《七米的Go语言之路》",
		Author:      "七米",
		PublishDate: time.Date(2023, 11, 15, 0, 0, 0, 0, time.UTC),
		Price:       100,
	}
	err := query.Book.WithContext(context.Background()).Create(&b1)
	if err != nil {
		fmt.Printf("create book fail, err:%v\n", err)
		return
	}

	// 更新
	ret, err := query.Book.WithContext(context.Background()).
		Where(query.Book.ID.Eq(1)).
		Update(query.Book.Price, 200)
	if err != nil {
		fmt.Printf("update book fail, err:%v\n", err)
		return
	}
	fmt.Printf("RowsAffected:%v\n", ret.RowsAffected)

	// 查询
	book, err := query.Book.WithContext(context.Background()).First()
	// 也可以使用全局Q对象查询
	//book, err := query.Q.Book.WithContext(context.Background()).First()
	if err != nil {
		fmt.Printf("query book fail, err:%v\n", err)
		return
	}
	fmt.Printf("book:%v\n", book)

	// 删除
	ret, err = query.Book.WithContext(context.Background()).Where(query.Book.ID.Eq(1)).Delete()
	if err != nil {
		fmt.Printf("delete book fail, err:%v\n", err)
		return
	}
	fmt.Printf("RowsAffected:%v\n", ret.RowsAffected)
}

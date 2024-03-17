package main

import (
	"database/sql"
	"errors"
	"fmt"
	"gen_demo/dal/model"
	"reflect"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

const MySQLDSN = "root:123456@tcp(127.0.0.1:3306)/db2?charset=utf8mb4&parseTime=True"

var db *sql.DB

func init() {
	conn, err := sql.Open("mysql", MySQLDSN)
	checkErr(err)
	db = conn
}

func main() {
	rows, err := db.Query("select * from book") // TODO
	checkErr(err)

	var b []model.Book
	err = scanRowsToSlice(rows, &b)
	checkErr(err)
	fmt.Println(b)

	var bs []model.Book
	err = Select(&bs, "")
	checkErr(err)
	fmt.Println(bs)
}

func Select(dest interface{}, where string, args ...interface{}) error {
	name, fields, err := getFields(dest, "json")
	if err != nil {
		return err
	}

	var sqlBuilder strings.Builder
	sqlBuilder.WriteString("select ")
	sqlBuilder.WriteString(strings.Join(fields, ","))
	sqlBuilder.WriteString(" from ")
	sqlBuilder.WriteString(name)
	sqlBuilder.WriteString(" where (1=1 ")
	sqlBuilder.WriteString(where)
	sqlBuilder.WriteString(")")

	fmt.Println(sqlBuilder.String())

	rows, err := db.Query(sqlBuilder.String(), args...)
	if err != nil {
		return err
	}

	return scanRowsToSlice(rows, dest)
}

func getFields(dest interface{}, tag string) (table string, fields []string, err error) {
	destType := reflect.TypeOf(dest)

	if destType.Kind() == reflect.Ptr {
		destType = destType.Elem()
	}

	if destType.Kind() == reflect.Slice {
		destType = destType.Elem()
	}

	if destType.Kind() == reflect.Ptr {
		destType = destType.Elem()
	}

	if destType.Kind() != reflect.Struct {
		return "", nil, errors.New("invalid dest type")
	}

	m := reflect.New(destType).MethodByName("TableName")
	if !m.IsValid() {
		return "", nil, errors.New("invlid TableName method")
	}
	resultValues := m.Call(nil)
	if len(resultValues) > 0 && resultValues[0].Kind() == reflect.String {
		table = resultValues[0].String()
	}

	for i := 0; i < destType.NumField(); i++ {
		field := destType.Field(i).Tag.Get(tag)
		fields = append(fields, field)
	}

	return
}

func scanRowsToSlice(rows *sql.Rows, dest interface{}) error {
	if rows == nil || dest == nil {
		return nil
	}

	destValue := reflect.ValueOf(dest).Elem()
	destType := destValue.Type().Elem()

	if destType.Kind() == reflect.Ptr {
		destType = destType.Elem()
	}

	if destType.Kind() != reflect.Struct {
		return errors.New("invalid dest type")
	}

	defer rows.Close()

	for rows.Next() {
		elem := reflect.New(destType).Interface()

		fields := make([]interface{}, destType.NumField())
		for i := 0; i < destType.NumField(); i++ {
			fields[i] = reflect.ValueOf(elem).Elem().Field(i).Addr().Interface()
		}

		err := rows.Scan(fields...)
		if err != nil {
			return fmt.Errorf("scan failed, err: %w", err)
		}

		destValue.Set(reflect.Append(destValue, reflect.ValueOf(elem).Elem()))
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("rows failed, err: %w", err)
	}

	return nil
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

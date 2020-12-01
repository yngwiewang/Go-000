package main

// sql.ErrNoRows 是 database/sql 包用于表示查询单行数据时是否实际查询到了结果，而并非代码错误。
// 在 gorm 包中用到了这个错误来判断数据库中是否存在某张表，

// func (s mysql) HasTable(tableName string) bool {
//         currentDatabase, tableName := currentDatabaseAndTable(&s, tableName)
//         var name string
//         if err := s.db.QueryRow(fmt.Sprintf("SHOW TABLES FROM `%s` WHERE `Tables_in_%s` = ?", currentDatabase, currentDatabase), tableName).Scan(&name); err != nil {
//                 if err == sql.ErrNoRows {
//                         return false
//                 }
//                 panic(err)
//         } else {
//                 return true
//         }
// }
//
// 我认为对这个错误的处理需要结合应用的场景，如果业务不需要处理未查到的情况需要显示地吞掉这个错误，
// 而如果业务需要对未查到做针对性处理，例如下面在查询一个人的别名时，可能没有这个人，也
// 可能有这个人但他没有别名，如果我们需要区别对待两种情况，就要处理 sql.ErrNoRows，
// 然而这个错误是在 dao 层抛出的，它只是 database/sql 包中的错误，在业务层做处理的时候
// 并不需要知道底层的错误是 sql.ErrNoRows 还是 mongo.ErrNoRows 亦或是 redis.ErrNoRows,
// 此时只需要知道是没查询到结果就可以了，因此应该声明一个统一的未查询到结果的 ErrNotFound，
// 将各种 ErrNoRows 转换成 ErrNotFound 供业务处理用。
// 以下代码假设 getUserAliasName 为 dao 层的查询方法，main 为 service/biz 层的处理函数。

import (
	"database/sql"
	"log"

	"github.com/pkg/errors"

	_ "github.com/go-sql-driver/mysql"
)

var (
	driverName = "mysql"
	dsn        = "root:111111@tcp(192.168.220.102:3306)/hello"
)

// ErrNotFound 代表没有查询到结果
var ErrNotFound = errors.New("not found")

func init() {

}

func getUserAliasName(name string) (string, error) {
	var aliasName string
	db, err := sql.Open(driverName, dsn)
	if err != nil {
		log.Fatalf("open db failed: %v", err)
	}
	err = db.QueryRow("select alias_name from user where name = ?", name).Scan(&aliasName)
	if err != nil {
		if err == sql.ErrNoRows {
			return aliasName, ErrNotFound
		}
	}
	return aliasName, errors.WithStack(err)
}

func main() {
	name := "张三1"
	aliasName, err := getUserAliasName(name)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			log.Printf("没有%s这个人\n", name)
			return
		} else {
			log.Fatalf("查询错误: %+v", err)
		}
	}
	log.Printf("%s的别名为: %v\n", name, aliasName)

}

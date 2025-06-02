package ConnectAndQuery

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func ConnectAndQuery() {
	// MySQL データベースへの接続文字列
	db, err := sql.Open("mysql", "test_user:nakanishi@tcp(localhost:3306)/test_db")

	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// データベースへの接続をテストする
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	//log.Println("Connected to MySQL database!")

	rows, err := db.Query("SELECT * FROM user")
	defer rows.Close()
	if err != nil {
		panic(err.Error())
	}

	for rows.Next() {
		var Name string
		var Age int
		var Sex int
		var Deleted_flag bool
		if err := rows.Scan(&Name, &Age, &Sex, &Deleted_flag); err != nil {
			panic(err.Error())
		}
		fmt.Println(Name, Age, Sex)
	}
}

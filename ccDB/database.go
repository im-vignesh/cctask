package ccDB

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
)

var (
	host     = os.Getenv("DB_HOST")
	port     = os.Getenv("PORT")
	user     = os.Getenv("USER")
	password = os.Getenv("PASSWORD")
	dbname   = os.Getenv("DATABASE")
)

func GetDBConnection() *sql.DB{
	mysqlInfo := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
		user, password,host,port, dbname)
	db, err := sql.Open("mysql", mysqlInfo)
	if err != nil {
		log.Fatal("Couldn't connect to Mysql Server ",err.Error())
		return nil
	}
	fmt.Println("Connection Established...")
	return db
}



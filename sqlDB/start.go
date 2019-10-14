package sqlDB
import(
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/zaddone/ctpSystem/config"
	"os"
)
var(
	DBsql *sql.DB
	DBName = config.Conf.SqlPath

)
func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
func create(){
	var err error
	DBsql, err = sql.Open("sqlite3", DBName)
	checkErr(err)
	fmt.Println("生成数据表")
	sql_table := `
CREATE TABLE IF NOT EXISTS "userinfo" (
   "uid" INTEGER PRIMARY KEY AUTOINCREMENT,
   "username" VARCHAR(64) NULL,
   "departname" VARCHAR(64) NULL,
   "created" TIMESTAMP default (datetime('now', 'localtime'))
);
CREATE TABLE IF NOT EXISTS "userdeatail" (
   "uid" INT(10) NULL,
   "intro" TEXT NULL,
   "profile" TEXT NULL,
   PRIMARY KEY (uid)
);
   `
	 DBsql.Exec(sql_table)
}
func init(){
	_,err := os.Stat(DBName)
	if err != nil {
		create()
	}
}

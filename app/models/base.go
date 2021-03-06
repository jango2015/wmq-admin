package models

import (
	"github.com/astaxie/beego"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/astaxie/beego/orm"
	"io/ioutil"
	"os"
	"database/sql"
	"fmt"
)

//初始化
func Init() {
	databaseType := beego.AppConfig.String("database.type");
	if(databaseType == "mysql") {
		mysqlConn();
	}
	if(databaseType == "sqlite") {
		sqliteConn();
	}

	orm.RegisterModel(new(User), new(Node), new(Notice));

	debug, _:= beego.AppConfig.Bool("database.debug");
	if (debug == true) {
		orm.Debug = true
	}
}

//mysql 连接
func mysqlConn()  {
	host := beego.AppConfig.String("database.mysql.host");
	port := beego.AppConfig.String("database.mysql.port");
	user := beego.AppConfig.String("database.mysql.user");
	password := beego.AppConfig.String("database.mysql.password");
	database := beego.AppConfig.String("database.mysql.name");
	charset := beego.AppConfig.String("database.mysql.charset");
	maxIdle := beego.AppConfig.DefaultInt("database.mysql.maxIdle", 30);
	maxOpenConn := beego.AppConfig.DefaultInt("database.mysql.maxOpenConn", 30);
	dsn := user + ":" + password + "@tcp(" + host + ":" + port + ")/" + database + "?charset=" + charset;

	//驱动
	orm.RegisterDriver("mysql", orm.DRMySQL)
	//注册
	orm.RegisterDataBase("default", "mysql", dsn);
	//数据库最大空闲连接
	orm.SetMaxIdleConns("default", maxIdle);
	//数据库最大连接
	orm.SetMaxOpenConns("default", maxOpenConn);
}

//sqlite 连接
func sqliteConn()  {

	sqlitePath := beego.AppConfig.String("database.sqlite.path");
	_, err := os.Stat(sqlitePath)
	if(os.IsNotExist(err)) {
		createSqlLite(sqlitePath)
	}

	orm.RegisterDriver("sqlite", orm.DRSqlite);
	orm.RegisterDataBase("default", "sqlite3", sqlitePath);
}

//创建sqllite
func createSqlLite(sqlitePath string)  {

	sqliteSqlPath := beego.AppConfig.String("database.sqlite.sql.path");
	db, _ := sql.Open("sqlite3", sqlitePath)

	sqlBytes, _ := ioutil.ReadFile(sqliteSqlPath);

	//创建表
	sqlTable := string(sqlBytes);
	res, err := db.Exec(sqlTable);
	fmt.Println(res);
	if(err != nil) {
		fmt.Println(err.Error());
	}

	db.Close();
}

func TableName(name string) string {
	return beego.AppConfig.String("database.prefix") + name;
}
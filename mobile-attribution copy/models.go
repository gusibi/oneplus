package handler

import (
	"fmt"
	"log"
	"time"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

// Config db config struct
type Config struct {
	DSN         string
	Active      int
	Idle        int
	IdleTimeout time.Duration
}

// NewConfig create new db config
func NewConfig(dsn string, active, idle int, idleTimeout time.Duration) *Config {
	return &Config{DSN: dsn, Active: active, Idle: idle, IdleTimeout: idleTimeout}
}

// NewMysql create new mysql client
func NewMysql(c *Config) (db *sql.DB) {
	db, err := Open(c)
	if err != nil {
		log.Panic(err)
	}
	return
}

// Open connect db
func Open(c *Config) (db *sql.DB, err error) {
	db, err = sql.Open("mysql", c.DSN)
	if err != nil {
		log.Panic("sql.Open() error(%v)", err)
		return nil, err
	}
	db.SetMaxOpenConns(c.Active)
	db.SetMaxIdleConns(c.Idle)
	db.SetConnMaxLifetime(c.IdleTimeout)
	return db, nil
}

// Dao DAO(Data Access Object)一个数据访问接口
type Dao struct {
	c  *Config
	db *sql.DB
}

// NewDao create new dao
func NewDao(c *Config) (d *Dao) {
	d = &Dao{
		c:  c,
		db: NewMysql(c),
	}
	return d
}

// Ping ping db
func (d *Dao) Ping() (err error) {
	return d.db.Ping()
}

// Close close db
func (d *Dao) Close() {
	d.db.Close()
}

func getDSN() string {
	db := EnvGet("C2HDB", "code2html")
	host := EnvGet("C2HHOST", "127.0.0.1")
	port := EnvGet("C2HPORT", "3306")
	user := EnvGet("C2HUSER", "code2html")
	pwd := EnvGet("C2HPWD", "pwd")
	dataSourceName := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8",
		user, pwd, host, port, db)
	log.Println(dataSourceName)
	return dataSourceName
}

var dao *Dao

// InitDB 初始化数据库
func InitDB() {
	dsn := getDSN()
	config := NewConfig(dsn, 20, 10, time.Minute)
	dao = NewDao(config)
	log.Println("init db succeed..")
}

// Code code struct
type Code struct {
	ID       string `db:"id"`
	Code     string `db:"code"`
	Language string `db:"language"`
}

// Get get code by id
func (c *Code) Get(id string) *Code {
	err := dao.db.QueryRow("SELECT code, language FROM code WHERE id=?", id).Scan(&c.Code, &c.Language)
	if err != nil {
		log.Println("query error: ", err)
		return nil
	}
	return c
}

// Create add new code content
func (c *Code) Create() *Code {
	_, err := dao.db.Exec("INSERT INTO code (id, code, language) value(?,?, ?)", c.ID, c.Code, c.Language)
	if err != nil {
		fmt.Println(c.Language, c.ID)
		log.Fatal("create error:", err)
	}
	return c
}

package handler

import (
	"fmt"
	"log"
	"time"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
)

var DB_SOURCE string

// MySQLConfig db config struct
type MySQLConfig struct {
	DSN         string
	Active      int
	Idle        int
	IdleTimeout time.Duration
}

// NewMySQLConfig create new mysql config
func NewMySQLConfig(dsn string, active, idle int, idleTimeout time.Duration) *MySQLConfig {
	return &MySQLConfig{DSN: dsn, Active: active, Idle: idle, IdleTimeout: idleTimeout}
}

// NewMySQL create new mysql client
func NewMySQL(c *MySQLConfig) (db *sql.DB) {
	db, err := OpenMySQL(c)
	if err != nil {
		log.Panic(err)
	}
	return
}

// OpenMySQL connect db
func OpenMySQL(c *MySQLConfig) (db *sql.DB, err error) {
	db, err = sql.Open("mysql", c.DSN)
	if err != nil {
		log.Fatal(fmt.Sprintf("sql.Open() error(%v)", err))
		return nil, err
	}
	db.SetMaxOpenConns(c.Active)
	db.SetMaxIdleConns(c.Idle)
	db.SetConnMaxLifetime(c.IdleTimeout)
	return db, nil
}

// SqliteConfig sqlite config
type SqliteConfig struct {
	DSN string
}

// NewSqliteConfig create new mysql config
func NewSqliteConfig(dsn string) *SqliteConfig {
	return &SqliteConfig{DSN: dsn}
}

// NewSqlite create new sqlite client
func NewSqlite(c *SqliteConfig) (db *sql.DB) {
	db, err := sql.Open("sqlite3", c.DSN)
	if err != nil {
		log.Panic(err)
	}
	return
}

// Dao DAO(Data Access Object)一个数据访问接口
type Dao struct {
	mc *MySQLConfig
	sc *SqliteConfig
	db *sql.DB
}

// NewMySQLDao create new dao
func NewMySQLDao(c *MySQLConfig) (d *Dao) {
	d = &Dao{
		mc: c,
		db: NewMySQL(c),
	}
	return d
}

// NewSqliteDao create new dao
func NewSqliteDao(c *SqliteConfig) (d *Dao) {
	d = &Dao{
		sc: c,
		db: NewSqlite(c),
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

func getMySQLDSN() string {
	db := EnvGet("C2HDB", "code2html")
	host := EnvGet("C2HHOST", "127.0.0.1")
	port := EnvGet("C2HPORT", "3306")
	user := EnvGet("C2HUSER", "code2html")
	pwd := EnvGet("C2HPWD", "pwd")
	dataSourceName := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8",
		user, pwd, host, port, db)
	return dataSourceName
}

func getSqliteDSN() string {
	dbfile := EnvGet("MOBILE_DBFILE", "./phone.db")
	return dbfile
}

var sdao, mdao *Dao

// InitDB 初始化数据库
func InitDB() {
	mdsn := getMySQLDSN()
	mconfig := NewMySQLConfig(mdsn, 20, 10, time.Minute)
	mdao = NewMySQLDao(mconfig)
	sdsn := getSqliteDSN()
	sconfig := NewSqliteConfig(sdsn)
	sdao = NewSqliteDao(sconfig)
	log.Println("init db succeed..")
}

func init() {
	DB_SOURCE = EnvGet("DB_SOURCE", "DB")
	if DB_SOURCE == "DB" {
		InitDB()
	}
}

// Phones phones struct
type Phones struct {
	ID       int `db:"id"`
	Number   int `db:"number"`
	Type     int `db:"type"`
	RegionID int `db:"region_id"`
}

// Create add new phone
func (p *Phones) Create() *Phones {
	_, err := mdao.db.Exec("INSERT INTO phones (number, type, region_id) value(?,?,?)", p.Number, p.Type, p.RegionID)
	if err != nil {
		// fmt.Println(p.RegionID, p.Number)
		log.Fatal("create error:", err)
	}
	return p
}

// Fetch 遍历全部数据
func (p *Phones) Fetch() (*sql.Rows, error) {
	rows, err := sdao.db.Query("SELECT * FROM phones ORDER BY number")
	return rows, err
}

// Phone 数据库连表查询结果
type Phone struct {
	Mobile   string `json:"mobile"`
	Type     int    `db:"type"` // 运营商类型
	Province string `json:"province"`
	City     string `json:"city"`
	ZipCode  string `json:"zipCode"`
	AreaCode string `json:"areaCode"`
}

// Get get code by id
func (p *Phone) Get(number int) *Phone {
	sql := `SELECT number, type, province, city, zip_code, area_code FROM phones 
			LEFT JOIN regions ON phones.region_id=regions.id
			WHERE number=?`
	// fmt.Println(sdao.db)
	err := sdao.db.QueryRow(sql, number).Scan(&p.Mobile, &p.Type, &p.Province, &p.City, &p.ZipCode, &p.AreaCode)
	if err != nil {
		// log.Println("query error: ", err)
		return nil
	}
	return p
}

// Region Region struct
type Region struct {
	ID       int    `db:"id"`
	Province string `json:"province"`
	City     string `json:"city"`
	ZipCode  string `json:"zip_code"`
	AreaCode string `json:"area_code"`
}

// Fetch 遍历全部数据
func (r *Region) Fetch() (*sql.Rows, error) {
	rows, err := sdao.db.Query("SELECT * FROM regions ORDER BY id")
	return rows, err
}

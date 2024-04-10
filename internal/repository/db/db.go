package db

import (
	"time"

	"go-template/pkg/zlog"

	"github.com/glebarez/sqlite"
	"github.com/pkg/errors"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var _ Repo = (*dbRepo)(nil)

type Repo interface {
	i()
	GetDb() *gorm.DB
	DbClose() error
}

type dbRepo struct {
	c  *Config
	Db *gorm.DB
}

func New(c *Config) (Repo, error) {
	if !c.Validate() {
		return nil, errors.New("db config is invalid")
	}
	zlog.Debug("db_new", zlog.Any("config", c))
	db, err := dbConnect(c)
	if err != nil {
		return nil, err
	}
	d := &dbRepo{
		Db: db,
	}
	if err := d.CreateTable(); err != nil {
		return nil, errors.Wrap(err, "create table failed")
	}
	// 只有特殊情况下需要手动添加触发器
	if err := d.CreateTrigger(); err != nil {
		return nil, errors.Wrap(err, "create trigger failed")
	}
	return d, err
}

func (d *dbRepo) CreateTable() error {
	tables := []interface{}{}
	return d.Db.Transaction(func(tx *gorm.DB) error {
		for _, table := range tables {
			if !tx.Migrator().HasTable(table) {
				if err := tx.Migrator().CreateTable(table); err != nil {
					return err
				}
			}
		}
		return nil
	})
}

func (d *dbRepo) CreateTrigger() error {
	// 表示使用的是sqlite，不需要执行以下sql
	if d.c.DbType == SQLite {
		return nil
	}
	tableNames := []string{}
	// 为了兼容sqlite3，所以使用sql语句设置表对应的updated_at字段设置为当前时间
	return d.Db.Transaction(func(tx *gorm.DB) error {
		for _, tableName := range tableNames {
			sql := "alter TABLE " + tableName + " change updated_at updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;"
			if err := tx.Exec(sql).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (d *dbRepo) GetDb() *gorm.DB {
	return d.Db
}

func (d *dbRepo) DbClose() error {
	sqlDB, err := d.Db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

func (d *dbRepo) i() {}

func dbConnect(c *Config) (*gorm.DB, error) {
	var dialect gorm.Dialector
	if c.DbType == MySQL {
		dialect = mysql.Open(c.Dsn)
	} else if c.DbType == SQLite {
		// note: 只有测试的时候才会使用sqlite
		dialect = sqlite.Open(c.Dsn)
	}
	db, err := gorm.Open(dialect, &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		Logger: NewGormLogger(time.Duration(c.SlowSqlThreshold) * time.Second),
	})
	if err != nil {
		return nil, errors.Wrap(err, "db connection failed")
	}

	db.Set("gorm:table_options", "CHARSET=utf8mb4")
	db.Set("gorm:table_options", "AUTO_INCREMENT=2")

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// 设置连接池 用于设置最大打开的连接数，默认值为0表示不限制.设置最大的连接数，可以避免并发太高导致连接mysql出现too many connections的错误。
	sqlDB.SetMaxOpenConns(c.MaxOpenConn)

	// 设置最大连接数 用于设置闲置的连接数.设置闲置的连接数则当开启的一个连接使用完成后可以放在池里等候下一次使用。
	sqlDB.SetMaxIdleConns(c.MaxIdleConn)

	// 设置最大连接超时
	sqlDB.SetConnMaxLifetime(time.Minute * time.Duration(c.ConnMaxLifeTime))

	// 使用插件
	if err = db.Use(&TracePlugin{}); err != nil {
		return nil, err
	}

	return db, nil
}

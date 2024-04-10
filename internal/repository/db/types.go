package db

type DatabaseType string

const (
	MySQL  DatabaseType = "mysql"
	SQLite DatabaseType = "sqlite"
)

func (d DatabaseType) Validate() bool {
	return d == MySQL || d == SQLite
}

type Config struct {
	DbType           DatabaseType `json:"db_type"`
	Dsn              string       `json:"dsn"`
	MaxOpenConn      int          `json:"max_open_conn"`
	MaxIdleConn      int          `json:"max_idle_conn"`
	SlowSqlThreshold int          `json:"slow_sql_threshold"`
	ConnMaxLifeTime  int          `json:"conn_max_life_time"`
}

func (c *Config) Validate() bool {
	if !c.DbType.Validate() ||
		c.MaxOpenConn <= 0 ||
		c.MaxIdleConn <= 0 ||
		c.ConnMaxLifeTime <= 0 {
		return false
	}
	return true
}

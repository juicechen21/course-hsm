package conf

type AppConf struct {
	MysqlConf `ini:"mysql"`
	RedisConf `ini:"redis"`
}

type MysqlConf struct {
	Dsn          string `ini:"dsn"`
	MaxIdleConns int    `ini:"maxIdleConns"`
	MaxOpenConns int    `ini:"maxOpenConns"`
}

type RedisConf struct {
	Address  string `ini:"address"`
	Password string `ini:"password"`
	Database int    `ini:"database"`
	PoolSize int    `ini:"poolSize"`
}
package model

import (
	"context"
	"fmt"
	"hsm/service/uavdata/conf"
	"time"

	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"gopkg.in/ini.v1"
)

var (
	cfg     = new(conf.AppConf)
	db      *sqlx.DB
	redisdb *redis.Client
)

func init() {
	// 1.加载配置文件
	err := ini.MapTo(cfg, "./conf/config.ini")
	if err != nil {
		fmt.Printf("load ini failed,err:%v\n", err)
		return
	}
	fmt.Println("load conf file success...")

	db, err = sqlx.Connect("mysql", cfg.MysqlConf.Dsn)
	if err != nil {
		fmt.Println("connect uav mysql err:", err)
		return
	}
	//defer db.Close()            // 注意这行代码要写在上面err判断的下面(关闭连接)
	db.SetMaxOpenConns(cfg.MysqlConf.MaxOpenConns) // 最大连接池
	db.SetMaxIdleConns(cfg.MysqlConf.MaxIdleConns) // 空闲时间最大连接池
	fmt.Println("connect uav mysql success...")

	redisdb = redis.NewClient(&redis.Options{
		Addr:     cfg.RedisConf.Address,
		Password: cfg.RedisConf.Password, // no password set
		DB:       cfg.RedisConf.Database, // use default DB
		PoolSize: cfg.RedisConf.PoolSize, // 连接池大小
	})
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err = redisdb.Ping(ctx).Result()
	if err != nil {
		fmt.Println("connect redis err:", err)
		return
	}
	fmt.Println("connect redis success.")
}

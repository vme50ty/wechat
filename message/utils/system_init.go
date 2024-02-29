package utils

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	DB  *gorm.DB
	Red *redis.Client
)

func InitConfig() {
	viper.SetConfigName("app")
	viper.AddConfigPath("message/config")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("config app inited")
}
func InitMySQL() {
	logs := logrus.WithFields(logrus.Fields{
		"func": "init DB",
	})
	dsn := viper.GetString("mysql.dns")
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Info,
			Colorful:      true,
		},
	)
	var err error
	if DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: newLogger}); err != nil {
		logs.Debugf("mysql inited failed: %v", err)
	} else {
		logs.Debugf("mysql inited")
	}
}

func InitRedis() {
	Red = redis.NewClient(&redis.Options{
		Addr:         viper.GetString("redis.addr"),
		Password:     viper.GetString("redis.password"),
		DB:           viper.GetInt("redis.DB"),
		PoolSize:     viper.GetInt("redis.poolSize"),
		MinIdleConns: viper.GetInt("redis.minIdleConn"),
	})
}

const (
	PublishKey = "websocket"
)

// 发布消息到redis
func Publish(ctx context.Context, channel string, msg string, log *logrus.Entry) error {
	err := Red.Publish(ctx, channel, msg).Err()
	log.Errorf("pubish error:%v", err)
	return err
}

func Subscribe(ctx context.Context, channel string, log *logrus.Entry) (string, error) {
	sub := Red.Subscribe(ctx, channel)
	msg, err := sub.ReceiveMessage(ctx)
	log.Debugf("Subscribe:%s", msg.Payload)
	return msg.Payload, err
}

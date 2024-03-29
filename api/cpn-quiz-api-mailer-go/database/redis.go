package database

import (
	"context"
	"cpn-quiz-api-mailer-go/logger"
	"fmt"

	"github.com/redis/go-redis/v9"
	config "github.com/spf13/viper"
)

var ctx = context.Background()

type RedisDatabase struct {
	Instance *redis.Client
	Log      logger.PatternLogger
}

func (rdb *RedisDatabase) GetConnectionRedisDB() *RedisDatabase {
	if rdb.Instance == nil {
		rdb.Instance = rdb.initRedisConnect()
	}
	return rdb
}

// =>Connection
func (rdb *RedisDatabase) initRedisConnect() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     config.GetString("cpn.quiz.redis.client.addr"),
		Password: config.GetString("cpn.quiz.redis.client.password"),
		DB:       0,
	})
}

func (rdb RedisDatabase) IsConnected() error {
	rdb = *rdb.GetConnectionRedisDB()
	//=>ตรวจสอบการเชื่อมต่อ
	return rdb.Instance.Ping(ctx).Err()
}

// =>Enqueue เพิ่มข้อมูลลงในคิว
func (rdb RedisDatabase) Enqueue(queueName string, value string) error {
	rdb = *rdb.GetConnectionRedisDB()
	return rdb.Instance.LPush(ctx, queueName, value).Err()
}

// =>Dequeue ดึงข้อมูลจากคิว
func (rdb RedisDatabase) Dequeue(queueName string) (string, error) {
	rdb = *rdb.GetConnectionRedisDB()
	// ดึงข้อมูลจากต้นคิวของคิว
	value, err := rdb.Instance.RPop(ctx, queueName).Result()
	if err == redis.Nil {
		// คิวว่าง
		return "", fmt.Errorf("queue is empty")
	} else if err != nil {
		// เกิดข้อผิดพลาดอื่นๆ
		return "", err
	}
	return value, nil
}

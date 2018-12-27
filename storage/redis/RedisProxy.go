package gothicredis

import (
	"time"
)

/**
 * @desc RedisProxy
 * @author zhaojiangwei
 * @date 2018-12-24 16:45
 */

import (
	"errors"
	"github.com/garyburd/redigo/redis"
)

type RedisClient struct{
	pool *redis.Pool
}

type RedisConn struct {
	conn redis.Conn
}

type RedisPoolConfig struct{
	Host 				string
	Port				string
	Password 			string

	MaxIdle 			int
	IdleTimeout 		int
	ConnTimeout 		int
	ReadTimeout			int
	WriteTimeout 		int
}

func NewRedisClient(config *RedisPoolConfig) *RedisClient {
	if config == nil{
		return nil
	}

	pool := &redis.Pool{
		MaxIdle:     config.MaxIdle,
		IdleTimeout: time.Duration(config.IdleTimeout) * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", config.Host + ":" + config.Port,
				redis.DialConnectTimeout(time.Duration(config.ConnTimeout)*time.Millisecond),
				redis.DialReadTimeout(time.Duration(config.ReadTimeout)*time.Millisecond),
				redis.DialWriteTimeout(time.Duration(config.WriteTimeout)*time.Millisecond),
			)

			if err != nil {
				return nil, err
			}

			if config.Password != "" {
				if _, err := c.Do("AUTH", config.Password); err != nil {
					c.Close()
					return nil, err
				}
			}

			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}

			_, err := c.Do("PING")
			return err
		},
	}

	return &RedisClient{pool: pool}
}

func (this *RedisClient) SelectDb(db int) error{
	if _, err := this.Do("SELECT", db); err != nil {
		return err
	}

	return nil
}

func (this *RedisClient) Do(command string, args ...interface{}) (reply interface{}, err error) {
	conn := this.pool.Get()
	defer conn.Close()

	return conn.Do(command, args...)
}

//###################################################Pipeline##################################################

//用完后需要手动调用rconn.Close()
func (this *RedisClient) StartPipe() (rconn *RedisConn, err error) {
	conn := this.pool.Get()
	if conn == nil {
		err = errors.New("get redis connection")
		return
	}

	rconn = &RedisConn{conn}
	return
}

func (this *RedisConn) Send(cmd string, args ...interface{}) error {
	if this.conn == nil {
		return errors.New("empty conn")
	}

	return this.conn.Send(cmd, args...)
}

func (this *RedisConn) Exec() (replay interface{}, err error) {
	if this.conn == nil {
		err = errors.New("empty conn")
		return
	}

	return this.conn.Do("")
}

func (this *RedisConn) Close() {
	if this.conn != nil {
		this.conn.Close()
		this.conn = nil
	}
}

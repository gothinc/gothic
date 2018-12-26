package gothic_redis

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
	"github.com/gothinc/gothic"
	logger "github.com/gothinc/gothic/logger"
)

type ConnType int

const (
	MASTER ConnType = 1
	SLAVE  ConnType = 2
)

//var RedisClient = &RedisProxy{}

type RedisProxy struct {
	master *redis.Pool
	slave  *redis.Pool
}

type RedisConn struct {
	conn redis.Conn
}

func NewPool(server, password string, max_idle int, max_idle_timeout int, conn_timeout, read_timeout, write_timeout int) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     max_idle,
		IdleTimeout: time.Duration(max_idle_timeout) * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server, redis.DialConnectTimeout(time.Duration(conn_timeout)*time.Millisecond),
				redis.DialReadTimeout(time.Duration(read_timeout)*time.Millisecond), redis.DialWriteTimeout(time.Duration(write_timeout)*time.Millisecond))
			if err != nil {
				log_data := map[string]interface{}{"service": "redis", "server_info": server, "errmsg": err, "type": "Redis NewConn Fail"}
				gothic.Logger.Error(log_data)
				return nil, err
			}

			if password != "" {
				if _, err := c.Do("AUTH", password); err != nil {
					log_data := map[string]interface{}{"service": "redis", "server_info": server, "password": password, "errmsg": err, "type": "Redis AUTH Fail"}
					gothic.Logger.Error(log_data)

					c.Close()
					return nil, err
				}
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			if err != nil {
				log_data := map[string]interface{}{"service": "redis", "server_info": server, "errmsg": err, "type": "Redis Test Conn Error"}
				gothic.Logger.Warn(log_data)
			}

			return err
		},
	}
}

func (this *RedisProxy) GetMaster() redis.Conn {
	if this.master != nil {
		conn := this.master.Get()
		return conn
	}

	return nil
}

func (this *RedisProxy) SetMaster(master *redis.Pool) {
	this.master = master
}

func (this *RedisProxy) SetSlave(slave *redis.Pool) {
	this.slave = slave
}

func (this *RedisProxy) Release(conn redis.Conn) {
	if conn != nil {
		conn.Close()
	}
}

func (this *RedisProxy) GetSlave() redis.Conn {
	if this.slave != nil {
		conn := this.slave.Get()
		return conn
	}

	return nil
}

func (this *RedisProxy) Do(ctype ConnType, command string, args ...interface{}) (reply interface{}, err error) {
	var conn redis.Conn
	if ctype == MASTER {
		conn = this.GetMaster()
	} else if ctype == SLAVE {
		conn = this.GetSlave()
	}

	if conn == nil {
		errinfo := "connect redis exception"
		log_data := map[string]interface{}{"service": "redis", "type": ctype, "errmsg": errinfo, "command": command, "args": args}
		gothic.Logger.Error(log_data)
		err = errors.New(errinfo)
		return
	}
	defer conn.Close()

	reply, err = conn.Do(command, args...)
	if err == nil && gothic.Logger.GetLogLevel() <= logger.LevelDebug {
		if command == "HGETALL" {
			logStr, _ := redis.StringMap(reply, nil)
			gothic.Logger.Debug("redis do ", command, args, logStr, err)
		} else if command == "SINTER" {
			logStr, _ := redis.Strings(reply, nil)
			gothic.Logger.Debug("redis do ", command, args, logStr, err)
		} else {
			logStr, _ := redis.String(reply, nil)
			gothic.Logger.Debug("redis do ", command, args, logStr, err)
		}
	}

	if err != nil {
		errinfo := "redis operate exception"
		log_data := map[string]interface{}{"service": "redis", "type": ctype, "errmsg": errinfo, "command": command, "rawerrinfo": err, "args": args}
		gothic.Logger.Error(log_data)
	}

	return
}

//用完后需要手动调用rconn.Close()
func (this *RedisProxy) StartMasterPipe() (rconn *RedisConn, err error) {
	conn := this.GetMaster()
	if conn == nil {
		errinfo := "redis connection exception"
		log_data := map[string]interface{}{"service": "redis", "type": "master", "errmsg": errinfo}
		gothic.Logger.Error(log_data)
		err = errors.New(errinfo)

		return
	}

	rconn = &RedisConn{conn}
	return
}

//用完后需要手动调用rconn.Close()
func (this *RedisProxy) StartSlavePipe() (rconn *RedisConn, err error) {
	conn := this.GetSlave()
	if conn == nil {
		errinfo := "redis connection exception"
		log_data := map[string]interface{}{"service": "redis", "type": "slave", "errmsg": errinfo}
		gothic.Logger.Error(log_data)
		err = errors.New(errinfo)

		return
	}

	rconn = &RedisConn{conn}
	return
}

//###################################################RedisConn##################################################
func (this *RedisConn) Send(cmd string, args ...interface{}) error {
	if this.conn == nil {
		return errors.New("empty conn")
	}

	err := this.conn.Send(cmd, args...)
	if err != nil {
		errinfo := "redis operate exception"
		log_data := map[string]interface{}{"service": "redis", "type": "Redis Trans Send Fail", "errmsg": errinfo, "command": cmd, "args": args}
		gothic.Logger.Error(log_data)
	}

	return err
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

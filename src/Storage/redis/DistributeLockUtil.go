package redis

import (
	"errors"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/satori/go.uuid"
	"log"

)

var (
	delScript = redis.NewScript(1, `
	if redis.call("get", KEYS[1]) == ARGV[1] then 
		return redis.call("del", KEYS[1]) 
	else 
		return 0
	end`)


)

const(
	LOCK_SUCCESS = "OK"
	SET_IF_NOT_EXIST = "NX"
	SET_WITH_EX = "EX"
	RELEASE_SUCCESS = int64(1)

	GET_LOCK_FAIL = 1
	UNLOCK_FAIL = 2
)


type Lock struct {
	resource string
	token string
	conn redis.Conn
	timeout int
}

func (lock *Lock) tryLock()(err error) {
	lockReply, err := lock.conn.Do("SET", lock.resource, lock.token, SET_WITH_EX, lock.timeout, SET_IF_NOT_EXIST)
	if err != nil {
		return errors.New("redis fail")
	}
	if lockReply == LOCK_SUCCESS {
		return nil
	} else {
		return errors.New("lock fail")
	}


}

func (lock *Lock) Unlock() (err error) {

	_, err = delScript.Do(lock.conn, lock.resource, lock.token)
	if err != nil {
		fmt.Println("unlock failed")
		fmt.Println("unlock err: ", err)
	}
	//_, err = lock.conn.Do("del", lock.resource)
	return
}

func (lock *Lock) key() string {
	//return fmt.Sprintf("redislock:%s", lock.resource)
	res := lock.resource
	return res
}

func (lock *Lock) AddTimeout(ex_time int64) (ok bool, err error) {
	ttl_time, err := redis.Int64(lock.conn.Do("TTL", lock.resource))

	if err != nil {
		log.Fatal("redis get failed: ", err)
	}

	if ttl_time > 0 {
		_, err := redis.String(lock.conn.Do("SET", lock.resource, lock.token, SET_WITH_EX, int(ttl_time + ex_time)))
		if err == redis.ErrNil {
			return false, nil

		}
		if err != nil {

			return false, err
		}
	}

	return true, nil
}

func TryLock(conn redis.Conn, resource string, token string, DefaultTimeout int) (lock *Lock, err error) {
	return TryLockWithTimeout(conn, resource, token, DefaultTimeout)
}

func TryLockWithTimeout(conn redis.Conn, resource string, token string, timeout int) (lock *Lock, err error) {
	lock = &Lock{resource, token, conn, timeout}

	err = lock.tryLock()
	if err != nil {
		lock = nil
		fmt.Println("try lock failed")
		fmt.Println("err: ", err)
	}
	return
}

type Function interface {
	Execute(conn redis.Conn, redisKey string, v... interface{})(m interface{} ,err error)
}

func (lock *Lock) DoWithLock(lockKey string, expire int, conn redis.Conn, function Function, redisKey string, v... interface{}) (m interface{}, err error) {

	requestId := uuid.Must(uuid.NewV4())
	lock, err = TryLock(conn, lockKey, fmt.Sprintf("%d", requestId), int(1))
	defer lock.Unlock()
	if err != nil {
		log.Fatal("Error while getting lock")
		m = nil
		return
	}

	fmt.Println("Distribute lock key: ", lockKey)
	fmt.Println("Business key: ", redisKey)
	res, err := function.Execute(conn, redisKey, v)
	fmt.Println("res: ", res)


	m = res
	return

}


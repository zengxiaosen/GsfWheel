package main

import (
	"fmt"
	"log"
	redigo "github.com/garyburd/redigo/redis"

	c "gsf/src/Storage/redis"
)

func main() {


	//pipeline
	//testPipeLine()

	//get
	testGet()

	//hmget
	testHmget()


}



func testHmget() {


	rds := c.New()
	conn, err := c.GetRedisConn(rds)
	defer conn.Close()

	if err != nil {
		fmt.Println("get redis conn failed")
	}


	//hmset
	_, err = conn.Do("hmset", "key", "hashkey1", "hashvalue1", "hashkey2", "hashvalue2", "hashkey3", "hashvalue3")
	if err != nil {
		log.Printf("hmset error", err)
		return
	}





	//hmget命令的参数，包括key
	args := make([]interface{}, 0)
	args = append(args, "key")
	args = append(args, "hashkey1")
	args = append(args, "hashkey2")
	args = append(args, "hashkey3")

	redisResult, err := conn.Do("hmget", args...)

	if err != nil {
		log.Printf("hmget command error, err %v\n", err)
	}

	//将redisResult转化称一个slice
	lists, err := redigo.Values(redisResult, nil)

	if err != nil {
		if err == redigo.ErrNil {
			log.Println("redisResult is an empty slice")
		} else {
			log.Printf("transform redisRsult to slice error, err : %v\n", err)
		}
	}


	for _, v := range lists {
		fmt.Printf("%s ", v.([]byte))
	}


}




func testGet() {

	rds := c.New()
	conn, err := c.GetRedisConn(rds)
	defer conn.Close()

	if err != nil {
		fmt.Println("get redis conn failed")
	}

	if redisResult, err := conn.Do("get", "k1"); err != nil {
		fmt.Printf("Get data from redis error, err %v", err)
		return
	} else {
		//使用redigo的转换函数进行类型转换，将结果从interface{}转换称string
		strResult, err := redigo.String(redisResult, nil)

		if err != nil {
			fmt.Printf("cast data type error, err : %v", err)
			return
		}
		fmt.Printf("redis get result : %v", strResult)
	}





}



func testPipeLine() {
	rds := c.New()
	conn, err := c.GetRedisConn(rds)
	defer conn.Close()

	if err != nil {
		fmt.Println("get redis conn failed")
	}

	rds.Add("set", "k1", "v1")
	rds.Add("set", "k2", "v2")
	rds.Add("set", "k3", "v3")


	values, err := rds.Pipeline(conn, "get", "username")
	log.Println(values)

	switch values.(type) {
	case []interface{}:
		fmt.Println("correct return type")
	default:
		fmt.Println("wrong return type")

	}




	rds.Add("get", "k1")
	rds.Add("get", "k2")
	rds.Add("get", "k3")

	values, err = rds.Pipeline(conn, "get", "username")
	log.Println(values)

	switch values.(type) {
	case []interface{}:
		fmt.Println("correct return type")

		for _, v := range values.([]interface{}) {
			b8 := v.([]uint8)
			ba := []byte{}
			for _, b := range b8 {
				ba = append(ba, byte(b))
			}
			bs := string(ba)
			log.Println(string(bs))

		}

	default:
		fmt.Println("wrong return type")

	}



	//

	rds.Add("set", "i1", 1)
	rds.Add("set", "i2", 2)
	rds.Add("set", "i3", 3)


	values, err = rds.Pipeline(conn, "get", "username")
	log.Println(values)

	switch values.(type) {
	case []interface{}:
		fmt.Println("correct return type")
	default:
		fmt.Println("wrong return type")

	}

	rds.Add("get", "i1")
	rds.Add("get", "i2")
	rds.Add("get", "i3")

	values, err = rds.Pipeline(conn, "get", "username")
	log.Println(values)

	switch values.(type) {
	case []interface{}:
		fmt.Println("correct return type")

		for _, v := range values.([]interface{}) {
			b8 := v.([]uint8)
			ba := []byte{}
			for _, b := range b8 {
				ba = append(ba, byte(b))
			}
			bs := string(ba)
			log.Println(string(bs))

		}

	default:
		fmt.Println("wrong return type")

	}
}
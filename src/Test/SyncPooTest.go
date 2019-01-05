package main

import (
	"log"
	"sync"
)
/*
sync.Pool的定位不是做类似连接池的东西，它的用途是减少gc
 */
func main() {

	poolTest1()


}



func poolTest1() {
	var pool = &sync.Pool{New: func() interface{} {
		return "default object"
	}}
	val := "Hello World!"
	//放入
	pool.Put(val)
	//取出
	log.Println(pool.Get())
	//再取就没有了，会自动调用NEW
	log.Println(pool.Get())
}

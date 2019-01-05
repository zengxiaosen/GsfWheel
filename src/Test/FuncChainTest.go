package main

import (
	"log"
	"runtime"
)

func main() {
	testAop()
}

func testAop() {
	getFuncName()
}

func getFuncName() {
	pc,file,line,ok := runtime.Caller(1)
	log.Println(pc)
	log.Println(file)
	log.Println(line)
	log.Println(ok)
	f := runtime.FuncForPC(pc)
	log.Println(f.Name())
}

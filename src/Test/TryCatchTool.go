package main

import (
	"errors"
	"fmt"
	"github.com/mongodb/mongo-go-driver/bson"
	"gsf/src/tryCatch"
	"reflect"
)

func main() {

	fun4()
}

func fun4() {
	usrinfo := bson.M{"a": "a", "b": "b"}
	c := usrinfo["c"]
	if c == nil {
		fmt.Println("...")
	}
}

func fun3() {
	err := errors.New("手机号不存在")
	fmt.Println(err.Error())
}

func fun2() {
	var num float64 = 1.2345
	tc := tryCatch.TryCatch{}
	pointer := reflect.ValueOf(num)
	var result interface{}
	count := 0
	for i := 1; i <= 10; i++ {
		excepFlag := false
		tc.Try(func() {
			convertPointer := pointer.Interface().(int)
			result = convertPointer

		}).Catch(tryCatch.MyError{}, func(err error) {
			println("catch MyError")
			excepFlag = true
		}).CatchAll(func(err error) {
			//其他运行时异常，比如 &runtime.TypeAssertionError{}
			println("catch error")
			excepFlag = true
		}).Finally(func() {
			println("finally do something")
		})

		if excepFlag == true {
			fmt.Println("出现运行时异常")
			continue

		}

		count += i
	}

	fmt.Println("count: ", count)
}

func fun1() {
	tc := tryCatch.TryCatch{}
	tc.Try(func() {
		println("do something buggy")

		merr := tryCatch.MyError{}
		merr.Code = 1
		merr.Desc = "自定义异常"
		panic(merr)

	}).Catch(tryCatch.MyError{}, func(err error) {
		println("catch MyError: ", err)
		code := err.(tryCatch.MyError).Code
		desc := err.(tryCatch.MyError).Desc
		println("code: ", code)
		println("desc: ", desc)
	}).CatchAll(func(err error) {
		println("catch error")
	}).Finally(func() {
		println("finally do something")
	})

	println("done")
}

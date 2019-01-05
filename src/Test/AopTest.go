package main

import (
	"encoding/json"
	"fmt"
	"grpc-go-demo/src/aop"
	"reflect"
)

type LogAspect struct {

}

func (l LogAspect) LogBefore() {
	fmt.Println("log before")
}

func (l LogAspect) LogAfter() {
	fmt.Println("log after")
}


type AopRefStruct struct {
	la LogAspect
}

func (ap AopRefStruct) DoAop(s string) {
	fmt.Println(s)
}
func main() {

	//非反射的aop
	AopTest1()

	//纯反射的aop
	//AopTest2()

	//半反射的aop，golang运行时绑定不了方法中传方法，暂时解决不了
	//AopTest3()
}



func AopTest2() {
	aopRef := AopRefStruct{LogAspect{}}

	//得到反射类型的对象
	getValue := reflect.ValueOf(aopRef)

	pointer := reflect.ValueOf(&aopRef.la)
	newValue := pointer.Elem().Interface().(LogAspect)
	newValue.LogBefore()

	//反射调用有参的方法,如果方法名错误，会抛出错误
	methodValue := getValue.MethodByName("DoAop")
	args := []reflect.Value{reflect.ValueOf("yaoming")}
	methodValue.Call(args)


	pointer = reflect.ValueOf(&aopRef.la)
	newValue = pointer.Elem().Interface().(LogAspect)
	newValue.LogAfter()
}

func AopTest3() {
	aopUtil := aop.AopStructure{}
	getValue := reflect.ValueOf(aopUtil)
	//反射调用有参的方法,如果方法名错误，会抛出错误
	methodValue := getValue.MethodByName("DoWithAop")
	args := []reflect.Value{reflect.ValueOf(TestAopFunc{}), reflect.ValueOf(33.3)}
	methodValue.Call(args)
}

func AopTest1() {

	aopUtil := &aop.AopStructure{}
	aopUtil.DoWithAop(TestAopFunc{}, 33.3)
}

type TestAopFunc struct {

}

//传入一个uin
func (testAopFunc TestAopFunc) Execute(v... interface{}) (m interface{}, err error) {
	uinStr := getUinStr(v)
	fmt.Println("uinStr: ", uinStr)
	return uinStr, nil
}

func getUinStr(v []interface{}) string {
	uin := v[0]
	juin, _ := json.Marshal(uin)
	uinStr := string(juin)
	uinStrLen := len(uinStr)
	uinStr = uinStr[1:uinStrLen-1]
	return uinStr
}


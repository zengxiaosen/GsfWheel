package main


import (
	"fmt"
	"reflect"
)

/*

 */

type User struct {
	Id int
	Name string
	Age int
}

func (u User) ReflectCallFunc() {
	fmt.Println("do ReflectCallFunc ")
}

func main() {
	//test1()
	//test2()
	test3()
}

//通过反射，进行方法的调用
func test3() {
	user := User{1, "andy", 25}
	//得到反射类型的对象
	getValue := reflect.ValueOf(user)

	//反射调用有参的方法,如果方法名错误，会抛出错误
	methodValue := getValue.MethodByName("ReflectCallFuncHasArgs")
	args := []reflect.Value{reflect.ValueOf("yaoming"), reflect.ValueOf(40)}
	methodValue.Call(args)

	//反射调用无参的方法
	methodValue = getValue.MethodByName("ReflectCallFuncNoArgs")
	args = make([]reflect.Value, 0)
	methodValue.Call(args)


}

func (u User) ReflectCallFuncHasArgs(name string, age int) {
	fmt.Println("ReflectCallFuncHasArgs name: ", name, ", age: ", age, " and original User.Name: ", u.Name)
}

func (u User) ReflectCallFuncNoArgs() {
	fmt.Println("ReflectCallFuncNoArgs")
}


//通过反射，修改值
func test2() {
	var num float64 = 1.2345
	fmt.Println("old value of pointer: ", num)

	//注意，参数必须时指针才能修改其值
	pointer := reflect.ValueOf(&num)
	newValue := pointer.Elem()
	fmt.Println("new Value: ", newValue)
	fmt.Println("type of pointer: ", newValue.Type())
	fmt.Println("settability of pointer: ", newValue.CanSet())

	//重新赋值
	newValue.SetFloat(77)
	fmt.Println("new Value of pointer: ", num)

}

//通过反射，获取变量名，方法名
func test1() {
	user := User{1, "andy", 25}
	DoAgentFieldAndMethod(user)
}

func DoAgentFieldAndMethod(input interface{}) {
	getType := reflect.TypeOf(input)
	fmt.Println("get Type is : ", getType.Name())

	/*
	go的劣势：
	getValue 的类型时Value，它是一个具体的值，而不是java反射那种可复用的反射对象，golang每次反射都需要malloc这个reflect.Value结构体，并且还设计gc
	java中：
	Field filed = clazz.getField("hello");
	field.get(obj1);
	field.get(obj2);
	这个取得的反反射对象类型是java.lang.reflect.Field。它是可以复用的。

	 */
	getValue := reflect.ValueOf(input)
	fmt.Println("get all Fields is: ", getValue)

	for i := 0; i< getType.NumField(); i++ {
		field := getType.Field(i)
		value := getValue.Field(i).Interface()
		fmt.Printf("%s: %v = %v\n", field.Name, field.Type, value)
	}

	// get func
	for i := 0; i< getType.NumMethod(); i++ {
		m := getType.Method(i)
		fmt.Printf("%s: %v\n", m.Name, m.Type)
	}
}
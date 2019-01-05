package aop

import "fmt"








type AopStructure struct {

}

func HappenBefore() {
	fmt.Println("happen before func")
}

func HappenAfter() {
	fmt.Println("happen after func")
}

type Function interface {
	Execute(v... interface{}) (m interface{}, err error)
}


func (ap *AopStructure) DoWithAop(function Function, v... interface{}) (m interface{}, err error) {
	HappenBefore()
	defer HappenAfter()

	res, err := function.Execute(v)

	m = res
	return
}


package main

import (
	"fmt"
	"time"
)

func main()  {

	ticker := time.NewTicker(10 * time.Second)
	time := <- ticker.C
	fmt.Println(time.String())
}

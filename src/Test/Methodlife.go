package main

import "fmt"

func main() {
	test1();
}

func test1(){
	{
		fmt.Println("0")
		{
			fmt.Println("1")
			{
				fmt.Println("2")
			}
		}
	}

}
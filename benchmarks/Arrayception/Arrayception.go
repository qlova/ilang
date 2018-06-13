package main

import "fmt"

func main() {
	var a []int
	
	for i:=0; i < 1000000; i++ {
		a = append(a, i*2)
	}
	
	fmt.Println(a[12317])
}

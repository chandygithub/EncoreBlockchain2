package main

import (
	"fmt"
	"strconv"
)

func main() {
	//fmt.Println("works")
	amt := int64(5000000000)
	//fmt.Println(amt)
	amtString := strconv.FormatInt(amt, 10)
	fmt.Println(amtString)
	amtByte := []byte(amtString)
	fmt.Println(amtByte)
	byteString := string(amtByte)
	fmt.Println(byteString)
}

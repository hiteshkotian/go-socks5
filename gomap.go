package main

import "fmt"

func main() {
	var myMap = make(map[string]string)
	myMap["first"] = "first"
	myMap["second"] = "second"
	myMap["third"] = "third"
	myMap["fourth"] = "fourth"

	for key, value := range myMap {
		fmt.Printf("%s : %s\n", key, value)
	}
}

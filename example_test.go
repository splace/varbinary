package varbinary_test

import "github.com/splace/varbinary"

import "fmt"


func ExampleUint64_String() {
	fmt.Println(varbinary.Uint64(0))
	fmt.Println(varbinary.Uint64(1))
	fmt.Println(varbinary.Uint64(257))
	// Output:
	// 1,2,3
	// 4,5,6,7
}


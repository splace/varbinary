package varbinary_test

import "github.com/splace/varbinary"

import "fmt"


func ExampleUint64_String() {
	fmt.Println(varbinary.Uint64(0))
	fmt.Println(varbinary.Uint64(1))
	fmt.Println(varbinary.Uint64(257))
	fmt.Println(varbinary.Uint64(65793))
	// Output:
	//
	// 00
	// 00 00
	// 00 00 00
}



package varbinary_test

import "github.com/splace/varbinary"

import "fmt"


func ExampleUint64_String() {
	fmt.Printf("'%s' '%s' '%s' '%s'",varbinary.Uint64(0),varbinary.Uint64(1),varbinary.Uint64(257),varbinary.Uint64(65793))
	// Output:
	// '' '00' '00 00' '00 00 00'
}


package varbinary_test

import "github.com/splace/varbinary"

import "fmt"
import "io/ioutil"
import "os"
import "log"

// persistance for unordered list/group of uint64's : puts all values with the same encoding length into the same file. Up too 8 files named after their contents individual lengths.
// could be a compact data representation if the uint64 values are biased toward smaller numbers.
func Example() {
	dir, err := ioutil.TempDir("", "Uint64")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(dir) // clean up temp files
	defer func() { func(dir string, _ error) { os.Chdir(dir) }(os.Getwd()) }()  // reset working dir
	os.Chdir(dir)

	set:=map[uint64]struct{}{1000:{}, 1002:{}, 1003:{}, 1004:{}, 1005:{}, 10000000:{}}
	Put(set)
	gotSet:=Get()

	fmt.Println(sum(set) == sum(gotSet)) 
	// Output:
	// true
}

func Put(p map[uint64]struct{}){
	var files [8]*os.File
	for v := range p {
		b,_ := (*varbinary.Uint64).MarshalBinary((*varbinary.Uint64)(&v))
		// use/make a file depending on encoding length
		if files[len(b)] == nil {
			var err error
			files[len(b)], err = os.Create(fmt.Sprintf("l%v", len(b)))
			if err != nil {
				log.Fatal(err)
			}
			defer files[len(b)].Close()
		}
		c, err := files[len(b)].Write(b) 
		if err != nil {
			log.Fatal(err)
		}
		if c != len(b) {
			log.Fatal(c,len(b))
		}
	}
}

func Get() (p map[uint64]struct{}){
	p=make(map[uint64]struct{})
	fileInfos, err := ioutil.ReadDir(".")
	if err != nil {
		log.Fatal(err)
	}
	var l int
	buf := make([]byte,8,8)
	for _, fi := range fileInfos {
		// scan file name for length held
		c, err := fmt.Sscanf(fi.Name(), "l%v", &l)
		if err != nil || c != 1 {
			continue
		}
		if l != 0 && fi.Size()%int64(l) != 0 {
			log.Fatal("Size indicates file isn't whole number of encodings.")
		}
		f, _ := os.Open(fi.Name())
		defer f.Close()
		var v varbinary.Uint64
		for {
			n,err := f.Read(buf[:l])
			if n != l {
				break
			}
			err = v.UnmarshalBinary(buf[:n])
			if err !=nil {
				log.Print(err)
			}
			p[uint64(v)]=struct{}{}
		}
	}
	return
}	

func sum(vs map[uint64]struct{}) (t uint64) {
	for v := range vs {
		t += v
	}
	return
}


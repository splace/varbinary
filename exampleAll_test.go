package varbinary_test

import "github.com/splace/varbinary"

import "fmt"
import "io/ioutil"
import "os"
import "log"
import "io"
import "path/filepath"

// persistance for unordered list/group of uint64's : puts all values with the same encoding length into the same file, to record the their length. Up too 8 files named after their contents individual lengths.
// could be a compact data representation if the uint64 values are biased toward smaller numbers.
func Example() {
	dir, err := ioutil.TempDir("", "Uint64")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(dir) // clean up temp files
	defer func() { func(dir string, _ error) { os.Chdir(dir) }(os.Getwd()) }()  // reset working dir
	os.Chdir(dir)

	put := []uint64{1000, 1002, 1003, 1004, 1005, 10000000}
	var files [8]*os.File
	buf := make([]byte, 8, 8)
	var l int
	for _, v := range put {
		l = varbinary.PutUint64(buf, varbinary.Uint64(v)) // encode
		if files[l] == nil {
			files[l], err = os.Create(fmt.Sprintf("l%v", l))
			if err != nil {
				log.Fatal(err)
			}
		}
		c, err := files[l].Write(buf[:l]) // different file depending on encoding length
		if err != nil || c != l {
			log.Fatal(err)
		}

	}
	for _, f := range files {
		if f != nil {
			f.Close()
		}
	}

	// read back in	
	var got []uint64
	fileInfos, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	for _, fi := range fileInfos {
		c, err := fmt.Sscanf(fi.Name(), "l%v", &l)
		if err != nil || c != 1 {
			log.Fatal(err)
		}
		if l != 0 && fi.Size()%int64(l) != 0 {
			log.Fatal("Wrong file size")
		}
		f, err := os.Open(filepath.Join(dir, fi.Name()))
		if err != nil {
			log.Fatal(err)
		}
		var i varbinary.Uint64
		var c64 int64
		for {
			c64, err = io.Copy(&i, &io.LimitedReader{f, int64(l)}) // read only the length of one encoding in this file.
			if c64 != int64(l) {
				break
			}
			got = append(got, uint64(i))
		}
	}
	
	fmt.Println(sum(put...) == sum(got...)) // same total even when order changed
	// Output:
	// true
}

func sum(vs ...uint64) (t uint64) {
	for _, v := range vs {
		t += v
	}
	return
}


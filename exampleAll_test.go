package varbinary_test

import "github.com/splace/varbinary"

import "fmt"
import "io/ioutil"
import "os"
import "log"
import "io"
import "path/filepath"

// slice encoder, puts blacks of the same length encodings into the same file, so with Uint64 up to 8 files, named after that length, slice order is lost.
// much smaller files if order isn't needed, bias to small numbers and indexing in the file not needed.
func Example() {
	dir, err := ioutil.TempDir("", "Uint64")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(dir) // clean up
	defer func(){func(dir string,_ error){os.Chdir(dir)}(os.Getwd())}()
	os.Chdir(dir)
	
	set:=[]uint64{1000,1002,1003,1004,1005,10000000}
	var files [8]*os.File
	b:=make([]byte,8,8)
	var l int
	for _,v :=range(set){
		l=varbinary.PutUint64(b,varbinary.Uint64(v))
		if files[l]==nil{
			files[l],err=os.Create(fmt.Sprintf("l%v",l))
			if err != nil {
				log.Fatal(err)
			}
		}
		c,err:=io.Copy(files[l],varbinary.Uint64(v))
		if err != nil || c!=int64(l){
			log.Fatal(err)
		}
	
	}
	for _,f:=range files{
		if f!=nil{f.Close()}
	}

	var get []uint64
	fileInfos,err:=ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}
	
	for _,fi:=range fileInfos{
		c,err:=fmt.Sscanf(fi.Name(),"l%v",&l)
		if err != nil || c!=1 {
			log.Fatal(err)
		}
		if l!=0 && fi.Size()%int64(l)!=0 {
			log.Fatal("Wrong file size")
		}
		f,err:=os.Open(filepath.Join(dir,fi.Name()))
		if err != nil {
			log.Fatal(err)
		}
		var i varbinary.Uint64
		var c64 int64
		for {
			c64,err=io.Copy(&i,&io.LimitedReader{f,int64(l)})
			if c64!=int64(l) {
				break
			}
			get=append(get,uint64(i))
		}
	}
	fmt.Println(sum(set...)==sum(get...))
	// Output:
	// true

}

func sum(vs ...uint64) (t uint64){
	for _,v :=range vs{
		t+=v
	}
	return
}


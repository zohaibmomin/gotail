package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"time"
)

func tail(filename string, out io.Writer) {
	f, err := os.Open(filename) //Open the file
	if err != nil {
		panic(err)
	}
	defer f.Close()
	r := bufio.NewReader(f) //read the file - like bufferedRead
	info, err := f.Stat()   // get file size using f stat
	if err != nil {
		panic(err)
	}
	oldSize := info.Size() // store existing file size
	for {
		//loop always-continous monitoring
		for line, prefix, err := r.ReadLine(); err != io.EOF; line, prefix, err = r.ReadLine() {
			// get the new lines until hit the err as end of file
			if prefix {
				fmt.Fprint(out, string(line))
			} else {
				fmt.Fprintln(out, string(line))
			}
		}
		pos, err := f.Seek(0, io.SeekCurrent) // helps get the position
		if err != nil {
			panic(err)
		}
		for {
			time.Sleep(time.Second)  //adds delay while concurrency
			newinfo, err := f.Stat() // gets new file size using f stat
			if err != nil {
				panic(err)
			}
			newSize := newinfo.Size()
			if newSize != oldSize {
				if newSize < oldSize { //if old file bigger then continue reading
					f.Seek(0, 0)
				} else {
					f.Seek(pos, io.SeekStart) //else go to new start
				}
				r = bufio.NewReader(f)
				oldSize = newSize // assign new size to old
				break
			}
		}
	}
}

func main() {
	tail("mylogs.log", os.Stdout)
}

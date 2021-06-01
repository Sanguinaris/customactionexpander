package main

import (
	"debug/pe"
	"fmt"
	"io"
	"os"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s SOURCEFILE\n", os.Args[0])
		return
	}

	f, err := pe.Open(os.Args[1])
	check(err)

	overlaySection := f.Sections[f.FileHeader.NumberOfSections-1].SectionHeader.Size + f.Sections[f.FileHeader.NumberOfSections-1].SectionHeader.Offset
	fmt.Println(os.Args[1], overlaySection)

	read, err := os.Open(os.Args[1])
	check(err)
	write, err := os.Create(os.Args[1] + ".cab")
	check(err)

	_, err = read.Seek(int64(overlaySection), 0)
	check(err)

	p := make([]byte, 4)

	// check that it's actually a cab and we're not at the end
	n, err := read.Read(p)
	if err == io.EOF {
		fmt.Println("No Overlay found!")
		return
	}
	if string(p[:n]) != "MSCF" {
		fmt.Println("Don't think this overlay data is a cab")
		return
	}
	fmt.Println("found cab!")
	write.Write(p[:n])

	// do the rest
	for {
		n, err = read.Read(p)
		if err == io.EOF {
			break
		}
		write.Write(p[:n])
	}
}

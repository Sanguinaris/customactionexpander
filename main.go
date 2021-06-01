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
	if len(os.Args) < 3 {
		fmt.Printf("Usage: %s SOURCEFILE OUTPUTDIR\n", os.Args[0])
		return
	}

	f, err := pe.Open(os.Args[1])
	check(err)

	overlaySection := f.Sections[f.FileHeader.NumberOfSections-1].SectionHeader.Size + f.Sections[f.FileHeader.NumberOfSections-1].SectionHeader.Offset
	fmt.Println(overlaySection)

	read, err := os.Open(os.Args[1])
	check(err)
	write, err := os.Create(os.Args[1] + ".cab")
	check(err)

	read.Seek(int64(overlaySection), 0)

	p := make([]byte, 4)
	for {
		n, err := read.Read(p)
		if err == io.EOF {
			break
		}
		write.Write(p[:n])
	}
}

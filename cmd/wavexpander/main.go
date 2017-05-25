package main

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/takuyaohashi/go-wav"
)

func usage() {
	fmt.Println("Usage: go-wavexpand [wav file]")
}

func expand(in <-chan []byte) <-chan []byte{
	out := make(chan []byte)

	go func() {
		defer close(out)
		for data := range in {
			buf := make([]byte, 4)
			for i := 0; i < len(data); i++ {
				buf[i+1] = data[i]
			}
			out <- buf
		}
	}()
	
	return out
}

func read(reader *bufio.Reader) <-chan []byte {
	out := make(chan []byte)

	go func() {
		defer close(out)
		for {
			buf := make([]byte, 3)
			_, err := reader.Read(buf)
			if err == io.EOF {
				break
			}
			out <- buf
		}
	}()
	return out
}

func main() {
	if len(os.Args) != 2 {
		usage()
	}

	f, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer f.Close()

	parser := wav.NewWav(f)
	err = parser.Parse()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	header := parser.GetHeader()

	wf, err2 := os.Create("hoge.wav")
	if err2 != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer wf.Close()
	writer := bufio.NewWriter(wf)

	fmt.Printf("size = %d\n", header.SubChunk2Size)

	reader := bufio.NewReaderSize(f, int(header.SubChunk2Size))
	out := read(reader)
	expand := expand(out)

	for i := range expand {
		writer.Write(i)
	}

	writer.Flush()
}

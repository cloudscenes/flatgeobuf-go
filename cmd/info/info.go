package main

import (
	"bufio"
	"flag"
	flatgeobuf_go "flatgeobuf-go"
	"fmt"
	"log"
	"os"
	"runtime/pprof"
)

func main() {
	quiet := flag.Bool("q", false, "quiet mode")

	f, err := os.Create("/tmp/fgb_info.pprof")
	if err != nil {
		log.Fatal("could not create CPU profile: ", err)
	}
	defer f.Close()
	if err = pprof.StartCPUProfile(f); err != nil {
		log.Fatal("could not start CPU profile: ", err)
	}
	defer pprof.StopCPUProfile()

	flag.Parse()
	if len(flag.Args()) != 1 {
		fmt.Println("please give a filename")
		return
	}

	file := flag.Args()[0]
	if stat, err := os.Stat(file); err != nil || stat.IsDir() {
		fmt.Printf("file %s not found", file)
		return
	}

	f, err = os.Open(file)
	if err != nil {
		log.Fatal(err)
	}

	br := bufio.NewReader(f)

	fgb, err := flatgeobuf_go.NewFGB(br)
	if err != nil {
		log.Fatal(err)
	}

	features := fgb.Features()

	if !*quiet {
		fmt.Print(features.Summary())
	}

	for features.Next() {
		feature := features.Read()

		if *quiet {
			continue
		}

		geom, _ := feature.Geometry()

		if geom == nil {
			continue
		}
	}
}

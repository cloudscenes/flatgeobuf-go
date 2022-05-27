package main

import (
	"bufio"
	"flag"
	flatgeobuf_go "flatgeobuf-go"
	"fmt"
	"github.com/twpayne/go-geom/encoding/wkt"
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

	fId := 0
	for features.Next() {
		feature := features.Read()

		if *quiet {
			continue
		}

		props := feature.Properties()
		geom, _ := feature.Geometry()

		fmt.Printf("Feature:%d\n", fId)
		for name, value := range props {
			fmt.Printf("  %s = %v \n", name, value)
		}

		geomWkt, _ := wkt.Marshal(geom)
		fmt.Printf("  %s \n", geomWkt)
		fmt.Printf("\n")
		fId++
	}
}

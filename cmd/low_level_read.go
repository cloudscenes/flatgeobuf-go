package main

import (
	"flatgeobuf-go"
	"flatgeobuf-go/FlatGeobuf"
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	f, err := os.Open("test/data/countries.fgb")
	if err != nil {
		log.Fatal(err)
	}

	b, err := io.ReadAll(f)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("total size: ", len(b))

	fgb, err := flatgeobuf_go.NewFGBReader(b)

	header := fgb.Header()
	crs := FlatGeobuf.Crs{}
	header.Crs(&crs)
	fmt.Println("crs", string(crs.Name()), string(crs.Wkt()))
	fmt.Println("name", string(header.Name()))
	fmt.Println("desc", header.Description())
	fmt.Println("feat count", header.FeaturesCount())

	colLen := header.ColumnsLength()
	fmt.Println("col len", colLen)
	for i := 0; i < colLen; i++ {
		var c FlatGeobuf.Column
		header.Columns(&c, i)
		fmt.Println("  col ", i, string(c.Name()), c.Description(), c.Type())
	}
	fmt.Println("index", header.IndexNodeSize())

	// READ features
	features := fgb.Features()

	for features.Next() {
		feature := features.Read()
		fmt.Println("col len", feature.ColumnsLength())
		fmt.Println("prop len", feature.PropertiesLength())
		fmt.Println("prop bytes", string(feature.PropertiesBytes()))

		// READ geometry
		var g FlatGeobuf.Geometry
		feature.Geometry(&g)
		fmt.Println(g.Type(), g.PartsLength(), ":")
		for i := 0; i < g.PartsLength(); i++ {
			var gp FlatGeobuf.Geometry
			g.Parts(&gp, i)
			fmt.Println("  ", i, gp.Type(), gp.XyLength())
			for j := 0; j < gp.XyLength(); j += 2 {
				fmt.Printf(" %f,%f ", gp.Xy(j), gp.Xy(j+1))
			}
			fmt.Println()
		}
	}
}

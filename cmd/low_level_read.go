package main

import (
	"flatgeobuf-go"
	"fmt"
	"github.com/twpayne/go-geom/encoding/geojson"
	"io"
	"log"
	"os"
)

func main() {
	f, err := os.Open("test/data/alldatatypes.fgb")
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
	crs := header.Crs(nil)

	if crs != nil {
		fmt.Println("crs", string(crs.Name()), string(crs.Wkt()))
	}
	fmt.Println("name", string(header.Name()))
	fmt.Println("desc", header.Description())
	fmt.Println("feat count", header.FeaturesCount())

	columns := flatgeobuf_go.NewColumns(header)
	fmt.Println(columns)
	fmt.Println("index", header.IndexNodeSize())

	propertyDecoder := flatgeobuf_go.NewPropertyDecoder(columns)
	// READ features
	features := fgb.Features()

	for features.Next() {
		feature := features.Read()
		fmt.Println("col len", feature.ColumnsLength())
		fmt.Println("prop len", feature.PropertiesLength())

		// TODO see proplength
		res := propertyDecoder.Decode(feature.PropertiesBytes())
		fmt.Println(res)

		// READ geometry
		geom, _ := flatgeobuf_go.ParseGeometry(header, feature.Geometry(nil))

		if geom == nil {
			continue
		}

		repr, _ := geojson.Marshal(geom)

		fmt.Printf("%s", repr)
	}
}

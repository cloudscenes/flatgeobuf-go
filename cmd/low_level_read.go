package main

import (
	"encoding/binary"
	"flatgeobuf-go/FlatGeobuf"
	"fmt"
	flatbuffers "github.com/google/flatbuffers/go"
	"io"
	"log"
	"math"
	"os"
)

const NODE_ITEM_LEN = 8*4 + 8

// adapted from the ts version
// TODO: rewrite
func calcTreeSize(numItems uint64, nodeSize uint16) uint32 {
	if nodeSize < 2 {
		nodeSize = 2
	}

	n := float64(numItems)
	numNodes := n
	for n != 1 {
		n = math.Ceil(float64(n) / float64(nodeSize))
		numNodes += n
	}

	return uint32(numNodes * NODE_ITEM_LEN)
}

func main() {
	f, err := os.Open("test/data/countries.fgb")
	if err != nil {
		log.Fatal(err)
	}

	b, err := io.ReadAll(f)
	if err != nil {
		log.Fatal(err)
	}

	// read format
	fmt.Printf("magic bytes %x\n", b[:8])
	//TODO: check format

	// READ Header
	calcLen := binary.LittleEndian.Uint32(b[8:12])
	header := FlatGeobuf.GetSizePrefixedRootAsHeader(b, 8)
	fmt.Println("len", calcLen)
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

	// READ Index
	//TODO: check if there is an index
	indexLength := calcTreeSize(header.FeaturesCount(), header.IndexNodeSize())
	fmt.Println("index length", indexLength)
	//TODO: parse index

	offset := 8 + 4 + flatbuffers.UOffsetT(calcLen) + flatbuffers.UOffsetT(indexLength)
	fmt.Println("offset", offset)

	// READ features
	featureLen := binary.LittleEndian.Uint32(b[offset : offset+4])
	fmt.Println("featlen", featureLen)
	features := FlatGeobuf.GetSizePrefixedRootAsFeature(b, offset)
	fmt.Println("col len", features.ColumnsLength())
	fmt.Println("prop len", features.PropertiesLength())
	fmt.Println("prop bytes", string(features.PropertiesBytes()))

	// READ geometry
	var g FlatGeobuf.Geometry
	features.Geometry(&g)
	fmt.Println(g.Type())
	fmt.Println(g.PartsLength())
}

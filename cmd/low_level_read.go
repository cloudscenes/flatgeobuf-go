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

type FGBReader struct {
	b              []byte
	indexOffset    uint32
	featuresOffset flatbuffers.UOffsetT
}

func NewFGBReader(b []byte) (*FGBReader, error) {
	fmt.Printf("magic bytes %x\n", b[:8])
	//TODO: check magic bytes

	res := FGBReader{b: b}

	headerLength := binary.LittleEndian.Uint32(b[8:12])
	var indexLength uint32

	header := FlatGeobuf.GetSizePrefixedRootAsHeader(b, 8)
	if header.IndexNodeSize() != 0 {
		indexLength = calcTreeSize(header.FeaturesCount(), header.IndexNodeSize())
	}

	res.indexOffset = 8 + 4 + headerLength
	res.featuresOffset = flatbuffers.UOffsetT(8 + 4 + headerLength + indexLength)

	// TODO: should we validate the size?
	//featureLen := binary.LittleEndian.Uint32(b[res.featuresOffset : res.featuresOffset+4])

	return &res, nil
}

func (fgb *FGBReader) Header() *FlatGeobuf.Header {
	return FlatGeobuf.GetSizePrefixedRootAsHeader(fgb.b, 8)
}

func (fgb *FGBReader) Feature() *FlatGeobuf.Feature {
	return FlatGeobuf.GetSizePrefixedRootAsFeature(fgb.b, fgb.featuresOffset)
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

	fgb, err := NewFGBReader(b)

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
	features := fgb.Feature()
	fmt.Println("col len", features.ColumnsLength())
	fmt.Println("prop len", features.PropertiesLength())
	fmt.Println("prop bytes", string(features.PropertiesBytes()))

	// READ geometry
	var g FlatGeobuf.Geometry
	features.Geometry(&g)
	fmt.Println(g.Type())
	fmt.Println(g.PartsLength())
}

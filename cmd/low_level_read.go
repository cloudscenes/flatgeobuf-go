package main

import (
	"encoding/binary"
	"flatgeobuf-go"
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

	res.indexOffset = 8 + flatbuffers.SizeUint32 + headerLength
	res.featuresOffset = flatbuffers.UOffsetT(8 + flatbuffers.SizeUint32 + headerLength + indexLength)

	// TODO: should we validate the size?
	//featureLen := binary.LittleEndian.Uint32(b[res.featuresOffset : res.featuresOffset+4])

	return &res, nil
}

func (fgb *FGBReader) Header() *FlatGeobuf.Header {
	return FlatGeobuf.GetSizePrefixedRootAsHeader(fgb.b, 8)
}

func (fgb *FGBReader) Features() *Features {
	return NewFeatures(fgb.b[fgb.featuresOffset:])
}

type Features struct {
	b       []byte
	pos     uint32
	started bool // TODO: this should not be needed
}

func NewFeatures(b []byte) *Features {
	return &Features{
		b:       b,
		pos:     0,
		started: false,
	}
}

func (fs *Features) featureLen() uint32 {
	return binary.LittleEndian.Uint32(fs.b[fs.pos : fs.pos+4])
}

func (fs *Features) Next() bool {
	if !fs.started {
		fs.started = true
		fs.pos = 0
		return true
	}

	if int(fs.pos+4) >= len(fs.b) {
		return false
	}

	fs.pos += fs.featureLen() + flatbuffers.SizeUint32
	if int(fs.pos) >= len(fs.b) {
		return false
	}
	return true
}

func (fs *Features) Read() *FlatGeobuf.Feature {
	return FlatGeobuf.GetSizePrefixedRootAsFeature(fs.b, flatbuffers.UOffsetT(fs.pos))
}

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

	fgb, err := NewFGBReader(b)

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
		g := feature.Geometry(nil)
		if g == nil {
			continue
		}
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

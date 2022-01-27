package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"flatgeobuf-go/FlatGeobuf"
	"fmt"
	flatbuffers "github.com/google/flatbuffers/go"
	"io"
	"log"
	"math"
	"os"
)

const NODE_ITEM_LEN = 8*4 + 8
const SupportedVersion uint8 = 3

// last byte is the patch level that is backwards compatible,
// so an implementation for a major version should accept any patch level version
var MagicBytes = []byte{0x66, 0x67, 0x62, SupportedVersion, 0x66, 0x67, 0x62, 0x00}

var (
	ErrUnsupportedVersion = errors.New("unsupported flatgeobuffer version")
	ErrInvalidFile        = errors.New("invalid flatgeobuffer file")
)

func Version(magicBytes []byte) (string, error) {
	if !bytes.Equal(magicBytes[0:3], MagicBytes[0:3]) || !bytes.Equal(magicBytes[4:7], MagicBytes[4:7]) {
		return "", ErrInvalidFile
	}

	var majorVersion = magicBytes[3]
	var patchVersion = magicBytes[7]

	if majorVersion != SupportedVersion {
		return "", ErrUnsupportedVersion
	}

	return fmt.Sprintf("%d.0.%d", majorVersion, patchVersion), nil
}

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
	version, err := Version(b[:8])

	if err != nil {
		panic(err)
	}

	fmt.Println(version)

	res := FGBReader{b: b}

	headerLength := flatbuffers.GetUint32(b[8:12])
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
	f, err := os.Open("test/data/countries.fgb")
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

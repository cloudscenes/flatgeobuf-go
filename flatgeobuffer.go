package flatgeobuf_go

import (
	"bytes"
	"errors"
	"flatgeobuf-go/FlatGeobuf"
	"fmt"
	"github.com/google/flatbuffers/go"
	"math"
)

const nodeItemLen = 8*4 + 8
const supportedVersion uint8 = 3

// last byte is the patch level that is backwards compatible,
// so an implementation for a major version should accept any patch level version
var magicBytes = []byte{0x66, 0x67, 0x62, supportedVersion, 0x66, 0x67, 0x62, 0x00}

var (
	ErrUnsupportedVersion = errors.New("unsupported flatgeobuffer version")
	ErrInvalidFile        = errors.New("invalid flatgeobuffer file")
)

func Version(fileMagicBytes []byte) (string, error) {
	if !bytes.Equal(fileMagicBytes[0:3], magicBytes[0:3]) || !bytes.Equal(fileMagicBytes[4:7], magicBytes[4:7]) {
		return "", ErrInvalidFile
	}

	var majorVersion = fileMagicBytes[3]
	var patchVersion = fileMagicBytes[7]

	if majorVersion != supportedVersion {
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

	return uint32(numNodes * nodeItemLen)
}

type FGBReader struct {
	b              []byte
	indexOffset    uint32
	featuresOffset flatbuffers.UOffsetT
}

func NewFGBReader(b []byte) (*FGBReader, error) {
	fmt.Printf("magic bytes %x\n", b[:8])

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

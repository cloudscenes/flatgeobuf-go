package flatgeobuf_go

import (
	"bytes"
	"errors"
	"flatgeobuf-go/FlatGeobuf"
	"flatgeobuf-go/index"
	"fmt"
	"github.com/google/flatbuffers/go"
)

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

type FGBReader struct {
	b              []byte
	indexOffset    uint32
	featuresOffset flatbuffers.UOffsetT
	prt            *index.PackedRTree
}

func NewFGBReader(b []byte) (*FGBReader, error) {
	_, err := Version(b[:8])
	if err != nil {
		panic(err)
	}

	res := FGBReader{b: b}

	headerLength := flatbuffers.GetUint32(b[8:12])
	res.indexOffset = 8 + flatbuffers.SizeUint32 + headerLength

	var indexLength uint64
	header := FlatGeobuf.GetSizePrefixedRootAsHeader(b, 8)
	if header.IndexNodeSize() != 0 {
		indexLength, err = index.CalcTreeSize(header.FeaturesCount(), header.IndexNodeSize())
		prt, err := index.ReadPackedRTreeBytes(header.FeaturesCount(), header.IndexNodeSize(), b[res.indexOffset:uint64(res.indexOffset)+indexLength])
		if err != nil {
			return nil, fmt.Errorf("could not read index: %w", err)
		}
		res.prt = prt
	}

	res.featuresOffset = flatbuffers.UOffsetT(8 + flatbuffers.SizeUint32 + headerLength + uint32(indexLength))

	// TODO: should we validate the size?
	//featureLen := binary.LittleEndian.Uint32(b[res.featuresOffset : res.featuresOffset+4])

	return &res, nil
}

func (fgb *FGBReader) Header() *FlatGeobuf.Header {
	return FlatGeobuf.GetSizePrefixedRootAsHeader(fgb.b, 8)
}

func (fgb *FGBReader) Features() *Features {
	return NewFeatures(fgb.b[fgb.featuresOffset:], fgb.Header())
}

func (fgb *FGBReader) Index() *index.PackedRTree {
	return fgb.prt
}

package flatgeobuf_go

import (
	"bytes"
	"errors"
	"flatgeobuf-go/FlatGeobuf"
	"flatgeobuf-go/index"
	"fmt"
	"github.com/google/flatbuffers/go"
	"io"
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
	reader io.Reader
	header *FlatGeobuf.HeaderT
	prt    *index.PackedRTree
}

func NewFGB(r io.Reader) (*FGBReader, error) {
	buffer := make([]byte, 12)
	n, err := io.ReadFull(r, buffer)
	if err != nil || n < 12 {
		return nil, fmt.Errorf("could not read first bytes: %w", err)
	}

	fgb := FGBReader{}

	headerLength := flatbuffers.GetUint32(buffer[8:12])

	buffer = make([]byte, headerLength)
	n, err = r.Read(buffer)
	if err != nil || n < int(headerLength) {
		return nil, fmt.Errorf("could not read header: %w", err)
	}

	var indexLength uint64
	header := FlatGeobuf.GetRootAsHeader(buffer, 0).UnPack()
	fgb.header = header

	if header.IndexNodeSize != 0 {
		indexLength, err = index.CalcTreeSize(header.FeaturesCount, header.IndexNodeSize)
		buffer = make([]byte, indexLength)
		n, err = io.ReadFull(r, buffer)
		if err != nil || n < int(indexLength) {
			return nil, fmt.Errorf("could not read index: %w", err)
		}

		prt, err := index.ReadPackedRTreeBytes(header.FeaturesCount, header.IndexNodeSize, buffer)
		if err != nil {
			return nil, fmt.Errorf("could not read index: %w", err)
		}
		fgb.prt = prt
	}

	fgb.reader = r

	return &fgb, nil
}

func (fgb *FGBReader) Header() *FlatGeobuf.HeaderT {
	return fgb.header
}

func (fgb *FGBReader) Index() *index.PackedRTree {
	return fgb.prt
}

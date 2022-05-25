package flatgeobuf_go

import (
	"encoding/binary"
	"flatgeobuf-go/FlatGeobuf"
	"io"
	"log"
)

type Features struct {
	header *FlatGeobuf.Header
	r      io.Reader
	fLen   uint32
}

func NewFeatures(r io.Reader, header *FlatGeobuf.Header) *Features {
	return &Features{
		header: header,
		r:      r,
		fLen:   0,
	}
}

func (fs *Features) featureLen() (uint32, error) {
	b := make([]byte, 4)
	_, err := fs.r.Read(b)
	if err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint32(b), nil
}

func (fs *Features) Next() bool {
	fLen, err := fs.featureLen()
	if err == io.EOF {
		return false
	} else if err != nil {
		log.Fatalf("unexpected error reading feature %v", err)
		return false
	}

	fs.fLen = fLen
	return true
}

func (fs *Features) Read() (*Feature, error) {
	b := make([]byte, fs.fLen)
	// TODO: handle errors
	io.ReadFull(fs.r, b)

	fgbFeature := FlatGeobuf.GetRootAsFeature(b, 0)

	feature, err := NewFeature(fgbFeature, fs.header)
	if err != nil {
		return nil, err
	}

	return feature, nil
}

func (fs *Features) ReadAt(pos uint32) *FlatGeobuf.Feature {
	//TODO: this cannot be implemented with a simple reader
	return nil
	//return FlatGeobuf.GetSizePrefixedRootAsFeature(fs.b, flatbuffers.UOffsetT(pos))
}

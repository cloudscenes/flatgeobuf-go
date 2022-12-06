package flatgeobuf_go

import (
	"encoding/binary"
	"flatgeobuf-go/FlatGeobuf"
	"github.com/twpayne/go-geom"
	"io"
	"log"
)

type Features struct {
	header          *FlatGeobuf.Header
	r               io.Reader
	propertyDecoder *PropertyDecoder
}

func NewFeatures(r io.Reader, header *FlatGeobuf.Header) *Features {
	columns := NewColumns(header)
	propertyDecoder := NewPropertyDecoder(columns)

	return &Features{
		header:          header,
		r:               r,
		propertyDecoder: propertyDecoder,
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

func (fs *Features) Read() (*Feature, error) {
	fLen, err := fs.featureLen()

	if err == io.EOF {
		return nil, err
	} else if err != nil {
		log.Fatalf("unexpected error reading feature %v", err)
		return nil, err
	}

	b := make([]byte, fLen)
	// TODO: handle errors
	io.ReadFull(fs.r, b)

	fgbFeature := FlatGeobuf.GetRootAsFeature(b, 0)

	feature := NewFeature(fgbFeature, fs)

	return feature, nil
}

func (fs *Features) ReadAt(pos uint32) *FlatGeobuf.Feature {
	//TODO: this cannot be implemented with a simple reader
	return nil
	//return FlatGeobuf.GetSizePrefixedRootAsFeature(fs.b, flatbuffers.UOffsetT(pos))
}

// TODO: check if header crs is never nil and defaults to 0 -> uknown as header.fbs
// 		 also check if 0 is an acceptable value for go geom srid
func (fs *Features) Crs() int {
	crs := fs.header.Crs(nil)

	if crs == nil {
		return 0
	}

	return int(crs.Code())
}

func (fs *Features) Layout() geom.Layout {
	return parseLayout(fs.header)
}

func (fs *Features) GeometryType() FlatGeobuf.GeometryType {
	return fs.header.GeometryType()
}

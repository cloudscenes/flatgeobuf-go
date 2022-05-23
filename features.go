package flatgeobuf_go

import (
	"encoding/binary"
	"flatgeobuf-go/FlatGeobuf"
	"github.com/google/flatbuffers/go"
	"github.com/twpayne/go-geom"
)

type Features struct {
	crs          *FlatGeobuf.Crs
	layout       geom.Layout
	geometryType FlatGeobuf.GeometryType
	b            []byte
	pos          uint32
	started      bool // TODO: this should not be needed
}

func NewFeatures(b []byte, header *FlatGeobuf.Header) *Features {
	return &Features{
		crs:          header.Crs(nil),
		layout:       ParseLayout(header),
		geometryType: header.GeometryType(),
		b:            b,
		pos:          0,
		started:      false,
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

func (fs *Features) ReadAt(pos uint32) *FlatGeobuf.Feature {
	return FlatGeobuf.GetSizePrefixedRootAsFeature(fs.b, flatbuffers.UOffsetT(pos))
}

func (fs *Features) ReadGeometry() (geom.T, error) {
	feature := fs.Read()

	return ParseGeometry(feature.Geometry(nil), fs.geometryType, fs.layout, fs.crs)
}

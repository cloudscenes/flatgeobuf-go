package flatgeobuf_go

import (
	"encoding/binary"
	"flatgeobuf-go/FlatGeobuf"
	"github.com/google/flatbuffers/go"
)

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

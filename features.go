package flatgeobuf_go

import (
	"encoding/binary"
	"flatgeobuf-go/FlatGeobuf"
	"github.com/twpayne/go-geom"
	"io"
	"log"
)

func (fgb *FGBReader) featureLen() (uint32, error) {
	b := make([]byte, 4)
	_, err := fgb.reader.Read(b)
	if err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint32(b), nil
}

func (fgb *FGBReader) ReadFeature() (*Feature, error) {
	fLen, err := fgb.featureLen()

	if err == io.EOF {
		return nil, err
	} else if err != nil {
		log.Fatalf("unexpected error reading feature %v", err)
		return nil, err
	}

	b := make([]byte, fLen)
	// TODO: handle errors
	io.ReadFull(fgb.reader, b)

	fgbFeature := FlatGeobuf.GetRootAsFeature(b, 0).UnPack()

	feature := &Feature{
		FeatureT:           *fgbFeature,
		crs:                fgb.Crs(),
		layout:             fgb.Layout(),
		headerGeometryType: fgb.GeometryType(),
		headerColumns:      fgb.header.Columns,
	}

	return feature, nil
}

func (fgb *FGBReader) ReadFeatureAt(pos uint32) *FlatGeobuf.FeatureT {
	//TODO: this cannot be implemented with a simple reader
	return nil
	//return FlatGeobuf.GetSizePrefixedRootAsFeature(fs.b, flatbuffers.UOffsetT(pos))
}

// TODO: check if header crs is never nil and defaults to 0 -> uknown as header.fbs
// 		 also check if 0 is an acceptable value for go geom srid
func (fgb *FGBReader) Crs() int {
	crs := fgb.header.Crs

	if crs == nil {
		return 0
	}

	return int(crs.Code)
}

func (fgb *FGBReader) Layout() geom.Layout {
	return parseLayout(fgb.header)
}

func (fgb *FGBReader) GeometryType() FlatGeobuf.GeometryType {
	return fgb.header.GeometryType
}

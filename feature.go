package flatgeobuf_go

import (
	"flatgeobuf-go/FlatGeobuf"
	"fmt"
	"github.com/twpayne/go-geom"
)

type Feature struct {
	fgbFeature   *FlatGeobuf.Feature
	geometryType FlatGeobuf.GeometryType
	crs          *FlatGeobuf.Crs
	layout       geom.Layout
	columns      *Columns
}

func NewFeature(fgbFeature *FlatGeobuf.Feature, header *FlatGeobuf.Header) *Feature {
	feature := Feature{
		fgbFeature:   fgbFeature,
		geometryType: header.GeometryType(),
		crs:          header.Crs(nil),
		layout:       parseLayout(header),
		columns:      NewColumns(header),
	}

	return &feature
}

func (f *Feature) Geometry() (geom.T, error) {
	geometry := f.fgbFeature.Geometry(nil)
	geometryType := f.geometryType
	if geometryType == FlatGeobuf.GeometryTypeUnknown {
		geometryType = geometry.Type()
	}

	var newGeom geom.T
	var err error

	if geometry.PartsLength() > 0 {
		newGeom, err = parseMultiGeometry(geometry, f.layout, geometryType)
	} else {
		newGeom, err = parseSimpleGeometry(geometry, f.layout, geometryType)
	}
	if err != nil {
		return nil, fmt.Errorf("unable to parse geometry: %w", err)
	}

	if f.crs == nil {
		return newGeom, nil
	}

	sridGeom, err := geom.SetSRID(newGeom, int(f.crs.Code()))
	if err != nil {
		return nil, fmt.Errorf("unable to set SRID: %w", err)
	}

	return sridGeom, nil
}

func (f *Feature) Properties() map[string]interface{} {
	propertyDecoder := NewPropertyDecoder(f.columns)
	props := propertyDecoder.Decode(f.fgbFeature.PropertiesBytes())

	return props
}

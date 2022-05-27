package flatgeobuf_go

import (
	"flatgeobuf-go/FlatGeobuf"
	"fmt"
	"github.com/twpayne/go-geom"
)

type Feature struct {
	fgbFeature *FlatGeobuf.Feature
	features   *Features
}

func NewFeature(fgbFeature *FlatGeobuf.Feature, features *Features) *Feature {
	feature := Feature{
		fgbFeature: fgbFeature,
		features:   features,
	}

	return &feature
}

func (f *Feature) Geometry() (geom.T, error) {
	geometry := f.fgbFeature.Geometry(nil)
	geometryType := f.features.GeometryType()
	if geometryType == FlatGeobuf.GeometryTypeUnknown {
		geometryType = geometry.Type()
	}

	var newGeom geom.T
	var err error

	layout := f.features.Layout()
	if geometry.PartsLength() > 0 {
		newGeom, err = parseMultiGeometry(geometry, layout, geometryType)
	} else {
		newGeom, err = parseSimpleGeometry(geometry, layout, geometryType)
	}
	if err != nil {
		return nil, fmt.Errorf("unable to parse geometry: %w", err)
	}

	crs := f.features.Crs()
	if crs == 0 {
		return newGeom, nil
	}

	sridGeom, err := geom.SetSRID(newGeom, crs)
	if err != nil {
		return nil, fmt.Errorf("unable to set SRID: %w", err)
	}

	return sridGeom, nil
}

func (f *Feature) Properties() map[string]interface{} {
	props := f.features.propertyDecoder.Decode(f.fgbFeature.PropertiesBytes())

	return props
}

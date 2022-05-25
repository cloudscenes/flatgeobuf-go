package flatgeobuf_go

import (
	"flatgeobuf-go/FlatGeobuf"
	"fmt"
	"github.com/twpayne/go-geom"
)

type Feature struct {
	geometry   geom.T
	properties map[string]interface{}
}

func NewFeature(fgbFeature *FlatGeobuf.Feature, header *FlatGeobuf.Header) (*Feature, error) {
	geometry, err := parseGeometry(fgbFeature, header)
	if err != nil {
		return nil, err
	}
	props := parseProperties(fgbFeature, header)
	feature := Feature{
		geometry:   geometry,
		properties: props,
	}

	return &feature, nil
}

func parseGeometry(fgbFeature *FlatGeobuf.Feature, header *FlatGeobuf.Header) (geom.T, error) {
	geometry := fgbFeature.Geometry(nil)
	geometryType := header.GeometryType()
	if geometryType == FlatGeobuf.GeometryTypeUnknown {
		geometryType = geometry.Type()
	}

	var newGeom geom.T
	var err error

	layout := parseLayout(header)
	if geometry.PartsLength() > 0 {
		newGeom, err = parseMultiGeometry(geometry, layout, geometryType)
	} else {
		newGeom, err = parseSimpleGeometry(geometry, layout, geometryType)
	}
	if err != nil {
		return nil, fmt.Errorf("unable to parse geometry: %w", err)
	}

	crs := header.Crs(nil)
	if crs == nil {
		return newGeom, nil
	}

	sridGeom, _ := geom.SetSRID(newGeom, int(crs.Code()))

	return sridGeom, nil
}

func parseProperties(fgbFeature *FlatGeobuf.Feature, header *FlatGeobuf.Header) map[string]interface{} {
	columns := NewColumns(header)
	propertyDecoder := NewPropertyDecoder(columns)
	props := propertyDecoder.Decode(fgbFeature.PropertiesBytes())

	return props
}

func (f *Feature) Geometry() geom.T {
	return f.geometry
}

func (f *Feature) Properties() map[string]interface{} {
	return f.properties
}

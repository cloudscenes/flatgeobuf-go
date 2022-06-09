package flatgeobuf_go

import (
	"flatgeobuf-go/FlatGeobuf"
	"fmt"
	"github.com/twpayne/go-geom"
	"reflect"
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

func (f *Feature) Unmarshal(v interface{}) error {
	rv := reflect.ValueOf(v)

	if rv.Kind() != reflect.Ptr || rv.IsNil() || rv.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("cannot unmarshal")
	}

	props := f.Properties()

	rt := rv.Elem().Type()
	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)

		target := field.Name
		if field.Tag != "" {
			target = field.Tag.Get("fgb")
		}

		value, exists := props[target]
		if !exists {
			continue
		}

		val := reflect.ValueOf(value)
		if val.Type() != field.Type {
			panic("cannot unmarshal " + val.Type().Name() + " into Go struct field " + rt.Name() + "." + field.Name + " of type " + field.Type.Name())
		}

		rv.Elem().Field(i).Set(reflect.ValueOf(value))
	}

	return nil
}

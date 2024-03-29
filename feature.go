package flatgeobuf_go

import (
	"flatgeobuf-go/FlatGeobuf"
	"fmt"
	"github.com/twpayne/go-geom"
	"reflect"
	"strings"
)

type Feature struct {
	FlatGeobuf.FeatureT
	crs                int
	layout             geom.Layout
	headerGeometryType FlatGeobuf.GeometryType
	headerColumns      []*FlatGeobuf.ColumnT
}

func (f *Feature) Geometry() (geom.T, error) {
	geometry := f.FeatureT.Geometry
	geometryType := f.headerGeometryType
	if geometryType == FlatGeobuf.GeometryTypeUnknown {
		geometryType = geometry.Type
	}

	var newGeom geom.T
	var err error

	layout := f.layout
	if len(geometry.Parts) > 0 {
		newGeom, err = parseMultiGeometry(geometry, layout, geometryType)
	} else {
		newGeom, err = parseSimpleGeometry(geometry, layout, geometryType)
	}
	if err != nil {
		return nil, fmt.Errorf("unable to parse geometry: %w", err)
	}

	crs := f.crs
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
	props := f.decode()

	return props
}

func (f *Feature) Unmarshal(v interface{}) error {
	rv := reflect.ValueOf(v)

	if rv.Kind() != reflect.Ptr || rv.IsNil() || rv.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("cannot unmarshal")
	}

	// TODO: handle error
	geometry, _ := f.Geometry()
	props := f.Properties()

	rt := rv.Elem().Type()
	gt := reflect.TypeOf(geometry)
	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		fieldType := field.Type
		fieldKind := fieldType.Kind()

		target := field.Name
		tag := field.Tag.Get("fgb")
		if tag != "" {
			target = tag
		}

		target = strings.ToLower(target)

		if target == "geom" {
			// only accept values that are geom.T
			if fieldKind != reflect.Ptr && fieldType != reflect.TypeOf((*geom.T)(nil)).Elem() {
				panic("cannot unmarshall geometry to a struct field that isn't a pointer or a geom.T interface!")
			}

			// don't try to match different types (e.g. store a Point in a Linestring var)
			if fieldKind == reflect.Ptr && fieldType.Elem() != gt.Elem() {
				panic("cannot unmarshal " + gt.Elem().Name() + " into Go struct field " + rt.Name() + "." + field.Name + " of type " + field.Type.Elem().Name())
			}

			rv.Elem().Field(i).Set(reflect.ValueOf(geometry))

			continue
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

package flatgeobuf_go

import (
	"flatgeobuf-go/FlatGeobuf"
	"fmt"
	"github.com/twpayne/go-geom"
)

func ParseGeometry(header *FlatGeobuf.Header, geometry *FlatGeobuf.Geometry) (geom.T, error) {
	geomType := header.GeometryType()

	if geomType == FlatGeobuf.GeometryTypeUnknown {
		geomType = geometry.Type()
	}

	layout := geom.XY

	if header.HasZ() && header.HasM() {
		layout = geom.XYZM
	} else if header.HasZ() {
		layout = geom.XYZ
	} else if header.HasM() {
		layout = geom.XYM
	}

	fmt.Printf("Layout: %v\n", layout)
	fmt.Printf("Geom Type: %v\n", geomType)
	fmt.Printf("XY Len: %v\n", geometry.XyLength())
	fmt.Printf("Z Len: %v\n", geometry.ZLength())
	fmt.Printf("M Len: %v\n\n", geometry.MLength())

	coords := make([]float64, 0, geometry.XyLength()/2*layout.Stride())

	for i := 0; i < geometry.XyLength(); i += 2 {
		coords = append(coords, geometry.Xy(i), geometry.Xy(i+1))

		if header.HasZ() {
			coords = append(coords, geometry.Z(i))
		}

		if header.HasM() {
			coords = append(coords, geometry.M(i))
		}
	}

	var newGeom geom.T

	switch geomType {
	case FlatGeobuf.GeometryTypePoint:
		newGeom = geom.NewPointFlat(layout, coords)
	case FlatGeobuf.GeometryTypeLineString:
		newGeom = geom.NewLineStringFlat(layout, coords)

		if v, ok := newGeom.(*geom.LineString); ok {
			fmt.Println(v.Coords())
		}
	case FlatGeobuf.GeometryTypePolygon, FlatGeobuf.GeometryTypeMultiPoint, FlatGeobuf.GeometryTypeMultiLineString, FlatGeobuf.GeometryTypeMultiPolygon, FlatGeobuf.GeometryTypeGeometryCollection:
		// TODO: implement methods
		newGeom = geom.NewPointFlat(geom.XY, []float64{0.0, 0.0})
	default:
		return nil, fmt.Errorf("unsupported geometry type")
	}

	crs := header.Crs(nil)

	if crs == nil {
		return newGeom, nil
	}

	sridGeom, _ := geom.SetSRID(newGeom, int(crs.Code()))

	return sridGeom, nil
}

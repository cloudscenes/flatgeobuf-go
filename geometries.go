package flatgeobuf_go

import (
	"flatgeobuf-go/FlatGeobuf"
	"fmt"
	"github.com/twpayne/go-geom"
)

func parseSimpleGeometry(geometry *FlatGeobuf.Geometry, layout geom.Layout, geomType FlatGeobuf.GeometryType) (geom.T, error) {
	coords := make([]float64, 0, geometry.XyLength()/2*layout.Stride())
	geomLen := geometry.XyLength() / 2
	for i := 0; i < geomLen; i += 1 {
		coords = append(coords, geometry.Xy(i*2), geometry.Xy(i*2+1))

		if layout == geom.XYZ || layout == geom.XYZM {
			coords = append(coords, geometry.Z(i))
		}
		if layout == geom.XYM || layout == geom.XYZM {
			coords = append(coords, geometry.M(i))
		}
	}

	var newGeom geom.T
	switch geomType {
	case FlatGeobuf.GeometryTypePoint:
		newGeom = geom.NewPointFlat(layout, coords)
	case FlatGeobuf.GeometryTypeLineString:
		newGeom = geom.NewLineStringFlat(layout, coords)
	case FlatGeobuf.GeometryTypePolygon:
		ends := getEnds(geometry, layout.Stride())
		// in flatgeobuf the ends doesn't include the last position, in go-geom it is needed
		ends = append(ends, len(coords))
		newGeom = geom.NewPolygonFlat(layout, coords, ends)
	case FlatGeobuf.GeometryTypeMultiPoint:
		newGeom = geom.NewMultiPointFlat(layout, coords)
	case FlatGeobuf.GeometryTypeMultiLineString:
		ends := getEnds(geometry, layout.Stride())
		newGeom = geom.NewMultiLineStringFlat(layout, coords, ends)
	default:
		return nil, fmt.Errorf("unsupported geometry type %s", geomType)
	}

	return newGeom, nil
}

func parseMultiGeometry(geometry *FlatGeobuf.Geometry, layout geom.Layout, geomType FlatGeobuf.GeometryType) (geom.T, error) {
	var newGeom geom.T
	var addToCollection func(createdGeom geom.T) error

	switch geomType {
	case FlatGeobuf.GeometryTypeMultiPolygon:
		multiPolygon := geom.NewMultiPolygon(layout)
		addToCollection = func(createdGeom geom.T) error {
			polygon, ok := createdGeom.(*geom.Polygon)
			if !ok {
				return fmt.Errorf("multipolygon cannot have geoms that are not polygons")
			}

			err := multiPolygon.Push(polygon)
			if err != nil {
				return fmt.Errorf("cannot push polygon to collection: %w", err)
			}

			return nil
		}

		newGeom = multiPolygon
	case FlatGeobuf.GeometryTypeGeometryCollection:
		geomCollection := geom.NewGeometryCollection()
		addToCollection = func(createdGeom geom.T) error {
			err := geomCollection.Push(createdGeom)
			if err != nil {
				return fmt.Errorf("cannot push geometry to collection: %w", err)
			}

			return nil
		}

		newGeom = geomCollection
	default:
		return nil, fmt.Errorf("unsupported geometry type %s", geomType)
	}

	for i := 0; i < geometry.PartsLength(); i++ {
		partGeom := FlatGeobuf.Geometry{}
		geometry.Parts(&partGeom, i)

		createdGeom, err := parseSimpleGeometry(&partGeom, layout, partGeom.Type())
		if err != nil {
			return nil, err
		}

		err = addToCollection(createdGeom)
		if err != nil {
			return nil, err
		}
	}

	return newGeom, nil
}

func getEnds(geometry *FlatGeobuf.Geometry, stride int) []int {
	ends := make([]int, geometry.EndsLength())
	for i := 0; i < geometry.EndsLength(); i++ {
		// in flatgeobuf ends is based on stride, in go-geom it's index based
		ends[i] = int(geometry.Ends(i)) * stride
	}

	return ends
}

func parseLayout(header *FlatGeobuf.Header) geom.Layout {
	layout := geom.XY
	if header.HasZ() && header.HasM() {
		layout = geom.XYZM
	} else if header.HasZ() {
		layout = geom.XYZ
	} else if header.HasM() {
		layout = geom.XYM
	}
	return layout
}

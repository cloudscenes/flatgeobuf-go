package flatgeobuf_go

import (
	"fmt"
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/encoding/wkt"
	"io"
	"log"
	"os"
	"reflect"
	"strings"
	"testing"
)

func readFile(path string) *FGBReader {
	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	b, err := io.ReadAll(f)
	if err != nil {
		log.Fatal(err)
	}
	fgb, err := NewFGBReader(b)
	if err != nil {
		log.Fatal(err)
	}

	return fgb
}

func featuresToMap(path string) map[string]string {
	fgb := readFile(path)

	header := fgb.Header()
	columns := NewColumns(header)
	features := fgb.Features()
	propertyDecoder := NewPropertyDecoder(columns)

	featuresMap := make(map[string]string)

	for features.Next() {
		feature := features.Read()
		props := propertyDecoder.Decode(feature.PropertiesBytes())
		geom, err := features.ReadGeometry()
		if err != nil {
			log.Fatal(err)
		}

		geomWKT, err := wkt.Marshal(geom, wkt.EncodeOptionWithMaxDecimalDigits(10))
		if err != nil {
			log.Fatal(err)
		}

		featuresMap[fmt.Sprint(props["type"])] = geomWKT
	}

	return featuresMap
}

func TestParseGeometry(t *testing.T) {
	tests := []struct {
		name string
		file string
		want map[string]string
	}{
		{
			name: "Parse XY",
			file: "test/data/all_geom_types.fgb",
			want: map[string]string{
				"point":               "POINT (-8.6291106307 41.1580318815)",
				"linestring":          "LINESTRING (-8.6107858821 41.145636564, -8.6108923285 41.1450754532, -8.6141070099 41.1427908808, -8.6144157045 41.1411234935, -8.6166617236 41.1403298664)",
				"polygon":             "POLYGON ((-8.6137876707 41.1459952717, -8.6137158194 41.1458069003, -8.6145780353 41.1455804532, -8.614705771 41.1456004928, -8.6147616553 41.1456986867, -8.6146339196 41.1457668212, -8.6137876707 41.1459952717))",
				"multipoint":          "MULTIPOINT (-8.6887871464 41.1733073983, -8.6898516104 41.1685638551)",
				"multilinestring":     "MULTILINESTRING ((-8.6111770726 41.1492896792, -8.6112728744 41.1478949867, -8.6116347922 41.1460513822), (-8.6112728744 41.1460513822, -8.6108896673 41.1478308622, -8.6107193531 41.1486243986, -8.6104213032 41.1491934945))",
				"multipolygon":        "MULTIPOLYGON (((-8.626590512 41.1469491439, -8.6265479335 41.1464842331, -8.6251960641 41.1463159026, -8.6250789731 41.1471575511, -8.6264095532 41.1473739733, -8.626590512 41.1469491439)), ((-8.6267076031 41.1488568468, -8.6271227441 41.1479510956, -8.6245467411 41.1473018327, -8.6241528894 41.1481755307, -8.6267076031 41.1488568468)))",
				"geometry collection": "GEOMETRYCOLLECTION (LINESTRING (-8.6107858821 41.145636564, -8.6108923285 41.1450754532, -8.6141070099 41.1427908808, -8.6144157045 41.1411234935, -8.6166617236 41.1403298664), MULTIPOINT (-8.6887871464 41.1733073983, -8.6898516104 41.1685638551))",
			},
		},
		{
			name: "Parse XYZ",
			file: "test/data/all_geom_types_z.fgb",
			want: map[string]string{
				"point":               "POINT Z (-8.6291106307 41.1580318815 10)",
				"linestring":          "LINESTRING Z (-8.6107858821 41.145636564 20, -8.6108923285 41.1450754532 21, -8.6141070099 41.1427908808 22, -8.6144157045 41.1411234935 23, -8.6166617236 41.1403298664 24)",
				"polygon":             "POLYGON Z ((-8.6137876707 41.1459952717 30, -8.6137158194 41.1458069003 31, -8.6145780353 41.1455804532 32, -8.614705771 41.1456004928 33, -8.6147616553 41.1456986867 34, -8.6146339196 41.1457668212 35, -8.6137876707 41.1459952717 30))",
				"multipoint":          "MULTIPOINT Z (-8.6887871464 41.1733073983 40, -8.6898516104 41.1685638551 41)",
				"multilinestring":     "MULTILINESTRING Z ((-8.6111770726 41.1492896792 50, -8.6112728744 41.1478949867 51, -8.6116347922 41.1460513822 52), (-8.6112728744 41.1460513822 53, -8.6108896673 41.1478308622 54, -8.6107193531 41.1486243986 55, -8.6104213032 41.1491934945 56))",
				"multipolygon":        "MULTIPOLYGON Z (((-8.626590512 41.1469491439 60, -8.6265479335 41.1464842331 61, -8.6251960641 41.1463159026 62, -8.6250789731 41.1471575511 63, -8.6264095532 41.1473739733 64, -8.626590512 41.1469491439 60)), ((-8.6267076031 41.1488568468 65, -8.6271227441 41.1479510956 66, -8.6245467411 41.1473018327 67, -8.6241528894 41.1481755307 68, -8.6267076031 41.1488568468 65)))",
				"geometry collection": "GEOMETRYCOLLECTION Z (LINESTRING Z (-8.6107858821 41.145636564 20, -8.6108923285 41.1450754532 21, -8.6141070099 41.1427908808 22, -8.6144157045 41.1411234935 23, -8.6166617236 41.1403298664 24), MULTIPOINT Z (-8.6887871464 41.1733073983 40, -8.6898516104 41.1685638551 41))",
			},
		},
		{
			name: "Parse XYM",
			file: "test/data/all_geom_types_m.fgb",
			want: map[string]string{
				"point":               "POINT M (-8.6291106307 41.1580318815 100)",
				"linestring":          "LINESTRING M (-8.6107858821 41.145636564 200, -8.6108923285 41.1450754532 201, -8.6141070099 41.1427908808 202, -8.6144157045 41.1411234935 203, -8.6166617236 41.1403298664 204)",
				"polygon":             "POLYGON M ((-8.6137876707 41.1459952717 300, -8.6137158194 41.1458069003 301, -8.6145780353 41.1455804532 302, -8.614705771 41.1456004928 303, -8.6147616553 41.1456986867 304, -8.6146339196 41.1457668212 305, -8.6137876707 41.1459952717 300))",
				"multipoint":          "MULTIPOINT M (-8.6887871464 41.1733073983 400, -8.6898516104 41.1685638551 401)",
				"multilinestring":     "MULTILINESTRING M ((-8.6111770726 41.1492896792 500, -8.6112728744 41.1478949867 501, -8.6116347922 41.1460513822 502), (-8.6112728744 41.1460513822 503, -8.6108896673 41.1478308622 504, -8.6107193531 41.1486243986 505, -8.6104213032 41.1491934945 506))",
				"multipolygon":        "MULTIPOLYGON M (((-8.626590512 41.1469491439 600, -8.6265479335 41.1464842331 601, -8.6251960641 41.1463159026 602, -8.6250789731 41.1471575511 603, -8.6264095532 41.1473739733 604, -8.626590512 41.1469491439 600)), ((-8.6267076031 41.1488568468 605, -8.6271227441 41.1479510956 606, -8.6245467411 41.1473018327 607, -8.6241528894 41.1481755307 608, -8.6267076031 41.1488568468 605)))",
				"geometry collection": "GEOMETRYCOLLECTION M (LINESTRING M (-8.6107858821 41.145636564 200, -8.6108923285 41.1450754532 201, -8.6141070099 41.1427908808 202, -8.6144157045 41.1411234935 203, -8.6166617236 41.1403298664 204), MULTIPOINT M (-8.6887871464 41.1733073983 400, -8.6898516104 41.1685638551 401))"},
		},
		{
			name: "Parse XYZM",
			file: "test/data/all_geom_types_zm.fgb",
			want: map[string]string{
				"point":               "POINT ZM (-8.6291106307 41.1580318815 10 100)",
				"linestring":          "LINESTRING ZM (-8.6107858821 41.145636564 20 200, -8.6108923285 41.1450754532 21 201, -8.6141070099 41.1427908808 22 202, -8.6144157045 41.1411234935 23 203, -8.6166617236 41.1403298664 24 204)",
				"polygon":             "POLYGON ZM ((-8.6137876707 41.1459952717 30 300, -8.6137158194 41.1458069003 31 301, -8.6145780353 41.1455804532 32 302, -8.614705771 41.1456004928 33 303, -8.6147616553 41.1456986867 34 304, -8.6146339196 41.1457668212 35 305, -8.6137876707 41.1459952717 30 300))",
				"multipoint":          "MULTIPOINT ZM (-8.6887871464 41.1733073983 40 400, -8.6898516104 41.1685638551 41 401)",
				"multilinestring":     "MULTILINESTRING ZM ((-8.6111770726 41.1492896792 50 500, -8.6112728744 41.1478949867 51 501, -8.6116347922 41.1460513822 52 502), (-8.6112728744 41.1460513822 53 503, -8.6108896673 41.1478308622 54 504, -8.6107193531 41.1486243986 55 505, -8.6104213032 41.1491934945 56 506))",
				"multipolygon":        "MULTIPOLYGON ZM (((-8.626590512 41.1469491439 60 600, -8.6265479335 41.1464842331 61 601, -8.6251960641 41.1463159026 62 602, -8.6250789731 41.1471575511 63 603, -8.6264095532 41.1473739733 64 604, -8.626590512 41.1469491439 60 600)), ((-8.6267076031 41.1488568468 65 605, -8.6271227441 41.1479510956 66 606, -8.6245467411 41.1473018327 67 607, -8.6241528894 41.1481755307 68 608, -8.6267076031 41.1488568468 65 605)))",
				"geometry collection": "GEOMETRYCOLLECTION ZM (LINESTRING ZM (-8.6107858821 41.145636564 20 200, -8.6108923285 41.1450754532 21 201, -8.6141070099 41.1427908808 22 202, -8.6144157045 41.1411234935 23 203, -8.6166617236 41.1403298664 24 204), MULTIPOINT ZM (-8.6887871464 41.1733073983 40 400, -8.6898516104 41.1685638551 41 401))",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := featuresToMap(tt.file); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseGeometry() got = %v, want = %v", got, tt.want)
			}
		})
	}
}

func TestUnsupportedGeometries(t *testing.T) {
	tests := []struct {
		name    string
		file    string
		want    geom.T
		wantErr error
	}{
		{
			name:    "Unsupported 2D types",
			file:    "test/data/unsupported2d_types.fgb",
			want:    nil,
			wantErr: nil,
		},
		{
			name:    "Unsupported 3D types",
			file:    "test/data/unsupported3d_types.fgb",
			want:    nil,
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fgb := readFile(tt.file)
			features := fgb.Features()

			for features.Next() {
				got, err := features.ReadGeometry()

				// TODO: create custom error to avoid weird handling
				if got != tt.want || !strings.HasPrefix(err.Error(), "unable to parse geometry: unsupported geometry type") {
					t.Errorf("Version() got = %v, want %v, error = %v, wantErr %v", got, tt.want, err, tt.wantErr)
				}
			}
		})
	}
}

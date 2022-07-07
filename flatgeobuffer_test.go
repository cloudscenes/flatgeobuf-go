package flatgeobuf_go

import (
	"fmt"
	"github.com/twpayne/go-geom"
	"log"
	"os"
	"reflect"
	"testing"
)

func TestVersion(t *testing.T) {
	tests := []struct {
		name    string
		args    []byte
		want    string
		wantErr error
	}{
		{name: "Valid", args: []byte{0x66, 0x67, 0x62, supportedVersion, 0x66, 0x67, 0x62, 0x01}, want: fmt.Sprintf("%d.0.1", supportedVersion), wantErr: nil},
		{name: "Valid 2", args: []byte{0x66, 0x67, 0x62, supportedVersion, 0x66, 0x67, 0x62, 0x08}, want: fmt.Sprintf("%d.0.8", supportedVersion), wantErr: nil},
		{name: "Invalid Bytes", args: []byte{0x99, 0x67, 0x62, supportedVersion, 0x66, 0x67, 0x62, 0x01}, want: "", wantErr: ErrInvalidFile},
		{name: "Invalid Bytes 2", args: []byte{0x66, 0x67, 0x62, supportedVersion, 0x66, 0x67, 0x31, 0x01}, want: "", wantErr: ErrInvalidFile},
		{name: "Unsupported Version", args: []byte{0x66, 0x67, 0x62, 2, 0x66, 0x67, 0x62, 0x01}, want: "", wantErr: ErrUnsupportedVersion},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Version(tt.args)
			if err != tt.wantErr || got != tt.want {
				t.Errorf("Version() got = %v, want %v, error = %v, wantErr %v", got, tt.want, err, tt.wantErr)
			}
		})
	}
}

// TODO: fix test when ReadAt is implemented
func searchFGB(file string, box []float64) ([]geom.T, []geom.T, error) {
	f, err := os.Open(file)
	defer f.Close()
	if err != nil {
		return nil, nil, err
	}

	fgb, err := NewFGB(f)
	if err != nil {
		return nil, nil, err
	}
	//header := fgb.Header()
	features := fgb.Features()
	searchResult := fgb.Index().Search(box[0], box[1], box[2], box[3])
	searchGeoms := make([]geom.T, len(searchResult))
	for i, _ := range searchResult {
		//feature := features.ReadAt(v.Offset)
		//g, _ := ParseGeometry(feature.Geometry(nil), header.GeometryType(), ParseLayout(header), header.Crs(nil))
		searchGeoms[i] = nil // g
	}

	seqGeoms := make([]geom.T, 0)
	filterBounds := geom.NewBounds(geom.XY).Set(box...)
	for features.Next() {
		feature := features.Read()

		g, err := feature.Geometry()
		if err != nil {
			log.Fatal(err)
		}
		gBounds := geom.NewBounds(geom.XY).Extend(g)
		if gBounds.Overlaps(geom.XY, filterBounds) {
			seqGeoms = append(seqGeoms, g)
		}
	}

	return searchGeoms, seqGeoms, nil
}

func Test_Search(t *testing.T) {
	t.Skip()

	tests := []struct {
		name      string
		file      string
		searchBox []float64
	}{
		{name: "small no results", file: "test/data/simple-small.fgb", searchBox: []float64{-5, -5, -2, -2}},
		{name: "small some results", file: "test/data/simple-small.fgb", searchBox: []float64{1, 1, 2, 2}},
		{name: "small all results", file: "test/data/simple-small.fgb", searchBox: []float64{-1, -1, 20, 20}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			searchRes, seqRes, err := searchFGB(tt.file, tt.searchBox)
			if err != nil {
				t.Errorf("got unexpect error %v", err)
				return
			}
			if len(searchRes) != len(seqRes) {
				t.Errorf("Search len() = %d, want %d", len(searchRes), len(seqRes))
				return
			}
			for i, v := range searchRes {
				if !reflect.DeepEqual(v, seqRes[i]) {
					t.Errorf("search geometries = %v, want %v", v, seqRes[i])
				}
			}
		})
	}
}

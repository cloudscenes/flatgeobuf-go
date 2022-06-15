package flatgeobuf_go

import (
	"github.com/twpayne/go-geom"
	"log"
	"os"
	"reflect"
	"testing"
	"time"
)

func TestUnmarshal(t *testing.T) {
	f, err := os.Open("test/data/alldatatypes.fgb")
	if err != nil {
		log.Fatal(err)
	}
	fgb, err := NewFGB(f)

	features := fgb.Features()
	features.Next()
	feature := features.Read()
	props := feature.Properties()

	type OutStruct struct {
		Geom     geom.T
		Binary   []uint8
		Boolean  bool `fgb:"bool"`
		Byte     int8
		Datetime time.Time
		Double   float64
		Float    float32
		Int      int32
		Json     string
		Long     int64
		Short    int16
		Str      string `fgb:"string"`
		Ubyte    uint8
		Uint     uint32
		Ulong    uint64
		Ushort   uint16
	}

	tests := []struct {
		name  string
		props map[string]interface{}
		want  *OutStruct
	}{
		{
			name:  "alldatatypes",
			props: props,
			want: &OutStruct{
				Geom:     geom.NewPointFlat(geom.XY, []float64{0, 0}),
				Binary:   []uint8{0x58},
				Boolean:  true,
				Byte:     int8(-1),
				Datetime: time.Date(2020, time.February, 29, 12, 34, 56, 0, time.UTC),
				Double:   float64(0),
				Float:    float32(0),
				Int:      int32(-1),
				Json:     "X",
				Long:     int64(-1),
				Short:    int16(-1),
				Str:      "X",
				Ubyte:    uint8(0xff),
				Uint:     uint32(0xffffffff),
				Ulong:    uint64(0xffffffffffffffff),
				Ushort:   uint16(0xffff),
			},
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := &OutStruct{}

			err := feature.Unmarshal(got)
			if err != nil {
				panic(err)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Unmarshal() = %v, want %v", got, tt.want)
			}
		})
	}
}

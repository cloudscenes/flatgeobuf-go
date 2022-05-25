package flatgeobuf_go

import (
	"log"
	"os"
	"reflect"
	"testing"
	"time"
)

func TestPropertyDecoder_Decode(t *testing.T) {
	f, err := os.Open("test/data/alldatatypes.fgb")
	if err != nil {
		log.Fatal(err)
	}
	fgb, err := NewFGB(f)

	header := fgb.Header()
	columns := NewColumns(header)

	features := fgb.Features()
	features.Next()
	feature, err := features.Read()
	if err != nil {
		log.Fatal(err)
	}
	props := feature.Properties()

	tests := []struct {
		name  string
		cols  *Columns
		props map[string]interface{}
		want  map[string]interface{}
	}{
		{
			name:  "alldatatypes",
			cols:  columns,
			props: props,
			want: map[string]interface{}{
				"binary":   []uint8{0x58},
				"bool":     true,
				"byte":     int8(-1),
				"datetime": time.Date(2020, time.February, 29, 12, 34, 56, 0, time.UTC),
				"double":   float64(0),
				"float":    float32(0),
				"int":      int32(-1),
				"json":     "X",
				"long":     int64(-1),
				"short":    int16(-1),
				"string":   "X",
				"ubyte":    uint8(0xff),
				"uint":     uint32(0xffffffff),
				"ulong":    uint64(0xffffffffffffffff),
				"ushort":   uint16(0xffff),
			},
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.props; !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Decode() = %v, want %v", got, tt.want)
			}
		})
	}
}

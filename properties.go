package flatgeobuf_go

import (
	"encoding/binary"
	"flatgeobuf-go/FlatGeobuf"
	"fmt"
	flatbuffers "github.com/google/flatbuffers/go"
)

type PropertyDecoder struct {
	c *Columns
}

func NewPropertyDecoder(c *Columns) *PropertyDecoder {
	return &PropertyDecoder{c}
}

func (pd *PropertyDecoder) Decode(b []byte) {
	pos := uint16(0)
	for i, v := range pd.c.ids {
		size := flatbuffers.GetUint16(b)
		pos += 2
		fmt.Printf("reading field %d(%s) with type %s, width: %d, size: %d\n", i, v.Name(), v.Type().String(), v.Width(), size)

		switch v.Type() {
		case FlatGeobuf.ColumnTypeByte:
		case FlatGeobuf.ColumnTypeUByte:
		case FlatGeobuf.ColumnTypeBool:
		case FlatGeobuf.ColumnTypeShort:
		case FlatGeobuf.ColumnTypeUShort:
		case FlatGeobuf.ColumnTypeInt:
		case FlatGeobuf.ColumnTypeUInt:
		case FlatGeobuf.ColumnTypeLong:
		case FlatGeobuf.ColumnTypeULong:
		case FlatGeobuf.ColumnTypeFloat:
		case FlatGeobuf.ColumnTypeDouble:
		case FlatGeobuf.ColumnTypeJson:
		case FlatGeobuf.ColumnTypeDateTime:
		case FlatGeobuf.ColumnTypeBinary:
			fmt.Println("TODO")
		case FlatGeobuf.ColumnTypeString:
			strSize := binary.LittleEndian.Uint32(b[pos:])
			fmt.Printf(" pos: %d, size: %d\n", pos, strSize)
			pos += 4
			fmt.Println(string(b[pos : pos+uint16(strSize)]))
			pos += uint16(strSize)
		}
		pos += size
	}
}

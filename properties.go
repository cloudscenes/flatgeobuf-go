package flatgeobuf_go

import (
	"flatgeobuf-go/FlatGeobuf"
	"fmt"
	flatbuffers "github.com/google/flatbuffers/go"
	"time"
)

type PropertyDecoder struct {
	c *Columns
}

func NewPropertyDecoder(c *Columns) *PropertyDecoder {
	return &PropertyDecoder{c}
}

func (pd *PropertyDecoder) Decode(b []byte) {
	if b == nil {
		return
	}

	pos := uint16(0)
	for i, v := range pd.c.ids {
		size := flatbuffers.GetUint16(b)
		pos += 2
		fmt.Printf("reading field %d(%s) with type %s, width: %d, size: %d\n", i, v.Name(), v.Type().String(), v.Width(), size)

		switch v.Type() {
		case FlatGeobuf.ColumnTypeByte:
			byteVal := int8(b[pos])
			fmt.Printf(" pos: %d = %d\n", pos, byteVal)
			pos += 1
		case FlatGeobuf.ColumnTypeUByte:
			ubyteVal := b[pos]
			fmt.Printf(" pos: %d = %d\n", pos, ubyteVal)
			pos += 1
		case FlatGeobuf.ColumnTypeBool:
			boolVal := b[pos] != 0
			fmt.Printf(" pos: %d = %v\n", pos, boolVal)
			pos += 1
		case FlatGeobuf.ColumnTypeShort:
			shortVal := flatbuffers.GetInt16(b[pos:])
			fmt.Printf(" pos: %d = %d\n", pos, shortVal)
			pos += flatbuffers.SizeInt16
		case FlatGeobuf.ColumnTypeUShort:
			ushortVal := flatbuffers.GetUint16(b[pos:])
			fmt.Printf(" pos: %d = %d\n", pos, ushortVal)
			pos += flatbuffers.SizeUint16
		case FlatGeobuf.ColumnTypeInt:
			intVal := flatbuffers.GetInt32(b[pos:])
			fmt.Printf(" pos: %d = %d\n", pos, intVal)
			pos += flatbuffers.SizeInt32
		case FlatGeobuf.ColumnTypeUInt:
			uintVal := flatbuffers.GetUint32(b[pos:])
			fmt.Printf(" pos: %d = %d\n", pos, uintVal)
			pos += flatbuffers.SizeUint32
		case FlatGeobuf.ColumnTypeLong:
			longVal := flatbuffers.GetInt64(b[pos:])
			fmt.Printf(" pos: %d = %d\n", pos, longVal)
			pos += flatbuffers.SizeInt64
		case FlatGeobuf.ColumnTypeULong:
			ulongVal := flatbuffers.GetUint64(b[pos:])
			fmt.Printf(" pos: %d = %d\n", pos, ulongVal)
			pos += flatbuffers.SizeUint64
		case FlatGeobuf.ColumnTypeFloat:
			floatVal := flatbuffers.GetFloat32(b[pos:])
			fmt.Printf(" pos: %d = %f\n", pos, floatVal)
			pos += flatbuffers.SizeFloat32
		case FlatGeobuf.ColumnTypeDouble:
			doubleVal := flatbuffers.GetFloat64(b[pos:])
			fmt.Printf(" pos: %d = %f\n", pos, doubleVal)
			pos += flatbuffers.SizeFloat64
		case FlatGeobuf.ColumnTypeString:
			strSize := flatbuffers.GetUint32(b[pos:])
			fmt.Printf(" pos: %d, size: %d\n", pos, strSize)
			pos += 4
			fmt.Println(string(b[pos : pos+uint16(strSize)]))
			pos += uint16(strSize)
		case FlatGeobuf.ColumnTypeJson:
			jsonSize := flatbuffers.GetUint32(b[pos:])
			fmt.Printf(" pos: %d, size: %d\n", pos, jsonSize)
			pos += 4
			fmt.Println(string(b[pos : pos+uint16(jsonSize)]))
			pos += uint16(jsonSize)
		case FlatGeobuf.ColumnTypeDateTime:
			datetimeSize := flatbuffers.GetUint32(b[pos:])
			fmt.Printf(" pos: %d, size: %d\n", pos, datetimeSize)
			pos += 4
			t, err := time.Parse(time.RFC3339, string(b[pos:pos+uint16(datetimeSize)]))
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println(t)
			}
			pos += uint16(datetimeSize)
		case FlatGeobuf.ColumnTypeBinary:
			strSize := flatbuffers.GetUint32(b[pos:])
			fmt.Printf(" pos: %d, size: %d\n", pos, strSize)
			pos += 4
			fmt.Println(string(b[pos : pos+uint16(strSize)]))
			pos += uint16(strSize)
		}
	}
}

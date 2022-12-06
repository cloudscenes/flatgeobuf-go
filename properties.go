package flatgeobuf_go

import (
	"flatgeobuf-go/FlatGeobuf"
	"fmt"
	flatbuffers "github.com/google/flatbuffers/go"
	"time"
)

// PropertyDecoder allows decoding of flatgeobuf properties
// according to the columns in the header
type PropertyDecoder struct {
	c *Columns
}

func NewPropertyDecoder(c *Columns) *PropertyDecoder {
	return &PropertyDecoder{c}
}

func (pd *PropertyDecoder) Decode(b []byte) map[string]interface{} {
	if b == nil {
		return nil
	}

	res := make(map[string]interface{})

	pos := uint16(0)
	for _, v := range pd.c.ids {
		pos += 2

		val, size := pd.decodeVal(b[pos:], v)

		res[v.Name] = val
		pos += size
	}

	return res
}

func (pd *PropertyDecoder) decodeVal(b []byte, v *FlatGeobuf.ColumnT) (interface{}, uint16) {
	var val interface{}
	var size uint16
	switch v.Type {
	case FlatGeobuf.ColumnTypeByte:
		val = int8(b[0])
		size = 1
	case FlatGeobuf.ColumnTypeUByte:
		val = b[0]
		size = 1
	case FlatGeobuf.ColumnTypeBool:
		val = b[0] != 0
		size = 1
	case FlatGeobuf.ColumnTypeShort:
		val = flatbuffers.GetInt16(b)
		size = flatbuffers.SizeInt16
	case FlatGeobuf.ColumnTypeUShort:
		val = flatbuffers.GetUint16(b)
		size = flatbuffers.SizeUint16
	case FlatGeobuf.ColumnTypeInt:
		val = flatbuffers.GetInt32(b)
		size = flatbuffers.SizeInt32
	case FlatGeobuf.ColumnTypeUInt:
		val = flatbuffers.GetUint32(b)
		size = flatbuffers.SizeUint32
	case FlatGeobuf.ColumnTypeLong:
		val = flatbuffers.GetInt64(b)
		size = flatbuffers.SizeInt64
	case FlatGeobuf.ColumnTypeULong:
		val = flatbuffers.GetUint64(b)
		size = flatbuffers.SizeUint64
	case FlatGeobuf.ColumnTypeFloat:
		val = flatbuffers.GetFloat32(b)
		size = flatbuffers.SizeFloat32
	case FlatGeobuf.ColumnTypeDouble:
		val = flatbuffers.GetFloat64(b)
		size = flatbuffers.SizeFloat64
	case FlatGeobuf.ColumnTypeString:
		strSize := flatbuffers.GetUint32(b)
		val = string(b[4 : 4+strSize])
		size = 4 + uint16(strSize)
	case FlatGeobuf.ColumnTypeJson:
		jsonSize := flatbuffers.GetUint32(b)
		val = string(b[4 : 4+jsonSize])
		size = 4 + uint16(jsonSize)
	case FlatGeobuf.ColumnTypeDateTime:
		datetimeSize := flatbuffers.GetUint32(b)
		t, err := time.Parse(time.RFC3339, string(b[4:4+datetimeSize]))
		if err != nil {
			// TODO: should we panic here?
			fmt.Println(err)
		} else {
			val = t
		}
		size = 4 + uint16(datetimeSize)
	case FlatGeobuf.ColumnTypeBinary:
		strSize := flatbuffers.GetUint32(b)
		val = b[4 : 4+strSize]
		size = 4 + uint16(strSize)
	}
	return val, size
}

package flatgeobuf_go

import (
	"flatgeobuf-go/FlatGeobuf"
	"fmt"
	flatbuffers "github.com/google/flatbuffers/go"
	"time"
)

func (f *Feature) decode() map[string]interface{} {
	fgbFeature := f.FeatureT
	if len(fgbFeature.Properties) == 0 {
		return nil
	}

	res := make(map[string]interface{})

	pos := uint16(0)

	columns := f.headerColumns
	if len(columns) == 0 {
		columns = fgbFeature.Columns
	}

	for _, col := range columns {
		pos += 2
		val, size := f.decodeVal(fgbFeature.Properties[pos:], col)

		res[col.Name] = val
		pos += size
	}

	return res
}

// constant that represents size of variable length data (string, json, dateTime, binary)
const varDataSize = flatbuffers.SizeUint32

func (f *Feature) decodeVal(b []byte, col *FlatGeobuf.ColumnT) (interface{}, uint16) {
	var val interface{}
	var size uint16
	switch col.Type {
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
	case FlatGeobuf.ColumnTypeString, FlatGeobuf.ColumnTypeJson:
		strSize := flatbuffers.GetUint32(b)
		val = string(b[varDataSize : varDataSize+strSize])
		size = varDataSize + uint16(strSize)
	case FlatGeobuf.ColumnTypeDateTime:
		datetimeSize := flatbuffers.GetUint32(b)
		t, err := time.Parse(time.RFC3339, string(b[varDataSize:varDataSize+datetimeSize]))
		if err != nil {
			// TODO: should we panic here?
			fmt.Println(err)
		} else {
			val = t
		}
		size = varDataSize + uint16(datetimeSize)
	case FlatGeobuf.ColumnTypeBinary:
		strSize := flatbuffers.GetUint32(b)
		val = b[varDataSize : varDataSize+strSize]
		size = varDataSize + uint16(strSize)
	}
	return val, size
}

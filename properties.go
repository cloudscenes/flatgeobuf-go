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

func (f *Feature) encodeVal(val interface{}, col *FlatGeobuf.ColumnT) ([]byte, error) {
	var res []byte
	var ok bool
	switch col.Type {
	case FlatGeobuf.ColumnTypeByte, FlatGeobuf.ColumnTypeUByte:
		var v byte
		v, ok = val.(byte)
		if !ok {
			break
		}
		res = make([]byte, flatbuffers.SizeByte)
		flatbuffers.WriteByte(res, v)
	case FlatGeobuf.ColumnTypeBool:
		var b bool
		b, ok = val.(bool)
		if !ok {
			break
		}
		res = make([]byte, flatbuffers.SizeBool)
		flatbuffers.WriteBool(res, b)
	case FlatGeobuf.ColumnTypeShort:
		var v int16
		v, ok = val.(int16)
		if !ok {
			break
		}
		res = make([]byte, flatbuffers.SizeInt16)
		flatbuffers.WriteInt16(res, v)
	case FlatGeobuf.ColumnTypeUShort:
		var v uint16
		v, ok = val.(uint16)
		if !ok {
			break
		}
		res = make([]byte, flatbuffers.SizeUint16)
		flatbuffers.WriteUint16(res, v)
	case FlatGeobuf.ColumnTypeInt:
		var v int32
		v, ok = val.(int32)
		if !ok {
			break
		}
		res = make([]byte, flatbuffers.SizeUint32)
		flatbuffers.WriteInt32(res, v)
	case FlatGeobuf.ColumnTypeUInt:
		var v uint32
		v, ok = val.(uint32)
		if !ok {
			break
		}
		res = make([]byte, flatbuffers.SizeUint32)
		flatbuffers.WriteUint32(res, v)
	case FlatGeobuf.ColumnTypeLong:
		var v int64
		v, ok = val.(int64)
		if !ok {
			break
		}
		res = make([]byte, flatbuffers.SizeInt64)
		flatbuffers.WriteInt64(res, v)
	case FlatGeobuf.ColumnTypeULong:
		var v uint64
		v, ok = val.(uint64)
		if !ok {
			break
		}
		res = make([]byte, flatbuffers.SizeUint64)
		flatbuffers.WriteUint64(res, v)
	case FlatGeobuf.ColumnTypeFloat:
		var v float32
		v, ok = val.(float32)
		if !ok {
			break
		}
		res = make([]byte, flatbuffers.SizeFloat64)
		flatbuffers.WriteFloat32(res, v)
	case FlatGeobuf.ColumnTypeDouble:
		var v float64
		v, ok = val.(float64)
		if !ok {
			break
		}
		res = make([]byte, flatbuffers.SizeFloat64)
		flatbuffers.WriteFloat64(res, v)
	case FlatGeobuf.ColumnTypeString, FlatGeobuf.ColumnTypeJson:
		var str string
		str, ok = val.(string)
		if !ok {
			break
		}
		strSize := uint32(len(str))
		res = make([]byte, varDataSize+strSize)
		flatbuffers.WriteUint32(res, strSize)
		copy(res[varDataSize:], str)
	case FlatGeobuf.ColumnTypeDateTime:
		var t time.Time
		t, ok = val.(time.Time)
		if !ok {
			break
		}
		str := t.Format(time.RFC3339)
		strSize := uint32(len(str))
		res = make([]byte, varDataSize+strSize)
		flatbuffers.WriteUint32(res, strSize)
		copy(res[varDataSize:], str)
	case FlatGeobuf.ColumnTypeBinary:
		var b []byte
		b, ok = val.([]byte)
		if !ok {
			break
		}
		bSize := uint32(len(b))
		flatbuffers.WriteUint32(res, bSize)
		copy(res[varDataSize:], b)
	}

	if !ok {
		return nil, fmt.Errorf("invalid value type")
	}
	return res, nil
}

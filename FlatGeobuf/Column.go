// Code generated by the FlatBuffers compiler. DO NOT EDIT.

package FlatGeobuf

import (
	flatbuffers "github.com/google/flatbuffers/go"
)

type ColumnT struct {
	Name string
	Type ColumnType
	Title string
	Description string
	Width int32
	Precision int32
	Scale int32
	Nullable bool
	Unique bool
	PrimaryKey bool
	Metadata string
}

func (t *ColumnT) Pack(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	if t == nil { return 0 }
	nameOffset := builder.CreateString(t.Name)
	titleOffset := builder.CreateString(t.Title)
	descriptionOffset := builder.CreateString(t.Description)
	metadataOffset := builder.CreateString(t.Metadata)
	ColumnStart(builder)
	ColumnAddName(builder, nameOffset)
	ColumnAddType(builder, t.Type)
	ColumnAddTitle(builder, titleOffset)
	ColumnAddDescription(builder, descriptionOffset)
	ColumnAddWidth(builder, t.Width)
	ColumnAddPrecision(builder, t.Precision)
	ColumnAddScale(builder, t.Scale)
	ColumnAddNullable(builder, t.Nullable)
	ColumnAddUnique(builder, t.Unique)
	ColumnAddPrimaryKey(builder, t.PrimaryKey)
	ColumnAddMetadata(builder, metadataOffset)
	return ColumnEnd(builder)
}

func (rcv *Column) UnPackTo(t *ColumnT) {
	t.Name = string(rcv.Name())
	t.Type = rcv.Type()
	t.Title = string(rcv.Title())
	t.Description = string(rcv.Description())
	t.Width = rcv.Width()
	t.Precision = rcv.Precision()
	t.Scale = rcv.Scale()
	t.Nullable = rcv.Nullable()
	t.Unique = rcv.Unique()
	t.PrimaryKey = rcv.PrimaryKey()
	t.Metadata = string(rcv.Metadata())
}

func (rcv *Column) UnPack() *ColumnT {
	if rcv == nil { return nil }
	t := &ColumnT{}
	rcv.UnPackTo(t)
	return t
}

type Column struct {
	_tab flatbuffers.Table
}

func GetRootAsColumn(buf []byte, offset flatbuffers.UOffsetT) *Column {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &Column{}
	x.Init(buf, n+offset)
	return x
}

func (rcv *Column) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *Column) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *Column) Name() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *Column) Type() ColumnType {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return ColumnType(rcv._tab.GetByte(o + rcv._tab.Pos))
	}
	return 0
}

func (rcv *Column) MutateType(n ColumnType) bool {
	return rcv._tab.MutateByteSlot(6, byte(n))
}

func (rcv *Column) Title() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *Column) Description() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(10))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *Column) Width() int32 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(12))
	if o != 0 {
		return rcv._tab.GetInt32(o + rcv._tab.Pos)
	}
	return -1
}

func (rcv *Column) MutateWidth(n int32) bool {
	return rcv._tab.MutateInt32Slot(12, n)
}

func (rcv *Column) Precision() int32 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(14))
	if o != 0 {
		return rcv._tab.GetInt32(o + rcv._tab.Pos)
	}
	return -1
}

func (rcv *Column) MutatePrecision(n int32) bool {
	return rcv._tab.MutateInt32Slot(14, n)
}

func (rcv *Column) Scale() int32 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(16))
	if o != 0 {
		return rcv._tab.GetInt32(o + rcv._tab.Pos)
	}
	return -1
}

func (rcv *Column) MutateScale(n int32) bool {
	return rcv._tab.MutateInt32Slot(16, n)
}

func (rcv *Column) Nullable() bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(18))
	if o != 0 {
		return rcv._tab.GetBool(o + rcv._tab.Pos)
	}
	return true
}

func (rcv *Column) MutateNullable(n bool) bool {
	return rcv._tab.MutateBoolSlot(18, n)
}

func (rcv *Column) Unique() bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(20))
	if o != 0 {
		return rcv._tab.GetBool(o + rcv._tab.Pos)
	}
	return false
}

func (rcv *Column) MutateUnique(n bool) bool {
	return rcv._tab.MutateBoolSlot(20, n)
}

func (rcv *Column) PrimaryKey() bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(22))
	if o != 0 {
		return rcv._tab.GetBool(o + rcv._tab.Pos)
	}
	return false
}

func (rcv *Column) MutatePrimaryKey(n bool) bool {
	return rcv._tab.MutateBoolSlot(22, n)
}

func (rcv *Column) Metadata() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(24))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func ColumnStart(builder *flatbuffers.Builder) {
	builder.StartObject(11)
}
func ColumnAddName(builder *flatbuffers.Builder, name flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(0, flatbuffers.UOffsetT(name), 0)
}
func ColumnAddType(builder *flatbuffers.Builder, type_ ColumnType) {
	builder.PrependByteSlot(1, byte(type_), 0)
}
func ColumnAddTitle(builder *flatbuffers.Builder, title flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(2, flatbuffers.UOffsetT(title), 0)
}
func ColumnAddDescription(builder *flatbuffers.Builder, description flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(3, flatbuffers.UOffsetT(description), 0)
}
func ColumnAddWidth(builder *flatbuffers.Builder, width int32) {
	builder.PrependInt32Slot(4, width, -1)
}
func ColumnAddPrecision(builder *flatbuffers.Builder, precision int32) {
	builder.PrependInt32Slot(5, precision, -1)
}
func ColumnAddScale(builder *flatbuffers.Builder, scale int32) {
	builder.PrependInt32Slot(6, scale, -1)
}
func ColumnAddNullable(builder *flatbuffers.Builder, nullable bool) {
	builder.PrependBoolSlot(7, nullable, true)
}
func ColumnAddUnique(builder *flatbuffers.Builder, unique bool) {
	builder.PrependBoolSlot(8, unique, false)
}
func ColumnAddPrimaryKey(builder *flatbuffers.Builder, primaryKey bool) {
	builder.PrependBoolSlot(9, primaryKey, false)
}
func ColumnAddMetadata(builder *flatbuffers.Builder, metadata flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(10, flatbuffers.UOffsetT(metadata), 0)
}
func ColumnEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}

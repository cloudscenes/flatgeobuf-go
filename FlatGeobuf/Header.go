// Code generated by the FlatBuffers compiler. DO NOT EDIT.

package FlatGeobuf

import (
	flatbuffers "github.com/google/flatbuffers/go"
)

type HeaderT struct {
	Name string
	Envelope []float64
	GeometryType GeometryType
	HasZ bool
	HasM bool
	HasT bool
	HasTm bool
	Columns []*ColumnT
	FeaturesCount uint64
	IndexNodeSize uint16
	Crs *CrsT
	Title string
	Description string
	Metadata string
}

func (t *HeaderT) Pack(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	if t == nil { return 0 }
	nameOffset := builder.CreateString(t.Name)
	envelopeOffset := flatbuffers.UOffsetT(0)
	if t.Envelope != nil {
		envelopeLength := len(t.Envelope)
		HeaderStartEnvelopeVector(builder, envelopeLength)
		for j := envelopeLength - 1; j >= 0; j-- {
			builder.PrependFloat64(t.Envelope[j])
		}
		envelopeOffset = builder.EndVector(envelopeLength)
	}
	columnsOffset := flatbuffers.UOffsetT(0)
	if t.Columns != nil {
		columnsLength := len(t.Columns)
		columnsOffsets := make([]flatbuffers.UOffsetT, columnsLength)
		for j := 0; j < columnsLength; j++ {
			columnsOffsets[j] = t.Columns[j].Pack(builder)
		}
		HeaderStartColumnsVector(builder, columnsLength)
		for j := columnsLength - 1; j >= 0; j-- {
			builder.PrependUOffsetT(columnsOffsets[j])
		}
		columnsOffset = builder.EndVector(columnsLength)
	}
	crsOffset := t.Crs.Pack(builder)
	titleOffset := builder.CreateString(t.Title)
	descriptionOffset := builder.CreateString(t.Description)
	metadataOffset := builder.CreateString(t.Metadata)
	HeaderStart(builder)
	HeaderAddName(builder, nameOffset)
	HeaderAddEnvelope(builder, envelopeOffset)
	HeaderAddGeometryType(builder, t.GeometryType)
	HeaderAddHasZ(builder, t.HasZ)
	HeaderAddHasM(builder, t.HasM)
	HeaderAddHasT(builder, t.HasT)
	HeaderAddHasTm(builder, t.HasTm)
	HeaderAddColumns(builder, columnsOffset)
	HeaderAddFeaturesCount(builder, t.FeaturesCount)
	HeaderAddIndexNodeSize(builder, t.IndexNodeSize)
	HeaderAddCrs(builder, crsOffset)
	HeaderAddTitle(builder, titleOffset)
	HeaderAddDescription(builder, descriptionOffset)
	HeaderAddMetadata(builder, metadataOffset)
	return HeaderEnd(builder)
}

func (rcv *Header) UnPackTo(t *HeaderT) {
	t.Name = string(rcv.Name())
	envelopeLength := rcv.EnvelopeLength()
	t.Envelope = make([]float64, envelopeLength)
	for j := 0; j < envelopeLength; j++ {
		t.Envelope[j] = rcv.Envelope(j)
	}
	t.GeometryType = rcv.GeometryType()
	t.HasZ = rcv.HasZ()
	t.HasM = rcv.HasM()
	t.HasT = rcv.HasT()
	t.HasTm = rcv.HasTm()
	columnsLength := rcv.ColumnsLength()
	t.Columns = make([]*ColumnT, columnsLength)
	for j := 0; j < columnsLength; j++ {
		x := Column{}
		rcv.Columns(&x, j)
		t.Columns[j] = x.UnPack()
	}
	t.FeaturesCount = rcv.FeaturesCount()
	t.IndexNodeSize = rcv.IndexNodeSize()
	t.Crs = rcv.Crs(nil).UnPack()
	t.Title = string(rcv.Title())
	t.Description = string(rcv.Description())
	t.Metadata = string(rcv.Metadata())
}

func (rcv *Header) UnPack() *HeaderT {
	if rcv == nil { return nil }
	t := &HeaderT{}
	rcv.UnPackTo(t)
	return t
}

type Header struct {
	_tab flatbuffers.Table
}

func GetRootAsHeader(buf []byte, offset flatbuffers.UOffsetT) *Header {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &Header{}
	x.Init(buf, n+offset)
	return x
}

func (rcv *Header) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *Header) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *Header) Name() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *Header) Envelope(j int) float64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		a := rcv._tab.Vector(o)
		return rcv._tab.GetFloat64(a + flatbuffers.UOffsetT(j*8))
	}
	return 0
}

func (rcv *Header) EnvelopeLength() int {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.VectorLen(o)
	}
	return 0
}

func (rcv *Header) MutateEnvelope(j int, n float64) bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		a := rcv._tab.Vector(o)
		return rcv._tab.MutateFloat64(a+flatbuffers.UOffsetT(j*8), n)
	}
	return false
}

func (rcv *Header) GeometryType() GeometryType {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		return GeometryType(rcv._tab.GetByte(o + rcv._tab.Pos))
	}
	return 0
}

func (rcv *Header) MutateGeometryType(n GeometryType) bool {
	return rcv._tab.MutateByteSlot(8, byte(n))
}

func (rcv *Header) HasZ() bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(10))
	if o != 0 {
		return rcv._tab.GetBool(o + rcv._tab.Pos)
	}
	return false
}

func (rcv *Header) MutateHasZ(n bool) bool {
	return rcv._tab.MutateBoolSlot(10, n)
}

func (rcv *Header) HasM() bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(12))
	if o != 0 {
		return rcv._tab.GetBool(o + rcv._tab.Pos)
	}
	return false
}

func (rcv *Header) MutateHasM(n bool) bool {
	return rcv._tab.MutateBoolSlot(12, n)
}

func (rcv *Header) HasT() bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(14))
	if o != 0 {
		return rcv._tab.GetBool(o + rcv._tab.Pos)
	}
	return false
}

func (rcv *Header) MutateHasT(n bool) bool {
	return rcv._tab.MutateBoolSlot(14, n)
}

func (rcv *Header) HasTm() bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(16))
	if o != 0 {
		return rcv._tab.GetBool(o + rcv._tab.Pos)
	}
	return false
}

func (rcv *Header) MutateHasTm(n bool) bool {
	return rcv._tab.MutateBoolSlot(16, n)
}

func (rcv *Header) Columns(obj *Column, j int) bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(18))
	if o != 0 {
		x := rcv._tab.Vector(o)
		x += flatbuffers.UOffsetT(j) * 4
		x = rcv._tab.Indirect(x)
		obj.Init(rcv._tab.Bytes, x)
		return true
	}
	return false
}

func (rcv *Header) ColumnsLength() int {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(18))
	if o != 0 {
		return rcv._tab.VectorLen(o)
	}
	return 0
}

func (rcv *Header) FeaturesCount() uint64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(20))
	if o != 0 {
		return rcv._tab.GetUint64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Header) MutateFeaturesCount(n uint64) bool {
	return rcv._tab.MutateUint64Slot(20, n)
}

func (rcv *Header) IndexNodeSize() uint16 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(22))
	if o != 0 {
		return rcv._tab.GetUint16(o + rcv._tab.Pos)
	}
	return 16
}

func (rcv *Header) MutateIndexNodeSize(n uint16) bool {
	return rcv._tab.MutateUint16Slot(22, n)
}

func (rcv *Header) Crs(obj *Crs) *Crs {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(24))
	if o != 0 {
		x := rcv._tab.Indirect(o + rcv._tab.Pos)
		if obj == nil {
			obj = new(Crs)
		}
		obj.Init(rcv._tab.Bytes, x)
		return obj
	}
	return nil
}

func (rcv *Header) Title() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(26))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *Header) Description() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(28))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *Header) Metadata() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(30))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func HeaderStart(builder *flatbuffers.Builder) {
	builder.StartObject(14)
}
func HeaderAddName(builder *flatbuffers.Builder, name flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(0, flatbuffers.UOffsetT(name), 0)
}
func HeaderAddEnvelope(builder *flatbuffers.Builder, envelope flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(1, flatbuffers.UOffsetT(envelope), 0)
}
func HeaderStartEnvelopeVector(builder *flatbuffers.Builder, numElems int) flatbuffers.UOffsetT {
	return builder.StartVector(8, numElems, 8)
}
func HeaderAddGeometryType(builder *flatbuffers.Builder, geometryType GeometryType) {
	builder.PrependByteSlot(2, byte(geometryType), 0)
}
func HeaderAddHasZ(builder *flatbuffers.Builder, hasZ bool) {
	builder.PrependBoolSlot(3, hasZ, false)
}
func HeaderAddHasM(builder *flatbuffers.Builder, hasM bool) {
	builder.PrependBoolSlot(4, hasM, false)
}
func HeaderAddHasT(builder *flatbuffers.Builder, hasT bool) {
	builder.PrependBoolSlot(5, hasT, false)
}
func HeaderAddHasTm(builder *flatbuffers.Builder, hasTm bool) {
	builder.PrependBoolSlot(6, hasTm, false)
}
func HeaderAddColumns(builder *flatbuffers.Builder, columns flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(7, flatbuffers.UOffsetT(columns), 0)
}
func HeaderStartColumnsVector(builder *flatbuffers.Builder, numElems int) flatbuffers.UOffsetT {
	return builder.StartVector(4, numElems, 4)
}
func HeaderAddFeaturesCount(builder *flatbuffers.Builder, featuresCount uint64) {
	builder.PrependUint64Slot(8, featuresCount, 0)
}
func HeaderAddIndexNodeSize(builder *flatbuffers.Builder, indexNodeSize uint16) {
	builder.PrependUint16Slot(9, indexNodeSize, 16)
}
func HeaderAddCrs(builder *flatbuffers.Builder, crs flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(10, flatbuffers.UOffsetT(crs), 0)
}
func HeaderAddTitle(builder *flatbuffers.Builder, title flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(11, flatbuffers.UOffsetT(title), 0)
}
func HeaderAddDescription(builder *flatbuffers.Builder, description flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(12, flatbuffers.UOffsetT(description), 0)
}
func HeaderAddMetadata(builder *flatbuffers.Builder, metadata flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(13, flatbuffers.UOffsetT(metadata), 0)
}
func HeaderEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}

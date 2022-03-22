package flatgeobuf_go

import (
	"flatgeobuf-go/FlatGeobuf"
	"fmt"
	"strings"
)

type Columns struct {
	names map[string]*FlatGeobuf.Column
	ids   []*FlatGeobuf.Column
}

func NewColumns(header *FlatGeobuf.Header) *Columns {

	names := make(map[string]*FlatGeobuf.Column)
	colLen := header.ColumnsLength()
	ids := make([]*FlatGeobuf.Column, colLen)

	for i := 0; i < colLen; i++ {
		var c FlatGeobuf.Column
		header.Columns(&c, i)
		name := string(c.Name())
		names[name] = &c
		ids[i] = &c
	}

	return &Columns{
		names: names,
		ids:   ids,
	}
}

func (cols *Columns) String() string {
	var b strings.Builder

	for i, v := range cols.ids {
		b.WriteString(fmt.Sprintf("%d: %s\n", i, string(v.Name())))
	}

	for k, v := range cols.names {
		b.WriteString(k)
		b.WriteString(" : ")
		b.WriteString(v.Type().String())
		b.WriteString("\n")
	}

	return b.String()
}

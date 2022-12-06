package flatgeobuf_go

import (
	"flatgeobuf-go/FlatGeobuf"
	"fmt"
	"strings"
)

type Columns struct {
	names map[string]*FlatGeobuf.ColumnT
	ids   []*FlatGeobuf.ColumnT
}

func NewColumns(header *FlatGeobuf.HeaderT) *Columns {
	names := make(map[string]*FlatGeobuf.ColumnT)
	ids := make([]*FlatGeobuf.ColumnT, len(header.Columns))

	for i, column := range header.Columns {
		name := column.Name
		names[name] = column
		ids[i] = column
	}

	return &Columns{
		names: names,
		ids:   ids,
	}
}

func (cols *Columns) String() string {
	var b strings.Builder

	for i, v := range cols.ids {
		b.WriteString(fmt.Sprintf("%d: %s\n", i, string(v.Name)))
	}

	for k, v := range cols.names {
		b.WriteString(k)
		b.WriteString(" : ")
		b.WriteString(v.Type.String())
		b.WriteString("\n")
	}

	return b.String()
}

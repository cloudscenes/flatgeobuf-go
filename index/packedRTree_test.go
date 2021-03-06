package index

import (
	"reflect"
	"testing"
)

func TestCalcTreeSize(t *testing.T) {
	type args struct {
		numItems uint64
		nodeSize uint16
	}
	tests := []struct {
		name    string
		args    args
		want    uint64
		wantErr bool
	}{
		{name: "small", args: args{numItems: 2, nodeSize: 16}, want: 120, wantErr: false},
		{name: "medium", args: args{numItems: 2000, nodeSize: 16}, want: 85360, wantErr: false},
		{name: "not enough items", args: args{numItems: 0, nodeSize: 16}, want: 123456, wantErr: true},
		{name: "node size too small", args: args{numItems: 20, nodeSize: 1}, want: 123456, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CalcTreeSize(tt.args.numItems, tt.args.nodeSize)
			if (err != nil) != tt.wantErr {
				t.Errorf("CalcTreeSize() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && got != tt.want {
				t.Errorf("CalcTreeSize() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNodeItem_intersects(t *testing.T) {
	tests := []struct {
		name   string
		target NodeItem
		args   NodeItem
		want   bool
	}{
		{
			name:   "no intersection",
			target: NodeItem{0, 0, 10, 10, 0},
			args:   NodeItem{11, 11, 20, 20, 0},
			want:   false,
		},
		{
			name:   "no intersection Y",
			target: NodeItem{0, 0, 10, 10, 0},
			args:   NodeItem{0, 11, 20, 20, 0},
			want:   false,
		},
		{
			name:   "no intersection 2",
			target: NodeItem{11, 11, 20, 20, 0},
			args:   NodeItem{0, 0, 10, 10, 0},
			want:   false,
		},
		{
			name:   "no intersection Y 2",
			target: NodeItem{0, 11, 20, 20, 0},
			args:   NodeItem{0, 0, 10, 10, 0},
			want:   false,
		},
		{
			name:   "intersects",
			target: NodeItem{0, 0, 10, 10, 0},
			args:   NodeItem{5, 5, 20, 20, 0},
			want:   true,
		},
		{
			name:   "contained",
			target: NodeItem{0, 0, 10, 10, 0},
			args:   NodeItem{2, 2, 8, 8, 0},
			want:   true,
		},
		{
			name:   "contains",
			target: NodeItem{2, 2, 8, 8, 0},
			args:   NodeItem{0, 0, 10, 10, 0},
			want:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.target.intersects(tt.args); got != tt.want {
				t.Errorf("intersects() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNodeItem_expand(t *testing.T) {
	tests := []struct {
		name   string
		target *NodeItem
		args   NodeItem
		want   NodeItem
	}{
		{
			name:   "equal",
			target: &NodeItem{0, 0, 10, 10, 0},
			args:   NodeItem{0, 0, 10, 10, 0},
			want:   NodeItem{0, 0, 10, 10, 0},
		},
		{
			name:   "contained",
			target: &NodeItem{0, 0, 10, 10, 0},
			args:   NodeItem{3, 3, 8, 8, 0},
			want:   NodeItem{0, 0, 10, 10, 0},
		},
		{
			name:   "contains",
			target: &NodeItem{3, 3, 8, 8, 0},
			args:   NodeItem{0, 0, 10, 10, 0},
			want:   NodeItem{0, 0, 10, 10, 0},
		},
		{
			name:   "disjoint",
			target: &NodeItem{0, 0, 3, 3, 0},
			args:   NodeItem{8, 8, 10, 10, 0},
			want:   NodeItem{0, 0, 10, 10, 0},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.target.expand(tt.args)
			if !reflect.DeepEqual(*tt.target, tt.want) {
				t.Errorf("expand() = %v, want %v", *tt.target, tt.want)
			}
		})
	}
}

func TestNodeItem_ReadFromBytes(t *testing.T) {
	tests := []struct {
		name string
		args []byte
		want NodeItem
	}{
		{
			name: "simple box",
			args: []byte{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x8, 0x40, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x8, 0x40, 0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
			want: NodeItem{0, 0, 3, 3, 1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ni := NodeItem{}
			ni.readFromBytes(tt.args)
			if !reflect.DeepEqual(ni, tt.want) {
				t.Errorf("read from bytes() = %v, want %v", ni, tt.want)
			}
		})
	}
}

func TestPackedRTree_generateLevelBounds(t *testing.T) {
	tests := []struct {
		name  string
		input PackedRTree
		want  []LevelBounds
	}{
		{
			name:  "simple",
			input: PackedRTree{numItems: 2, nodeSize: 16, levelBounds: nil},
			want:  []LevelBounds{{0, 1}, {1, 3}},
		},
		{
			name:  "bigger",
			input: PackedRTree{numItems: 200, nodeSize: 16, levelBounds: nil},
			want:  []LevelBounds{{0, 1}, {1, 14}, {14, 214}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.input.generateLevelBounds()
			if !reflect.DeepEqual(tt.input.levelBounds, tt.want) {
				t.Errorf("generate level bounds() = %v, want %v", tt.input.levelBounds, tt.want)
			}
		})
	}
}

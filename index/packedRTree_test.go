package index

import "testing"

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

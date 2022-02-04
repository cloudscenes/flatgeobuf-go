package flatgeobuf_go

import (
	"fmt"
	"testing"
)

func TestVersion(t *testing.T) {
	type args struct {
		fileMagicBytes []byte
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr error
	}{
		{name: "Valid", args: args{fileMagicBytes: []byte{0x66, 0x67, 0x62, supportedVersion, 0x66, 0x67, 0x62, 0x01}}, want: fmt.Sprintf("%d.0.1", supportedVersion), wantErr: nil},
		{name: "Valid 2", args: args{fileMagicBytes: []byte{0x66, 0x67, 0x62, supportedVersion, 0x66, 0x67, 0x62, 0x08}}, want: fmt.Sprintf("%d.0.8", supportedVersion), wantErr: nil},
		{name: "Invalid Bytes", args: args{fileMagicBytes: []byte{0x99, 0x67, 0x62, supportedVersion, 0x66, 0x67, 0x62, 0x01}}, want: "", wantErr: ErrInvalidFile},
		{name: "Invalid Bytes 2", args: args{fileMagicBytes: []byte{0x66, 0x67, 0x62, supportedVersion, 0x66, 0x67, 0x31, 0x01}}, want: "", wantErr: ErrInvalidFile},
		{name: "Unsupported Version", args: args{fileMagicBytes: []byte{0x66, 0x67, 0x62, 2, 0x66, 0x67, 0x62, 0x01}}, want: "", wantErr: ErrUnsupportedVersion},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Version(tt.args.fileMagicBytes)
			if err != tt.wantErr || got != tt.want {
				t.Errorf("Version() got = %v, want %v, error = %v, wantErr %v", got, tt.want, err, tt.wantErr)
			}
		})
	}
}
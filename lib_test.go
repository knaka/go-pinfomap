package pinfomap

import (
	"github.com/knaka/go-pinfomap/stringer/testdata"
	"golang.org/x/tools/go/packages"
	"testing"
)

func TestNewIntTypeInfo(t *testing.T) {
	type args struct {
		target any
	}
	tests := []struct {
		name            string
		args            args
		wantIntConstNum int
		wantErr         bool
	}{
		{"Test 1", args{packages.LoadMode(0)}, 16, false},
		{"Test 2", args{"./stringer/testdata.Pill"}, 4, false},
		{"Test 3", args{testdata.Fruit(0)}, 3, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotIntType, err := NewIntTypeInfo(tt.args.target)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewIntTypeInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(gotIntType.Consts) != tt.wantIntConstNum {
				t.Errorf("NewIntTypeInfo() gotIntType.Consts = %v, want %v", len(gotIntType.Consts), tt.wantIntConstNum)
				return
			}
		})
	}
}

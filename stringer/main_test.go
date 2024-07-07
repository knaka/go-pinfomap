package stringer

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"
)

func TestNewIntTypeInfo(t *testing.T) {
	type args struct {
		packageDirOrPath string
		typeName         string
	}
	filepathSeparator := fmt.Sprintf("%c", filepath.Separator)
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"Test 1", args{strings.Join([]string{".", "testdata"}, filepathSeparator), "Pill"}, false},
		{"Test 2", args{"golang.org/x/tools/go/packages", "LoadMode"}, false},
		{"Test 3", args{"example.com/not/found", "FooBar"}, true},
		{"Test 4", args{"golang.org/x/tools/go/packages", "FooBar"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := NewIntTypeInfo(tt.args.packageDirOrPath, tt.args.typeName); (err != nil) != tt.wantErr {
				t.Errorf("NewIntTypeInfo() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

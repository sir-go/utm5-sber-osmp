package utm

import (
	"testing"
)

func TestHasDigitsOnly(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"empty", args{""}, false},
		{"yes", args{"112141"}, true},
		{"no", args{"112 141 "}, false},
		{"no", args{"- rr141 "}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HasDigitsOnly(tt.args.s); got != tt.want {
				t.Errorf("HasDigitsOnly() = %v, want %v", got, tt.want)
			}
		})
	}
}

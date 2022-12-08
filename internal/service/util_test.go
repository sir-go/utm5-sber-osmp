package service

import (
	"testing"
)

func TestRoundBalance(t *testing.T) {
	type args struct {
		b float64
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{"zero", args{0.0}, 0.0},
		{"", args{-3999.0212813620074}, -4000},
		{"", args{-3999.36895}, -4000},
		{"", args{-3999.01234}, -4000},
		{"", args{3999.36895}, 3999},
		{"", args{3999.01234}, 3999},
		{"", args{-0.00021}, -1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RoundBalance(tt.args.b); got != tt.want {
				t.Errorf("RoundBalance() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRoundRecSum(t *testing.T) {
	type args struct {
		b float64
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{"zero", args{0.0}, 0.0},
		{"", args{-3999.0212813620074}, 4000},
		{"", args{-3999.36895}, 4000},
		{"", args{-3999.01234}, 4000},
		{"", args{3999.36895}, 4000},
		{"", args{3999.01234}, 4000},
		{"", args{-0.00021}, 1},
		{"", args{0.00021}, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RoundRecSum(tt.args.b); got != tt.want {
				t.Errorf("RoundRecSum() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRightPadID(t *testing.T) {
	type args struct {
		prefix   int
		id       int
		totalLen int
	}
	tests := []struct {
		name    string
		args    args
		want    int64
		wantErr bool
	}{
		{"empty", args{0, 0, 0}, 0, true},
		{"errTooBig", args{0, 123, 4}, 0, true},
		{"errTooBig", args{1, 123, 4}, 0, true},
		{"ok8", args{666, 732, 8}, 66600732, false},
		{"ok14", args{666, 732, 14}, 66600000000732, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := RightPadID(tt.args.prefix, tt.args.id, tt.args.totalLen)
			if (err != nil) != tt.wantErr {
				t.Errorf("RightPadID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("RightPadID() got = %v, want %v", got, tt.want)
			}
		})
	}
}

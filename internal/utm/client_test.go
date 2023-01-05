package utm

import (
	"testing"
)

func TestClient_GetPrefixByExtID(t *testing.T) {
	Prefixes := map[string]*Prefix{"tih": {Api: "tih"}, "kor": {Api: "kor"}}
	tests := []struct {
		name       string
		extId      string
		wantApi    string
		wantAidInt int
	}{
		{"empty", "", "tih", 0},
		{"tih", "1342011", "tih", 0},
		{"tih", "ext134", "tih", 0},
		{"kor", "134", "kor", 134},
		{"kor", "11234", "kor", 11234},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{cfg: Config{Prefixes: Prefixes}}
			gotPref, gotAidInt := c.GetPrefixByExtID(tt.extId)
			if gotPref.Api != tt.wantApi {
				t.Errorf("GetPrefixByExtID() gotPref.api = %v, want %v", gotPref.Api, tt.wantApi)
			}
			if gotAidInt != tt.wantAidInt {
				t.Errorf("GetPrefixByExtID() gotAidInt = %v, want %v", gotAidInt, tt.wantAidInt)
			}
		})
	}
}

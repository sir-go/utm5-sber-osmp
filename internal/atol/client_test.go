package atol

import (
	"bytes"
	"testing"
	"text/template"
)

func Test_makeUUID(t *testing.T) {
	tests := []struct {
		name    string
		wantLen int
	}{
		{"lenTest", 32},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := makeUUID(); len(got) != tt.wantLen {
				t.Errorf("len(makeUUID()) = %v, want %v", got, tt.wantLen)
			}
		})
	}
}

func TestClient_makeBody(t *testing.T) {
	tests := []struct {
		name       string
		tmplString string
		vars       map[string]interface{}
		wantBody   []byte
		wantErr    bool
	}{
		{"empty",
			"", nil, nil, true},
		{"plain",
			"some template text", nil, []byte("some template text"), false},
		{"templated empty",
			"some {{.w}} text", nil, []byte("some <no value> text"), false},
		{"templated",
			"some {{.w}} text", map[string]interface{}{
				"w": "filled",
			}, []byte("some filled text"), false,
		},
		{"e2e",
			`
{
	"uuid": "{{.taskId}}",
	"request": {
		"type": "sell",
		"items": [{
			"tax": {"type": "none"},
			"name": "абон. плата за Интернет [{{.payId}}]",
			"type": "position",
			"price": "{{printf "%.2f" .amount}}",
			"amount": "{{printf "%.2f" .amount}}",
			"quantity": 1
		}],
		"total": "{{printf "%.2f" .amount}}",
		"operator": {"name": "{{.operator}}"},
		"payments": [{"sum": {{printf "%.2f" .amount}}, "type": "electronically"}],
		"clientInfo": {"emailOrPhone": "{{.contact}}"},
		"paymentsPlace": "{{.place}}",
		"electronically": true
	}
}
`, map[string]interface{}{
				"taskId":   "fb129654-4b8a-4c1f-a748-b06a874111a3",
				"payId":    "e03079a23f4622d0f1f80d1ff929e3a5",
				"amount":   123.456,
				"operator": "John Smith",
				"contact":  "some@email.com",
				"place":    "office",
			},
			[]byte(
				`
{
	"uuid": "fb129654-4b8a-4c1f-a748-b06a874111a3",
	"request": {
		"type": "sell",
		"items": [{
			"tax": {"type": "none"},
			"name": "абон. плата за Интернет [e03079a23f4622d0f1f80d1ff929e3a5]",
			"type": "position",
			"price": "123.46",
			"amount": "123.46",
			"quantity": 1
		}],
		"total": "123.46",
		"operator": {"name": "John Smith"},
		"payments": [{"sum": 123.46, "type": "electronically"}],
		"clientInfo": {"emailOrPhone": "some@email.com"},
		"paymentsPlace": "office",
		"electronically": true
	}
}
`), false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpl, err := template.New(tt.name).Parse(tt.tmplString)
			if err != nil {
				t.Errorf("template parsing error, %v", err)
				return
			}

			c := &Client{reqTmpl: tmpl}
			gotBody, err := c.makeBody(tt.vars)
			if (err != nil) != tt.wantErr {
				t.Errorf("makeBody() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				return
			}
			if !bytes.Equal(gotBody.Bytes(), tt.wantBody) {
				t.Errorf("makeBody() gotBody = %v, want %v", gotBody.String(), string(tt.wantBody))
			}
		})
	}
}

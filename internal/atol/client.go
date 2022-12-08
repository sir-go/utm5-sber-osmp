package atol

import (
	"bytes"
	"encoding/hex"
	"net/http"
	"text/template"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type Client struct {
	Config
	reqTmpl *template.Template
}

func makeUUID() string {
	strUUID := uuid.NewString()
	if id, err := uuid.NewRandom(); err != nil {
		return strUUID
	} else {
		if binID, err := id.MarshalBinary(); err != nil {
			return strUUID
		} else {
			return hex.EncodeToString(binID)
		}
	}
}

func (c *Client) makeBody(vars map[string]interface{}) (body *bytes.Buffer, err error) {
	body = bytes.NewBuffer(make([]byte, 0))
	err = c.reqTmpl.Execute(body, vars)
	return
}

func NewClient(config Config) (*Client, error) {
	tmpl, err := template.New("atol-request").Parse(config.RequestTemplate)
	if err != nil {
		return nil, err
	}
	return &Client{config, tmpl}, nil
}

func (c *Client) MakeFiscal(payId int64, amount float64, operator string, contact string) (taskId string, err error) {
	var resp *http.Response
	taskId = makeUUID()

	body, err := c.makeBody(map[string]interface{}{
		"taskId":   taskId,
		"payId":    payId,
		"amount":   amount,
		"operator": operator,
		"contact":  contact,
		"place":    operator,
	})
	if err != nil {
		return "", err
	}

	resp, err = http.Post(c.ApiURL, "application/json", body)
	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		err = errors.New("atol api response status: " + resp.Status)
		return
	}

	return
}

package main

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

const DefaultContact = "cashbox@tih.ru"

func MakeAtolFiscal(utmPayId int64, amount float64, operator string, contact string) (taskId string, err error) {
	var resp *http.Response
	taskId = MakeUUID()

	if contact == "" {
		contact = DefaultContact
	}

	jsonStr := fmt.Sprintf(CFG.Atol.RequestTemplate, taskId, utmPayId,
		amount, amount, amount, operator, amount, contact, operator)

	if resp, err = http.Post(CFG.Atol.ApiUrl, "application/json",
		bytes.NewBuffer([]byte(jsonStr))); err != nil {
		return
	}
	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		err = errors.New("atol api response status: " + resp.Status)
		return
	}

	return
}

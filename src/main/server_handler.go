package main

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	ErrOk = iota
	ErrTemporary
	ErrUnknownRequest
	ErrAccountNotFound
	_ // ErrAccountIdWrongFormat
	_ // ErrAccountDisabled
	ErrPayIdWrongFormat
	_ // ErrMaintenance
	ErrPaymentAlreadyExists
	ErrAmountWrongFormat
	ErrAmountTooSmall
	ErrAmountTooBig
	ErrPayDateWrongFormat
	_ // ErrOther = 300
)

type (
	Response interface{ XML() string }

	BaseResponse struct {
		Response
		XMLName xml.Name
		Code    int
		Msg     string
	}
)

func (r *BaseResponse) XML() string {
	return fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<response>
  <CODE>%d</CODE>
  <MESSAGE>%s</MESSAGE>
</response>
`, r.Code, r.Msg)
}

func (s *Server) Handler(c *gin.Context) {
	var (
		ok        bool
		action    string
		account   string
		uid       int
		aid       int
		p         string
		amount    float64
		bankPayId string
		payDate   time.Time
		r         Response
		err       error
	)

	user := c.GetString("user")
	LOG.SetPrefix(fmt.Sprintf("[%s][@%s] ", c.Request.Header.Get("X-Request-Id"), user))

	defer func() {
		header := c.Writer.Header()
		header["Content-Type"] = []string{"application/xml; charset=utf-8"}
		c.Status(http.StatusOK)
		if _, err = c.Writer.WriteString(r.XML()); err != nil {
			ehSkip(err)
			c.AbortWithStatus(http.StatusInternalServerError)
		}
		LOG.SetPrefix("")
	}()

	if action, ok = c.GetQuery("action"); !ok {
		r = &BaseResponse{Code: ErrUnknownRequest, Msg: "action param not found"}
		return
	}

	if account, ok = c.GetQuery("account"); !ok {
		r = &BaseResponse{Code: ErrUnknownRequest, Msg: "account param not found"}
		return
	}
	if account == "" || account == "0" || len(account) > 8 {
		r = &BaseResponse{Code: ErrAccountNotFound, Msg: "account not found"}
		return
	}

	utmPrefix, aid := GetBillingByExtID(account)
	utm := &UtmApi{BillingPrefix: utmPrefix}

	if utmPrefix.Api == "tih" {
		if aid, err = utm.GetAidByExtID(account); err != nil {
			ehSkip(err)
			r = &BaseResponse{Code: ErrAccountNotFound, Msg: "account not found: " + account}
			return
		}
	}
	if aid == 0 {
		r = &BaseResponse{Code: ErrAccountNotFound, Msg: "account not found: " + account}
		return
	}

	if uid, err = utm.GetUidByAid(aid); err != nil {
		ehSkip(err)
		return
	}

	if uid == 0 {
		r = &BaseResponse{Code: ErrAccountNotFound, Msg: "account not found: " + account}
		return
	}

	switch action {
	case "check":
		r = Check(utm, account, uid, aid)
	case "payment":
		if p, ok = c.GetQuery("amount"); !ok {
			r = &BaseResponse{Code: ErrUnknownRequest, Msg: "amount param not found"}
			return
		}
		if amount, err = strconv.ParseFloat(p, 64); err != nil {
			r = &BaseResponse{Code: ErrAmountWrongFormat, Msg: "amount format is wrong"}
			return
		}
		if amount < float64(CFG.OSMP.PayAmountMin) {
			r = &BaseResponse{Code: ErrAmountTooSmall, Msg: "amount is too small"}
			return
		}
		if amount > float64(CFG.OSMP.PayAmountMax) {
			r = &BaseResponse{Code: ErrAmountTooBig, Msg: "amount is too big"}
			return
		}

		if bankPayId, ok = c.GetQuery("pay_id"); !ok {
			r = &BaseResponse{Code: ErrUnknownRequest, Msg: "pay_id param not found"}
			return
		}
		if len(bankPayId) > 255 || len(strings.TrimSpace(bankPayId)) < 1 {
			r = &BaseResponse{Code: ErrPayIdWrongFormat, Msg: "pay_id format is wrong"}
			return
		}

		if p, ok = c.GetQuery("pay_date"); !ok {
			r = &BaseResponse{Code: ErrUnknownRequest, Msg: "pat_date param not found"}
			return
		}
		if payDate, err = time.ParseInLocation(CFG.OSMP.TimeLayout, p, LOC); err != nil {
			r = &BaseResponse{Code: ErrPayDateWrongFormat, Msg: "pat_date format is wrong"}
			return
		}

		if user == "test" {
			payDate = time.Now().Local()
			r = &PayResponse{
				base:     BaseResponse{Code: ErrOk, Msg: "successful"},
				UtmPayID: 1001001,
				RegDate:  payDate,
				Amount:   amount}
			return
		}
		r = Pay(utm, uid, aid, amount, bankPayId, payDate, user, c.Query("contact"))
	default:
		r = &BaseResponse{Code: ErrUnknownRequest, Msg: "action is unknown"}
	}

	if r == nil {
		r = &BaseResponse{Code: ErrTemporary, Msg: "temporary error"}
	}
}

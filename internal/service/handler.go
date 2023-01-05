package service

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	zlog "github.com/rs/zerolog/log"

	"utm5-sber-osmp/internal/atol"
	"utm5-sber-osmp/internal/utm"
)

const (
	ErrOk = iota
	ErrTemporary
	ErrUnknownRequest
	ErrAccountNotFound
	ErrPayIdWrongFormat
	ErrPaymentAlreadyExists
	ErrAmountWrongFormat
	ErrAmountTooSmall
	ErrAmountTooBig
	ErrPayDateWrongFormat
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

func handle(c *gin.Context) {
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
	LOG := zlog.With().Str("user", user).Str("req-id", c.Request.Header.Get("X-Request-Id")).Logger()
	LOG.Info().Msg("request")

	defer func() {
		header := c.Writer.Header()
		header["Content-Type"] = []string{"application/xml; charset=utf-8"}
		c.Status(http.StatusOK)
		if _, err = c.Writer.WriteString(r.XML()); err != nil {
			LOG.Err(err).Msg("send response")
			c.AbortWithStatus(http.StatusInternalServerError)
		}
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

	configValue, _ := c.Get("config")
	config := configValue.(Config)

	utmClient := utm.NewClient(config.Billing)

	pref, aid := utmClient.GetPrefixByExtID(account)
	utmClient.SetActivePrefix(pref)

	if pref.Api == "tih" {
		if aid, err = utmClient.GetAidByExtID(account); err != nil {
			LOG.Err(err).Msg("get account")
			r = &BaseResponse{Code: ErrAccountNotFound, Msg: "account not found: " + account}
			return
		}
	}
	if aid == 0 {
		r = &BaseResponse{Code: ErrAccountNotFound, Msg: "account not found: " + account}
		return
	}

	if uid, err = utmClient.GetUidByAid(aid); err != nil {
		LOG.Err(err).Msg("get userId")
		return
	}

	if uid == 0 {
		r = &BaseResponse{Code: ErrAccountNotFound, Msg: "account not found: " + account}
		return
	}

	switch action {
	case "check":
		r = Check(utmClient, account, uid, aid, config.OSMP.CheckInfo)

	case "payment":
		if p, ok = c.GetQuery("amount"); !ok {
			r = &BaseResponse{Code: ErrUnknownRequest, Msg: "amount param not found"}
			return
		}
		if amount, err = strconv.ParseFloat(p, 64); err != nil {
			r = &BaseResponse{Code: ErrAmountWrongFormat, Msg: "amount format is wrong"}
			return
		}
		if amount < float64(config.OSMP.PayAmountMin) {
			r = &BaseResponse{Code: ErrAmountTooSmall, Msg: "amount is too small"}
			return
		}
		if amount > float64(config.OSMP.PayAmountMax) {
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

		locationVal, _ := c.Get("tzLocation")
		tzLocation := locationVal.(*time.Location)
		if payDate, err = time.ParseInLocation("02.01.2006_15:04:05", p, tzLocation); err != nil {
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

		atolClient, err := atol.NewClient(config.Atol)
		if err != nil {
			LOG.Err(err).Msg("atol client init")
			return
		}

		r = Pay(
			utmClient,
			atolClient,
			uid,
			aid,
			amount,
			bankPayId,
			payDate,
			user,
			c.Query("contact"),
			config.OSMP.IdMaxLen)
	default:
		r = &BaseResponse{Code: ErrUnknownRequest, Msg: "action is unknown"}
	}

	if r == nil {
		r = &BaseResponse{Code: ErrTemporary, Msg: "temporary error"}
	}
}

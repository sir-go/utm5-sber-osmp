package service

import (
	"fmt"
	"time"

	zlog "github.com/rs/zerolog/log"

	"utm5-sber-osmp/internal/atol"
	"utm5-sber-osmp/internal/utm"
)

type PayResponse struct {
	base     BaseResponse
	UtmPayID int64
	RegDate  time.Time
	Amount   float64
}

func (r *PayResponse) XML() string {
	return fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<response>
	<CODE>%d</CODE>
	<MESSAGE>%s</MESSAGE>
	<EXT_ID>%d</EXT_ID>
	<REG_DATE>%s</REG_DATE>
	<AMOUNT>%.2f</AMOUNT>
</response>
`, r.base.Code, r.base.Msg, r.UtmPayID, r.RegDate.Format("02.01.2006_15:04:05"), r.Amount)
}

func Pay(utmClient *utm.Client, atolClient *atol.Client, uid int, aid int,
	amount float64, bankPayId string, dt time.Time,
	comment string, contact string, idMaxLen int) (resp Response) {

	var (
		formattedUtmPayId int64
		utmPayId          int
		err               error
		atolTaskId        string
	)

	LOG := zlog.With().Str("action", "pay").Logger()

	payReport, err := utmClient.GetPayments(uid, aid, dt)
	if err != nil {
		LOG.Err(err).Msg("get payments")
		return
	}

	for _, reportRow := range payReport.Rows {
		if !utmClient.IsPayMethodBack(reportRow.MethodId) &&
			reportRow.BankPayId == bankPayId &&
			time.Unix(reportRow.ActualDate, 0).Local().Equal(dt) &&
			amount == reportRow.Amount {

			formattedUtmPayId, err = RightPadID(utmClient.GetActivePrefix().PayId, reportRow.UtmPayId, idMaxLen)
			if err != nil {
				LOG.Err(err).Msg("padding Id")
				return
			}

			return &PayResponse{
				base:     BaseResponse{Code: ErrPaymentAlreadyExists, Msg: "Дублирование транзакции"},
				UtmPayID: formattedUtmPayId,
				RegDate:  time.Unix(reportRow.EnterDate, 0).Local(),
				Amount:   amount}
		}
	}

	if utmPayId, err = utmClient.AddPayment(aid, amount, dt, comment, bankPayId); err != nil {
		LOG.Err(err).Msg("add payment")
		return
	}

	formattedUtmPayId, err = RightPadID(utmClient.GetActivePrefix().PayId, utmPayId, idMaxLen)
	if err != nil {
		LOG.Err(err).Msg("padding Id")
		return
	}

	atolTaskId, err = atolClient.MakeFiscal(formattedUtmPayId, amount, comment, contact)
	if err != nil {
		LOG.Err(err).Msg("make atol fiscal")
	}
	LOG.Info().Str("taskId", atolTaskId).Msg("make atol fiscal")

	return &PayResponse{
		base:     BaseResponse{Code: ErrOk, Msg: "successful"},
		UtmPayID: formattedUtmPayId,
		RegDate:  time.Now().Local(),
		Amount:   amount}
}

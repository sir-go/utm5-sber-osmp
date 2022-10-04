package main

import (
	"fmt"
	"time"
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
`, r.base.Code, r.base.Msg, r.UtmPayID, r.RegDate.Format(CFG.OSMP.TimeLayout), r.Amount)
}

func Pay(utm *UtmApi, uid int, aid int,
	amount float64, bankPayId string, dt time.Time,
	comment string, contact string) (resp Response) {

	var (
		formattedUtmPayId int64
		utmPayId          int
		err               error
		atolTaskId        string
	)

	payReport, err := utm.GetPayments(uid, aid, dt.Add(-CFG.Billing.PaymentReportRetro.Duration))
	if err != nil {
		ehSkip(err)
		return
	}

	for _, reportRow := range payReport.Rows {
		if reportRow.MethodId != CFG.Billing.PaymentBackMethod &&
			reportRow.BankPayId == bankPayId &&
			time.Unix(reportRow.ActualDate, 0).Local().Equal(dt) &&
			amount == reportRow.Amount {
			if formattedUtmPayId, err = RightPadID(
				utm.BillingPrefix.PayId,
				reportRow.UtmPayId,
				CFG.OSMP.IdMaxLen); err != nil {
				ehSkip(err)
				return
			}
			return &PayResponse{
				base:     BaseResponse{Code: ErrPaymentAlreadyExists, Msg: "Дублирование транзакции"},
				UtmPayID: formattedUtmPayId,
				RegDate:  time.Unix(reportRow.EnterDate, 0).Local(),
				Amount:   amount}
		}
	}

	if utmPayId, err = utm.AddPayment(aid, amount, dt, comment, bankPayId); err != nil {
		ehSkip(err)
		return
	}

	if formattedUtmPayId, err = RightPadID(
		utm.BillingPrefix.PayId,
		utmPayId,
		CFG.OSMP.IdMaxLen); err != nil {
		ehSkip(err)
		return
	}

	if atolTaskId, err = MakeAtolFiscal(formattedUtmPayId, amount, comment, contact); err != nil {
		ehSkip(err)
	}
	LOG.Println("made atol task id: " + atolTaskId)

	return &PayResponse{
		base:     BaseResponse{Code: ErrOk, Msg: "successful"},
		UtmPayID: formattedUtmPayId,
		RegDate:  time.Now().Local(),
		Amount:   amount}
}

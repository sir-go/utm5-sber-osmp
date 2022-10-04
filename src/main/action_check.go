package main

import (
	"fmt"
)

type CheckResponse struct {
	base    BaseResponse
	Name    string
	Address string
	Balance int
	RecSum  int
	Info    string
}

func (r *CheckResponse) XML() string {
	return fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<response>
	<CODE>%d</CODE>
	<MESSAGE>%s</MESSAGE>
	<FIO>%s</FIO>
	<ADDRESS>%s</ADDRESS>
	<BALANCE>%d</BALANCE>
	<REC_SUM>%d</REC_SUM>
	<INFO>%s</INFO>
</response>
`, r.base.Code, r.base.Msg, r.Name, r.Address, r.Balance, r.RecSum, r.Info)
}

func Check(utm *UtmApi, extId string, uid int, aid int) (resp Response) {
	var (
		err     error
		balance float64
		cost    float64
	)

	a := &CheckResponse{
		base: BaseResponse{Code: ErrOk, Msg: "account exist"},
		Info: CFG.OSMP.CheckInfo}

	errChUserInfo := make(chan error, 0)
	go func() {
		if a.Name, a.Address, err = utm.GetUserInfo(uid); err != nil {
			ehSkip(err)
			errChUserInfo <- err
			return
		}
		errChUserInfo <- nil
	}()

	errChBalance := make(chan error, 0)
	go func() {
		if balance, err = utm.GetBalance(aid); err != nil {
			errChBalance <- err
			return
		}
		errChBalance <- nil
	}()

	errChCost := make(chan error, 0)
	go func() {
		if cost, err = utm.GetServicesCost(aid); err != nil {
			errChCost <- err
			return
		}
		errChCost <- nil
	}()

	if err = <-errChUserInfo; err != nil {
		return &BaseResponse{Code: ErrAccountNotFound, Msg: "account not found: " + extId}
	}

	if err = <-errChBalance; err != nil {
		balance = 0
		ehSkip(err)
	}
	a.Balance = RoundBalance(balance)

	if err = <-errChCost; err != nil {
		ehSkip(err)
		cost = 0
	}

	if balance < cost {
		a.RecSum = RoundRecSum(cost - balance)
	}

	return a
}

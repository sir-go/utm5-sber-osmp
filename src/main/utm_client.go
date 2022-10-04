package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/KeisukeYamashita/jsonrpc"
)

func GetBillingByExtID(extId string) (billing *CfgBillingPrefix, aidInt int) {
	var err error
	if HasDigitsOnly(extId) {
		if aidInt, err = strconv.Atoi(extId); err != nil {
			ehSkip(err)
			return &CFG.Billing.PrefixTih, 0
		}
		if aidInt < 12000 {
			return &CFG.Billing.PrefixKor, aidInt
		}
	}
	return &CFG.Billing.PrefixTih, 0
}

const UTMOnceServiceType = 1

type (
	rAid struct {
		Aid int `json:"aid"`
	}

	rUid struct {
		Uid int `json:"user_id"`
	}

	rUInfo struct {
		Name       string `json:"full_name"`
		Address    string `json:"act_address"`
		FlatNumber string `json:"flat_number"`
	}

	rBalance struct {
		Balance float64 `json:"balance"`
	}

	rAccountServices struct {
		Slinks []struct {
			ServiceID   int     `json:"service_id"`
			ServiceType int     `json:"service_type"`
			ServiceCost float64 `json:"service_cost"`
		} `json:"slinks"`
	}

	rOnceServiceCost struct {
		Cost float64 `json:"cost"`
	}

	rPayReport struct {
		Rows []struct {
			UtmPayId   int     `json:"id"`
			BankPayId  string  `json:"payment_ext_number"`
			ActualDate int64   `json:"actual_date"`
			EnterDate  int64   `json:"payment_enter_date"`
			Amount     float64 `json:"payment"`
			MethodId   int     `json:"method"`
		} `json:"rows"`
	}

	rPayId struct {
		UtmPayId int `json:"payment_transaction_id"`
	}

	UtmApi struct {
		BillingPrefix *CfgBillingPrefix
	}

	UtmArgs map[string]interface{}
)

func (u *UtmApi) call(method string, args UtmArgs, target interface{}) (err error) {
	LOG.Printf("%s.%s %v", u.BillingPrefix.Api, method, args)
	var res *jsonrpc.RPCResponse
	client := jsonrpc.NewRPCClient(CFG.Billing.ApiURL)
	client.SetBasicAuth(CFG.Billing.Username, CFG.Billing.Password)

	if res, err = client.CallNamed(u.BillingPrefix.Api+"."+method, args); err != nil {
		return
	}

	if res.Error != nil {
		err = fmt.Errorf("urfa api error: %d : %s : %v",
			res.Error.Code, res.Error.Message, res.Error.Data)
		return
	}

	return res.GetObject(target)
}

func (u *UtmApi) GetAidByExtID(extId string) (int, error) {
	o := new(rAid)
	if err := u.call("rpcf_is_account_external_id_used",
		UtmArgs{"external_id": extId}, o); err != nil {
		return 0, err
	}
	return o.Aid, nil
}

func (u *UtmApi) GetUidByAid(aid int) (int, error) {
	o := new(rUid)
	if err := u.call("rpcf_get_user_by_account",
		UtmArgs{"account_id": aid}, o); err != nil {
		return 0, err
	}
	return o.Uid, nil
}

func (u *UtmApi) GetUserInfo(aid int) (string, string, error) {
	o := new(rUInfo)
	if err := u.call("rpcf_get_userinfo",
		UtmArgs{"user_id": aid}, o); err != nil {
		return "", "", err
	}
	if o.FlatNumber != "" {
		return o.Name, fmt.Sprintf("%s, кв. %s", o.Address, o.FlatNumber), nil
	}
	return o.Name, o.Address, nil
}

func (u *UtmApi) GetBalance(aid int) (float64, error) {
	o := new(rBalance)
	if err := u.call("rpcf_get_accountinfo",
		UtmArgs{"account_id": aid}, o); err != nil {
		return 0.0, err
	}

	return o.Balance, nil
}

func (u *UtmApi) GetPayments(uid, aid int, timeStart time.Time) (*rPayReport, error) {
	o := new(rPayReport)
	if err := u.call("rpcf_payments_report_new",
		UtmArgs{"user_id": uid, "account_id": aid, "time_start": timeStart.Unix()}, o); err != nil {
		return nil, err
	}

	return o, nil
}

func (u *UtmApi) GetServices(aid int) (*rAccountServices, error) {
	o := new(rAccountServices)
	if err := u.call("rpcf_get_all_services_for_user",
		UtmArgs{"account_id": aid}, o); err != nil {
		return nil, err
	}

	return o, nil
}

func (u *UtmApi) GetOnceServiceCost(sid int) (float64, error) {
	o := new(rOnceServiceCost)
	if err := u.call("rpcf_get_once_service",
		UtmArgs{"sid": sid}, o); err != nil {
		return 0.0, err
	}

	return o.Cost, nil
}

func (u *UtmApi) GetServicesCost(aid int) (cost float64, err error) {
	var (
		services *rAccountServices
		onesCost float64
	)

	if services, err = u.GetServices(aid); err != nil {
		return
	}

	for _, slink := range services.Slinks {
		if slink.ServiceType != UTMOnceServiceType {
			cost += slink.ServiceCost
			continue
		}

		if onesCost, err = u.GetOnceServiceCost(slink.ServiceID); err != nil {
			return
		}
		cost += onesCost

	}
	return
}

func (u *UtmApi) AddPayment(aid int, amount float64, dt time.Time,
	comment string, bankPayId string) (int, error) {
	o := new(rPayId)
	if err := u.call("rpcf_add_payment_for_account",
		UtmArgs{
			"account_id":         aid,
			"payment":            amount,
			"payment_date":       dt.Unix(),
			"payment_method":     CFG.Billing.PaymentMethod,
			"admin_comment":      comment,
			"payment_ext_number": bankPayId,
		}, o); err != nil {
		return 0, err
	}

	return o.UtmPayId, nil
}

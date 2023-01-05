package utm

import "C"
import (
	"fmt"
	"strconv"
	"time"

	"github.com/KeisukeYamashita/jsonrpc"
	zlog "github.com/rs/zerolog/log"
)

const OnceServiceType = 1

type (
	Client struct {
		cfg          Config
		activePrefix *Prefix
	}

	RetAid struct {
		Aid int `json:"aid"`
	}

	RetUid struct {
		Uid int `json:"user_id"`
	}

	RetUInfo struct {
		Name       string `json:"full_name"`
		Address    string `json:"act_address"`
		FlatNumber string `json:"flat_number"`
	}

	RetBalance struct {
		Balance float64 `json:"balance"`
	}

	RetAccountServices struct {
		Slinks []struct {
			ServiceID   int     `json:"service_id"`
			ServiceType int     `json:"service_type"`
			ServiceCost float64 `json:"service_cost"`
		} `json:"slinks"`
	}

	RetOnceServiceCost struct {
		Cost float64 `json:"cost"`
	}

	RetPayReport struct {
		Rows []struct {
			UtmPayId   int     `json:"id"`
			BankPayId  string  `json:"payment_ext_number"`
			ActualDate int64   `json:"actual_date"`
			EnterDate  int64   `json:"payment_enter_date"`
			Amount     float64 `json:"payment"`
			MethodId   int     `json:"method"`
		} `json:"rows"`
	}

	RetPayId struct {
		UtmPayId int `json:"payment_transaction_id"`
	}

	Args map[string]interface{}
)

func NewClient(config Config) *Client {
	return &Client{cfg: config}
}

func (c *Client) GetActivePrefix() Prefix {
	return *c.activePrefix
}

func (c *Client) SetActivePrefix(pref *Prefix) {
	c.activePrefix = pref
}

func (c *Client) GetPrefixByExtID(extId string) (pref *Prefix, aidInt int) {
	var err error
	if HasDigitsOnly(extId) {
		if aidInt, err = strconv.Atoi(extId); err != nil {
			return c.cfg.Prefixes["tih"], 0
		}
		if aidInt < 12000 {
			return c.cfg.Prefixes["kor"], aidInt
		}
	}
	return c.cfg.Prefixes["tih"], 0
}

func (c *Client) call(method string, args Args, target interface{}) (err error) {
	prefMethod := fmt.Sprintf("%s.%s", c.activePrefix.Api, method)
	zlog.Debug().Str("method", prefMethod).Interface("args", args)

	var res *jsonrpc.RPCResponse
	client := jsonrpc.NewRPCClient(c.cfg.ApiURL)
	client.SetBasicAuth(c.cfg.Username, c.cfg.Password)

	if res, err = client.CallNamed(prefMethod, args); err != nil {
		return
	}

	if res.Error != nil {
		err = fmt.Errorf("urfa api error: %d : %s : %v",
			res.Error.Code, res.Error.Message, res.Error.Data)
		return
	}

	return res.GetObject(target)
}

func (c *Client) IsPayMethodBack(payMethodId int) bool {
	return c.cfg.PaymentBackMethod == payMethodId
}

func (c *Client) GetAidByExtID(extId string) (int, error) {
	o := new(RetAid)
	if err := c.call("rpcf_is_account_external_id_used",
		Args{"external_id": extId}, o); err != nil {
		return 0, err
	}
	return o.Aid, nil
}

func (c *Client) GetUidByAid(aid int) (int, error) {
	o := new(RetUid)
	if err := c.call("rpcf_get_user_by_account",
		Args{"account_id": aid}, o); err != nil {
		return 0, err
	}
	return o.Uid, nil
}

func (c *Client) GetUserInfo(aid int) (string, string, error) {
	o := new(RetUInfo)
	if err := c.call("rpcf_get_userinfo",
		Args{"user_id": aid}, o); err != nil {
		return "", "", err
	}
	if o.FlatNumber != "" {
		return o.Name, fmt.Sprintf("%s, кв. %s", o.Address, o.FlatNumber), nil
	}
	return o.Name, o.Address, nil
}

func (c *Client) GetBalance(aid int) (float64, error) {
	o := new(RetBalance)
	if err := c.call("rpcf_get_accountinfo",
		Args{"account_id": aid}, o); err != nil {
		return 0.0, err
	}

	return o.Balance, nil
}

func (c *Client) GetPayments(uid, aid int, timeEnd time.Time) (*RetPayReport, error) {
	o := new(RetPayReport)
	if err := c.call("rpcf_payments_report_new",
		Args{
			"user_id":    uid,
			"account_id": aid,
			"time_start": timeEnd.Add(-c.cfg.PaymentReportRetro).Unix(),
		}, o); err != nil {
		return nil, err
	}

	return o, nil
}

func (c *Client) GetServices(aid int) (*RetAccountServices, error) {
	o := new(RetAccountServices)
	if err := c.call("rpcf_get_all_services_for_user",
		Args{"account_id": aid}, o); err != nil {
		return nil, err
	}

	return o, nil
}

func (c *Client) GetOnceServiceCost(sid int) (float64, error) {
	o := new(RetOnceServiceCost)
	if err := c.call("rpcf_get_once_service",
		Args{"sid": sid}, o); err != nil {
		return 0.0, err
	}

	return o.Cost, nil
}

func (c *Client) GetServicesCost(aid int) (cost float64, err error) {
	var (
		services *RetAccountServices
		onesCost float64
	)

	if services, err = c.GetServices(aid); err != nil {
		return
	}

	for _, slink := range services.Slinks {
		if slink.ServiceType != OnceServiceType {
			cost += slink.ServiceCost
			continue
		}

		if onesCost, err = c.GetOnceServiceCost(slink.ServiceID); err != nil {
			return
		}
		cost += onesCost

	}
	return
}

func (c *Client) AddPayment(aid int, amount float64, dt time.Time,
	comment string, bankPayId string) (int, error) {
	o := new(RetPayId)
	if err := c.call("rpcf_add_payment_for_account",
		Args{
			"account_id":         aid,
			"payment":            amount,
			"payment_date":       dt.Unix(),
			"payment_method":     c.cfg.PaymentMethod,
			"admin_comment":      comment,
			"payment_ext_number": bankPayId,
		}, o); err != nil {
		return 0, err
	}

	return o.UtmPayId, nil
}

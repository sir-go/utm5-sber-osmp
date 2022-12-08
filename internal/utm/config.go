package utm

import (
	"time"
)

type (
	Prefix struct {
		Api    string `toml:"api_prefix"`
		PayId  int    `toml:"pay_id_prefix"`
		Office string `toml:"office"`
	}

	Config struct {
		ApiURL             string             `toml:"api_url"`
		Username           string             `toml:"username"`
		Password           string             `toml:"password"`
		PaymentMethod      int                `toml:"payment_method"`
		PaymentBackMethod  int                `toml:"payment_back_method"`
		PaymentReportRetro time.Duration      `toml:"payment_report_retro"`
		Prefixes           map[string]*Prefix `toml:"prefixes"`
	}
)

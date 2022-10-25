package main

import (
	"flag"
	"os"
	"time"

	"github.com/BurntSushi/toml"
)

type (
	Duration struct {
		time.Duration
	}

	CfgService struct {
		Host     string `toml:"host"`
		Port     int    `toml:"port"`
		Location string `toml:"location"`
		Timeouts struct {
			Write *Duration `toml:"write"`
			Read  *Duration `toml:"read"`
			Idle  *Duration `toml:"idle"`
		} `toml:"timeouts"`
		Users map[string]string `toml:"users"`
	}

	CfgBillingPrefix struct {
		Api    string `toml:"api_prefix"`
		PayId  int    `toml:"pay_id_prefix"`
		Office string `toml:"office"`
	}

	CfgBilling struct {
		ApiURL             string           `toml:"api_url"`
		Username           string           `toml:"username"`
		Password           string           `toml:"password"`
		PaymentMethod      int              `toml:"payment_method"`
		PaymentBackMethod  int              `toml:"payment_back_method"`
		PaymentReportRetro *Duration        `toml:"payment_report_retro"`
		PrefixTih          CfgBillingPrefix `toml:"tih"`
		PrefixKor          CfgBillingPrefix `toml:"kor"`
	}

	CfgOSMP struct {
		CheckInfo    string `toml:"check_info"`
		TimeLayout   string `toml:"time_layout"`
		IdMaxLen     int    `toml:"id_max_len"`
		PayAmountMin int    `toml:"pay_amount_min"`
		PayAmountMax int    `toml:"pay_amount_max"`
	}

	CfgAtol struct {
		ApiUrl          string `toml:"api_url"`
		RequestTemplate string `toml:"request_template"`
	}

	Config struct {
		Service CfgService `toml:"service"`
		Billing CfgBilling `toml:"billing"`
		OSMP    CfgOSMP    `toml:"osmp"`
		Atol    CfgAtol    `toml:"atol"`
		Path    string
	}
)

func (d *Duration) UnmarshalText(text []byte) error {
	var err error
	d.Duration, err = time.ParseDuration(string(text))
	return err
}

func ConfigInit() *Config {
	fCfgPath := flag.String("c", DefaultConfFile, "path to conf file")
	flag.Parse()

	conf := new(Config)
	file, err := os.Open(*fCfgPath)
	if err != nil {
		panic(err)
	}

	defer func() {
		if file == nil {
			return
		}
		if err = file.Close(); err != nil {
			panic(err)
		}
	}()

	if _, err = toml.DecodeFile(*fCfgPath, &conf); err != nil {
		panic(err)
	}
	conf.Path = *fCfgPath
	return conf
}

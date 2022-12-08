package service

import (
	"time"

	"utm5-sber-osmp/internal/atol"
	"utm5-sber-osmp/internal/utm"
)

type (
	CfgOSMP struct {
		CheckInfo    string `toml:"check_info"`
		IdMaxLen     int    `toml:"id_max_len"`
		PayAmountMin int    `toml:"pay_amount_min"`
		PayAmountMax int    `toml:"pay_amount_max"`
	}

	CfgService struct {
		Host     string `toml:"host"`
		Port     int    `toml:"port"`
		Location string `toml:"location"`
		Timeouts struct {
			Write time.Duration `toml:"write"`
			Read  time.Duration `toml:"read"`
			Idle  time.Duration `toml:"idle"`
		} `toml:"timeouts"`
		Users map[string]string `toml:"users"`
	}

	Config struct {
		Billing utm.Config  `toml:"billing"`
		Atol    atol.Config `toml:"atol"`
		Service CfgService  `toml:"service"`
		OSMP    CfgOSMP     `toml:"osmp"`
	}
)

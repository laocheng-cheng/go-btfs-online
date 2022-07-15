package config

import (
	"fmt"
	"time"

	"github.com/tron-us/go-common/v2/env"
	"github.com/tron-us/go-common/v2/log"
	"github.com/tron-us/status-server/common"

	"go.uber.org/zap"
)

var (
	Host = ""
	Port = ":50051"

	StagRequsterPid     = "16Uiu2HAkuszaXgpfttTG1eTvVns2Cdtwp1KikeC7828N6Q1bV4T1"
	ProdRequsterPid     = "16Uiu2HAmGhz9BpSpHRXV5ppDnCMzsxLh5pjvvFZaBn2oHe6M1T1B"
	TempProdRequsterPid = "16Uiu2HAkvRaYNTtdoWjsp796ZDrzUiRwvAGTBcY7bbmfLKPZoVTV"
	OtherRequsterPid    = "16Uiu2HAm9BZSAT1xAcGq2o63zzL6k23dCzzjE7zRCQDEzqTR5T8D"

	MinRefreshInterval = 7 * 60 * time.Minute
	MaxRefreshTimeout  = 24 * 60 * time.Minute
)

func init() {
	if env, h := env.GetEnv("HOST"); h != "" {
		Host = h
		log.Debug(common.ServerConfig, zap.String(env, fmt.Sprint(Host)))
	}
	if env, p := env.GetEnv("PORT"); p != "" {
		Port = p
		log.Debug(common.ServerConfig, zap.String(env, fmt.Sprint(Port)))
	}

	if _, r := env.GetEnv("ENV"); r != "production" {
		MinRefreshInterval = 30 * time.Minute
		MaxRefreshTimeout = 5 * time.Minute
	}
}

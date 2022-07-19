package config

import (
	"fmt"
	"github.com/tron-us/go-common/v2/env"
	"github.com/tron-us/go-common/v2/log"
	"go.uber.org/zap"
)

var (
	Host = ""
	Port = ":50051"

	SignAddr       = "0x22df207EC3C8D18fEDeed87752C5a68E5b4f6FbD"
	SignPrivateKey = "744ba22387c27cf73dff283a37f0a7e63054a86be15965be97c807816d79da39"
)

func init() {
	if env, h := env.GetEnv("HOST"); h != "" {
		Host = h
		log.Debug("config, HOST ", zap.String(env, fmt.Sprint(Host)))
	}
	if env, p := env.GetEnv("PORT"); p != "" {
		Port = p
		log.Debug("config, PORT ", zap.String(env, fmt.Sprint(Port)))
	}

	if env, sa := env.GetEnv("SIGN_ADDR"); sa != "" {
		SignAddr = sa
		log.Debug("config, SIGN_ADDR ", zap.String(env, SignAddr))
	}
	if env, spk := env.GetEnv("SIGN_PRIVATE_KEY"); spk != "" {
		SignPrivateKey = spk
		log.Debug("config, SIGN_PRIVATE_KEY ", zap.String(env, SignPrivateKey))
	}
}

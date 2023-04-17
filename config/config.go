package config

import (
	"time"

	"github.com/spf13/viper"
)

const (
	// ServerAddr порт на котором будет сервис
	ServerAddr        = "SERVER_ADDR"
	serverAddrDefault = ":3333"

	// Store путь к хранилищу
	Store        = "STORE"
	storeDefault = `./users.json`

	// CtxTimeout timeout для контекста
	CtxTimeout        = "TIMEOUT"
	ctxTimeoutDefault = 60 * time.Second
)

func SetDefaults() {
	// ставим default значения
	viper.SetDefault(Store, storeDefault)
	viper.SetDefault(CtxTimeout, ctxTimeoutDefault)
	viper.SetDefault(ServerAddr, serverAddrDefault)
}

func Configure() {
	SetDefaults()
}

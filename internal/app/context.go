package app

import (
	"context"
	"sync"

	"github.com/kondohiroki/go-boilerplate/config"
	"github.com/kondohiroki/go-boilerplate/internal/discord"
	"go.uber.org/zap"
)

var appCtx *AppContext // Read-only global variable
var m sync.Mutex

type AppContext struct {
	Ctx     context.Context
	Config  *config.Config
	Discord *discord.Discord
	Logger  *zap.Logger
}

func GetAppContext() *AppContext {
	return appCtx
}

func SetAppContext(ctx *AppContext) {
	m.Lock()
	defer m.Unlock()
	appCtx = ctx
}

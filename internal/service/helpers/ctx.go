package helpers

import (
	"context"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/tokend/course-certificates/ccp/internal/config"
	"net/http"
)

type ctxKey int

const (
	logCtxKey ctxKey = iota
	configCtxKey
)

func CtxLog(entry *logan.Entry) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, logCtxKey, entry)
	}
}

func Log(r *http.Request) *logan.Entry {
	return r.Context().Value(logCtxKey).(*logan.Entry)
}

func CtxConfig(entry config.Config) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, configCtxKey, entry)
	}
}

func Config(r *http.Request) config.Config {
	return r.Context().Value(configCtxKey).(config.Config)
}

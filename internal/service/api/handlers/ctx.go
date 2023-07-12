package handlers

import (
	"context"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/tokend/course-certificates/ccp/internal/config"
	"gitlab.com/tokend/course-certificates/ccp/internal/data"
	"gitlab.com/tokend/course-certificates/ccp/internal/service/core/pdf"
	"net/http"
)

type ctxKey int

const (
	logCtxKey ctxKey = iota
	configCtxKey
	masterqCtxKey
	pdfCreatorCtxKey
	staticConfigerCtxKey
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

func CtxMasterQ(entry data.MasterQ) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, masterqCtxKey, entry)
	}
}

func MasterQ(r *http.Request) data.MasterQ {
	return r.Context().Value(masterqCtxKey).(data.MasterQ).New()
}

func CtxPdfCreator(entry pdf.CreatorPDF) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, pdfCreatorCtxKey, entry)
	}
}

func PdfCreator(r *http.Request) pdf.CreatorPDF {
	return r.Context().Value(pdfCreatorCtxKey).(pdf.CreatorPDF)
}

func CtxStaticConfiger(entry *config.StaticConfig) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, staticConfigerCtxKey, entry)
	}
}

func StaticConfiger(r *http.Request) *config.StaticConfig {
	return r.Context().Value(staticConfigerCtxKey).(*config.StaticConfig)
}

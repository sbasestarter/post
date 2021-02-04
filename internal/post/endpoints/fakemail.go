package endpoints

import (
	"context"
	"github.com/jiuzhou-zhao/go-fundamental/loge"
)

type fakeMail struct {
}

func NewFakeMail() EndPoint {
	return &fakeMail{}
}

func (fm *fakeMail) Send(ctx context.Context, to []string, templateID string, vars []string) error {
	loge.Infof(ctx, "fake mail: to %v, detail: %+v", to, vars)
	return nil
}

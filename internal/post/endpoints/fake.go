package endpoints

import (
	"context"

	"github.com/sbasestarter/post/internal/config"
	"github.com/sgostarter/i/l"
)

type fake struct {
	endPoint config.EndPoint
	logger   l.WrapperWithContext
}

func NewFake(endPoint config.EndPoint, logger l.Wrapper) EndPoint {
	if logger == nil {
		logger = l.NewNopLoggerWrapper()
	}

	return &fake{
		endPoint: endPoint,
		logger:   logger.WithFields(l.StringField(l.ClsKey, "fake")).GetWrapperWithContext(),
	}
}

func (fm *fake) Send(ctx context.Context, to []string, templateID string, vars []string) error {
	fm.logger.Infof(ctx, "fake %s: to %v, detail: %+v", fm.endPoint.Name, to, vars)

	return nil
}

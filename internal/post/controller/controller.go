package controller

import (
	"context"
	"fmt"
	"strings"

	"github.com/sbasestarter/post/internal/config"
	"github.com/sbasestarter/post/internal/post/endpoints"
	"github.com/sbasestarter/post/pkg/post"
	"github.com/sgostarter/i/l"
)

type Controller struct {
	cfg       *config.Config
	endPoints map[string]map[string]endpoints.EndPoint
	logger    l.WrapperWithContext
}

// nolint: gocognit
func NewController(ctx context.Context, cfg *config.Config, logger l.Wrapper) *Controller {
	if logger == nil {
		logger = l.NewNopLoggerWrapper()
	}

	lc := logger.WithFields(l.StringField(l.ClsKey, "Controller")).GetWrapperWithContext()

	controller := &Controller{
		cfg:       cfg,
		endPoints: make(map[string]map[string]endpoints.EndPoint),
		logger:    lc,
	}

	for protocol, providers := range cfg.ProtocolProviders {
		protocol = strings.ToLower(protocol)

		for provider, endPoints := range providers {
			provider = strings.ToLower(provider)

			for _, endPoint := range endPoints {
				if endPoint.Name == "" {
					lc.Fatal(ctx, "empty endpoint name")

					continue
				}

				endPoint.Name = strings.ToLower(endPoint.Name)

				var ep endpoints.EndPoint

				switch protocol {
				case post.ProtocolTypeEmail:
					switch provider {
					case post.ProviderGoMail:
						ep = endpoints.NewGoMailEndPoint(ctx, endPoint, logger)
					case post.ProviderFake:
						ep = endpoints.NewFake(endPoint, logger)
					}
				case post.ProtocolTypeSMS:
					switch provider {
					case post.ProviderTencentSMS:
						ep = endpoints.NewTencentSMSEndPoint(ctx, endPoint, logger)
					case post.ProviderFake:
						ep = endpoints.NewFake(endPoint, logger)
					}
				}

				if ep == nil {
					lc.Errorf(ctx, "create endPoint from %v, %v failed", protocol, provider)

					continue
				}

				if _, ok := controller.endPoints[protocol]; !ok {
					controller.endPoints[protocol] = make(map[string]endpoints.EndPoint)
				}

				if _, ok := controller.endPoints[protocol][endPoint.Name]; ok {
					lc.Fatalf(ctx, "dup endpoint name: %v", endPoint.Name)

					continue
				}

				controller.endPoints[protocol][endPoint.Name] = ep
			}
		}
	}

	return controller
}

func (controller *Controller) SendTemplate(ctx context.Context, to []string, protocolType, templateID string, vars []string) (err error) {
	endPoints, ok := controller.endPoints[protocolType]
	if !ok {
		err = fmt.Errorf("unknown protocol: %v", protocolType)

		return
	}

	for name, endPoint := range endPoints {
		err = endPoint.Send(ctx, to, templateID, vars)
		if err == nil {
			break
		}

		controller.logger.Errorf(ctx, "send by %v failed: %v", name, err)
	}

	return
}

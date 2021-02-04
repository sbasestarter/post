package controller

import (
	"context"
	"fmt"
	"strings"

	"github.com/jiuzhou-zhao/go-fundamental/loge"
	"github.com/sbasestarter/post/internal/config"
	"github.com/sbasestarter/post/internal/post/endpoints"
	"github.com/sbasestarter/post/pkg/post"
)

type Controller struct {
	cfg       *config.Config
	endPoints map[string]map[string]endpoints.EndPoint
}

func NewController(ctx context.Context, cfg *config.Config) *Controller {
	controller := &Controller{
		cfg:       cfg,
		endPoints: make(map[string]map[string]endpoints.EndPoint),
	}

	for protocol, providers := range cfg.ProtocolProviders {
		protocol = strings.ToLower(protocol)
		for provider, endPoints := range providers {
			provider = strings.ToLower(provider)
			for _, endPoint := range endPoints {
				if endPoint.Name == "" {
					loge.Fatal(ctx, "empty endpoint name")
					continue
				}
				endPoint.Name = strings.ToLower(endPoint.Name)
				var ep endpoints.EndPoint
				switch protocol {
				case post.ProtocolTypeEmail:
					switch provider {
					case post.ProviderGoMail:
						ep = endpoints.NewGoMailEndPoint(ctx, endPoint)
					case post.ProviderFakeMail:
						ep = endpoints.NewFakeMail()
					}
				case post.ProtocolTypeSMS:
					switch provider {
					case post.ProviderTencentSMS:
						ep = endpoints.NewTencentSMSEndPoint(ctx, endPoint)
					}
				}
				if ep == nil {
					loge.Errorf(ctx, "create endPoint from %v, %v failed", protocol, provider)
					continue
				}
				if _, ok := controller.endPoints[protocol]; !ok {
					controller.endPoints[protocol] = make(map[string]endpoints.EndPoint)
				}
				if _, ok := controller.endPoints[protocol][endPoint.Name]; ok {
					loge.Fatalf(ctx, "dup endpoint name: %v", endPoint.Name)
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
		loge.Errorf(ctx, "send by %v failed: %v", name, err)
	}
	return
}

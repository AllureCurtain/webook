package ratelimite

import (
	"context"
	"fmt"
	"webook/internal/service/sms"
	"webook/pkg/ratelimit"
)

var errLimited = fmt.Errorf("触发了限流")

type RatelimitSMSService struct {
	svc     sms.Service
	limiter ratelimit.Limiter
}

func NewRatelimitSMSService(svc sms.Service, limiter ratelimit.Limiter) *RatelimitSMSService {
	return &RatelimitSMSService{
		svc:     svc,
		limiter: limiter,
	}
}

func (s RatelimitSMSService) Send(ctx context.Context, tpl string, args []string, numbers ...string) error {
	limited, err := s.limiter.Limit(ctx, "sms:tencent")
	if err != nil {
		return fmt.Errorf("短信服务判断是否限流出现问题, %w", err)
	}
	if limited {
		return errLimited
	}

	// 在这里加一些代码，新特性
	err = s.svc.Send(ctx, tpl, args, numbers...)
	// 这里也可以加
	return err
}

package retryable

import (
	"context"
	"errors"
	"webook/internal/service/sms"
)

type Service struct {
	svc sms.Service
	// 重试
	retryMax int
}

func (s Service) Send(ctx context.Context, biz string, args []string, numbers ...string) error {
	err := s.svc.Send(ctx, biz, args, numbers...)
	cnt := 1
	for err != nil && cnt < s.retryMax {
		err = s.svc.Send(ctx, biz, args, numbers...)
		if err == nil {
			return nil
		}
		cnt++
	}
	return errors.New("重试都失败了")
}

// 设计并实现一个高可用的短信平台
// 1. 提高可用性：重试机制、客户端限流、failover（轮询、实时检测）
// 1.1 实时检测：
// 1.1.1 基于超时的实时检测（连续超时）
// 1.1.2 基于响应时间的实时检测（例如，平均响应时间上升 20%）
// 1.1.3 基于长尾请求的实时检测（例如，响应时间超 1s 的请求占比超 10%）
// 1.1.4 错误率
// 2. 提高安全性：
// 2.1 完整的资源申请与审批流程
// 2.2 鉴权：
// 2.2.1 静态 token
// 2.2.2 动态 token
// 3. 提高可观测性：日志、metrics、tracing，丰富完善的排查手段
// 4. 提高性能：
// （注：当前暂未明确高并发实现方案，可扩展异步队列、连接池、缓存等方向）

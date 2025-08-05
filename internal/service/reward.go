package service

import (
	"context"
	"errors"
	"fmt"
	"webook/internal/domain"
	"webook/internal/repository"
	"webook/pkg/logger"
)

type RewardService interface {
	// PreReward 准备打赏，
	// 你也可以直接理解为对标到创建一个打赏的订单
	// 因为目前我们只支持微信扫码支付，所以实际上直接把接口定义成这个样子就可以了
	PreReward(ctx context.Context, r domain.Reward) (domain.CodeURL, error)
	GetReward(ctx context.Context, rid, uid int64) (domain.Reward, error)
	UpdateReward(ctx context.Context, bizTradeNO string, status domain.RewardStatus) error
}

type rewardService struct {
	repo repository.RewardRepository
	// 这里应该有支付服务，但目前看起来还没有实现
	// paymentSvc PaymentService
	l logger.LoggerV1
}

func NewRewardService(repo repository.RewardRepository, l logger.LoggerV1) RewardService {
	return &rewardService{
		repo: repo,
		l:    l,
	}
}
func (s *rewardService) PreReward(ctx context.Context, r domain.Reward) (domain.CodeURL, error) {
	// 首先检查缓存中是否已有二维码
	codeURL, err := s.repo.GetCachedCodeURL(ctx, r)
	if err == nil {
		// 缓存命中，直接返回
		return codeURL, nil
	}

	// 设置初始状态
	r.Status = domain.RewardStatusInit

	// 创建打赏记录
	rid, err := s.repo.CreateReward(ctx, r)
	if err != nil {
		s.l.Error("创建打赏记录失败",
			logger.String("biz", r.Target.Biz),
			logger.Int64("bizId", r.Target.BizId),
			logger.Int64("uid", r.Uid),
			logger.Error(err))
		return domain.CodeURL{}, err
	}

	// TODO: 这里应该调用支付服务生成支付二维码
	// 目前先返回一个模拟的二维码URL
	codeURL = domain.CodeURL{
		Rid: rid,
		URL: fmt.Sprintf("https://pay.example.com/qr?rid=%d&amt=%d", rid, r.Amt),
	}

	// 缓存二维码
	r.Id = rid
	err = s.repo.CachedCodeURL(ctx, codeURL, r)
	if err != nil {
		s.l.Error("缓存二维码失败",
			logger.Int64("rid", rid),
			logger.Error(err))
		// 缓存失败不影响主流程，继续返回结果
	}

	return codeURL, nil
}

func (s *rewardService) GetReward(ctx context.Context, rid, uid int64) (domain.Reward, error) {
	reward, err := s.repo.GetReward(ctx, rid)
	if err != nil {
		s.l.Error("获取打赏记录失败",
			logger.Int64("rid", rid),
			logger.Int64("uid", uid),
			logger.Error(err))
		return domain.Reward{}, err
	}

	// 验证用户权限：只有打赏者本人可以查看
	if reward.Uid != uid {
		s.l.Warn("用户尝试访问非本人的打赏记录",
			logger.Int64("rid", rid),
			logger.Int64("requestUid", uid),
			logger.Int64("rewardUid", reward.Uid))
		return domain.Reward{}, errors.New("无权访问该打赏记录")
	}

	return reward, nil
}

func (s *rewardService) UpdateReward(ctx context.Context, bizTradeNO string, status domain.RewardStatus) error {
	// TODO: 这里需要根据 bizTradeNO 查找对应的打赏记录
	// 目前的 DAO 层还没有提供根据 bizTradeNO 查询的方法
	// 这通常是支付回调时使用的方法

	// 临时实现：记录日志并返回错误，提示需要完善
	s.l.Error("UpdateReward 方法需要完善",
		logger.String("bizTradeNO", bizTradeNO),
		logger.String("status", fmt.Sprintf("%d", status)))

	return errors.New("UpdateReward 方法暂未完全实现，需要在 DAO 层添加根据 bizTradeNO 查询的方法")
}

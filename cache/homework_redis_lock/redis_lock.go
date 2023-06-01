package homework_redis_lock

import (
	"context"
	_ "embed"
	"errors"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"time"
)

var (
	errLockErr      = errors.New("超过重试次数")
	errLockNotExist = errors.New("锁不存在")
	//注入redis 脚本
	//go:embed lua/lock.lua
	luaLock string
	//go:embed lua/unlock.lua
	luaUnLock string
	//go:embed lua/refresh.lua
	luaRefresh string
)

type Client struct {
	client redis.Cmdable
}
type Lock struct {
	//想要在lock上解锁需要客户端
	client     redis.Cmdable
	key        string
	val        string
	expiration time.Duration
	stopCh     chan struct{}
}

//Unlock 定义在Lock，便于解锁
func (l *Lock) Unlock(ctx context.Context) error {
	res, err := l.client.Eval(ctx, luaUnLock, []string{l.key}, l.val).Int64()
	if err != nil {
		return err
	}
	if res != 0 {
		return errLockNotExist
	}
	return nil
}
func NewClient(client redis.Cmdable) *Client {
	return &Client{client: client}
}
func (l *Lock) Refresh(ctx context.Context) error {
	res, err := l.client.Eval(ctx, luaRefresh, []string{l.key}, l.val, l.expiration.Seconds()).Int64()
	if err != nil {
		return err
	}
	if res != 1 {
		return errLockNotExist
	}
	l.stopCh <- struct{}{}
	return nil
}
func (l *Lock) AutoRefresh(interval, timeout time.Duration) error {
	ticker := time.NewTicker(interval)
	timeoutCh := make(chan struct{}, 1)
	for {
		select {
		//到了续约时间,进行续约
		case <-ticker.C:
			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			res, err := l.client.Eval(ctx, luaRefresh, []string{l.key}, l.val, l.expiration.Seconds()).Int64()
			cancel()
			if err == context.DeadlineExceeded {
				timeoutCh <- struct{}{}
				continue
			}
			if err != nil {
				return err
			}
			if res != 1 {
				return errLockNotExist
			}
		//有续约超时的，可以继续重试
		case <-timeoutCh:
			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			res, err := l.client.Eval(ctx, luaRefresh, []string{l.key}, l.val, l.expiration.Seconds()).Int64()
			cancel()
			if err != nil && err == context.DeadlineExceeded {
				timeoutCh <- struct{}{}
			} else {
				return err
			}
			if res != 1 {
				return errLockNotExist
			}
		//如果解锁了怎么办，就需要在lock中加入一个stopCh让在解锁时通知不用再续约
		case <-l.stopCh:
			return nil
		}
	}
}

//RetryStrategy 重试策略
type RetryStrategy interface {
	Next() (time.Duration, bool)
}

type MaxCntRetry struct {
	//重试超时时间
	interval time.Duration
	maxCnt   int
	cnt      int
}

//Next 返回重试时间间隔和是否重试
func (m *MaxCntRetry) Next() (time.Duration, bool) {
	if m.cnt < m.maxCnt {
		m.cnt++
		return m.interval, true
	}
	return 0, false
}

//Lock 带有重试的加锁
func (c *Client) Lock(ctx context.Context, key string, expiration time.Duration, timeout time.Duration, retry RetryStrategy) (*Lock, error) {
	val := uuid.New().String()
	//创建一个计时器给重试策略的计时使用
	var timer *time.Timer
	//使用for循环来重试
	for {
		//控制Lock超时的ctx
		tctx, cancel := context.WithTimeout(ctx, timeout)
		res, err := c.client.Eval(tctx, luaLock, []string{key}, val, expiration.Seconds()).Result()
		cancel()
		//超时错误是可以解决的，所以判断一下
		if err != nil && err != context.DeadlineExceeded {
			return nil, err
		}
		//超时不能直接重试，要看超重试策略
		//if err == context.DeadlineExceeded {
		//	//超时，进行下一次
		//	continue
		//}
		//如果成功，返回Lock
		if res == "Ok" {
			return &Lock{
				key:        key,
				val:        val,
				expiration: expiration,
				stopCh:     make(chan struct{}, 1),
			}, nil
		}
		//执行重试策略，看还要不要重试
		interval, ok := retry.Next()
		if !ok {
			return nil, errLockErr
		}
		//还要重试,为计时器赋值或重置
		if timer == nil {
			timer = time.NewTimer(interval)
		} else {
			timer.Reset(interval)
		}
		select {
		//等待重试时间间隔
		case <-timer.C:
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
}

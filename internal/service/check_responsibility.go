package service

import (
	"context"
	"errors"
	"net/http"
	"sync"
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"go.uber.org/zap"

	"github.com/nktch1/tank/internal/domain"
)

type ready struct {
	host string
	rps  int
}

func (t *Tank) CheckResponsibility(ctx context.Context, searchResults *SearchResults) (*domain.Response, error) {
	if searchResults == nil {
		return nil, nil
	}

	resp := &domain.Response{
		HostToOptimalRPS: map[string]int{},
	}

	ctx, cancelByTimeout := context.WithTimeout(ctx, t.conf.Timeout)
	defer cancelByTimeout()

	ctx, done := context.WithCancel(ctx)

	//go func() {
	//	for {
	//		select {
	//		case <-ctx.Done():
	//			return
	//		default:
	//		}
	//		println("\n", runtime.NumGoroutine(), "\n")
	//		time.Sleep(time.Millisecond * 100)
	//	}
	//}()

	go t.processHost(ctx, done, resp, searchResults)

	<-ctx.Done()

	ctxzap.Extract(ctx).Info(
		"results",
		zap.Int("completed", len(resp.HostToOptimalRPS)),
		zap.Int("from", len(searchResults.Items)),
	)

	return resp, nil
}

func (t *Tank) processHost(ctx context.Context, done context.CancelFunc,
	resp *domain.Response, queue *SearchResults) {
	defer done()

	var (
		logger    = ctxzap.Extract(ctx)
		rChannel  = make(chan ready, len(queue.Items))
		waitGroup = &sync.WaitGroup{}
	)

	for _, host := range queue.Items {
		logger.Info("queue", zap.String("host", host.Host))

		waitGroup.Add(1)
		go t.benchmark(ctx, host, rChannel, waitGroup)
	}

	go func() {
		waitGroup.Wait()
		close(rChannel)
	}()

	for r := range rChannel {
		t.Lock()
		resp.HostToOptimalRPS[r.host] = r.rps
		t.Unlock()
	}
}

func (t *Tank) benchmark(ctx context.Context, host responseItem, rChannel chan ready, wg *sync.WaitGroup) {
	defer wg.Done()

	var (
		logger            = ctxzap.Extract(ctx)
		waitWorkers       = &sync.WaitGroup{}
		currentRPS        = t.conf.StartRPS
		hostStatusChannel = make(chan error, 1)
		ctxPerHost, done  = context.WithCancel(ctx)
	)

	waitWorkers.Add(1)
	go func() {
		defer waitWorkers.Done()

		err := <-hostStatusChannel
		if err == nil {
			return
		}

		logger.Info(
			"dequeue",
			zap.String("host", host.Host),
			zap.Int("rps", currentRPS),
			zap.String("err", err.Error()),
		)

		rChannel <- ready{
			host: host.Host,
			rps:  currentRPS,
		}

		done()
	}()

	//for {
	start := time.Now()

	for i := 0; i < currentRPS; i++ {
		select {
		case <-ctxPerHost.Done():
			return
		default:
		}

		waitWorkers.Add(1)
		go t.get(ctxPerHost, host.Url, hostStatusChannel, waitWorkers)
	}

	waitWorkers.Wait()
	currentRPS += t.conf.IncreasingStepRPS

	logger.Debug(
		"rps increased",
		zap.String("host", host.Host),
		zap.Int("rps", currentRPS),
		zap.Int64("time", time.Since(start).Milliseconds()),
	)

	//	time.Sleep(time.Second)
	//}
}

func (t *Tank) get(ctx context.Context, url string, statusChan chan error, wg *sync.WaitGroup) {
	defer wg.Done()

	passErr := func(err error) {
		select {
		case statusChan <- err:
		default:
		}
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		passErr(err)
		return
	}

	resp, err := t.client.Do(req.WithContext(ctx))
	if err != nil {
		passErr(err)
		return
	}

	if resp != nil && resp.StatusCode == http.StatusTooManyRequests {
		err = errors.New("invalid status code: " + resp.Status)
		passErr(err)
		return
	}

	passErr(nil)
}

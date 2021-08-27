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

type e struct {
	err error
	rps int
}

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

	ctx, cancelByTimeout := context.WithTimeout(ctx, t.Conf.Timeout)
	defer cancelByTimeout()

	ctx, done := context.WithCancel(ctx)

	go t.processHost(ctx, done, resp, searchResults)

	<-ctx.Done()

	ctxzap.Extract(ctx).Debug(
		"results",
		zap.Int("completed", len(resp.HostToOptimalRPS)),
	)

	return resp, nil
}

func (t *Tank) processHost(ctx context.Context, done context.CancelFunc,
	resp *domain.Response, queue *SearchResults) {
	defer done()

	var (
		logger   = ctxzap.Extract(ctx)
		rChannel = make(chan ready, len(queue.Items))
		wg       = &sync.WaitGroup{}
	)

	for _, host := range queue.Items {
		logger.Info("queue", zap.String("host", host.Host))

		wg.Add(1)
		go t.benchmark(ctx, host, rChannel, wg)
	}

	go func() {
		wg.Wait()
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
		logger           = ctxzap.Extract(ctx)
		waitWorkers      = &sync.WaitGroup{}
		currentRPS       = t.Conf.StartRPS
		statusChannel    = make(chan error)
		ctxPerHost, done = context.WithCancel(ctx)
	)

	go func() {
		status := <-statusChannel

		logger.Info(
			"dequeue",
			zap.String("host", host.Host),
			zap.Int("rps", currentRPS),
			zap.String("status", status.Error()),
		)

		rChannel <- ready{
			host: host.Host,
			rps:  currentRPS,
		}

		done()
	}()

	for {
		start := time.Now()

		for i := 0; i < currentRPS; i++ {
			select {
			case <-ctxPerHost.Done():
				return
			case <-ctx.Done():
				return
			default:
			}

			waitWorkers.Add(1)
			go t.get(ctxPerHost, host.Url, statusChannel, waitWorkers)
		}

		waitWorkers.Wait()
		currentRPS += t.Conf.IncreasingStepRPS

		logger.Debug(
			"rps increased",
			zap.String("host", host.Host),
			zap.Int("rps", currentRPS),
			zap.Int64("time", time.Since(start).Milliseconds()),
		)

		time.Sleep(time.Second)
	}
}

func (t *Tank) get(ctx context.Context, url string, statusChan chan error, wg *sync.WaitGroup) {
	defer wg.Done()

	errC := make(chan error)

	go func() {
		resp, err := http.Get(url)
		if err != nil {
			errC <- err
			return
		}

		if resp != nil && resp.StatusCode == http.StatusTooManyRequests {
			errC <- errors.New("invalid status code: " + resp.Status)
			return
		}

		errC <- nil
	}()

	select {
	case <-ctx.Done():
		return

	case err := <-errC:
		if err != nil {
			statusChan <- err
		}

	case <-time.After(t.Conf.TimeoutPerHost):
		statusChan <- errors.New("waiting time is exceeded")
	}
}
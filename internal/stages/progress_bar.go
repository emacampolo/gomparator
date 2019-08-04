package stages

import (
	"sync"

	"gopkg.in/cheggaaa/pb.v1"
)

type ProgressBar struct {
	okPb    *pb.ProgressBar
	errorPb *pb.ProgressBar
	pool    *pb.Pool
}

func NewProgressBar(total int) *ProgressBar {
	okPb := makeProgressBar(total, "ok")
	errorPb := makeProgressBar(total, "error")

	return &ProgressBar{
		okPb,
		errorPb,
		nil,
	}
}

func (p *ProgressBar) IncrementOk() {
	p.okPb.Add(1)
}

func (p *ProgressBar) IncrementError() {
	p.errorPb.Add(1)
}

func (p *ProgressBar) Start() {
	pool, err := pb.StartPool(p.okPb, p.errorPb)
	if err != nil {
		panic(err)
	}
	p.pool = pool
	p.okPb.Start()
}

func (p *ProgressBar) Stop() {
	wg := new(sync.WaitGroup)
	for _, bar := range []*pb.ProgressBar{p.okPb, p.errorPb} {
		wg.Add(1)
		go func(cb *pb.ProgressBar) {
			cb.Finish()
			wg.Done()
		}(bar)
	}
	wg.Wait()
	// close pool
	_ = p.pool.Stop()
}

func makeProgressBar(total int, prefix string) *pb.ProgressBar {
	bar := pb.New(total)
	bar.Prefix(prefix)
	bar.SetMaxWidth(120)
	bar.ShowElapsedTime = true
	bar.ShowTimeLeft = false
	return bar
}

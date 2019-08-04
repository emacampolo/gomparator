package pipeline

import (
	"context"

	"github.com/ecampolo/gomparator/internal/stages"
)

type Reader interface {
	Read() <-chan stages.URLPair
}

type Producer interface {
	Produce(in <-chan stages.URLPair) <-chan stages.HostsPair
}

type Consumer interface {
	Consume(in stages.HostsPair)
}

func New(reader Reader, producer Producer, ctx context.Context, consumer Consumer) *Pipeline {
	return &Pipeline{
		reader:   reader,
		producer: producer,
		consumer: consumer,
		ctx:      ctx,
	}
}

type Pipeline struct {
	reader   Reader
	producer Producer
	consumer Consumer
	ctx      context.Context
}

func (p *Pipeline) Run() {
	readStream := p.reader.Read()
	producerStream := p.producer.Produce(readStream)

	orDone := func(ctx context.Context, c <-chan stages.HostsPair) <-chan stages.HostsPair {
		valStream := make(chan stages.HostsPair)
		go func() {
			defer close(valStream)
			for {
				select {
				case <-ctx.Done():
					return
				case v, ok := <-c:
					if !ok {
						return
					}
					select {
					case valStream <- v:
					case <-ctx.Done():
					}
				}
			}
		}()
		return valStream
	}

	for val := range orDone(p.ctx, producerStream) {
		p.consumer.Consume(val)
	}
}

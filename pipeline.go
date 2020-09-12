package main

import (
	"context"
)

type Reader interface {
	Read() <-chan URLPair
}

type Producer interface {
	Produce(in <-chan URLPair) <-chan HostsPair
}

type Consumer interface {
	Consume(in HostsPair)
}

func New(reader Reader, producer Producer, consumer Consumer) *Pipeline {
	return &Pipeline{
		reader:   reader,
		producer: producer,
		consumer: consumer,
	}
}

type Pipeline struct {
	reader   Reader
	producer Producer
	consumer Consumer
}

func (p *Pipeline) Run(ctx context.Context) {
	readStream := p.reader.Read()
	producerStream := p.producer.Produce(readStream)

	orDone := func(ctx context.Context, c <-chan HostsPair) <-chan HostsPair {
		valStream := make(chan HostsPair)
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

	for val := range orDone(ctx, producerStream) {
		p.consumer.Consume(val)
	}
}

package service

import (
	"fundinvest/internal/task004/cache"
	"fundinvest/internal/task004/exporter"
	"fundinvest/internal/task004/parser"
)

type Imports struct {
	decoder parser.Decoder
	cache   *cache.Store
	queue   *exporter.Queue
}

func New(cache *cache.Store, queue *exporter.Queue) *Imports {
	return &Imports{cache: cache, queue: queue}
}

func (i *Imports) Accept(id, payload string) {
	record := i.decoder.Parse(id, payload)
	i.cache.Save(record)
	i.queue.Add(record)
}

package exporter

import "fundinvest/internal/task004/parser"

type Queue struct{ pending []parser.Record }

func (q *Queue) Add(record parser.Record) { q.pending = append(q.pending, record) }

func (q *Queue) Flush() []string {
	result := make([]string, 0, len(q.pending))
	for _, record := range q.pending {
		result = append(result, string(record.Payload))
	}
	q.pending = nil
	return result
}

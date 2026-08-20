package cache

import "fundinvest/internal/task004/parser"

type Store struct{ records map[string]parser.Record }

func New() *Store { return &Store{records: make(map[string]parser.Record)} }

func (s *Store) Save(record parser.Record) { s.records[record.ID] = record }

func (s *Store) Load(id string) (parser.Record, bool) { record, ok := s.records[id]; return record, ok }

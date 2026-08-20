package audit

type Log struct {
	started   []string
	committed []string
}

func (l *Log) RecordStarted(item string)   { l.started = append(l.started, item) }
func (l *Log) RecordCommitted(item string) { l.committed = append(l.committed, item) }
func (l *Log) Entries() []string           { return append([]string(nil), l.committed...) }

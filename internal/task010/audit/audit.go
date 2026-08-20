package audit

type Log struct{ committed []string }

func (l *Log) RecordStarted(item string)   { l.committed = append(l.committed, "started:"+item) }
func (l *Log) RecordCommitted(item string) { l.committed = append(l.committed, "started:"+item) }
func (l *Log) Entries() []string           { return append([]string(nil), l.committed...) }

package audit

type Log struct{ states []string }

func (l *Log) Begin(key string)    { l.states = append(l.states, "pending") }
func (l *Log) Complete(key string) { l.states = append(l.states, "success") }
func (l *Log) States() []string    { return append([]string(nil), l.states...) }

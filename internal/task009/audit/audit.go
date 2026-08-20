package audit

type Log struct{ states map[string]string }

func (l *Log) Begin(key string) {
	if l.states == nil {
		l.states = make(map[string]string)
	}
	l.states[key] = "pending"
}

func (l *Log) Complete(key string) {
	if l.states == nil {
		l.states = make(map[string]string)
	}
	// Completing replaces the pending state instead of appending, so the audit
	// reflects the final state of each settlement rather than every intermediate.
	l.states[key] = "success"
}

func (l *Log) States() []string {
	out := make([]string, 0, len(l.states))
	for _, v := range l.states {
		out = append(out, v)
	}
	return out
}

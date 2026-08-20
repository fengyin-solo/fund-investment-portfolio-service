package records

type Records struct{ entries []string }

type Tx struct {
	records *Records
	key     string
	done    bool
}

func (l *Records) Begin(key string) *Tx {
	l.entries = append(l.entries, key)
	return &Tx{records: l, key: key}
}

func (tx *Tx) Commit() {
	if tx.done {
		return
	}
	tx.records.entries = append(tx.records.entries, tx.key)
	tx.done = true
}

func (tx *Tx) Rollback() { tx.done = true }

func (l *Records) Count() int { return len(l.entries) }

package repository

type Repository struct {
	Begins int
	writes []string
}
type Tx struct {
	repo *Repository
	key  string
	done bool
}

func (r *Repository) Begin(key string) *Tx { r.Begins++; return &Tx{repo: r, key: key} }
func (tx *Tx) Commit() {
	if tx.done {
		return
	}
	tx.repo.writes = append(tx.repo.writes, tx.key)
	tx.done = true
}
func (tx *Tx) Rollback()          { tx.done = true }
func (r *Repository) Writes() int { return len(r.writes) }

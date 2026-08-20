package model

type Job struct {
	ID, State string
	Version   int
}

func New(id string) Job     { return Job{ID: id, State: "running", Version: 1} }
func (j Job) Complete() Job { j.State = "succeeded"; return j }

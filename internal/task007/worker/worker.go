package worker

import (
	"fundinvest/internal/task007/cache"
	"fundinvest/internal/task007/model"
)

type Worker struct{ Cache *cache.Store }

func (w Worker) Publish(job model.Job)      { w.Cache.Apply(job) }
func (w Worker) Delayed(snapshot model.Job) { w.Cache.Apply(snapshot) }

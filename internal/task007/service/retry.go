package service

import (
	"fmt"
	"fundinvest/internal/task007/effects"
	"fundinvest/internal/task007/model"
	"fundinvest/internal/task007/worker"
)

func Retry(job model.Job, client *effects.Client, w worker.Worker) (model.Job, error) {
	// The idempotency key is stable across attempts: the same logical
	// settlement must carry the same key so the gateway can dedup a retry
	// instead of charging the side effect a second time.
	key := fmt.Sprintf("settlement:%s", job.ID)
	for attempt := 0; attempt < 2; attempt++ {
		if err := client.Execute(key); err != nil {
			if attempt == 1 {
				return job, err
			}
			continue
		}
		completed := job.Complete()
		w.Publish(completed)
		return completed, nil
	}
	return job, nil
}

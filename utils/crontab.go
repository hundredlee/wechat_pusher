package utils

import "github.com/robfig/cron"

var (
	cronInstance *cron.Cron
)

func getCronInstance() *cron.Cron {
	if cronInstance != nil {
		return cronInstance
	}
	cronInstance = cron.New()
	return cronInstance
}

func StartCron() {
	getCronInstance().Start()
	select {}
}

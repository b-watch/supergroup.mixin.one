package plugin

import "github.com/robfig/cron"

var c = cron.New()

func (*PluginContext) RegisterCronJob(spec string, cmd func()) error {
	mutex.Lock()
	defer mutex.Unlock()
	err := c.AddFunc(spec, cmd)
	return err
}

func RunCron() {
	c.Start()
}

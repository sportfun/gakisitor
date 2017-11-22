package kernel

import (
	"fmt"

	"github.com/xunleii/fantastic-broccoli/common/types"
	"github.com/xunleii/fantastic-broccoli/common/types/service"
	"github.com/xunleii/fantastic-broccoli/constant"
	"github.com/xunleii/fantastic-broccoli/log"
	"github.com/xunleii/fantastic-broccoli/properties"
)

type Core struct {
	services   []service.Service
	logger     log.Logger
	properties *properties.Properties

	notifications *service.NotificationQueue
	internal      error
	state         types.StateType
}

var (
	InfoStartServices   = log.NewArgumentBinder("start services")
	InfoServicesStarted = log.NewArgumentBinder("services successfully started (%d services)")
)

func (core *Core) Configure(services []service.Service, props *properties.Properties, logger log.Logger) error {
	// Property file can be not loaded (props.IsLoaded = false) if file not found or invalid
	if !props.IsLoaded() {
		return fmt.Errorf("properties not loaded")
	}

	core.services = services
	core.logger = logger
	core.notifications = service.NewNotificationQueue()

	core.internal = nil
	logger.Info(InfoStartServices)
	for _, s := range services {
		if !core.checkIf(s, s.Start(core.notifications, logger), IsStarted) ||
			!core.checkIf(s, s.Configure(props), IsConfigured) {
			return core.internal
		}
	}
	logger.Info(InfoServicesStarted.Bind(len(services)))

	core.state = constant.States.Idle
	return nil
}

func (core *Core) Run() error {
	core.state = constant.States.Working
	for _, s := range core.services {
		if !core.checkIf(s, s.Process(), IsProcessed) {
			return core.internal
		}

		for _, n := range core.notifications.Notifications(constant.EntityNames.Core) {
			core.handle(n)
		}
	}
	core.state = constant.States.Idle
	return nil
}

func (core *Core) Stop() error {
	for _, s := range core.services {
		if core.checkIf(s, s.Stop(), IsStopped) {
			return core.internal
		}
	}
	core.state = constant.States.Stopped
	return nil
}

func (core *Core) State() types.StateType {
	return core.state
}

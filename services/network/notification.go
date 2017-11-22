package network

import (
	"fmt"

	"github.com/xunleii/fantastic-broccoli/common/types/notification"
	"github.com/xunleii/fantastic-broccoli/common/types/notification/object"
	"github.com/xunleii/fantastic-broccoli/constant"
	"github.com/xunleii/fantastic-broccoli/log"
)

var (
	DebugNotificationHandled    = log.NewArgumentBinder("notification handled")
	UnhandledNotificationOrigin = log.NewArgumentBinder("unhandled notification origin (%s)")
	UnknownContentType          = log.NewArgumentBinder("unknown content type")
)

func (service *Service) handle(n *notification.Notification) error {
	service.logger.Debug(DebugNotificationHandled.More("notification", *n))

	switch string(n.From()) {
	case constant.EntityNames.Services.Module:
		fallthrough
	case constant.EntityNames.Core:
		return defaultNotificationHandler(service, n)
	default:
		service.logger.Warn(UnhandledNotificationOrigin.Bind(n.From()).More("content", n.Content()))
	}
	return nil
}

func defaultNotificationHandler(service *Service, n *notification.Notification) error {
	var succeed = true

	switch o := n.Content().(type) {
	case *object.CommandObject:
		succeed = service.emit(constant.Channels.Command.String(), *o)
	case *object.DataObject:
		succeed = service.emit(constant.Channels.Data.String(), *o)
	case *object.ErrorObject:
		succeed = service.emit(constant.Channels.Error.String(), *o)
	default:
		service.logger.Warn(UnknownContentType.More("packet", n.Content()))
	}

	if !succeed {
		return fmt.Errorf("failed to emit message")
	}
	return nil
}

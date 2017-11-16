package module

import (
	"fantastic-broccoli/properties"
	"go.uber.org/zap"
)

// TODO: [v1.x] Retourner une erreur avec type precis (pour la gestion d'erreur)
type Module interface {
	Start(queue *NotificationQueue, logger *zap.Logger) error
	Configure(properties *properties.Properties) error
	Process() error
	Stop() error

	StartSession() error
	StopSession() error

	Name() string
	State() int
}

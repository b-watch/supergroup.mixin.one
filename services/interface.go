package services

import (
	"context"

	"github.com/MixinNetwork/supergroup.mixin.one/models"
)

type Service interface {
	Run(context.Context, chan models.WsBroadcastMessage) error
}

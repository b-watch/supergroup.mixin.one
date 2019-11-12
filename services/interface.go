package services

import "context"

type Service interface {
	Run(context.Context, chan WsBroadcastMessage) error
}

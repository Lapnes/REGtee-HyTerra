package auth

import "backend/service/auth"

type Controller struct {
	svc auth.Servicer
}

func NewController(
	svc auth.Servicer,
) *Controller {
	return &Controller{
		svc: svc,
	}
}

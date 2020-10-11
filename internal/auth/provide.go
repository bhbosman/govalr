package auth

import "go.uber.org/fx"

func Provide() fx.Option {
	return fx.Provide(
		func() *ValrAuth {
			return &ValrAuth{}
		})
}

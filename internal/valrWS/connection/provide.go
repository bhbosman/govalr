package connection

import (
	"github.com/bhbosman/gocomms/impl"
	"github.com/bhbosman/gocomms/intf"
	"github.com/bhbosman/gocomms/netDial"
	"github.com/bhbosman/govalr/internal/auth"
	"github.com/bhbosman/govalr/internal/keys"
	"github.com/cskr/pubsub"
	"go.uber.org/fx"
)

const FactoryName = "ValrWSS"

func ProvideValrWsDialer(
	canDial netDial.ICanDial) fx.Option {
	var canDials []netDial.ICanDial
	if canDial != nil {
		canDials = append(canDials, canDial)
	}

	const ValrDialerConst = "ValrDialer"
	return fx.Options(
		fx.Provide(
			fx.Annotated{
				Group: impl.ConnectionReactorFactoryConst,
				Target: func(
					params struct {
						fx.In
						PubSub   *pubsub.PubSub `name:"Application"`
						Settings *keys.ValrConnectionSettings
						ValrAuth *auth.ValrAuth
					}) (intf.IConnectionReactorFactory, error) {
					return NewFactory(
						ValrDialerConst,
						params.PubSub,
						params.Settings,
						params.ValrAuth), nil
				},
			}),
		fx.Provide(
			fx.Annotated{
				Group: "Apps",
				Target: netDial.NewNetDialApp(
					"valr",
					"wss://api.valr.com:443/ws/trade",
					impl.WebSocketName,
					ValrDialerConst,
					netDial.MaxConnectionsSetting(1),
					netDial.CanDial(canDials...)),
			}))
}

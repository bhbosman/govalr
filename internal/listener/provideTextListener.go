package listener

import (
	"encoding/json"
	"github.com/bhbosman/gocomms/impl"
	"github.com/bhbosman/gocomms/intf"
	"github.com/bhbosman/gocomms/netListener"
	"github.com/bhbosman/gomessageblock"
	"github.com/bhbosman/goprotoextra"
	"github.com/bhbosman/govalr/internal/ConsumerCounter"
	"github.com/cskr/pubsub"
	"go.uber.org/fx"
)

import (
	"google.golang.org/protobuf/proto"
)

func TextListener(maxConnections int, url string) fx.Option {
	const TextListenerConnection = "TextListenerConnection"
	return fx.Options(
		fx.Provide(
			fx.Annotated{
				Group: impl.ConnectionReactorFactoryConst,
				Target: func(params struct {
					fx.In
					ConsumerCounter *ConsumerCounter.ConsumerCounter
					PubSub          *pubsub.PubSub `name:"Application"`
				}) (intf.IConnectionReactorFactory, error) {
					return NewFactory(
						TextListenerConnection,
						params.PubSub,
						func(m proto.Message) (goprotoextra.IReadWriterSize, error) {
							bytes, err := json.MarshalIndent(m, "", "\t")
							if err != nil {
								return nil, err
							}
							return gomessageblock.NewReaderWriterBlock(bytes), nil
						},
						params.ConsumerCounter), nil
				},
			}),
		fx.Provide(
			fx.Annotated{
				Group: "Apps",
				Target: netListener.NewNetListenApp(
					TextListenerConnection,
					url,
					impl.TransportFactoryEmptyName,
					TextListenerConnection,
					netListener.MaxConnectionsSetting(maxConnections)),
			}),
	)
}

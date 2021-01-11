package listener

import (
	"github.com/bhbosman/gocommon/stream"
	"github.com/bhbosman/gocomms/impl"
	"github.com/bhbosman/gocomms/intf"
	"github.com/bhbosman/gocomms/netListener"
	"github.com/bhbosman/goprotoextra"
	"github.com/bhbosman/govalr/internal/ConsumerCounter"
	"github.com/cskr/pubsub"
	"go.uber.org/fx"
	"google.golang.org/protobuf/proto"
)

func CompressedListener(maxConnections int, url string) fx.Option {
	const CompressedListenerConnection = "CompressedListenerConnection"
	return fx.Options(
		fx.Provide(
			fx.Annotated{
				Group: impl.ConnectionReactorFactoryConst,
				Target: func(params struct {
					fx.In
					PubSub          *pubsub.PubSub `name:"Application"`
					ConsumerCounter *ConsumerCounter.ConsumerCounter
				}) (intf.IConnectionReactorFactory, error) {
					return NewFactory(
						CompressedListenerConnection,
						params.PubSub,
						func(data proto.Message) (goprotoextra.IReadWriterSize, error) {
							return stream.Marshall(data)
						},
						params.ConsumerCounter), nil
				},
			}),
		fx.Provide(
			fx.Annotated{
				Group: "Apps",
				Target: netListener.NewNetListenApp(
					CompressedListenerConnection,
					url,
					impl.TransportFactoryCompressedName,
					CompressedListenerConnection,
					netListener.MaxConnectionsSetting(maxConnections)),
			}),
	)
}

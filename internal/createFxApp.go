package internal

import (
	app2 "github.com/bhbosman/gocommon/app"
	"github.com/bhbosman/gocommon/fxHelper"
	"github.com/bhbosman/gocomms/connectionManager"
	"github.com/bhbosman/gocomms/connectionManager/endpoints"
	"github.com/bhbosman/gocomms/connectionManager/view"
	"github.com/bhbosman/gocomms/impl"
	"github.com/bhbosman/gocomms/provide"
	"github.com/bhbosman/gologging"
	"github.com/bhbosman/govalr/internal/ConsumerCounter"
	"github.com/bhbosman/govalr/internal/auth"
	"github.com/bhbosman/govalr/internal/keys"
	"github.com/bhbosman/govalr/internal/listener"
	"github.com/bhbosman/govalr/internal/valrWS/connection"
	"go.uber.org/fx"
	"log"
	"os"
	//"path"
)

func CreateFxApp() (*fx.App, fx.Shutdowner) {
	settings := &AppSettings{
		Logger:                log.New(os.Stderr, "", log.LstdFlags),
		textListenerUrl:       "tcp4://127.0.0.1:3020",
		compressedListenerUrl: "tcp4://127.0.0.1:3021",
		HttpListenerUrl:       "http://127.0.0.1:8082",
	}
	ConsumerCounter := &ConsumerCounter.ConsumerCounter{}
	var shutDowner fx.Shutdowner
	fxApp := fx.New(
		fx.Supply(settings, ConsumerCounter),
		fx.Logger(settings.Logger),
		gologging.ProvideLogFactory(settings.Logger, nil),
		fx.Populate(&shutDowner),
		app2.RegisterRootContext(),
		connectionManager.RegisterDefaultConnectionManager(),
		provide.RegisterHttpHandler(settings.HttpListenerUrl),
		auth.Provide(),
		endpoints.RegisterConnectionManagerEndpoint(),
		view.RegisterConnectionsHtmlTemplate(),
		impl.RegisterAllConnectionRelatedServices(),
		keys.ProvideValrConnectionSettings(),
		connection.ProvideValrWsDialer(ConsumerCounter),
		listener.TextListener(1024, settings.textListenerUrl),
		listener.CompressedListener(1024, settings.compressedListenerUrl),
		fxHelper.InvokeApps(),
	)
	return fxApp, shutDowner
}

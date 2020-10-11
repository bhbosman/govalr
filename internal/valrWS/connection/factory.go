package connection

import (
	"context"
	"fmt"
	"github.com/bhbosman/gocomms/intf"
	"github.com/bhbosman/gologging"
	"github.com/bhbosman/govalr/internal/auth"
	"github.com/bhbosman/govalr/internal/keys"
	"github.com/cskr/pubsub"
	"net/url"
	"strconv"
	"time"
)

type Factory struct {
	name     string
	PubSub   *pubsub.PubSub
	Settings *keys.ValrConnectionSettings
	ValrAuth *auth.ValrAuth
}

func (self *Factory) Values(inputValues map[string]interface{}) (map[string]interface{}, error) {
	intfValue, ok := inputValues["url"]
	if !ok {
		return nil, fmt.Errorf("url required, but not found")
	}
	urlObject, ok := intfValue.(*url.URL)
	if !ok {
		return nil, fmt.Errorf(`object market as "url", is not of type "*url.URL"`)
	}
	timeStamp := strconv.FormatInt(time.Now().UnixNano()/1000000, 10)
	connectionHeaderMap := make(map[string][]string)
	connectionHeaderMap["X-VALR-API-KEY"] = []string{
		self.Settings.ApiKey,
	}
	connectionHeaderMap["X-VALR-SIGNATURE"] = []string{
		self.ValrAuth.Hash(
			self.Settings.SecretKey,
			timeStamp,
			urlObject.Path,
			"GET",
			""),
	}
	connectionHeaderMap["X-VALR-TIMESTAMP"] = []string{
		timeStamp,
	}

	result := make(map[string]interface{})
	result["connectionHeader"] = connectionHeaderMap
	return result, nil
}

func (self Factory) Name() string {
	return self.name
}

func (self Factory) Create(
	name string,
	cancelCtx context.Context,
	cancelFunc context.CancelFunc,
	logger *gologging.SubSystemLogger,
	userContext interface{}) intf.IConnectionReactor {
	return NewReactor(
		logger,
		name,
		cancelCtx,
		cancelFunc,
		userContext,
		self.PubSub)
}

func NewFactory(
	name string,
	pubSub *pubsub.PubSub,
	settings *keys.ValrConnectionSettings,
	ValrAuth *auth.ValrAuth) *Factory {
	return &Factory{
		name:     name,
		PubSub:   pubSub,
		Settings: settings,
		ValrAuth: ValrAuth,
	}
}

package connection

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/bhbosman/gocommon/messageRouter"
	"github.com/bhbosman/gocommon/stream"
	"github.com/bhbosman/gocomms/connectionManager"
	"github.com/bhbosman/gocomms/impl"
	"github.com/bhbosman/gocomms/stacks/websocket/wsmsg"
	"github.com/bhbosman/gologging"
	"github.com/bhbosman/gomessageblock"
	"github.com/bhbosman/goprotoextra"
	valrWsStream "github.com/bhbosman/govalr/internal/valrWS/internal/stream"
	"github.com/cskr/pubsub"
	"github.com/golang/protobuf/jsonpb"
	"github.com/reactivex/rxgo/v2"
	"net"
	"net/url"
)

type RePublishMessage struct {
}

type registrationKey struct {
	pair string
	name string
}

func newRegistrationKey(pair string, name string) registrationKey {
	return registrationKey{pair: pair, name: name}
}

type outstandingSubscription struct {
	Reqid uint32
	Pair  string
	Name  string
}

type registeredSubscription struct {
	channelName string
	channelId   uint32
	Reqid       uint32
	Pair        string
	Name        string
}

type Reactor struct {
	impl.BaseConnectionReactor
	messageRouter            *messageRouter.MessageRouter
	connectionID             uint64
	status                   string
	version                  string
	PubSub                   *pubsub.PubSub
	outstandingSubscriptions map[uint32]outstandingSubscription
	registeredSubscriptions  map[uint32]registeredSubscription
	reqid                    uint32
	republishChannelName     string
	publishChannelName       string
}

type WebsocketDataResponse []interface{}

func (self *Reactor) handleWebSocketMessageWrapper(inData *wsmsg.WebSocketMessageWrapper) error {
	switch inData.Data.OpCode {
	case wsmsg.WebSocketMessage_OpText:
		if len(inData.Data.Message) > 0 && inData.Data.Message[0] == '[' { //type WebsocketDataResponse []interface{}
			var dataResponse WebsocketDataResponse
			err := json.Unmarshal(inData.Data.Message, &dataResponse)
			if err != nil {
				return err
			}
			_, _ = self.messageRouter.Route(dataResponse)
			return nil

		} else {
			krakenMessage := &valrWsStream.KrakenWsMessageIncoming{}
			unMarshaler := jsonpb.Unmarshaler{
				AllowUnknownFields: true,
				AnyResolver:        nil,
			}
			err := unMarshaler.Unmarshal(bytes.NewBuffer(inData.Data.Message), krakenMessage)
			if err != nil {
				return err
			}
			_, _ = self.messageRouter.Route(krakenMessage)
			return nil
		}
	case wsmsg.WebSocketMessage_OpEndLoop:
		return nil

	case wsmsg.WebSocketMessage_OpStartLoop:
		return nil
	default:
		return nil
	}
}

func (self Reactor) handleMessageBlockReaderWriter(inData *gomessageblock.ReaderWriter) error {
	marshal, err := stream.UnMarshal(inData, self.CancelCtx, self.CancelFunc, self.ToReactor, self.ToConnection)
	if err != nil {
		println(err.Error())
		return err
	}

	_, err = self.messageRouter.Route(marshal)
	if err != nil {
		return err
	}

	return nil
}

func (self *Reactor) Init(
	conn net.Conn,
	url *url.URL,
	connectionId string,
	connectionManager connectionManager.IConnectionManager,
	onSend goprotoextra.ToConnectionFunc,
	toConnectionReactor goprotoextra.ToReactorFunc) (rxgo.NextExternalFunc, error) {
	_, err := self.BaseConnectionReactor.Init(
		conn,
		url,
		connectionId,
		connectionManager,
		onSend,
		toConnectionReactor)
	if err != nil {
		return nil, err
	}

	self.republishChannelName = "republishChannel"
	self.publishChannelName = "publishChannel"

	republishChannel := self.PubSub.Sub(self.republishChannelName)
	go func(ch chan interface{}, topics ...string) {
		<-self.CancelCtx.Done()
		self.PubSub.Unsub(ch, topics...)
	}(republishChannel, self.republishChannelName)

	go func(ch chan interface{}, topics ...string) {
		for range ch {
			if self.CancelCtx.Err() == nil {
				_ = self.ToReactor(false, &RePublishMessage{})
			}
		}
	}(republishChannel, self.republishChannelName)

	return self.doNext, nil
}

func (self *Reactor) Close() error {
	println("close")
	return self.BaseConnectionReactor.Close()
}
func (self *Reactor) Open() error {
	err := self.BaseConnectionReactor.Open()
	if err != nil {
		return err
	}

	return nil
}

func (self *Reactor) doNext(b bool, i interface{}) {
	_, _ = self.messageRouter.Route(i)
}

func (self *Reactor) handleSystemStatus(data valrWsStream.ISystemStatus) error {
	self.status = data.GetStatus()
	self.connectionID = data.GetConnectionID()
	self.version = data.GetVersion()
	return nil
}

func (self Reactor) handlePing(data valrWsStream.IPing) error {
	outgoing := &valrWsStream.KrakenWsMessageOutgoing{}
	outgoing.Event = "pong"
	outgoing.Reqid = data.GetReqid()

	return SendTextOpMessage(outgoing, self.ToConnection)
}

func (self Reactor) handleHeartbeat(data interface{}) error {
	return nil
}

func (self *Reactor) handleSubscriptionStatus(inData valrWsStream.ISubscriptionStatus) error {
	if data, ok := self.outstandingSubscriptions[inData.GetReqid()]; ok {
		if inData.GetStatus() == "error" {
			return self.Logger.ErrorWithDescription("subscription failed", fmt.Errorf(inData.GetErrorMessage()))
		}
		delete(self.outstandingSubscriptions, data.Reqid)

		self.registeredSubscriptions[inData.GetChannelID()] = registeredSubscription{
			channelName: inData.GetChannelName(),
			channelId:   inData.GetChannelID(),
			Reqid:       inData.GetReqid(),
			Pair:        inData.GetPair(),
			Name:        inData.GetSubscription().Name,
		}
	}
	return nil
}

func (self *Reactor) HandlePublishMessage(msg *RePublishMessage) error {
	return self.publishData(true)
}

func (self *Reactor) HandleEmptyQueue(msg *rxgo.EmptyQueue) error {
	return self.publishData(false)
}

func (self *Reactor) publishData(forcePublish bool) error {
	return nil
}

func NewReactor(
	logger *gologging.SubSystemLogger,
	name string,
	cancelCtx context.Context,
	cancelFunc context.CancelFunc,
	userContext interface{},
	PubSub *pubsub.PubSub) *Reactor {
	result := &Reactor{
		BaseConnectionReactor: impl.NewBaseConnectionReactor(
			logger,
			name,
			cancelCtx,
			cancelFunc,
			userContext),
		messageRouter:            messageRouter.NewMessageRouter(),
		connectionID:             0,
		status:                   "",
		version:                  "",
		PubSub:                   PubSub,
		outstandingSubscriptions: make(map[uint32]outstandingSubscription),
		registeredSubscriptions:  make(map[uint32]registeredSubscription),
	}
	result.messageRouter.Add(result.handleMessageBlockReaderWriter)
	result.messageRouter.Add(result.handleWebSocketMessageWrapper)
	result.messageRouter.Add(result.HandleEmptyQueue)
	result.messageRouter.Add(result.HandlePublishMessage)

	return result
}

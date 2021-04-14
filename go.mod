module github.com/bhbosman/govalr

go 1.15

require (
	github.com/bhbosman/goMessages v0.0.0-20201004192822-66d168b4b744
	github.com/bhbosman/gocommon v0.0.0-20201004145117-eae02ab42c9a
	github.com/bhbosman/gocomms v0.0.0-20210108094235-212b4e8c628c
	github.com/bhbosman/goerrors v0.0.0-20210201065523-bb3e832fa9ab
	github.com/bhbosman/gologging v0.0.0-20200921180328-d29fc55c00bc
	github.com/bhbosman/gomessageblock v0.0.0-20200921180725-7cd29a998aa3
	github.com/cskr/pubsub v1.0.2
	github.com/golang/protobuf v1.4.2
	github.com/reactivex/rxgo/v2 v2.1.0
	go.uber.org/fx v1.13.1
	google.golang.org/protobuf v1.25.0
	github.com/bhbosman/goprotoextra v0.0.2-0.20210414124526-a342e2a9e82f
)


replace (
	github.com/reactivex/rxgo/v2  => ../../reactivex/rxgo
	github.com/bhbosman/goMessages => ../goMessages
	github.com/bhbosman/gocommon => ../gocommon
	github.com/bhbosman/gocomms => ../gocomms
	github.com/bhbosman/gologging => ../gologging
	github.com/bhbosman/gomessageblock => ../gomessageblock
)

module github.com/embedscript/backend

go 1.15

require (
	github.com/golang/protobuf v1.4.3
	github.com/google/uuid v1.1.2
	github.com/gosimple/slug v1.9.0
	github.com/m3o/services v0.0.0-20201118173211-acf31ee96432
	github.com/micro/dev v0.0.0-20201117163752-d3cfc9788dfa
	github.com/micro/micro/v3 v3.0.5-0.20210205114115-75aad3b94f08
	github.com/micro/services v0.16.1-0.20210217104759-f62bcadf262e
	github.com/onsi/gomega v1.7.1
	github.com/patrickmn/go-cache v2.1.0+incompatible
	github.com/teris-io/shortid v0.0.0-20171029131806-771a37caa5cf
	golang.org/x/crypto v0.0.0-20201012173705-84dcc777aaee
	golang.org/x/net v0.0.0-20201010224723-4f7140c49acb
	google.golang.org/protobuf v1.25.0
)

// This can be removed once etcd becomes go gettable, version 3.4 and 3.5 is not,
// see https://github.com/etcd-io/etcd/issues/11154 and https://github.com/etcd-io/etcd/issues/11931.
replace google.golang.org/grpc => google.golang.org/grpc v1.26.0

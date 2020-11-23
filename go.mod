module github.com/embedscript/backend

go 1.15

require (
	github.com/golang/protobuf v1.4.3
	github.com/google/uuid v1.1.2
	github.com/m3o/services v0.0.0-20201118173211-acf31ee96432
	github.com/micro/micro/v3 v3.0.0
	github.com/patrickmn/go-cache v2.1.0+incompatible
)

// This can be removed once etcd becomes go gettable, version 3.4 and 3.5 is not,
// see https://github.com/etcd-io/etcd/issues/11154 and https://github.com/etcd-io/etcd/issues/11931.
replace google.golang.org/grpc => google.golang.org/grpc v1.26.0

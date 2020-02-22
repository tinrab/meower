package event

import "github.com/tinrab/meower/schema"

type Store interface {
	Close()
	PublishMeowCreated(meow schema.Meow) error
	SubscribeMeowCreated() (<-chan MeowCreatedMessage, error)
	OnMeowCreated(f func(MeowCreatedMessage)) error
}

var impl Store

func SetEventStore(es Store) {
	impl = es
}

func Close() {
	impl.Close()
}

func PublishMeowCreated(meow schema.Meow) error {
	return impl.PublishMeowCreated(meow)
}

func SubscribeMeowCreated() (<-chan MeowCreatedMessage, error) {
	return impl.SubscribeMeowCreated()
}

func OnMeowCreated(f func(MeowCreatedMessage)) error {
	return impl.OnMeowCreated(f)
}

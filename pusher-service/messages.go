package main

const (
	KindMeowCreated = iota + 1
)

type MeowCreatedMessage struct {
	Kind uint32 `json:"kind"`
	ID   string `json:"id"`
	Body string `json:"body"`
}

func newMeowCreatedMessage(id string, body string) *MeowCreatedMessage {
	return &MeowCreatedMessage{
		Kind: KindMeowCreated,
		ID:   id,
		Body: body,
	}
}

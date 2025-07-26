package types

//go:generate mockgen -destination=../../mocks/types.go -package=mocks github/erastusk/tracer/internal/types KafkaApps,RedisApps,TCPDial
type KafkaApps interface {
	Connectivity(PromptOptions) error
	ListTopics(PromptOptions) error
	Produce(PromptOptions, string) error
	Consume(PromptOptions, string) error
}
type RedisApps interface {
	Connectivity(PromptOptions) error
}
type TCPDial interface {
	TCPDial(server string) error
}

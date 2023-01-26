package controller

type KafkaController interface {
	Manager() chan error
	Worker()
}

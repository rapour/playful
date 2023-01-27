package controller

type HttpController interface {
	Manager() chan error
}

type KafkaController interface {
	Manager() chan error
	Worker()
}

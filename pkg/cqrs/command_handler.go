package cqrs

type CommandHandler[TCommand Command, TResult interface{}] interface {
	Handle(Command) (TResult, error)
}

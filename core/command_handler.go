package core

type ForthCommandHandler interface {
	Execute(command string) string
}

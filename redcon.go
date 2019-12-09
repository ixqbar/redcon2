package redcon2

import (
	"github.com/tidwall/redcon"
	"strings"
)

type RedconHandlerFunc func(conn redcon.Conn, cmd redcon.Command)
type RedconAcceptFunc func(conn redcon.Conn) bool
type RedconClosedFunc func(conn redcon.Conn, err error)

type RedconServeMux struct {
	handlers map[string]RedconHandlerFunc
	accept   func(conn redcon.Conn) bool
	closed   func(conn redcon.Conn, err error)
}

func NewRedconServeMux() *RedconServeMux {
	return &RedconServeMux{
		handlers: make(map[string]RedconHandlerFunc),
		accept: func(conn redcon.Conn) bool {
			return true
		},
		closed: func(conn redcon.Conn, err error) {},
	}
}

func (m *RedconServeMux) HandleFunc(command string, handler RedconHandlerFunc) {
	if handler == nil {
		panic("redcon: nil handler")
	}
	m.Handle(command, handler)
}

func (m *RedconServeMux) Handle(command string, handler RedconHandlerFunc) {
	if command == "" {
		panic("redcon: invalid command")
	}
	if handler == nil {
		panic("redcon: nil handler")
	}
	if _, exist := m.handlers[command]; exist {
		panic("redcon: multiple registrations for " + command)
	}

	m.handlers[command] = handler
}

func (m *RedconServeMux) do(conn redcon.Conn, cmd redcon.Command) {
	command := strings.ToLower(string(cmd.Args[0]))

	if handler, ok := m.handlers[command]; ok {
		handler(conn, cmd)
	} else {
		conn.WriteError("ERR unknown command '" + command + "'")
	}
}

func (m *RedconServeMux) Accept(f RedconAcceptFunc) {
	m.accept = f
}

func (m *RedconServeMux) Closed(f RedconClosedFunc) {
	m.closed = f
}

func (m *RedconServeMux) Run(address string) error {
	return redcon.ListenAndServe(
		address,
		m.do,
		m.accept,
		m.closed,
	)
}

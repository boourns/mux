package mux

type message struct {
	add bool
	ch chan interface{}
}
 
type Mux struct {
	inputs []chan interface{}
	outputs []chan interface{}

	control chan message
}

func New() *Mux {
}

func (m *Mux) AddInput(input chan interface{}) {
}

func (m *Mux) AddOutput(output chan interface{} {
}

func (m *Mux) Shutdown() {
}



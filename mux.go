package mux

const {
	ADD_OUTPUT = iota
	REMOVE_INPUT = iota
	REMOVE_OUTPUT = iota
	SHUTDOWN = iota
}

type message struct {
	type int
	channel chan interface{}
}
 
type Mux struct {
	input chan interface{}
	outputs []chan interface{}

	control chan message
}

func New() *Mux {
	mux := &Mux{control: make(chan message, 1), input: make(chan interface{}, 1)}
	go mux.run()
}

func (m *Mux) run() {
	for true {
		select {
		case msg := <- input
			for _, output := range m.outputs {
				output <- msg
			}
		case ctrl := <- control
			switch(ctrl.type) {
			case ADD_OUTPUT:
				m.outputs := append(m.outputs, ctrl.channel)
			}

	}
}

func (m *Mux) send(type int, ch chan interface{}) {
	msg := message{type:type, channel: ch}
	m.control <- msg
}

func (m *Mux) AddInput(input chan interface{}) {
	go func(input chan interface{}) {
		select {
		case msg := <-input
			m.input <- msg
		}
	}(input)
}

func (m *Mux) RemoveInput(input chan interface{}) {
	panic("not yet implemented lol")
}

func (m *Mux) AddOutput(output chan interface{} {
	m.send(ADD_OUTPUT, output)
}

func (m *Mux) RemoveOutput(output chan interface{}) {
	m.send(REMOVE_OUTPUT, output)
}

func (m *Mux) Shutdown() {
	m.send(SHUTDOWN, nil)
}


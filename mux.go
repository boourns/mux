package mux

import (
	"fmt"
)

const (
	ADD_OUTPUT    = iota
	ADD_INPUT     = iota
	REMOVE_INPUT  = iota
	REMOVE_OUTPUT = iota
	SHUTDOWN      = iota
)

type message struct {
	msgType int
	channel chan interface{}
}

type input struct {
	ch       chan interface{}
	shutdown chan bool
}

type Mux struct {
	input   chan interface{}
	outputs []chan interface{}
	inputs  []input
	control chan message
	filter  func(interface{}) []interface{}
}

func New(filter func(interface{}) []interface{}) *Mux {
	mux := &Mux{control: make(chan message, 1), input: make(chan interface{}, 1), filter: filter}
	go mux.run()
	return mux
}

func (m *Mux) run() {
	for true {
		select {
		case msg := <-m.input:
			if m.filter != nil {
				messages := m.filter(msg)
				for _, message := range messages {
					for _, output := range m.outputs {
						output <- message
					}
				}
			} else {
				for _, output := range m.outputs {
					output <- msg
				}
			}
		case ctrl := <-m.control:
			switch ctrl.msgType {
			case ADD_INPUT:
				in := input{ch: ctrl.channel, shutdown: make(chan bool, 1)}
				go func(in input) {
					for true {
						select {
						case msg := <-in.ch:
							m.input <- msg
						case <-in.shutdown:
							return
						}
					}
				}(in)
				m.inputs = append(m.inputs, in)
			case REMOVE_INPUT:
				for idx, in := range m.inputs {
					if in.ch == ctrl.channel {
						in.shutdown <- true
						m.inputs[idx] = m.inputs[len(m.inputs)-1]
						m.inputs = m.inputs[0 : len(m.inputs)-1]
						break
					}
				}
			case ADD_OUTPUT:
				m.outputs = append(m.outputs, ctrl.channel)
			case REMOVE_OUTPUT:
				for i, ch := range m.outputs {
					if ctrl.channel == ch {
						m.outputs[i] = m.outputs[len(m.outputs)-1]
						m.outputs = m.outputs[0 : len(m.outputs)-1]
						break
					}
				}
			case SHUTDOWN:
				for _, in := range m.inputs {
					in.shutdown <- true
				}
				return
			}
		}
	}
}

func (m *Mux) send(msgType int, ch chan interface{}) {
	msg := message{msgType: msgType, channel: ch}
	m.control <- msg
}

func (m *Mux) AddInput(input chan interface{}) {
	m.send(ADD_INPUT, input)
}

func (m *Mux) RemoveInput(input chan interface{}) {
	m.send(REMOVE_INPUT, input)
}

func (m *Mux) AddOutput(output chan interface{}) {
	m.send(ADD_OUTPUT, output)
}

func (m *Mux) RemoveOutput(output chan interface{}) {
	m.send(REMOVE_OUTPUT, output)
}

func (m *Mux) Shutdown() {
	m.send(SHUTDOWN, nil)
}

func String(in interface{}) string {
	return fmt.Sprintf("%v", in)
}

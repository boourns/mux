package mux

import (
	"testing"
)

func TestDoubleSendAndReceive(t *testing.T) {
	m := New(nil)
	out1 := make(chan interface{}, 1)
	out2 := make(chan interface{}, 1)

	in1 := make(chan interface{}, 1)
	in2 := make(chan interface{}, 1)

	m.AddInput(in1)
	m.AddInput(in2)
	m.AddOutput(out1)
	m.AddOutput(out2)

	in1 <- 1
	o := <-out1
	if String(o) != "1" {
		t.Fatalf("did not receive in1 output back on out1, got %v\n", o)
	}

	o = <-out2
	if String(o) != "1" {
		t.Fatalf("did not receive in1 output back on out2, got %v\n", o)
	}

	in2 <- 1
	o = <-out1
	if String(o) != "1" {
		t.Fatalf("did not receive in2 output back on out1, got %v\n", o)
	}

	o = <-out2
	if String(o) != "1" {
		t.Fatalf("did not receive in2 output back on out2, got %v\n", o)
	}
	m.Shutdown()
}

func TestRemovingOutput(t *testing.T) {
	m := New(nil)
	out1 := make(chan interface{}, 1)
	out2 := make(chan interface{}, 1)
	out3 := make(chan interface{}, 1)

	in1 := make(chan interface{}, 1)

	m.AddOutput(out1)
	m.AddOutput(out2)
	m.AddOutput(out3)
	m.RemoveOutput(out2)

	m.AddInput(in1)

	in1 <- 1
	o := <-out1
	if String(o) != "1" {
		t.Fatalf("Did not receive 1 back on in1, got %v\n", o)
	}

	o = <-out3
	if String(o) != "1" {
		t.Fatalf("Did not receive 1 back on in1, got %v\n", o)
	}
	m.Shutdown()
}

func TestRemovingInput(t *testing.T) {
	m := New(nil)
	out1 := make(chan interface{}, 1)

	in1 := make(chan interface{}, 1)
	in2 := make(chan interface{}, 1)
	in3 := make(chan interface{}, 1)

	m.AddOutput(out1)
	m.AddInput(in1)
	m.AddInput(in2)
	m.AddInput(in3)

	m.RemoveInput(in2)
	in1 <- 1
	o := <-out1
	if String(o) != "1" {
		t.Fatalf("did not receive expected from input, got %v\n", o)
	}
	m.Shutdown()
}

func TestFilter(t *testing.T) {
	m := New(func(msg interface{}) []interface{} {
		output := make([]interface{}, 0)
		if String(msg) == "ok" {
			output = append(output, "no way")
		}
	 	return output
        })

	out1 := make(chan interface{}, 1)
	in1 := make(chan interface{}, 1)
	m.AddInput(in1)
	m.AddOutput(out1)

	in1 <- "lol"
	in1 <- "ok"
	o := <- out1
	if (String(o) != "no way") {
		t.Fatalf("Filtering did not work, received %v\n", o)
	}
}

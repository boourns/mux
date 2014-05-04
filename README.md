mux
===

broadcast from many input channels to many output channels in Go

This package facilitates multicasting messages from one or more go channels (inputs) to one or more go channels (outputs).

You can optionally specify a 'filter' method which will process all incoming messages and return a slice of messages that should be sent to all listeners.

Each mux spawns a control goroutine that acts as the router/controller.  This serializes message passing and mux control into a single goroutine so mux modification is concurrency-safe.


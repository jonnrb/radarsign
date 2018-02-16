package main

// Probes should be linearized as to not affect each other. This just runs
// anonymous functions received from the returned channel in sequence until the
// returned channel is closed. This is effectively a mutex.
//
func Linearizer() chan func() {
	c := make(chan func())
	go func() {
		for action := range c {
			action()
		}
	}()
	return c
}

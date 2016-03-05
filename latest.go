// Package latest provides an untyped thread-safe single-value container.
package latest

import "sync"

// Latest acts like an untyped thread-safe variable. Get() returns the
// last thing given to Put(). It is safe to call both methods from
// multiple goroutines.
type Latest struct {
	in, out chan interface{}
	once    sync.Once
}

// Put replaces the current thing with the given thing.
func (l *Latest) Put(thing interface{}) {
	l.once.Do(l.start)
	l.in <- thing
}

// Get returns the last thing passed to Put. If Get is called first,
// it blocks until the first Put.
func (l *Latest) Get() interface{} {
	l.once.Do(l.start)
	return <-l.out
}

// Stop releases resources. Do not call Get or Put after calling Stop.
func (l *Latest) Stop() {
	l.once.Do(l.start)
	close(l.in)
}

func (l *Latest) start() {
	l.in = make(chan interface{})
	l.out = make(chan interface{})
	in := l.in
	go func() {
		defer close(l.out)
		thing := <-l.in
		for ok := true; ok; {
			select {
			case thing, ok = <-in:
			case l.out <- thing:
			}
		}
	}()
}

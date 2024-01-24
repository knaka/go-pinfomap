package main

type Bar struct{}

func (rcv *Bar) GreetRef(name string) string {
	return "hello" + name
}

func (rcv Bar) Greet() (int, error) {
	return 123, nil
}

type Baz struct {
	bar Bar
}

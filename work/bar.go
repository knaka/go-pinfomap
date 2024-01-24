package main

import "log"

type Bar struct {
	N int
}

func (rcv *Bar) GreetRef(name string) string {
	rcv.N++
	return "hello" + name
}

func (rcv Bar) Greet() (int, error) {
	rcv.N++
	return 123, nil
}

func bar() {
	bar := Bar{}
	log.Println(bar.N)
	bar.Greet()
	log.Println(bar.N)
	bar.GreetRef("foo")
	log.Println(bar.N)
	barRef := &Bar{}
	log.Println(barRef.N)
	barRef.Greet()
	log.Println(barRef.N)
	barRef.GreetRef("foo")
	log.Println(barRef.N)
}

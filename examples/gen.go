package examples

//go:generate_input foo.go gen_foo_accessor_go/*.go
//go:generate_output foo_accessor.go
//go:generate go run ./gen_foo_accessor_go/

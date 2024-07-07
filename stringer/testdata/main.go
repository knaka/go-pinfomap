package testdata

type Pill int

const (
	Placebo Pill = iota
	Aspirin
	Ibuprofen
	Paracetamol
	Acetaminophen = Paracetamol
)

type Fruit int8

const (
	Apple Fruit = iota
	Banana
	Cherry
)

package pinfomap

import (
	plu "github.com/gertd/go-pluralize"
	"sync"
)

var vPlu struct {
	client *plu.Client
	once   sync.Once
}

func Plural(s string) string {
	vPlu.once.Do(func() {
		vPlu.client = plu.NewClient()
	})
	return vPlu.client.Plural(s)
}

func Singular(s string) string {
	vPlu.once.Do(func() {
		vPlu.client = plu.NewClient()
	})
	return vPlu.client.Singular(s)
}

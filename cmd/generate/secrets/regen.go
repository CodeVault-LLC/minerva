package secrets

import (
	"github.com/lucasjones/reggen"
)

func NewSecret(regex string) string {
	g, err := reggen.NewGenerator(regex)
	if err != nil {
		panic(err)
	}
	return g.Generate(1)
}

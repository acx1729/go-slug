package slug_test

import (
	"fmt"
	"log"
	"math/rand/v2"
	"slices"

	"github.com/acx1729/go-slug"
)

func Example() {
	fmt.Println(slug.New()) // e.g. clever-falcon-89
}

func ExampleGenerator() {
	g := slug.Generator{Nouns: slug.Food(), Digits: 4, Separator: "_"}
	fmt.Println(g.New()) // e.g. sunny_mango_2048
}

func ExampleGenerator_seeded() {
	g := slug.Generator{
		Adjectives: []string{"clever"},
		Nouns:      []string{"falcon"},
		Digits:     2,
		Rand:       rand.New(rand.NewPCG(78, 0)),
	}
	fmt.Println(g.New())
	// Output: clever-falcon-89
}

func ExampleGenerator_allNouns() {
	g := slug.Generator{Nouns: slices.Concat(slug.Animals(), slug.Food(), slug.Nature(), slug.Space())}
	fmt.Println(g.New()) // e.g. brave-comet
}

func ExampleUnique() {
	slugs, err := slug.Unique(3)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(len(slugs))
	// Output: 3
}

func ExampleGenerator_Unique() {
	g := slug.Generator{Adjectives: []string{"happy", "sunny"}, Nouns: []string{"otter"}}
	slugs, err := g.Unique(2)
	if err != nil {
		log.Fatal(err)
	}
	slices.Sort(slugs)
	fmt.Println(slugs)
	// Output: [happy-otter sunny-otter]
}

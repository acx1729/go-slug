package slug

import (
	_ "embed"
	"slices"
	"strings"
)

// The built-in words live in words/*.txt, one word per line, and are
// compiled into the package. To add a list, add its file, embed it below,
// and give it a function like the ones at the end of this file.

//go:embed words/adjectives.txt
var adjectivesText string

//go:embed words/animals.txt
var animalsText string

//go:embed words/food.txt
var foodText string

//go:embed words/nature.txt
var natureText string

//go:embed words/space.txt
var spaceText string

var (
	adjectives = strings.Fields(adjectivesText)
	animals    = strings.Fields(animalsText)
	food       = strings.Fields(foodText)
	nature     = strings.Fields(natureText)
	space      = strings.Fields(spaceText)
)

// Adjectives returns the built-in adjectives, such as "clever" and "sunny".
// Every call returns a new slice.
func Adjectives() []string { return slices.Clone(adjectives) }

// Animals returns the built-in animals, such as "falcon" and "otter". They
// are the default nouns. Every call returns a new slice.
func Animals() []string { return slices.Clone(animals) }

// Food returns the built-in foods, such as "mango" and "waffle". Every
// call returns a new slice.
func Food() []string { return slices.Clone(food) }

// Nature returns the built-in nature words, such as "meadow" and "canyon".
// Every call returns a new slice.
func Nature() []string { return slices.Clone(nature) }

// Space returns the built-in space words, such as "comet" and "nebula".
// Every call returns a new slice.
func Space() []string { return slices.Clone(space) }

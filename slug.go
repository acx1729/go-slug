// Package slug makes short, readable, random names such as
// "clever-falcon-89".
//
//	s := slug.New() // "clever-falcon-89"
//
// A [Generator] changes the words, the digits, or the separator:
//
//	g := slug.Generator{Nouns: slug.Food(), Digits: 4, Separator: "_"}
//	s := g.New() // "sunny_mango_2048"
//
// Slugs are random, not secret. Do not use them as passwords or tokens.
package slug

import (
	"fmt"
	"math"
	"math/rand/v2"
	"strings"
)

// std makes the slugs returned by New and Unique.
var std = Generator{Digits: 2}

// New returns a random slug made of an adjective, an animal, and two
// digits, such as "clever-falcon-89". It is safe for concurrent use.
func New() string { return std.New() }

// Unique returns n different slugs made the way New makes them. It is safe
// for concurrent use. See [Generator.Unique] for when it fails.
func Unique(n int) ([]string, error) { return std.Unique(n) }

// A Generator makes slugs from an adjective, a noun, and optional digits.
// Every field is optional: the zero value makes slugs like "clever-falcon".
type Generator struct {
	// Adjectives and Nouns are the words slugs are made from. When empty,
	// they default to [Adjectives] and [Animals].
	Adjectives []string
	Nouns      []string

	// Digits is how many random digits end each slug: 2 ends it like
	// "-89", 4 like "-0423". Zero adds none. A negative value panics.
	Digits int

	// Separator joins the parts. When empty, it defaults to "-".
	Separator string

	// Rand is the source of randomness; set it for repeatable output. When
	// nil, the Generator uses the global math/rand/v2 source and is safe
	// for concurrent use. A *rand.Rand is not, so neither is a Generator
	// that holds one.
	Rand *rand.Rand
}

// New returns a random slug.
func (g Generator) New() string {
	adjs, nouns, sep := g.resolve()

	var b strings.Builder
	b.WriteString(adjs[g.intN(len(adjs))])
	b.WriteString(sep)
	b.WriteString(nouns[g.intN(len(nouns))])
	if g.Digits > 0 {
		b.WriteString(sep)
		for range g.Digits {
			b.WriteByte('0' + byte(g.intN(10)))
		}
	}
	return b.String()
}

// Unique returns n slugs that all differ. It fails if n is negative, if n
// is more than the Generator can make, or if the Separator shares a
// character with the words or digits, because two different slugs could
// then read the same. It slows down as n nears the number of possible
// slugs.
func (g Generator) Unique(n int) ([]string, error) {
	adjs, nouns, sep := g.resolve()
	if n < 0 {
		return nil, fmt.Errorf("slug: negative count %d", n)
	}
	if containsAny(adjs, sep) || containsAny(nouns, sep) || (g.Digits > 0 && strings.ContainsAny(sep, "0123456789")) {
		return nil, fmt.Errorf("slug: separator %q also appears in the words or digits", sep)
	}
	if possible := g.count(adjs, nouns); n > possible {
		return nil, fmt.Errorf("slug: cannot make %d unique slugs, only %d are possible", n, possible)
	}

	seen := make(map[string]bool, n)
	out := make([]string, 0, n)
	for len(out) < n {
		if s := g.New(); !seen[s] {
			seen[s] = true
			out = append(out, s)
		}
	}
	return out, nil
}

// resolve returns the words and separator g uses, applying the defaults.
func (g Generator) resolve() (adjs, nouns []string, sep string) {
	if g.Digits < 0 {
		panic(fmt.Sprintf("slug: negative Digits %d", g.Digits))
	}
	adjs, nouns, sep = g.Adjectives, g.Nouns, g.Separator
	if len(adjs) == 0 {
		adjs = adjectives
	}
	if len(nouns) == 0 {
		nouns = animals
	}
	if sep == "" {
		sep = "-"
	}
	return adjs, nouns, sep
}

func (g Generator) intN(n int) int {
	if g.Rand != nil {
		return g.Rand.IntN(n)
	}
	return rand.IntN(n)
}

// count returns how many different slugs g can make from adjs and nouns,
// capped at math.MaxInt.
func (g Generator) count(adjs, nouns []string) int {
	n := mul(distinct(adjs), distinct(nouns))
	for i := 0; i < g.Digits && n < math.MaxInt; i++ {
		n = mul(n, 10)
	}
	return n
}

// mul returns a*b for non-negative a and b, capped at math.MaxInt.
func mul(a, b int) int {
	if b != 0 && a > math.MaxInt/b {
		return math.MaxInt
	}
	return a * b
}

func distinct(words []string) int {
	set := make(map[string]struct{}, len(words))
	for _, w := range words {
		set[w] = struct{}{}
	}
	return len(set)
}

func containsAny(words []string, chars string) bool {
	for _, w := range words {
		if strings.ContainsAny(w, chars) {
			return true
		}
	}
	return false
}

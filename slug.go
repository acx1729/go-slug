// Package slug makes short, readable, random names such as
// "clever-falcon-89".
//
//	s := slug.New() // "clever-falcon-89"
//
// A [Generator] changes the words, the ending, or the separator:
//
//	g := slug.Generator{Nouns: slug.Food(), Suffix: "?#?#?", Separator: "_"}
//	s := g.New() // "sunny_mango_k7m2q"
//
// Slugs are random, not secret. Do not use them as passwords or tokens.
package slug

import (
	"fmt"
	"math"
	"math/rand/v2"
	"strings"

	"github.com/acx1729/go-slug/internal/words"
)

// std makes the slugs returned by New and Unique.
var std = Generator{Suffix: "##"}

// New returns a random slug made of an adjective, an animal, and two
// digits, such as "clever-falcon-89". It is safe for concurrent use.
func New() string { return std.New() }

// Unique returns n different slugs made the way New makes them. It is safe
// for concurrent use. See [Generator.Unique] for when it fails.
func Unique(n int) ([]string, error) { return std.Unique(n) }

// A Generator makes slugs from an adjective, a noun, and an optional
// random ending. Every field is optional: the zero value makes slugs like
// "clever-falcon".
type Generator struct {
	// Adjectives and Nouns are the words slugs are made from. When empty,
	// they default to [Adjectives] and [Animals].
	Adjectives []string
	Nouns      []string

	// Suffix is a template for the random ending, one symbol per
	// character: # is a digit 0-9 and ? is a letter a-z, the convention
	// of gofakeit's Numerify and Lexify. "##" ends a slug like "-89" and
	// "?#?#?" like "-k7m2q". Empty means no ending. Any other character
	// is an error: New panics and Unique returns it.
	Suffix string

	// Separator joins the parts. When empty, it defaults to "-".
	Separator string

	// Rand is the source of randomness; set it for repeatable output. When
	// nil, the Generator uses the global math/rand/v2 source and is safe
	// for concurrent use. A *rand.Rand is not, so neither is a Generator
	// that holds one.
	Rand *rand.Rand
}

// New returns a random slug. It panics if Suffix is not a valid template.
func (g Generator) New() string {
	adjs, nouns, sep, err := g.resolve()
	if err != nil {
		panic(err)
	}

	var b strings.Builder
	b.WriteString(adjs[g.intN(len(adjs))])
	b.WriteString(sep)
	b.WriteString(nouns[g.intN(len(nouns))])
	if g.Suffix != "" {
		b.WriteString(sep)
		for i := 0; i < len(g.Suffix); i++ {
			if g.Suffix[i] == '#' {
				b.WriteByte('0' + byte(g.intN(10)))
			} else {
				b.WriteByte('a' + byte(g.intN(26)))
			}
		}
	}
	return b.String()
}

// Unique returns n slugs that all differ. It fails if Suffix is not a valid
// template, if n is negative, if n is more than the Generator can make, or
// if the Separator shares a character with the words or the ending, because
// two different slugs could then read the same. It slows down as n nears
// the number of possible slugs.
func (g Generator) Unique(n int) ([]string, error) {
	adjs, nouns, sep, err := g.resolve()
	if err != nil {
		return nil, err
	}
	if n < 0 {
		return nil, fmt.Errorf("slug: negative count %d", n)
	}
	if g.clashes(sep, adjs, nouns) {
		return nil, fmt.Errorf("slug: separator %q also appears in the words or the ending", sep)
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

// resolve checks Suffix and returns the words and separator g uses,
// applying the defaults.
func (g Generator) resolve() (adjs, nouns []string, sep string, err error) {
	for i, r := range g.Suffix {
		if r != '#' && r != '?' {
			return nil, nil, "", fmt.Errorf("slug: Suffix %q: %q at byte %d is not # (a digit) or ? (a letter)", g.Suffix, r, i)
		}
	}
	adjs, nouns, sep = g.Adjectives, g.Nouns, g.Separator
	if len(adjs) == 0 {
		adjs = words.Adjectives
	}
	if len(nouns) == 0 {
		nouns = words.Animals
	}
	if sep == "" {
		sep = "-"
	}
	return adjs, nouns, sep, nil
}

func (g Generator) intN(n int) int {
	if g.Rand != nil {
		return g.Rand.IntN(n)
	}
	return rand.IntN(n)
}

// clashes reports whether sep shares a character with the words or with
// what the ending can contain.
func (g Generator) clashes(sep string, adjs, nouns []string) bool {
	if strings.Contains(g.Suffix, "#") && strings.ContainsAny(sep, "0123456789") {
		return true
	}
	if strings.Contains(g.Suffix, "?") && strings.ContainsAny(sep, "abcdefghijklmnopqrstuvwxyz") {
		return true
	}
	return containsAny(adjs, sep) || containsAny(nouns, sep)
}

// count returns how many different slugs g can make from adjs and nouns,
// capped at math.MaxInt. Suffix holds only ASCII symbols once resolve has
// accepted it.
func (g Generator) count(adjs, nouns []string) int {
	n := mul(distinct(adjs), distinct(nouns))
	for i := 0; i < len(g.Suffix) && n < math.MaxInt; i++ {
		if g.Suffix[i] == '#' {
			n = mul(n, 10)
		} else {
			n = mul(n, 26)
		}
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

package slug_test

import (
	"math/rand/v2"
	"regexp"
	"slices"
	"strings"
	"sync"
	"testing"

	"github.com/acx1729/go-slug"
)

func seeded(seed uint64) *rand.Rand { return rand.New(rand.NewPCG(seed, seed)) }

func TestNew(t *testing.T) {
	shape := regexp.MustCompile(`^[a-z]+-[a-z]+-[0-9]{2}$`)
	adjectives, animals := slug.Adjectives(), slug.Animals()
	for range 1000 {
		s := slug.New()
		if !shape.MatchString(s) {
			t.Fatalf("New() = %q, want adjective-animal-NN", s)
		}
		parts := strings.Split(s, "-")
		if !slices.Contains(adjectives, parts[0]) || !slices.Contains(animals, parts[1]) {
			t.Fatalf("New() = %q uses a word outside the built-in lists", s)
		}
	}
}

// TestCleverFalcon89 proves New can return "clever-falcon-89": both words
// are built in, and two digits can read "89".
func TestCleverFalcon89(t *testing.T) {
	if !slices.Contains(slug.Adjectives(), "clever") {
		t.Fatal(`"clever" is not a built-in adjective`)
	}
	if !slices.Contains(slug.Animals(), "falcon") {
		t.Fatal(`"falcon" is not a built-in animal`)
	}
	g := slug.Generator{Adjectives: []string{"clever"}, Nouns: []string{"falcon"}, Digits: 2, Rand: seeded(1)}
	for range 10_000 {
		if g.New() == "clever-falcon-89" {
			return
		}
	}
	t.Fatal("never made clever-falcon-89")
}

func TestGenerator(t *testing.T) {
	for _, tc := range []struct {
		name  string
		g     slug.Generator
		shape string
		nouns []string
	}{
		{"zero value", slug.Generator{}, `^[a-z]+-[a-z]+$`, slug.Animals()},
		{"food, four digits, underscore", slug.Generator{Nouns: slug.Food(), Digits: 4, Separator: "_"}, `^[a-z]+_[a-z]+_[0-9]{4}$`, slug.Food()},
		{"space with dots", slug.Generator{Nouns: slug.Space(), Separator: "."}, `^[a-z]+\.[a-z]+$`, slug.Space()},
		{"own words", slug.Generator{Adjectives: []string{"happy"}, Nouns: []string{"bella", "tommy"}, Digits: 1}, `^happy-(bella|tommy)-[0-9]$`, []string{"bella", "tommy"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			shape := regexp.MustCompile(tc.shape)
			sep := tc.g.Separator
			if sep == "" {
				sep = "-"
			}
			for range 500 {
				s := tc.g.New()
				if !shape.MatchString(s) {
					t.Fatalf("New() = %q, want %s", s, tc.shape)
				}
				if noun := strings.Split(s, sep)[1]; !slices.Contains(tc.nouns, noun) {
					t.Fatalf("New() = %q: noun %q is not in the list", s, noun)
				}
			}
		})
	}
}

func TestNegativeDigitsPanics(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("Digits: -1 did not panic")
		}
	}()
	slug.Generator{Digits: -1}.New()
}

func TestRandIsRepeatable(t *testing.T) {
	a := slug.Generator{Digits: 4, Rand: seeded(42)}
	b := slug.Generator{Digits: 4, Rand: seeded(42)}
	for range 100 {
		if x, y := a.New(), b.New(); x != y {
			t.Fatalf("same seed gave %q and %q", x, y)
		}
	}
}

func TestUnique(t *testing.T) {
	g := slug.Generator{Adjectives: []string{"red", "blue"}, Nouns: []string{"cat", "dog", "fish"}}

	got, err := g.Unique(6)
	if err != nil {
		t.Fatal(err)
	}
	slices.Sort(got)
	want := []string{"blue-cat", "blue-dog", "blue-fish", "red-cat", "red-dog", "red-fish"}
	if !slices.Equal(got, want) {
		t.Fatalf("Unique(6) = %v, want every combination %v", got, want)
	}

	if _, err := g.Unique(7); err == nil {
		t.Fatal("Unique(7) from 6 possible slugs did not fail")
	}
	if _, err := g.Unique(-1); err == nil {
		t.Fatal("Unique(-1) did not fail")
	}
	if got, err := g.Unique(0); err != nil || len(got) != 0 {
		t.Fatalf("Unique(0) = %v, %v; want no slugs and no error", got, err)
	}

	repeats := slug.Generator{Adjectives: []string{"red", "red"}, Nouns: []string{"cat"}}
	if _, err := repeats.Unique(2); err == nil {
		t.Fatal("repeated words were counted twice")
	}
}

func TestUniqueSeparatorClash(t *testing.T) {
	for _, tc := range []struct {
		name string
		g    slug.Generator
	}{
		{"in an adjective", slug.Generator{Adjectives: []string{"a-b"}, Nouns: []string{"cat"}}},
		{"in a noun", slug.Generator{Nouns: []string{"x_y"}, Separator: "_"}},
		{"a letter of the words", slug.Generator{Separator: "e"}},
		{"a digit", slug.Generator{Digits: 2, Separator: "0"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := tc.g.Unique(1); err == nil {
				t.Fatal("Unique did not fail")
			}
		})
	}
	if _, err := (slug.Generator{Separator: "7"}).Unique(1); err != nil {
		t.Fatalf("a digit separator without digits failed: %v", err)
	}
}

func TestUniqueDefault(t *testing.T) {
	got, err := slug.Unique(5000)
	if err != nil {
		t.Fatal(err)
	}
	seen := map[string]bool{}
	for _, s := range got {
		if seen[s] {
			t.Fatalf("Unique repeated %q", s)
		}
		seen[s] = true
	}
	if len(seen) != 5000 {
		t.Fatalf("Unique(5000) returned %d slugs", len(seen))
	}
}

// TestConcurrentUse is meaningful under -race.
func TestConcurrentUse(t *testing.T) {
	g := slug.Generator{Nouns: slug.Nature(), Digits: 3}
	var wg sync.WaitGroup
	for range 8 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range 1000 {
				_ = slug.New()
				_ = g.New()
			}
			if _, err := slug.Unique(100); err != nil {
				t.Error(err)
			}
		}()
	}
	wg.Wait()
}

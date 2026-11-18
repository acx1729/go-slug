package slug_test

import (
	"regexp"
	"testing"

	"github.com/acx1729/go-slug"
)

// lists names every built-in list. Add a new list here too.
var lists = map[string]func() []string{
	"Adjectives":       slug.Adjectives,
	"Animals":          slug.Animals,
	"Architecture":     slug.Architecture,
	"Art":              slug.Art,
	"Astronomy":        slug.Astronomy,
	"Biology":          slug.Biology,
	"Chemistry":        slug.Chemistry,
	"ComputerScience":  slug.ComputerScience,
	"Countries":        slug.Countries,
	"Economy":          slug.Economy,
	"Food":             slug.Food,
	"General":          slug.General,
	"Geography":        slug.Geography,
	"History":          slug.History,
	"Literature":       slug.Literature,
	"Math":             slug.Math,
	"Medicine":         slug.Medicine,
	"Microbiology":     slug.Microbiology,
	"MolecularBiology": slug.MolecularBiology,
	"Music":            slug.Music,
	"Nature":           slug.Nature,
	"Physics":          slug.Physics,
	"Plants":           slug.Plants,
	"Scientists":       slug.Scientists,
	"Space":            slug.Space,
	"Sports":           slug.Sports,
	"Technology":       slug.Technology,
}

var wordRule = regexp.MustCompile(`^[a-z0-9]+$`)

// TestWordLists keeps every list safe to put in a slug: lower-case letters
// a-z and digits only, sorted, with no repeats.
func TestWordLists(t *testing.T) {
	for name, list := range lists {
		words := list()
		if len(words) == 0 {
			t.Errorf("%s is empty", name)
		}
		for i, w := range words {
			switch {
			case !wordRule.MatchString(w):
				t.Errorf("%s[%d] %q: use lower-case letters a-z and digits only", name, i, w)
			case i > 0 && w <= words[i-1]:
				t.Errorf("%s[%d] %q must come after %q: keep the list sorted, with no repeats", name, i, w, words[i-1])
			}
		}
	}
}

func TestWordListsAreCopies(t *testing.T) {
	a := slug.Animals()
	a[0] = "changed"
	if slug.Animals()[0] == "changed" {
		t.Fatal("changing a returned list changed the built-in list")
	}
}

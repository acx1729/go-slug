package slug_test

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"testing"

	"github.com/acx1729/go-slug"
)

// lists maps every file in words/ to the function that returns it.
var lists = map[string]func() []string{
	"adjectives.txt": slug.Adjectives,
	"animals.txt":    slug.Animals,
	"food.txt":       slug.Food,
	"nature.txt":     slug.Nature,
	"space.txt":      slug.Space,
}

var wordRule = regexp.MustCompile(`^[a-z]{3,10}$`)

// TestWordFiles keeps every word file easy to read and safe to use: one
// word per line, 3 to 10 lower-case letters a-z, sorted, and each word in
// only one file.
func TestWordFiles(t *testing.T) {
	entries, err := os.ReadDir("words")
	if err != nil {
		t.Fatal(err)
	}
	owner := map[string]string{} // word -> file it first appeared in
	for _, e := range entries {
		path := filepath.Join("words", e.Name())
		load, ok := lists[e.Name()]
		if !ok {
			t.Errorf("%s is not loaded: embed it in words.go and add it to lists in this test", path)
			continue
		}
		b, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		text := string(b)
		if !strings.HasSuffix(text, "\n") {
			t.Errorf("%s: must end with a newline", path)
		}
		lines := strings.Split(strings.TrimSuffix(text, "\n"), "\n")
		for i, w := range lines {
			at := fmt.Sprintf("%s:%d", path, i+1)
			switch {
			case !wordRule.MatchString(w):
				t.Errorf("%s: %q must be 3 to 10 lower-case letters a-z", at, w)
			case i > 0 && w <= lines[i-1]:
				t.Errorf("%s: %q must come after %q (keep the file sorted, no repeats)", at, w, lines[i-1])
			}
			if first, dup := owner[w]; dup {
				t.Errorf("%s: %q is already in %s", at, w, first)
			} else {
				owner[w] = path
			}
		}
		if got := load(); !slices.Equal(got, lines) {
			t.Errorf("%s: the loaded list does not match the file", path)
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

# go-slug

Short, readable, random names for Go.

```
clever-falcon-89
```

Use them for preview environments, job and branch names, temporary
directories, or any label a person has to read and say out loud.

```bash
go get github.com/acx1729/go-slug
```

Go 1.22 or later. No dependencies.

## Usage

```go
import "github.com/acx1729/go-slug"

s := slug.New() // "clever-falcon-89"
```

That's the whole default: an adjective, an animal, and two digits. `slug.New`
is safe to call from many goroutines.

### Change the shape

A `Generator` is a plain struct. Every field is optional:

```go
g := slug.Generator{
	Nouns:  slug.Literature(),
	Suffix: "?#?#?",
}
g.New() // "brave-odyssey-k7m2q"
```

| Field | Empty means | Example |
| --- | --- | --- |
| `Adjectives` | `slug.Adjectives()` | `[]string{"happy", "sleepy"}` |
| `Nouns` | `slug.Animals()` | `slug.Space()` |
| `Suffix` | no ending | `"?#?#?"` → `-k7m2q` |
| `Separator` | `"-"` | `"_"` |
| `Rand` | the global random source | `rand.New(rand.NewPCG(1, 2))` for repeatable output |

The zero `Generator{}` makes slugs like `clever-falcon`.

### The ending

`Suffix` is a template with one symbol per character:

| Symbol | Becomes |
| --- | --- |
| `#` | a random digit, 0–9 |
| `?` | a random letter, a–z |

```go
slug.Generator{Suffix: "##"}    // brave-falcon-89   (what slug.New uses)
slug.Generator{Suffix: "?#?#?"} // brave-falcon-k7m2q
slug.Generator{Suffix: "??###"} // brave-falcon-km482
```

`#` and `?` mean the same as in gofakeit's `Numerify` and `Lexify`. Any other
character is an error: `New` panics, the way `regexp.MustCompile` does, and
`Unique` returns the error. Put fixed text outside the template, for example
`"run-" + g.New()`.

### No repeats

```go
slugs, err := slug.Unique(100) // 100 different slugs
```

`Unique` fails instead of returning a repeat. That happens when you ask for
more slugs than the generator can make, or when the separator also appears
inside the words or the ending, because two different slugs could then read
the same.

### Word lists

`slug.Adjectives()` holds the adjectives. Every other function returns nouns,
and `slug.Animals()` is the default:

| Theme | Functions |
| --- | --- |
| Everyday | `General`, `Animals`, `Food`, `Plants`, `Nature` |
| Science | `Astronomy`, `Space`, `Biology`, `Microbiology`, `MolecularBiology`, `Chemistry`, `Physics`, `Math`, `Medicine` |
| Technology | `ComputerScience`, `Technology` |
| Culture | `Art`, `Architecture`, `Literature`, `Music`, `History`, `Sports` |
| World | `Countries`, `Geography`, `Economy` |
| People | `Scientists` |

Each call returns a fresh copy. Mix the lists freely, or use your own words:

```go
science := slices.Concat(slug.Astronomy(), slug.Chemistry(), slug.Physics())
g := slug.Generator{Nouns: science, Suffix: "##"} // "quiet-quasar-07"
```

With the lists as shipped, `slug.New()` can make about 106 million different
slugs (adjectives × animals × 100). Each `#` in the ending multiplies the
count by 10, and each `?` by 26.

## Updating the words

Each list is a Go slice in its own file under
[internal/words](internal/words), for example `animals.go`:

```go
var Animals = []string{
	"aardvark",
	"aardwolf",
	…
}
```

To add or remove a word, edit its file. Then run `go test ./...`, which names
the list, position and word that break the rules:

- lower-case letters a–z and digits only
- sorted, with no repeats

To add a list, add a file in `internal/words`, a function in
[words.go](words.go), and an entry in `lists` in
[words_test.go](words_test.go).

Slugs are random, not secret. Don't use them as passwords or tokens.

## Credits

go-slug is a Go rewrite of
[unique-namer](https://github.com/aziele/unique-namer) by Andrzej Zielezinski,
a Python package that makes names like `talented-toucan`.

- The Go code is new.
- The word lists are translated from unique-namer's
  [data](https://github.com/aziele/unique-namer/tree/main/src/namer/data) and
  extended with new words, including the `Nature` and `Space` lists.
- `c++` is written `cpp` and `nadp+` is written `nadp`, so every word is safe
  in a URL.

unique-namer is MIT-licensed. Its copyright and license notice is in
[NOTICE](NOTICE).

## License

[MIT](LICENSE)

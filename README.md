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
	Nouns:     slug.Food(),
	Digits:    4,
	Separator: "_",
}
g.New() // "sunny_mango_2048"
```

| Field | Empty means | Example |
| --- | --- | --- |
| `Adjectives` | `slug.Adjectives()` | `[]string{"happy", "sleepy"}` |
| `Nouns` | `slug.Animals()` | `slug.Space()` |
| `Digits` | no digits | `4` → `-0423` |
| `Separator` | `"-"` | `"_"` |
| `Rand` | the global random source | `rand.New(rand.NewPCG(1, 2))` for repeatable output |

The zero `Generator{}` makes slugs like `clever-falcon`.

### No repeats

```go
slugs, err := slug.Unique(100) // 100 different slugs
```

`Unique` fails instead of returning a repeat. That happens when you ask for
more slugs than the generator can make, or when the separator also appears
inside the words, because two different slugs could then read the same.

### Word lists

| Function | Words |
| --- | --- |
| `slug.Adjectives()` | clever, sunny, brave, … |
| `slug.Animals()` | falcon, otter, panda, … |
| `slug.Food()` | mango, waffle, pretzel, … |
| `slug.Nature()` | meadow, canyon, willow, … |
| `slug.Space()` | comet, nebula, orbit, … |

Mix them freely, or use your own:

```go
all := slices.Concat(slug.Animals(), slug.Food(), slug.Nature(), slug.Space())
g := slug.Generator{Nouns: all, Digits: 2} // "brave-comet-07"
```

With the lists as shipped, `slug.New()` can make about 6.7 million different
slugs (adjectives × animals × 100). Each extra digit multiplies that by ten.

## Updating the words

The lists are plain text files in [words/](words), one word per line,
compiled into the package:

```
words/
  adjectives.txt
  animals.txt
  food.txt
  nature.txt
  space.txt
```

To add or remove a word, edit its file. Then run `go test ./...`, which names
the file and line of any word that breaks the rules:

- 3 to 10 lower-case letters, a–z only
- sorted, with no repeats
- each word in only one file

A new list is a new file plus a few lines in [words.go](words.go).

Slugs are random, not secret. Don't use them as passwords or tokens.

## Credits

go-slug is a Go rewrite of
[unique-namer](https://github.com/aziele/unique-namer) by Andrzej Zielezinski,
a Python package that makes names like `talented-toucan`. The idea of
adjective-noun names with a random suffix comes from there. The code and the
word lists in this repository are new.

## License

[MIT](LICENSE)

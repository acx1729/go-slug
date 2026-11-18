package slug

import (
	"slices"

	"github.com/acx1729/go-slug/internal/words"
)

// Each function returns a new copy of a built-in list, so callers may
// change it freely. The lists live in internal/words.

// Adjectives returns the adjectives, such as "clever" and "sunny". They are
// the default adjectives.
func Adjectives() []string { return slices.Clone(words.Adjectives) }

// Animals returns the animals, such as "falcon" and "otter". They are the
// default nouns.
func Animals() []string { return slices.Clone(words.Animals) }

// Architecture returns architecture nouns.
func Architecture() []string { return slices.Clone(words.Architecture) }

// Art returns art nouns, including artists.
func Art() []string { return slices.Clone(words.Art) }

// Astronomy returns astronomy nouns.
func Astronomy() []string { return slices.Clone(words.Astronomy) }

// Biology returns general biology nouns. [Microbiology] and
// [MolecularBiology] hold the more specialised terms.
func Biology() []string { return slices.Clone(words.Biology) }

// Chemistry returns chemistry nouns.
func Chemistry() []string { return slices.Clone(words.Chemistry) }

// ComputerScience returns computer science nouns.
func ComputerScience() []string { return slices.Clone(words.ComputerScience) }

// Countries returns country names.
func Countries() []string { return slices.Clone(words.Countries) }

// Economy returns economy nouns.
func Economy() []string { return slices.Clone(words.Economy) }

// Food returns foods and drinks, such as "mango" and "waffle".
func Food() []string { return slices.Clone(words.Food) }

// General returns everyday nouns, the largest list.
func General() []string { return slices.Clone(words.General) }

// Geography returns geography nouns.
func Geography() []string { return slices.Clone(words.Geography) }

// History returns history nouns, including historical figures.
func History() []string { return slices.Clone(words.History) }

// Literature returns literature nouns, including authors.
func Literature() []string { return slices.Clone(words.Literature) }

// Math returns mathematics nouns, including mathematicians.
func Math() []string { return slices.Clone(words.Math) }

// Medicine returns medicine nouns.
func Medicine() []string { return slices.Clone(words.Medicine) }

// Microbiology returns microbiology nouns.
func Microbiology() []string { return slices.Clone(words.Microbiology) }

// MolecularBiology returns molecular biology nouns.
func MolecularBiology() []string { return slices.Clone(words.MolecularBiology) }

// Music returns music nouns.
func Music() []string { return slices.Clone(words.Music) }

// Nature returns nature nouns, such as "meadow" and "canyon".
func Nature() []string { return slices.Clone(words.Nature) }

// Physics returns physics nouns.
func Physics() []string { return slices.Clone(words.Physics) }

// Plants returns plant nouns.
func Plants() []string { return slices.Clone(words.Plants) }

// Scientists returns scientists' names.
func Scientists() []string { return slices.Clone(words.Scientists) }

// Space returns space nouns, such as "comet" and "nebula".
func Space() []string { return slices.Clone(words.Space) }

// Sports returns sports nouns.
func Sports() []string { return slices.Clone(words.Sports) }

// Technology returns technology nouns.
func Technology() []string { return slices.Clone(words.Technology) }

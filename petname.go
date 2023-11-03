// Package petname is a library for generating human-readable, random
// names for objects (e.g. hostnames, containers, blobs).
package petname

import (
	"math/rand"
	"time"

	"github.com/Bios-Marcel/go-petname/short"
)

var names, adjectives, adverbs []string

func init() {
	SetNames(short.Names)
	SetAdjectives(short.Adjectives)
	SetAdverbs(short.Adverbs)
}

// SetNames specifies which names to use when generating petnames. Note that
// these are expected to be all lowercase characters and not containing any
// trailing or leading whitespace.
func SetNames(n []string) {
	names = n
}

// SetAdjectives specifies which adjectives to use when generating petnames.
// Note that these are expected to be all lowercase characters and not
// containing any trailing or leading whitespace.
func SetAdjectives(a []string) {
	adjectives = a
}

// SetAdverbs specifies which adverbs to use when generating petnames. Note
// that these are expected to be all lowercase characters and not containing
// any trailing or leading whitespace.
func SetAdverbs(a []string) {
	adverbs = a
}

var localRand *rand.Rand = rand.New(rand.NewSource(1))

// NonDeterministicMode configures the local random generator used internally
// to provide non deterministic results, instead of a pre-defined order that is
// reproducible even after a process restart. If you wish to specify a custom
// contant, call [Seed(int64)].
func NonDeterministicMode() {
	localRand.Seed(time.Now().UnixNano())
}

// Seed configures the local random generator, allowing you to specify a
// constant for reproducible "randomness" or provide a custom value for
// "true" randomness.
func Seed(seed int64) {
	localRand.Seed(seed)
}

// Adverb returns a random adverb from a list of petname adverbs.
func Adverb() string {
	return adverbs[localRand.Intn(len(adverbs))]
}

// Adjective returns a random adjective from a list of petname adjectives.
func Adjective() string {
	return adjectives[localRand.Intn(len(adjectives))]
}

// Name returns a random name from a list of petname names.
func Name() string {
	return names[localRand.Intn(len(names))]
}

// Casing specifies the casing of the overall petname, not its separate words.
type Casing int

const (
	// Lower will keep all letters as they are. Since the requirement for the
	// word lists is that they are all lowercase, this will result in a NO-OP.
	Lower Casing = iota
	// Upper will uppercase all letters.
	Upper
	// Title will uppercase the first letter and keep the rest lowercased..
	Title
)

// Separator specifies the separator between words in a petname.
type Separator = byte

const (
	None       Separator = 0
	Hyphen     Separator = '-'
	Underscore Separator = '_'
)

func asciiByteToUpper(b byte) byte {
	if b >= 'a' && b <= 'z' {
		return b - 32
	}

	return b
}

// Generate will create a petname using the given configuration. Casing and
// word separation are different, allowing things such as `Word-Word-Word` or
// `WORD_WORD_WORD`.
func Generate(wordCount uint, casing Casing, separator Separator) string {
	var words []string
	switch wordCount {
	case 0:
		// Without this case, passing 0 as the wordCount will cause a very
		// slow call. This is because we are using a unit, which we'll subtract
		// 2 from, causing MaxUint - 2 iterations.
		return ""
	case 1:
		words = []string{Name()}
	case 2:
		words = []string{Adjective(), Name()}
	case 3:
		// Potentially common cases have shortcut implementations to
		// reduce allocations and CPU usage, even though default: would handle
		// them correctly.
		words = []string{Adverb(), Adjective(), Name()}
	case 4:
		words = []string{Adverb(), Adverb(), Adjective(), Name()}
	default:
		words = make([]string, 0, wordCount)
		for i := uint(0); i < wordCount-2; i++ {
			words = append(words, Adverb())
		}

		words = append(words, Adjective(), Name())
	}

	var byteLen int
	for _, word := range words {
		byteLen += len(word)
	}

	if separator != None {
		byteLen += len(words) - 1
	}

	var buffer []byte
	if byteLen <= 64 {
		buffer = make([]byte, 0, 64)
	} else {
		buffer = make([]byte, 0, byteLen)
	}

	for _, word := range words {
		if separator != None && len(buffer) > 0 {
			buffer = append(buffer, separator)
		}

		switch casing {
		case Upper:
			for i := 0; i < len(word); i++ {
				buffer = append(buffer, asciiByteToUpper(word[i]))
			}
		case Title:
			buffer = append(buffer, asciiByteToUpper(word[0]))
			buffer = append(buffer, word[1:]...)
		case Lower:
			fallthrough
		default:
			buffer = append(buffer, word...)
		}
	}

	return string(buffer)
}

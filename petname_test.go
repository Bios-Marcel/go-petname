package petname

import (
	"strings"
	"testing"
	"time"

	"github.com/Bios-Marcel/go-petname/long"
	"github.com/Bios-Marcel/go-petname/medium"
	"github.com/Bios-Marcel/go-petname/short"
	"github.com/stretchr/testify/require"
)

// oldGenerate is the implementation i initially wrote for dustinkirklands
// golang-petname library. I use it as a baseline for performance comparisons.
func oldGenerate(wordCount int, separator string) string {
	switch wordCount {
	case 1:
		return Name()
	case 2:
		return Adjective() + separator + Name()
	case 3:
		// Potentially common cases have shortcut implementations to
		// reduce allocations and CPU usage, even though default: would handle
		// them correctly.
		return Adverb() + separator + Adjective() + separator + Name()
	case 4:
		return Adverb() + separator + Adverb() + separator + Adjective() + separator + Name()
	default:
		words := make([]string, 0, wordCount)
		for i := 0; i < wordCount-2; i++ {
			words = append(words, Adverb())
		}

		return strings.Join(append(words, Adjective(), Name()), separator)
	}
}

func BenchmarkOldGenerate(b *testing.B) {
	for i := 0; i < b.N; i++ {
		oldGenerate(3, "-")
	}
}

func BenchmarkGenerate(b *testing.B) {
	b.Run("lowercase", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			Generate(3, Lower, Hyphen)
		}
	})
	b.Run("titlecase", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			Generate(3, Title, Hyphen)
		}
	})
	b.Run("uppercase", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			Generate(3, Upper, Hyphen)
		}
	})
}

func TestGenerateDifferentWordLists(t *testing.T) {
	SetNames(long.Names)
	SetAdjectives(long.Adjectives)
	SetAdverbs(long.Adverbs)

	Seed(1)
	require.NotEmpty(t, Generate(3, Lower, None))

	SetNames(medium.Names)
	SetAdjectives(medium.Adjectives)
	SetAdverbs(medium.Adverbs)

	Seed(1)
	require.NotEmpty(t, Generate(3, Lower, None))

	// Do short last, so everything is configured correctly for the other tests.

	SetNames(short.Names)
	SetAdjectives(short.Adjectives)
	SetAdverbs(short.Adverbs)

	Seed(1)
	require.NotEmpty(t, Generate(3, Lower, None))
}

func TestGenerate(t *testing.T) {
	Seed(1)
	require.Equal(t, "likelygivingmagpie", Generate(3, Lower, None))
	Seed(1)
	require.Equal(t, "LIKELYGIVINGMAGPIE", Generate(3, Upper, None))
	Seed(1)
	require.Equal(t, "LikelyGivingMagpie", Generate(3, Title, None))
}

func TestGenerateDifferentLengths(t *testing.T) {
	for i := 0; i < 100; i++ {
		Seed(1)
		require.NotEmpty(t, Generate(3, Lower, None))
		Seed(1)
		require.NotEmpty(t, Generate(3, Upper, None))
		Seed(1)
		require.NotEmpty(t, Generate(3, Title, None))
	}
}

func TestGenerateZeroWords(t *testing.T) {
	done := make(chan struct{})

	go func() {
		// Passing 0 words should not cause an "infinite" loop.
		Generate(0, 0, 0)
		done <- struct{}{}
	}()

	select {
	case <-done:
	// Success
	case <-time.NewTimer(2 * time.Second).C:
		t.Fatal("Generate(0, 0, 0) did not return")
	}
}

func FuzzGenerate(f *testing.F) {
	f.Add(uint(3), 1, None)
	f.Fuzz(func(t *testing.T, wordCount uint, casing int, separator Separator) {
		Generate(wordCount, Casing(casing), separator)
	})
}

## Alternative to dustinkirklands golang-petname.

[![Verify](https://github.com/Bios-Marcel/go-petname/actions/workflows/verify.yml/badge.svg)](https://github.com/Bios-Marcel/go-petname/actions/workflows/verify.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/Bios-Marcel/go-petname.svg)](https://pkg.go.dev/github.com/Bios-Marcel/go-petname)

**!The API is incompatible with dustinkirkland/golang-petname!**

The goal is to provide a low overhead alternative that'S well-maintained and
simple to work with.

## Usage

Add to your dependencies:

```sh
go get github.com/Bios-Marcel/go-petname
```

Call the `Generate` function:

```go
// Results in word_word_word
petname.Generate(3, petname.Lower, petname.Underscore)
```

You can change the wordlists by calling `SetNames`, `SetAdjectives` and
`SetAdverbs`. Note that by default, the `short` package is used as the source of
words for all groups. The other packages available are `medium` and `long`.

Technically you can provide your own lists.

## Generate wordlists

The wordlists are generated like this:

```sh
go run ./cmd/generate folder > target/words.go
```

Replace `folder` with the folder that contains `adjectives.txt`, `adverbs.txt`
and `names.txt`.


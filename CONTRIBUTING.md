# Contributing

Thanks for helping. The most useful contribution is a **review of the Nepali word lists by a native speaker**:
spotting entries that are ordinary words or names, adding common spellings, and saying whether an entry is at the
right strictness.

The full guide, including the rules for adding a word, is at
[mukhxadnahunna.com/go/contributing](https://mukhxadnahunna.com/go/contributing).

## Reporting a problem

All the ports share one issue tracker:
[github.com/PG-Momik/no-nepali-profanity/issues](https://github.com/PG-Momik/no-nepali-profanity/issues). Include:

- the **exact input text**,
- which words were found, and what you expected,
- that you used the Go package, and the options you passed, if any.

For a false positive, say whether the word is a name, a place or an ordinary word.

## Development setup

```sh
git clone https://github.com/PG-Momik/no-nepali-profanity-go.git
cd no-nepali-profanity-go
go test -race ./...
go vet ./...
gofmt -l .           # should print nothing
```

`lexicon.go` holds the word lists, `filter.go` the matching logic and `types.go` the exported types.

## Changing the word lists

The word lists are shared by every port, and they're maintained in the
[JavaScript repository](https://github.com/PG-Momik/no-nepali-profanity). Please send word-list changes there rather
than here, so every port gets them. This repository's copy is updated to match.

## Changing the matching rules

Every port gives the same result for the same text, and this one is checked against the JavaScript package's output.
A change to how text is matched (normalization, tokenizing, suffixes, wildcards, phrases or censoring) has to be made
in the JavaScript package first, so please open an issue there before changing it here. Bug fixes that make this port
match the JavaScript package are always welcome.

## Pull requests

- Keep each pull request to one change, and explain why it's needed.
- Add a test to `filter_test.go`: a "catches" case for something that should be flagged, and a "does not flag" case for a
  word or name that shouldn't be.
- Make sure the tests pass before you open it.
- Don't add runtime dependencies without discussing it in an issue first.

By contributing, you agree that your contribution is released under the [MIT License](LICENSE).

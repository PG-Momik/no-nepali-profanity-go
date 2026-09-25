# no-nepali-profanity-go

A small profanity matcher for **English**, **Romanized (Latin) Nepali** and **Devanagari Nepali**, plus the Hindi
slang common in Nepal. Built for moderating user-written text — names, comments, reviews — on Nepali sites, where
false positives on real names are more damaging than a missed swear.

The Go port of [no-nepali-profanity](https://github.com/PG-Momik/no-nepali-profanity). It uses the same word lists
and matching rules as the JavaScript package, so a piece of text gets the same result in every language. Its only
dependency is `golang.org/x/text`, for Unicode normalization.

**Documentation: [mukhxadnahunna.com/go](https://mukhxadnahunna.com/go/)**

## Install

```sh
go get github.com/PG-Momik/no-nepali-profanity-go
```

Requires Go 1.22+.

```go
import nepaliprofanity "github.com/PG-Momik/no-nepali-profanity-go"
```

## Usage

```go
nepaliprofanity.ContainsProfanity("Great teacher!")   // false
nepaliprofanity.ContainsProfanity("muji")             // true
nepaliprofanity.ContainsProfanity("मुजीको कक्षा")      // true  (Devanagari + postposition)

nepaliprofanity.FindProfanity("f.u.c.k this sh1t")    // [fuck shit]
nepaliprofanity.FindProfanity("*ss teacher")          // [*ss]
nepaliprofanity.FindProfanity("Randip Thapa")         // []

nepaliprofanity.Tokenize("Great teacher!")            // [great teacher]
```

### Censoring

```go
nepaliprofanity.Censor("you muji")                                              // "you ****"
nepaliprofanity.Censor("F.U.C.K this Sh1t!")                                    // "******* this ****!"
nepaliprofanity.Censor("you muji", nepaliprofanity.CensorOptions{Mask: "#"})   // "you ####"
nepaliprofanity.Censor("you muji", nepaliprofanity.CensorOptions{
	Replace: func(m nepaliprofanity.ProfanityMatch) string { return "[censored]" },
})                                                                              // "you [censored]"

nepaliprofanity.FindProfanityMatches("you muji")
// [{Text:muji Normalized:muji Start:4 End:8}]
```

`Start` and `End` are byte offsets, so `text[m.Start:m.End] == m.Text`.

Check and censor in one pass:

```go
result := nepaliprofanity.Check("you muji")
result.HasProfanity   // true
result.Words          // [muji]
result.Censor()       // "you ****"
```

### Options

The package-level functions check all three languages at `Standard` strictness. For anything else, build a
`Filter` once and reuse it. It is safe for concurrent use.

```go
f, err := nepaliprofanity.NewFilter(nepaliprofanity.FilterOptions{
	Languages:  []nepaliprofanity.Language{nepaliprofanity.Romanized},
	Strictness: nepaliprofanity.Lenient,
})
f.FindProfanity("fuck muji murkha")   // [muji]
```

- `Languages`: any of `English`, `Romanized`, `Devanagari`. `nil` checks all three; an empty slice checks none.
- `Strictness`: `Lenient` (severe words only), `Standard` (the default; adds milder insults like `idiot`, `murkha`)
  or `Strict` (adds entries that are also ordinary words, like `damn`, and the stems `rand`, `cond`, `kand`, `lund`;
  names they would hit, like `Randip`, are on a built-in allow list).
- `ExtraWords`: more words to flag. `AllowWords`: words never to flag, such as names on your site.

`NewFilter` returns an error for an unknown language or strictness; `MustNewFilter` panics instead.

## API

| Function | Returns | Notes |
|---|---|---|
| `ContainsProfanity(text)` | `bool` | |
| `FindProfanity(text)` | `[]string` | Matching words, **normalized** (leet decoded, lower-cased) and deduplicated. |
| `FindProfanityMatches(text)` | `[]ProfanityMatch` | Every occurrence with its byte offsets, sorted by position. |
| `Censor(text, ...CensorOptions)` | `string` | The text with each match masked. Options: `Mask` (default `"*"`), `Replace`. |
| `Check(text)` | `*ProfanityCheck` | Scans once: `Text`, `HasProfanity`, `Words`, `Matches`, `Censor(...CensorOptions)`. |
| `NewFilter(FilterOptions)` | `(*Filter, error)` | A filter with the same methods, for other options. |
| `Tokenize(text)` | `[]string` | The tokens the matcher sees. Useful for debugging. |

The word lists are exported as `Words`, `Stems` and `Phrases` (tagged with language and strictness), flat per-script
lists (`LatinWords`, `DevanagariWords`…) and `LatinSuffixes` / `DevanagariSuffixes`.

## What it catches

- **Case and Unicode forms**: `IDIOT`, full-width letters.
- **Leetspeak**: `sh1t`, `@ss` (`0 1 3 4 5 7 @ $`).
- **`!` for `i` between letters**: `sh!t`, `b!tch`. Sentence-final `Great teacher!` is left alone.
- **`*` for a hidden letter**: `f*ck`, `sh*t`, `*ss`. Markdown emphasis like `*sh*t*` still reads as the word, while
  `*is*` stays clean.
- **Stretched letters**: `fuuuuck`, for words of 4+ letters.
- **Spelled-out letters**: `f.u.c.k`, `f u c k`.
- **Nepali postpositions and plurals glued on**: `mujiko`, `randiharu`, `मुजीको`, `…हरू`.
- **Devanagari spelling variants**: nukta, chandrabindu vs anusvara, zero-width joiners.
- **Stems** where no ordinary word starts the same way: `fucking`, `bitches`, `machiknee`.
- **Multi-word phrases**: `chaak ko pwal`, `pesa garne`, `sasto manche` (Latin and Devanagari).

## What it deliberately doesn't

- **Short words match exactly**, so `as`, `class`, `assignment` and `Assam` are fine.
- **Name collisions**: `shit` is a whole word only, because **Shitij / शितिज** is a name. Names like **Randip**,
  **Putali / पुतली** and **Asha** are checked in the test suite.
- **No caste names, surnames or ordinary words that are only offensive in context** (e.g. *kami*, *kukur*). A word
  list can't tell a slur from someone's name; that needs human moderation.
- **No judgement of context, sarcasm or meaning.** This is a first-pass filter, not a moderator.

## Development

```sh
go test -race ./...
go run ./example
```

The word lists live in `lexicon.go`, kept in step with `src/lexicon.ts` in the JavaScript package. Native-speaker
review of the Nepali lists is the most valuable contribution; please send word-list changes to the
[JavaScript repository](https://github.com/PG-Momik/no-nepali-profanity) so every port gets them.

## License

MIT

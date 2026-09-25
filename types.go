package nepaliprofanity

// Language is one of the languages the filter checks.
type Language string

const (
	English    Language = "english"
	Romanized  Language = "romanized"
	Devanagari Language = "devanagari"
)

// Strictness sets how much the filter catches.
//   - Lenient: severe profanity only.
//   - Standard: adds milder insults (idiot, murkha, sala…). The default.
//   - Strict: adds stems that also start ordinary words or names (rand → Randip, cond → condition).
type Strictness string

const (
	Lenient  Strictness = "lenient"
	Standard Strictness = "standard"
	Strict   Strictness = "strict"
)

// FilterOptions configures a Filter.
type FilterOptions struct {
	// Languages to check. Nil checks all three; an empty, non-nil slice checks none.
	Languages []Language
	// Strictness sets how much to catch. Empty means Standard.
	Strictness Strictness
}

// ProfanityMatch is one place where profanity was found in the original text.
type ProfanityMatch struct {
	// Text is the matched text exactly as it appears in the input, e.g. "F.U.C.K".
	Text string
	// Normalized is the normalized form that was matched, e.g. "fuck". The same value FindProfanity returns.
	Normalized string
	// Start is the byte offset of the match in the input, so Text == input[Start:End].
	Start int
	// End is the byte offset just past the match.
	End int
}

// CensorOptions configures how matches are masked.
type CensorOptions struct {
	// Mask replaces each visible character of a match, except whitespace. Empty means "*".
	Mask string
	// Replace returns the replacement for a whole match. It takes precedence over Mask.
	Replace func(match ProfanityMatch) string
}

// ProfanityCheck is the result of scanning one text: inspect it, then censor it without scanning again.
type ProfanityCheck struct {
	// Text is the text that was checked.
	Text string
	// HasProfanity reports whether any profanity was found.
	HasProfanity bool
	// Words are the normalized words found, without duplicates. The same as FindProfanity.
	Words []string
	// Matches are every match with its position, sorted by position. The same as FindProfanityMatches.
	Matches []ProfanityMatch

	found []ProfanityMatch
}

// Censor returns the checked text with every match masked.
func (c *ProfanityCheck) Censor(options ...CensorOptions) string {
	return censorMatches(c.Text, c.found, firstOr(options))
}

func firstOr(options []CensorOptions) CensorOptions {
	if len(options) > 0 {
		return options[0]
	}
	return CensorOptions{}
}

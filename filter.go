package nepaliprofanity

import (
	"fmt"
	"sort"
	"strings"
	"sync"
	"unicode"
	"unicode/utf8"

	"golang.org/x/text/unicode/norm"
)

var allLanguages = []Language{English, Romanized, Devanagari}

var strictnessLevel = map[Strictness]int{Lenient: 0, Standard: 1, Strict: 2}

var leet = map[rune]rune{'0': 'o', '1': 'i', '3': 'e', '4': 'a', '5': 's', '7': 't', '@': 'a', '$': 's'}

// Words shorter than this after collapsing repeated letters must match exactly, so "as" never matches "ass".
const minCollapse = 4

func isDevanagariRune(r rune) bool { return r >= 0x0900 && r <= 0x097F }

func isDevanagari(s string) bool {
	for _, r := range s {
		if isDevanagariRune(r) {
			return true
		}
	}
	return false
}

func isZeroWidth(r rune) bool { return (r >= 0x200B && r <= 0x200D) || r == 0x2060 || r == 0xFEFF }

func isLetterOrNumber(r rune) bool { return unicode.IsLetter(r) || unicode.IsNumber(r) }

// isSpace matches the same characters as \s in JavaScript, which the other ports follow.
func isSpace(r rune) bool {
	switch r {
	case '\t', '\n', '\v', '\f', '\r', ' ', 0x00A0, 0x1680, 0x2028, 0x2029, 0x202F, 0x205F, 0x3000, 0xFEFF:
		return true
	}
	return r >= 0x2000 && r <= 0x200A
}

func runeLen(s string) int { return utf8.RuneCountInString(s) }

// collapse turns every run of a repeated letter into one: "fuuuuck" → "fuck".
func collapse(s string) string { return shortenRuns(s, 1, 2) }

// squeeze turns every run of three or more of a letter into two: "fuuuuck" → "fuuck", but "mooji" stays.
func squeeze(s string) string { return shortenRuns(s, 2, 3) }

// shortenRuns cuts each run of one rune that is at least minRun long down to keep runes.
func shortenRuns(s string, keep, minRun int) string {
	var b strings.Builder
	runes := []rune(s)
	for i := 0; i < len(runes); {
		j := i
		for j < len(runes) && runes[j] == runes[i] {
			j++
		}
		n := j - i
		if n >= minRun {
			n = keep
		}
		for k := 0; k < n; k++ {
			b.WriteRune(runes[i])
		}
		i = j
	}
	return b.String()
}

func normalizeRune(r rune) string {
	if isZeroWidth(r) {
		return ""
	}
	if isDevanagariRune(r) {
		// Decompose so a precomposed nukta letter (ऩ) loses its nukta too, and fold chandrabindu into anusvara.
		s := norm.NFD.String(string(r))
		s = strings.ReplaceAll(s, "़", "")
		return strings.ReplaceAll(s, "ँ", "ं")
	}
	var b strings.Builder
	for _, c := range norm.NFKC.String(string(r)) {
		if c == 'İ' {
			// Lower-case it the way JavaScript and Python do, keeping the dot as a combining mark.
			b.WriteString("i̇")
			continue
		}
		c = unicode.ToLower(c)
		if l, ok := leet[c]; ok {
			c = l
		}
		b.WriteRune(c)
	}
	return b.String()
}

// normalized is the normalized text, plus the span of the original text that each byte of it came from, so a
// match found in the normalized text can be traced back to what the user typed.
type normalized struct {
	text   string
	starts []int
	ends   []int
}

func normalize(input string) normalized {
	var b strings.Builder
	var starts, ends []int
	for offset := 0; offset < len(input); {
		r, size := utf8.DecodeRuneInString(input[offset:])
		out := normalizeRune(r)
		b.WriteString(out)
		for k := 0; k < len(out); k++ {
			starts = append(starts, offset)
			ends = append(ends, offset+size)
		}
		offset += size
	}
	return replaceBangs(normalized{text: b.String(), starts: starts, ends: ends})
}

// replaceBangs turns each run of "!" between two letters or digits into one "i" spanning the whole run, so
// "sh!t" reads as "shit" while a sentence-final "!" stays punctuation.
func replaceBangs(n normalized) normalized {
	text := n.text
	var b strings.Builder
	var starts, ends []int
	last := 0
	for i := 0; i < len(text); {
		if text[i] != '!' {
			_, size := utf8.DecodeRuneInString(text[i:])
			i += size
			continue
		}
		j := i
		for j < len(text) && text[j] == '!' {
			j++
		}
		before, _ := utf8.DecodeLastRuneInString(text[:i])
		after, _ := utf8.DecodeRuneInString(text[j:])
		if i > 0 && j < len(text) && isLetterOrNumber(before) && isLetterOrNumber(after) {
			b.WriteString(text[last:i])
			b.WriteByte('i')
			starts = append(append(starts, n.starts[last:i]...), n.starts[i])
			ends = append(append(ends, n.ends[last:i]...), n.ends[j-1])
			last = j
		}
		i = j
	}
	if last == 0 {
		return n
	}
	b.WriteString(text[last:])
	starts = append(starts, n.starts[last:]...)
	ends = append(ends, n.ends[last:]...)
	return normalized{text: b.String(), starts: starts, ends: ends}
}

func normalizeText(s string) string { return normalize(s).text }

type tables struct {
	latinExact     map[string]bool
	latinCollapsed map[string]bool
	latinWords     [][]rune
	latinStems     []string
	devWords       map[string]bool
	devStems       []string
	phrases        [][]string
}

func buildTables(options FilterOptions) (*tables, error) {
	languages := map[Language]bool{}
	if options.Languages == nil {
		for _, l := range allLanguages {
			languages[l] = true
		}
	}
	for _, l := range options.Languages {
		if _, ok := map[Language]bool{English: true, Romanized: true, Devanagari: true}[l]; !ok {
			return nil, fmt.Errorf("nepaliprofanity: unknown language %q, use one of: english, romanized, devanagari", l)
		}
		languages[l] = true
	}
	strictness := options.Strictness
	if strictness == "" {
		strictness = Standard
	}
	level, ok := strictnessLevel[strictness]
	if !ok {
		return nil, fmt.Errorf("nepaliprofanity: unknown strictness %q, use one of: lenient, standard, strict", strictness)
	}

	active := func(entries []LexiconEntry, devanagari bool) []string {
		var out []string
		for _, e := range entries {
			if languages[e.Language] && strictnessLevel[e.Strictness] <= level && (e.Language == Devanagari) == devanagari {
				out = append(out, normalizeText(e.Text))
			}
		}
		return out
	}

	t := &tables{
		latinExact:     map[string]bool{},
		latinCollapsed: map[string]bool{},
		devWords:       map[string]bool{},
	}
	for _, w := range active(Words, false) {
		t.latinExact[squeeze(w)] = true
		t.latinWords = append(t.latinWords, []rune(squeeze(w)))
		if c := collapse(w); runeLen(c) >= minCollapse {
			t.latinCollapsed[c] = true
		}
	}
	for _, s := range active(Stems, false) {
		t.latinStems = append(t.latinStems, collapse(s))
	}
	for _, w := range active(Words, true) {
		t.devWords[w] = true
	}
	t.devStems = active(Stems, true)
	for _, p := range append(active(Phrases, false), active(Phrases, true)...) {
		t.phrases = append(t.phrases, strings.Fields(p))
	}
	return t, nil
}

// wildcardEqual reports whether pattern matches word, where each "*" in pattern stands for one hidden letter.
func wildcardEqual(pattern, word []rune) bool {
	if len(pattern) != len(word) {
		return false
	}
	for i, p := range pattern {
		if p != '*' && p != word[i] {
			return false
		}
	}
	return true
}

func (t *tables) wildcardTokenMatches(token string) bool {
	if !strings.Contains(token, "*") {
		return false
	}
	clean := strings.Trim(token, "*")
	if strings.IndexFunc(clean, unicode.IsLetter) < 0 {
		return false
	}

	// "*" on both ends is markdown emphasis ("*sh*t*"). On one end only, it may also hide a first or last letter ("*ss").
	tokens := []string{clean}
	emphasis := strings.HasPrefix(token, "*") && strings.HasSuffix(token, "*")
	if !emphasis && clean != token {
		tokens = append(tokens, token)
	}

	for _, tok := range tokens {
		for _, f := range []string{squeeze(tok), collapse(tok)} {
			form := []rune(f)
			for _, w := range t.latinWords {
				if wildcardEqual(form, w) {
					return true
				}
			}
			for _, stem := range t.latinStems {
				s := []rune(stem)
				if len(form) >= len(s) && wildcardEqual(form[:len(s)], s) {
					return true
				}
			}
		}
	}
	return false
}

func (t *tables) latinTokenMatches(token string) bool {
	candidates := []string{token}
	for _, s := range LatinSuffixes {
		if strings.HasSuffix(token, s) && runeLen(token)-runeLen(s) >= 3 {
			candidates = append(candidates, strings.TrimSuffix(token, s))
			break
		}
	}
	for _, c := range candidates {
		collapsed := collapse(c)
		if t.latinExact[squeeze(c)] || (runeLen(collapsed) >= minCollapse && t.latinCollapsed[collapsed]) {
			return true
		}
		for _, stem := range t.latinStems {
			if strings.HasPrefix(collapsed, stem) {
				return true
			}
		}
		if t.wildcardTokenMatches(c) {
			return true
		}
	}
	return false
}

func (t *tables) devanagariTokenMatches(token string) bool {
	candidates := []string{token}
	for _, s := range DevanagariSuffixes {
		if strings.HasSuffix(token, s) && runeLen(token) > runeLen(s)+1 {
			candidates = append(candidates, strings.TrimSuffix(token, s))
			break
		}
	}
	for _, c := range candidates {
		if t.devWords[c] {
			return true
		}
		for _, stem := range t.devStems {
			if strings.HasPrefix(c, stem) {
				return true
			}
		}
	}
	return false
}

// span is a token with the span of the original text it came from.
type span struct {
	value      string
	start, end int
}

func isTokenRune(r rune) bool { return unicode.IsLetter(r) || unicode.IsMark(r) || r == '*' }

func tokenSpans(n normalized) []span {
	var tokens, run []span

	// Three or more single letters in a row ("f.u.c.k", "f u c k") are read as one word.
	flush := func() {
		if len(run) >= 3 {
			var b strings.Builder
			for _, t := range run {
				b.WriteString(t.value)
			}
			tokens = append(tokens, span{b.String(), run[0].start, run[len(run)-1].end})
		}
		run = nil
	}

	text := n.text
	for i := 0; i < len(text); {
		r, size := utf8.DecodeRuneInString(text[i:])
		if !isTokenRune(r) {
			i += size
			continue
		}
		j := i
		for j < len(text) {
			r, size := utf8.DecodeRuneInString(text[j:])
			if !isTokenRune(r) {
				break
			}
			j += size
		}
		t := span{text[i:j], n.starts[i], n.ends[j-1]}
		switch {
		case isDevanagari(t.value):
			flush()
			tokens = append(tokens, t)
		case runeLen(t.value) == 1:
			run = append(run, t)
		default:
			flush()
			tokens = append(tokens, t)
		}
		i = j
	}
	flush()
	return tokens
}

// phraseSpans finds each non-overlapping occurrence of a phrase in the normalized text, as byte ranges. A phrase
// starts the text or follows a character that is not a letter or digit, and is not followed by a letter or digit.
func phraseSpans(text string, words []string) [][2]int {
	var spans [][2]int

	// matchAt returns where the phrase ends if it starts at byte i, or -1.
	matchAt := func(i int) int {
		for k, w := range words {
			if k > 0 {
				j := i
				for j < len(text) {
					r, size := utf8.DecodeRuneInString(text[j:])
					if !isSpace(r) {
						break
					}
					j += size
				}
				if j == i {
					return -1
				}
				i = j
			}
			if !strings.HasPrefix(text[i:], w) {
				return -1
			}
			i += len(w)
		}
		if next, _ := utf8.DecodeRuneInString(text[i:]); i < len(text) && isLetterOrNumber(next) {
			return -1
		}
		return i
	}

	for i := 0; i < len(text); {
		if i == 0 {
			if end := matchAt(0); end >= 0 {
				spans = append(spans, [2]int{0, end})
				i = end
				continue
			}
		}
		r, size := utf8.DecodeRuneInString(text[i:])
		if !isLetterOrNumber(r) {
			if end := matchAt(i + size); end >= 0 {
				spans = append(spans, [2]int{i + size, end})
				i = end
				continue
			}
		}
		i += size
	}
	return spans
}

// scan returns every match, words first and then phrases, in the order found.
func (t *tables) scan(text string) []ProfanityMatch {
	if text == "" {
		return nil
	}
	n := normalize(text)
	match := func(normalized string, start, end int) ProfanityMatch {
		return ProfanityMatch{Text: text[start:end], Normalized: normalized, Start: start, End: end}
	}

	var found []ProfanityMatch
	for _, tok := range tokenSpans(n) {
		var hit bool
		if isDevanagari(tok.value) {
			hit = t.devanagariTokenMatches(tok.value)
		} else {
			hit = t.latinTokenMatches(tok.value)
		}
		if hit {
			found = append(found, match(tok.value, tok.start, tok.end))
		}
	}
	for _, words := range t.phrases {
		for _, s := range phraseSpans(n.text, words) {
			found = append(found, match(n.text[s[0]:s[1]], n.starts[s[0]], n.ends[s[1]-1]))
		}
	}
	return found
}

func unique(matches []ProfanityMatch) []string {
	seen := map[string]bool{}
	words := []string{}
	for _, m := range matches {
		if !seen[m.Normalized] {
			seen[m.Normalized] = true
			words = append(words, m.Normalized)
		}
	}
	return words
}

// byPosition returns the matches sorted by start, and a longer match before a shorter one at the same start.
func byPosition(matches []ProfanityMatch) []ProfanityMatch {
	sorted := append([]ProfanityMatch{}, matches...)
	sort.SliceStable(sorted, func(i, j int) bool {
		if sorted[i].Start != sorted[j].Start {
			return sorted[i].Start < sorted[j].Start
		}
		return sorted[i].End > sorted[j].End
	})
	return sorted
}

// Viramas of the Indic scripts, which join the consonants on either side into one visible character.
const viramas = "्্੍્୍்్್്්"

func isGraphemeExtend(r rune) bool {
	return unicode.In(r, unicode.Mn, unicode.Mc, unicode.Me) || (r >= 0xFE00 && r <= 0xFE0F) || r == 0x200C || r == 0x200D
}

// endsWithVirama reports whether a consonant joins this cluster: it ends in a virama, perhaps followed by other
// marks or a ZWJ, so a conjunct such as "ण्ड" masks as one visible character.
func endsWithVirama(cluster string) bool {
	for len(cluster) > 0 {
		r, size := utf8.DecodeLastRuneInString(cluster)
		if strings.ContainsRune(viramas, r) {
			return true
		}
		if r != 0x200D && !unicode.Is(unicode.Mn, r) {
			return false
		}
		cluster = cluster[:len(cluster)-size]
	}
	return false
}

// graphemes splits s into visible characters: a letter with its combining marks, and an Indic conjunct, stay whole.
func graphemes(s string) []string {
	var clusters []string
	for _, r := range s {
		if last := len(clusters) - 1; last >= 0 &&
			(isGraphemeExtend(r) || (unicode.IsLetter(r) && endsWithVirama(clusters[last]))) {
			clusters[last] += string(r)
			continue
		}
		clusters = append(clusters, string(r))
	}
	return clusters
}

func isAllSpace(s string) bool {
	for _, r := range s {
		if !isSpace(r) {
			return false
		}
	}
	return s != ""
}

func censorMatches(text string, matches []ProfanityMatch, options CensorOptions) string {
	mask := options.Mask
	if mask == "" {
		mask = "*"
	}

	// Merge overlapping matches, such as the word "chaak" inside the phrase "chaak ko pwal", into one.
	var merged []ProfanityMatch
	for _, m := range byPosition(matches) {
		if last := len(merged) - 1; last >= 0 && m.Start < merged[last].End {
			if m.End > merged[last].End {
				merged[last].End = m.End
				merged[last].Text = text[merged[last].Start:m.End]
			}
			continue
		}
		merged = append(merged, m)
	}

	var b strings.Builder
	last := 0
	for _, m := range merged {
		b.WriteString(text[last:m.Start])
		if options.Replace != nil {
			b.WriteString(options.Replace(m))
		} else {
			for _, g := range graphemes(m.Text) {
				if isAllSpace(g) {
					b.WriteString(g)
				} else {
					b.WriteString(mask)
				}
			}
		}
		last = m.End
	}
	b.WriteString(text[last:])
	return b.String()
}

// Filter checks text against the word lists for one set of options. Build it once with NewFilter and reuse it;
// it is safe for concurrent use.
type Filter struct {
	tables *tables
}

// NewFilter builds a filter for the given options. It returns an error for an unknown language or strictness.
func NewFilter(options FilterOptions) (*Filter, error) {
	t, err := buildTables(options)
	if err != nil {
		return nil, err
	}
	return &Filter{tables: t}, nil
}

// MustNewFilter is like NewFilter but panics on an unknown language or strictness.
func MustNewFilter(options FilterOptions) *Filter {
	f, err := NewFilter(options)
	if err != nil {
		panic(err)
	}
	return f
}

// Check scans the text once. Use the result to check for profanity and to censor it, e.g. f.Check(text).Censor().
func (f *Filter) Check(text string) *ProfanityCheck {
	found := f.tables.scan(text)
	return &ProfanityCheck{
		Text:         text,
		HasProfanity: len(found) > 0,
		Words:        unique(found),
		Matches:      byPosition(found),
		found:        found,
	}
}

// ContainsProfanity reports whether the text contains any profanity.
func (f *Filter) ContainsProfanity(text string) bool { return len(f.tables.scan(text)) > 0 }

// FindProfanity returns the normalized words found, without duplicates, e.g. ["fuck", "shit"] for "f.u.c.k sh1t".
func (f *Filter) FindProfanity(text string) []string { return unique(f.tables.scan(text)) }

// FindProfanityMatches returns every occurrence with its position in the input, sorted by position.
func (f *Filter) FindProfanityMatches(text string) []ProfanityMatch {
	return byPosition(f.tables.scan(text))
}

// Censor returns the text with every match masked, e.g. "you muji" → "you ****".
func (f *Filter) Censor(text string, options ...CensorOptions) string {
	return censorMatches(text, f.tables.scan(text), firstOr(options))
}

var defaultFilter = sync.OnceValue(func() *Filter { return MustNewFilter(FilterOptions{}) })

// Check scans the text once with the default options: all languages, Standard strictness.
func Check(text string) *ProfanityCheck { return defaultFilter().Check(text) }

// ContainsProfanity reports whether the text contains any profanity, with the default options.
func ContainsProfanity(text string) bool { return defaultFilter().ContainsProfanity(text) }

// FindProfanity returns the normalized words found, without duplicates, with the default options.
func FindProfanity(text string) []string { return defaultFilter().FindProfanity(text) }

// FindProfanityMatches returns every occurrence with its position in the input, with the default options.
func FindProfanityMatches(text string) []ProfanityMatch {
	return defaultFilter().FindProfanityMatches(text)
}

// Censor returns the text with every match masked, with the default options.
func Censor(text string, options ...CensorOptions) string {
	return defaultFilter().Censor(text, options...)
}

// Tokenize returns the tokens the matcher sees. Useful for debugging why a word is or isn't caught.
func Tokenize(text string) []string {
	tokens := []string{}
	if text == "" {
		return tokens
	}
	for _, s := range tokenSpans(normalize(text)) {
		tokens = append(tokens, s.value)
	}
	return tokens
}

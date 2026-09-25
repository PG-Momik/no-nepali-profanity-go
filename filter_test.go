package nepaliprofanity

import (
	"reflect"
	"strings"
	"testing"
)

func TestCatches(t *testing.T) {
	cases := []struct{ name, text string }{
		{"English", "what a bitch"},
		{"an inflected English word (stem)", "fucking useless"},
		{"Romanized Nepali", "muji teacher"},
		{"Romanized Nepali with a postposition", "machikneko class"},
		{"Devanagari Nepali", "यो मुजी हो"},
		{"Devanagari with a postposition", "मुजीको कक्षा"},
		{"Devanagari with a nukta or zero-width joiner", "मु‍जी"},
		{"leetspeak", "sh1t lecturer"},
		{"leetspeak with @", "@ss"},
		{"leetspeak with $", "a$$"},
		{"dodging with ! for i inside a word", "sh!t lecturer"},
		{"dodging with a wildcard for a hidden letter", "f*ck this"},
		{"dodging with a wildcard for the hidden i", "sh*t"},
		{"dodging with a wildcard for the hidden first letter", "that *ss"},
		{"dodging with a wildcard for the hidden last letter", "fuc* off"},
		{"markdown emphasis still reads the word", "*sh*t* is bad"},
		{"stretched letters", "fuuuuck"},
		{"letters spelled out with dots", "f.u.c.k"},
		{"letters spelled out with spaces", "m u j i"},
		{"upper case", "IDIOT"},
		{"full-width letters", "ＦＵＣＫ"},
		{"Hindi slang common in Nepal", "chutiya"},
		{"Romanized invective", "murkha"},
		{"Dodged word stays caught with a suffix", "gandako budi"},
		{"exact word, not a long place name", "look at that gand"},
		{"stem catches the -ne inflected form", "chodne manche"},
		{"Devanagari invective", "मुर्ख"},
		{"Devanagari with a doubled consonant", "थुक्क"},
		{"Devanagari slang", "कमिना"},
		{"Devanagari vulgar term", "लुंड"},
		{"multi-word phrase", "chaak ko pwal"},
		{"multi-word phrase with leetspeak", "p3sa g@rne taba"},
		{"multi-word phrase ending a sentence", "thulo sasto manche!"},
		{"review-supplied term", "chhakka lai hami mukhulla bhanchha"},
		{"short spelling", "mji"},
		{"leet spelling", "m00ji"},
		{"English leet spelling", "f4ck off"},
		{"Devanagari spelling", "राण्डी"},
		{"Latin phrase", "khatako choro"},
		{"Devanagari phrase", "राण्डीको बान"},
		{"spelling variant with -ey", "yo khatey payment app kahiley chaley po"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if !ContainsProfanity(c.text) {
				t.Errorf("ContainsProfanity(%q) = false, want true", c.text)
			}
		})
	}
}

func TestDoesNotFlag(t *testing.T) {
	texts := []string{
		// Ordinary words that contain or resemble a listed word
		"The class assignment was as hard as expected",
		"Computing and data structures",
		"Dickson explained Scunthorpe problems",
		"Assam and Gandaki are places",
		// Real Nepali names, in both scripts
		"Kshitij Shrestha",
		"Shitij Adhikari",
		"Putali Gurung",
		"पुतली गुरुङ",
		"Randip Thapa",
		"Asha Sharma",
		"Machindra Karki",
		"Harimaya Tamang",
		"सीता कार्की",
		"क्षितिज श्रेष्ठ",
		// Sentence punctuation
		"Great teacher!",
		"No way! That can't be right",
		"the starred items are on page 12*",
		"feed ** me ** the list",
		// Markdown emphasis around an ordinary word must not read the asterisks as hidden letters
		"this *is* good",
		"*and* then",
		"**hi** there",
		// Ordinary words that must not start matching the stems
		"chicken biryani is good",
		"the salaam greeting sounded nice",
		"Gandaki river is in Nepal",
		// Words separated so a phrase must not match
		"sasto ra manche duitai ho",
		// Ordinary words that the Strict-only stems would catch
		"terms and conditions",
		"the conductor",
		"bhrastachar kanda",
		"Kandel sir",
		"Lundberg",
		// Ordinary words removed from the lexicon
		"fohor pani",
		"kano manche",
		"lato keta",
		"फोहोर पानी",
		"लाटो केटा",
		"",
	}
	for _, text := range texts {
		if got := FindProfanity(text); len(got) != 0 {
			t.Errorf("FindProfanity(%q) = %q, want none", text, got)
		}
	}
}

func TestTokenize(t *testing.T) {
	cases := map[string][]string{
		"f.u.c.k this": {"fuck", "this"},
		"सीता कार्की":  {"सीता", "कार्की"},
		"f*ck this!":   {"f*ck", "this"},
		"":             {},
	}
	for text, want := range cases {
		if got := Tokenize(text); !reflect.DeepEqual(got, want) {
			t.Errorf("Tokenize(%q) = %q, want %q", text, got, want)
		}
	}
}

func TestLanguages(t *testing.T) {
	romanized := MustNewFilter(FilterOptions{Languages: []Language{Romanized}})
	for text, want := range map[string]bool{"muji": true, "fuck": false, "मुजी": false, "sasto manche": true} {
		if got := romanized.ContainsProfanity(text); got != want {
			t.Errorf("romanized ContainsProfanity(%q) = %v, want %v", text, got, want)
		}
	}

	english := MustNewFilter(FilterOptions{Languages: []Language{English}})
	if got := english.FindProfanity("fuck muji मुजी"); !reflect.DeepEqual(got, []string{"fuck"}) {
		t.Errorf("english FindProfanity = %q", got)
	}
	devanagari := MustNewFilter(FilterOptions{Languages: []Language{Devanagari}})
	if got := devanagari.FindProfanity("fuck muji मुजी"); !reflect.DeepEqual(got, []string{"मुजी"}) {
		t.Errorf("devanagari FindProfanity = %q", got)
	}
	none := MustNewFilter(FilterOptions{Languages: []Language{}})
	if got := none.FindProfanity("fuck muji मुजी"); len(got) != 0 {
		t.Errorf("no languages FindProfanity = %q, want none", got)
	}
	if _, err := NewFilter(FilterOptions{Languages: []Language{"hindi"}}); err == nil {
		t.Error("NewFilter accepted an unknown language")
	}
}

func TestStrictness(t *testing.T) {
	lenient := MustNewFilter(FilterOptions{Strictness: Lenient})
	for text, want := range map[string]bool{"idiot": false, "murkha": false, "sasto manche": false, "muji": true} {
		if got := lenient.ContainsProfanity(text); got != want {
			t.Errorf("lenient ContainsProfanity(%q) = %v, want %v", text, got, want)
		}
	}

	standard := MustNewFilter(FilterOptions{Strictness: Standard})
	if !ContainsProfanity("idiot") || !standard.ContainsProfanity("idiot") || ContainsProfanity("Randip Thapa") {
		t.Error("Standard is not the default")
	}

	strict := MustNewFilter(FilterOptions{Strictness: Strict})
	for text, want := range map[string][]string{
		"randikoban":           {"randikoban"},
		"terms and conditions": {"conditions"},
		"Randip Thapa":         {"randip"},
	} {
		if got := strict.FindProfanity(text); !reflect.DeepEqual(got, want) {
			t.Errorf("strict FindProfanity(%q) = %q, want %q", text, got, want)
		}
	}

	if _, err := NewFilter(FilterOptions{Strictness: "max"}); err == nil {
		t.Error("NewFilter accepted an unknown strictness")
	}
}

func TestCombinedOptions(t *testing.T) {
	f := MustNewFilter(FilterOptions{Languages: []Language{Romanized}, Strictness: Lenient})
	if !f.ContainsProfanity("muji") || f.ContainsProfanity("murkha") {
		t.Error("combined options not applied")
	}
	if got := f.FindProfanity("fuck muji"); !reflect.DeepEqual(got, []string{"muji"}) {
		t.Errorf("FindProfanity = %q", got)
	}
}

func TestFindProfanityMatches(t *testing.T) {
	got := FindProfanityMatches("F.U.C.K this sh1t, muji. MUJI")
	want := []ProfanityMatch{
		{Text: "F.U.C.K", Normalized: "fuck", Start: 0, End: 7},
		{Text: "sh1t", Normalized: "shit", Start: 13, End: 17},
		{Text: "muji", Normalized: "muji", Start: 19, End: 23},
		{Text: "MUJI", Normalized: "muji", Start: 25, End: 29},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("FindProfanityMatches = %+v, want %+v", got, want)
	}

	got = FindProfanityMatches("chaak ko pwal")
	want = []ProfanityMatch{
		{Text: "chaak ko pwal", Normalized: "chaak ko pwal", Start: 0, End: 13},
		{Text: "chaak", Normalized: "chaak", Start: 0, End: 5},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("a phrase should come before a word it contains: %+v", got)
	}

	// Offsets are in bytes, so they slice the input.
	text := "यो मुजीको कक्षा"
	for _, m := range FindProfanityMatches(text) {
		if text[m.Start:m.End] != m.Text {
			t.Errorf("text[%d:%d] = %q, want %q", m.Start, m.End, text[m.Start:m.End], m.Text)
		}
	}

	if got := FindProfanityMatches("Great teacher!"); len(got) != 0 {
		t.Errorf("clean text matched: %+v", got)
	}
}

func TestCensor(t *testing.T) {
	cases := map[string]string{
		"you muji":           "you ****",
		"F.U.C.K this Sh1t!": "******* this ****!",
		"sh!!t happens":      "***** happens",
		"ＦＵＣＫ off":           "**** off",
		"fuuuuck yeah":       "******* yeah",
		"f*ck and *sh*t*":    "**** and ******",
		"muji muji":          "**** ****",
		"you 😀 muji 😀":       "you 😀 **** 😀",
		// Devanagari is masked by visible character, and a conjunct is one character.
		"मुजीको कक्षा": "*** कक्षा",
		"गाण्ड":        "**",
		// The spaces inside a phrase are kept, and a word inside a phrase is merged with it.
		"sasto   manche":  "*****   ******",
		"chaak ko pwal":   "***** ** ****",
		"Great teacher!":  "Great teacher!",
		"Shitij is great": "Shitij is great",
		"":                "",
	}
	for text, want := range cases {
		if got := Censor(text); got != want {
			t.Errorf("Censor(%q) = %q, want %q", text, got, want)
		}
	}
}

func TestCensorOptions(t *testing.T) {
	if got := Censor("you muji", CensorOptions{Mask: "#"}); got != "you ####" {
		t.Errorf("mask: %q", got)
	}
	replaced := Censor("you muji fuck", CensorOptions{Mask: "#", Replace: func(ProfanityMatch) string { return "[censored]" }})
	if replaced != "you [censored] [censored]" {
		t.Errorf("replace: %q", replaced)
	}
	firstLetter := Censor("you muji", CensorOptions{Replace: func(m ProfanityMatch) string {
		return m.Text[:1] + strings.Repeat("*", len(m.Text)-1)
	}})
	if firstLetter != "you m***" {
		t.Errorf("replace with match: %q", firstLetter)
	}

	if got := MustNewFilter(FilterOptions{Languages: []Language{Romanized}}).Censor("fuck muji"); got != "fuck ****" {
		t.Errorf("languages: %q", got)
	}
	if got := MustNewFilter(FilterOptions{Strictness: Lenient}).Censor("you idiot"); got != "you idiot" {
		t.Errorf("lenient: %q", got)
	}
	if got := MustNewFilter(FilterOptions{Strictness: Strict}).Censor("Randip Thapa"); got != "****** Thapa" {
		t.Errorf("strict: %q", got)
	}
}

func TestCheck(t *testing.T) {
	c := Check("you muji, F.U.C.K")
	if c.Text != "you muji, F.U.C.K" || !c.HasProfanity || !reflect.DeepEqual(c.Words, []string{"muji", "fuck"}) {
		t.Errorf("Check = %+v", c)
	}
	if len(c.Matches) != 2 || c.Matches[0].Text != "muji" || c.Matches[1].Text != "F.U.C.K" {
		t.Errorf("Check matches = %+v", c.Matches)
	}
	if got := c.Censor(); got != "you ****, *******" {
		t.Errorf("Check censor = %q", got)
	}
	if got := c.Censor(CensorOptions{Mask: "#"}); got != "you ####, #######" {
		t.Errorf("Check censor with mask = %q", got)
	}

	clean := Check("Great teacher!")
	if clean.HasProfanity || len(clean.Words) != 0 || len(clean.Matches) != 0 || clean.Censor() != "Great teacher!" {
		t.Errorf("clean Check = %+v", clean)
	}

	if got := MustNewFilter(FilterOptions{Languages: []Language{Romanized}}).Check("fuck muji").Censor(); got != "fuck ****" {
		t.Errorf("Check with options = %q", got)
	}
}

func TestConcurrentUse(t *testing.T) {
	done := make(chan bool)
	for i := 0; i < 8; i++ {
		go func() {
			for j := 0; j < 100; j++ {
				if Censor("you muji") != "you ****" {
					t.Error("concurrent Censor gave a different result")
				}
			}
			done <- true
		}()
	}
	for i := 0; i < 8; i++ {
		<-done
	}
}

package nepaliprofanity

// The word lists, kept in step with lexicon.ts in the JavaScript package.

// LexiconEntry is one word, stem or phrase, with the language it belongs to and the lowest strictness that turns
// it on.
type LexiconEntry struct {
	Text       string
	Language   Language
	Strictness Strictness
}

func tag(language Language, strictness Strictness, texts []string) []LexiconEntry {
	entries := make([]LexiconEntry, len(texts))
	for i, text := range texts {
		entries[i] = LexiconEntry{Text: text, Language: language, Strictness: strictness}
	}
	return entries
}

func concat(groups ...[]LexiconEntry) []LexiconEntry {
	var all []LexiconEntry
	for _, g := range groups {
		all = append(all, g...)
	}
	return all
}

// Words are matched as whole words, after Nepali postpositions are taken off.
var Words = concat(
	tag(English, Lenient, []string{
		"fuck", "fuk", "fck", "phuck", "shit", "shitty", "shithead", "bullshit", "bitch", "bastard", "ass",
		"asshole", "arsehole", "dumbass", "dick", "dickhead", "cunt", "whore", "slut", "cock", "pussy", "twat",
		"wanker", "retard", "fack",
	}),
	tag(English, Standard, []string{
		"piss", "idiot", "stupid", "moron",
	}),
	tag(Romanized, Lenient, []string{
		"muji", "mujhi", "muzi", "machikne", "machhikne", "mchikne", "mcikne", "machikney", "randi", "raandi",
		"rando", "rande", "radi", "lado", "lodo", "puti", "geda", "jatha", "jantha", "jathya", "chikne", "chikney",
		"bhalu", "khate", "khatey", "harami", "gandu", "chutiya", "chutia", "bhosdi", "bhosadi", "bhosdike", "bsdk",
		"madarchod", "behenchod", "bhenchod", "chhakka", "lauro", "gukhane", "gand", "gaand", "gandako", "lund",
		"lundra", "lundri", "chod", "chodna", "beshya", "hijada", "kamina", "haramzada", "turi", "pakhe", "condo",
		"kando", "chaak", "gula", "bajiya", "mji", "mzi", "mujj", "mooji", "moozi", "mcne", "mechikne", "laado",
		"puuti", "zatya", "chodeko", "toori", "kundo",
	}),
	tag(Romanized, Standard, []string{
		"kutta", "kutti", "kuttiya", "murkha", "badmas", "sala", "saley", "sali", "chhucho", "chhuchi", "gu",
		"thukk", "nalayak", "beijjat", "nikamma", "ghinlagdo", "nindaniya", "paji", "moot", "bhate", "chhura",
		"torpe", "mukhulla", "gobre", "bhusya", "dhurt",
	}),
	tag(Devanagari, Lenient, []string{
		"मुजी", "मुजि", "माचिक्ने", "मचिक्ने", "रण्डी", "रन्डी", "रंडी", "रांडी", "रण्डो", "राण्डे", "राडी", "लाडो",
		"लांडो", "पुती", "गेडा", "जाठा", "जांठा", "जाठ्या", "चिक्ने", "भालु", "खाते", "हरामी", "गान्डु", "गांडु",
		"चुतिया", "भोस्डी", "भोसडी", "भोस्डीके", "मादरचोद", "बहनचोद", "भेनचोद", "छक्का", "लौरो", "गुखाने", "गान्ड",
		"गाण्ड", "गान्डको", "लुंड", "लुण्ड", "लुन्ड्रा", "लुन्ड्री", "लोडो", "चोद", "चोद्ना", "चोदेको", "वेश्या",
		"हिजडा", "कमिना", "हरामजादा", "तुरी", "पाखे", "कोंडो", "काण्डो", "कुन्डो", "चाक", "गुला", "बजिया", "राण्डी",
		"गाण्डु", "भोसडीके", "गाण्डको", "कोन्डो",
	}),
	tag(Devanagari, Standard, []string{
		"कुत्ता", "कुत्ती", "कुत्तिया", "मुर्ख", "मूर्ख", "बदमास", "साला", "साले", "साली", "छुच्चो", "छुच्ची", "गु",
		"किचकिच", "थुक", "थुक्क", "नालायक", "बेइज्जत", "निकम्मा", "घिनलाग्दो", "निन्दनीय", "पाजी", "मूत", "भाते",
		"छुरा", "टोर्पे", "मुखुल्ला", "गोबरे", "भुस्या", "धूर्त",
	}),
)

// Stems are matched at the start of a word, so "fuck" also catches "fucking". Entries tagged Strict also start
// ordinary words or names: rand → Randip, cond → condition, kand → kanda / Kandel, lund → Lundberg.
var Stems = concat(
	tag(English, Lenient, []string{
		"fuck", "motherfuck", "bitch", "bastard", "asshol", "cunt", "whore", "slut", "wank", "retard",
	}),
	tag(Romanized, Lenient, []string{
		"machikn", "mchikn", "chutiy", "bhosd", "madarch", "behench", "bhench", "chikn", "chickn", "chod", "jath",
	}),
	tag(Romanized, Strict, []string{
		"rand", "cond", "kand", "lund",
	}),
	tag(Devanagari, Lenient, []string{
		"माचिक्न", "मचिक्न", "चुतिय", "भोस्ड", "मादरच", "बहनच", "भेनच", "चोद", "चिक्न", "रण्ड", "लुण्ड", "जाठ",
	}),
	tag(Devanagari, Strict, []string{
		"कोन्ड", "कान्ड",
	}),
)

// Phrases are matched as a run of words separated by whitespace.
var Phrases = concat(
	tag(Romanized, Lenient, []string{
		"chaak ko pwal", "tero aama ko", "muji jasto", "lado khaye", "lado khos", "randi ko choro", "randi ko ban",
		"gand mara", "gand fatchya", "geda jasto", "geda khaya", "khatako choro", "bhaluko ban", "machikne khate",
	}),
	tag(Romanized, Standard, []string{
		"pesa garne", "sasto manche",
	}),
	tag(Devanagari, Lenient, []string{
		"चाकको प्वाल", "चाक को प्वाल", "तेरो आमाको", "मुजी जस्तो", "लाडो खाए", "लाडो खोस्", "राण्डीको छोरो",
		"राण्डीको बान", "गाण्ड मरा", "गाण्ड फाट्या", "गेडा जस्तो", "गेडा खाया", "खातेको छोरो", "भालुको बान",
		"माचिक्ने खाते",
	}),
	tag(Devanagari, Standard, []string{
		"पेसा गर्ने", "सस्तो मान्छे",
	}),
)

// LatinSuffixes are the Nepali postpositions and plural endings taken off a Latin-script word, longest first.
var LatinSuffixes = []string{
	"haruko", "harule", "sanga", "haru", "bata", "lai", "ko", "ki", "ka", "le", "ma", "ni", "ne", "yo",
}

// DevanagariSuffixes are the Nepali postpositions and plural endings taken off a Devanagari word, longest first.
var DevanagariSuffixes = []string{
	"हरूको", "हरुको", "हरूले", "हरुले", "हरू", "हरु", "बाट", "सँग", "संग", "लाई", "को", "की", "का", "ले", "मा",
	"नि", "ने", "यो",
}

func texts(entries []LexiconEntry, latin bool) []string {
	var out []string
	for _, e := range entries {
		if (e.Language != Devanagari) == latin {
			out = append(out, e.Text)
		}
	}
	return out
}

// Flat lists of every entry at every strictness, by script.
var (
	LatinWords        = texts(Words, true)
	LatinStems        = texts(Stems, true)
	LatinPhrases      = texts(Phrases, true)
	DevanagariWords   = texts(Words, false)
	DevanagariStems   = texts(Stems, false)
	DevanagariPhrases = texts(Phrases, false)
)

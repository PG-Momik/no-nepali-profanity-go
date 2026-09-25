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
		"fag", "fags", "faggot", "faggots", "fagot", "faggy", "nigger", "niggers", "nigga", "niggas", "tranny", "kike",
		"shitface", "shithole", "dipshit", "horseshit", "batshit", "apeshit", "jackass", "asshat", "asswipe", "cocksucker",
		"jizz", "dildo", "skank", "douchebag", "motherfucker", "bollocks", "bellend",
	}),
	tag(English, Standard, []string{
		"piss", "idiot", "stupid", "moron",
		"crap", "crappy", "bugger", "douche", "jerk", "scumbag", "dumb", "wtf", "stfu", "tits", "boobs", "porn", "porno",
		"horny", "bimbo", "thot",
	}),
	// Each of these is also an ordinary word: spic (span), chink (in the armour), dyke (a wall), hoe (a tool), cum
	// (laude), prick (a pin), damn.
	tag(English, Strict, []string{
		"spic", "chink", "dyke", "hoe", "cum", "prick", "damn", "rape",
	}),
	tag(Romanized, Lenient, []string{
		"muji", "mujhi", "muzi", "machikne", "machhikne", "mchikne", "mcikne", "machikney", "randi", "raandi",
		"rando", "rande", "radi", "lado", "lodo", "puti", "geda", "jatha", "jantha", "jathya", "chikne", "chikney",
		"bhalu", "khate", "khatey", "harami", "gandu", "chutiya", "chutia", "bhosdi", "bhosadi", "bhosdike", "bsdk",
		"madarchod", "behenchod", "bhenchod", "chhakka", "lauro", "gukhane", "gand", "gaand", "gandako", "lund",
		"lundra", "lundri", "chod", "chodna", "beshya", "hijada", "kamina", "haramzada", "turi", "pakhe", "condo",
		"kando", "chaak", "gula", "bajiya", "mji", "mzi", "mujj", "mooji", "moozi", "mcne", "mechikne", "laado",
		"puuti", "zatya", "chodeko", "toori", "kundo",
		"chhakke", "maxikne", "mxikne", "xutiya", "xutia", "chhutiya", "chootiya", "xikne", "jhant", "jhaant",
		"jhantu", "lauda", "lavda", "lawda", "loda", "lwado", "lwada", "bhosda", "bhosri", "chhinal", "chhinar", "besya",
		"maachod", "madarchood",
	}),
	tag(Romanized, Standard, []string{
		"kutta", "kutti", "kuttiya", "murkha", "badmas", "sala", "saley", "sali", "chhucho", "chhuchi", "gu",
		"thukk", "nalayak", "beijjat", "nikamma", "ghinlagdo", "nindaniya", "paji", "moot", "bhate", "chhura",
		"torpe", "mukhulla", "gobre", "bhusya", "dhurt",
		"gadha", "ullu", "badmash", "thukka", "haramkhor", "fataha",
	}),
	tag(Devanagari, Lenient, []string{
		"मुजी", "मुजि", "माचिक्ने", "मचिक्ने", "रण्डी", "रन्डी", "रंडी", "रांडी", "रण्डो", "राण्डे", "राडी", "लाडो",
		"लांडो", "पुती", "गेडा", "जाठा", "जांठा", "जाठ्या", "चिक्ने", "भालु", "खाते", "हरामी", "गान्डु", "गांडु",
		"चुतिया", "भोस्डी", "भोसडी", "भोस्डीके", "मादरचोद", "बहनचोद", "भेनचोद", "छक्का", "लौरो", "गुखाने", "गान्ड",
		"गाण्ड", "गान्डको", "लुंड", "लुण्ड", "लुन्ड्रा", "लुन्ड्री", "लोडो", "चोद", "चोद्ना", "चोदेको", "वेश्या",
		"हिजडा", "कमिना", "हरामजादा", "तुरी", "पाखे", "कोंडो", "काण्डो", "कुन्डो", "चाक", "गुला", "बजिया", "राण्डी",
		"गाण्डु", "भोसडीके", "गाण्डको", "कोन्डो",
		"छिनाल", "झांट", "झाँट", "लौडा", "लवडा",
	}),
	tag(Devanagari, Standard, []string{
		"कुत्ता", "कुत्ती", "कुत्तिया", "मुर्ख", "मूर्ख", "बदमास", "साला", "साले", "साली", "छुच्चो", "छुच्ची", "गु",
		"किचकिच", "थुक", "थुक्क", "नालायक", "बेइज्जत", "निकम्मा", "घिनलाग्दो", "निन्दनीय", "पाजी", "मूत", "भाते",
		"छुरा", "टोर्पे", "मुखुल्ला", "गोबरे", "भुस्या", "धूर्त",
		"गधा", "उल्लु", "उल्लू", "बदमाश", "थुक्का", "हरामखोर", "फटाहा",
	}),
)

// Stems are matched at the start of a word, so "fuck" also catches "fucking". Entries tagged Strict also start
// ordinary words or names: rand → Randip, cond → condition, kand → kanda / Kandel, lund → Lundberg. The ones they
// hit most are in Allowed.
var Stems = concat(
	tag(English, Lenient, []string{
		"fuck", "motherfuck", "bitch", "bastard", "asshol", "cunt", "whore", "slut", "wank", "retard",
		"shit", "nigger", "cocksuck", "bullshit", "dickhead", "douchebag", "jizz", "dildo",
	}),
	tag(Romanized, Lenient, []string{
		"machikn", "mchikn", "chutiy", "bhosd", "madarch", "behench", "bhench", "chikn", "chickn", "chod", "jath",
		"xutiy", "chhutiy", "maxikn", "jhant",
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
		"teri maa ki", "teri ma ki",
	}),
	tag(Romanized, Standard, []string{
		"pesa garne", "sasto manche",
		"gu khane", "gu khaa",
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

// Infixes are Latin roots caught anywhere inside a word, not just at its start: dumbfuck, sonofabitch. Only roots
// that no ordinary word contains are here; the few that do, like Scunthorpe, are in Allowed.
var Infixes = concat(
	tag(English, Lenient, []string{
		"fuck", "cunt", "bitch", "whore", "nigger", "faggot", "jizz",
	}),
)

// Allowed are ordinary words and names that a stem or an embedded root would otherwise flag. A word here is never
// flagged, with or without a postposition (Randipko, Shitijlai), at any strictness.
var Allowed = []string{
	"shitij", "shitiz", "shital", "shitala", "shitalpati", "shitanshu", "shiitake", "shitake", "shiite", "shiites",
	"shiitic", "niger", "nigeria", "nigerian", "nigerians", "nigerien", "snigger", "sniggers", "sniggered",
	"sniggering", "sniggerer", "scunthorpe", "randip", "randeep", "randhir", "randhawa", "random", "randomly",
	"randomness", "randomize", "randomized", "randy", "condition", "conditions", "conditional", "conditionally",
	"conditioner", "conditioning", "conditioned", "conduct", "conducts", "conducted", "conducting", "conductor",
	"conductors", "conduction", "conductive", "condense", "condensed", "condenser", "condemn", "condemned",
	"condolence", "condolences", "condiment", "kanda", "kandel", "kandu", "kandahar", "lundberg", "lundup",
}

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
	LatinInfixes      = texts(Infixes, true)
	DevanagariWords   = texts(Words, false)
	DevanagariStems   = texts(Stems, false)
	DevanagariPhrases = texts(Phrases, false)
)

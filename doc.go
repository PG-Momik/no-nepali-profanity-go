// Package nepaliprofanity finds and censors profanity in English, Romanized (Latin) Nepali and Devanagari Nepali,
// plus the Hindi slang common in Nepal.
//
// It is built for moderating user-written text, such as names, comments and reviews, on Nepali sites, where a false
// positive on a real name does more damage than a missed swear. It uses the same word lists and matching rules as
// the JavaScript package, so a piece of text gets the same result in every language.
//
//	nepaliprofanity.ContainsProfanity("Great teacher!")    // false
//	nepaliprofanity.FindProfanity("f.u.c.k this sh1t")    // [fuck shit]
//	nepaliprofanity.Censor("you muji")                     // "you ****"
//
//	result := nepaliprofanity.Check("you muji")           // scan once...
//	result.Words                                           // [muji]
//	result.Censor()                                        // ...then censor: "you ****"
//
// The package-level functions check all three languages at Standard strictness. For other options, build a Filter
// once and reuse it:
//
//	f, err := nepaliprofanity.NewFilter(nepaliprofanity.FilterOptions{
//		Languages:  []nepaliprofanity.Language{nepaliprofanity.Romanized},
//		Strictness: nepaliprofanity.Lenient,
//	})
//
// # What it catches
//
//   - Case and Unicode forms: IDIOT, full-width letters.
//   - Leetspeak: sh1t, @ss.
//   - "!" for "i" between letters, "*" for a hidden letter: sh!t, f*ck, *ss.
//   - Stretched and spelled-out letters: fuuuuck, f.u.c.k, f u c k.
//   - Nepali postpositions and plurals: mujiko, मुजीको, …हरू.
//   - Devanagari spelling variants: nukta, chandrabindu for anusvara, zero-width joiners.
//   - Multi-word phrases, in both scripts.
//
// Short words match exactly, so "class" and "Assam" are left alone, and names such as Shitij, Randip and Putali are
// in the test suite.
package nepaliprofanity

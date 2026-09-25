// Command example shows the main functions of the package.
package main

import (
	"fmt"
	"strings"

	nepaliprofanity "github.com/PG-Momik/no-nepali-profanity-go"
)

func main() {
	fmt.Println(nepaliprofanity.ContainsProfanity("Great teacher!")) // false
	fmt.Println(nepaliprofanity.ContainsProfanity("मुजीको कक्षा"))   // true
	fmt.Println(nepaliprofanity.FindProfanity("f.u.c.k this sh1t"))  // [fuck shit]
	fmt.Println(nepaliprofanity.FindProfanity("Randip Thapa"))       // []

	for _, m := range nepaliprofanity.FindProfanityMatches("hello muji world") {
		fmt.Printf("%q (%s) at bytes %d-%d\n", m.Text, m.Normalized, m.Start, m.End) // "muji" (muji) at bytes 6-10
	}

	fmt.Println(nepaliprofanity.Censor("you muji"))                                           // you ****
	fmt.Println(nepaliprofanity.Censor("you muji", nepaliprofanity.CensorOptions{Mask: "#"})) // you ####
	fmt.Println(nepaliprofanity.Censor("you muji", nepaliprofanity.CensorOptions{
		Replace: func(m nepaliprofanity.ProfanityMatch) string {
			return m.Text[:1] + strings.Repeat("*", len(m.Text)-1)
		},
	})) // you m***

	// Scan once, then inspect and censor the result.
	result := nepaliprofanity.Check("you muji")
	fmt.Println(result.HasProfanity, result.Words, result.Censor()) // true [muji] you ****

	// Build a filter once for other options.
	f := nepaliprofanity.MustNewFilter(nepaliprofanity.FilterOptions{
		Languages:  []nepaliprofanity.Language{nepaliprofanity.Romanized},
		Strictness: nepaliprofanity.Lenient,
	})
	fmt.Println(f.FindProfanity("fuck muji murkha")) // [muji]
}

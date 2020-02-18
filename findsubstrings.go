package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
)

/* This program searches standard input for patterns that are provided in a
 * file. It will list where each exact substring match is found in each line.
 */

type Args struct {
	Patterns string
	Help     bool
	K        int
}

var args = Args{}

type Seed struct {
	Pattern string
	Seed    string
	Offset  int
}

func init() {
	log.SetFlags(0)
	flag.StringVar(&args.Patterns, "patterns", "", "filename containing substrings to look for (required)")
	flag.BoolVar(&args.Help, "help", false, "show this help message")
	flag.IntVar(&args.K, "k", 8, "minimum substring length, used for finding match seeds")

	flag.Usage = func() {
		log.Println("usage: findsubstrings [options]")
		flag.PrintDefaults()
	}
}

func main() {
	flag.Parse()

	if len(os.Args) == 1 {
		flag.Usage()
		os.Exit(1)
	}
	if args.Help {
		flag.Usage()
		os.Exit(0)
	}

	if args.Patterns == "" {
		log.Println("must specify -patterns FILE parameter")
		flag.Usage()
		os.Exit(1)
	}

	// We'll store and count patterns using this map
	patterns := make(map[string]int, 0)

	// Read through the list of patterns
	fp, err := os.Open(args.Patterns)
	if err != nil {
		log.Fatalf("Failed to open list of patterns%s: %v", args.Patterns, err)
	}
	patternScanner := bufio.NewScanner(fp)
	lineNum := 0
	for patternScanner.Scan() {
		line := patternScanner.Text()
		lineNum++
		_, present := patterns[line]
		if present {
			log.Fatalf("Patterns must be unique. Already saw pattern %s", line)
		}
		// We require patterns to be at least k long, though it might be possible to relax this.
		if len(line) < args.K {
			log.Fatalf("pattern %s is too short, must at least k (%d) long", line, args.K)
		}
		patterns[line] = 0
	}

	// We build a data structure that stores k substrings for each pattern
	// consisting of the k-mers starting at offsets of [0:k]. This is to make
	// lookup of patterns quick. We thus only need to check len/k times for each string.
	seeds := make(map[string][]Seed, 0)

	seedCount := 0
	for pattern, _ := range patterns {
		for i := 0; i < args.K; i++ {
			if len(pattern) < i+args.K {
				break
			}
			substr := pattern[i:(i + args.K)]
			seed := Seed{
				Pattern: pattern,
				Offset:  i,
				Seed:    substr,
			}
			seeds[substr] = append(seeds[substr], seed)
			seedCount++
			//log.Println("Stored seed", seed.Seed)
		}
	}
	log.Println("stored", seedCount, "seeds")

	// Iterate over the input from stdin
	inputScanner := bufio.NewScanner(os.Stdin)
	lineNum = 0
	totalMatches := 0
	for inputScanner.Scan() {
		line := inputScanner.Text()
		//log.Printf("line %d: %s\n", lineNum, line)
		for offset := 0; offset+args.K <= len(line); offset += args.K {
			substr := line[offset:(offset + args.K)]
			//log.Printf("considering substring %s at offset %d\n", substr, offset)
			seedList, ok := seeds[substr]
			if ok {
				//log.Printf("found %d seeds\n", len(seedList))
				// Now check for the full patterns implied by each seed
				for _, seed := range seedList {
					patternStart := offset - seed.Offset
					// Make sure the pattern fits in the string at the location implied by the seed
					if patternStart >= 0 && len(line)-offset >= len(seed.Pattern)-seed.Offset {
						// Check if the string actually contains the exact pattern
						candidateMatch := line[patternStart:(patternStart + len(seed.Pattern))]
						//log.Printf("candidate match: %s\n", candidateMatch)
						if candidateMatch == seed.Pattern {
							// line number, offset into line, pattern, full line
							fmt.Printf("%d\t%d\t%s\t%s\n", lineNum, patternStart, seed.Pattern, line)
							patterns[seed.Pattern]++
							totalMatches++
						}
					}
				}
			}
		}
		lineNum++
	}

	log.Println("found", totalMatches, "matches")
	for p, n := range patterns {
		log.Println(p, n)
	}
}

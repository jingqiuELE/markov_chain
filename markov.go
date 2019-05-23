package main

import (
	"bufio"
	"fmt"
	"io"
	"math/rand"
	"os"
)

const NPREF = 2
const NHASH = 4093
const MAXGEN = 10000

type State struct {
	pref [2]string
	suf  *Suffix
	next *State
}

type Suffix struct {
	word string
	next *Suffix
}

var statetab [NHASH]*State

/* hash: compute hash value for array of NPREF strings using djb2 algorithm */
func hash(s [NPREF]string) uint {
	var h uint = 5381
	for i := 0; i < NPREF; i++ {
		for _, p := range s[i] {
			h = ((h << 5) + h) + uint(p)
		}
	}
	return h % NHASH
}

/* lookup: search for prefix; create if requested. */
/* returns pointer if present or created, nil if not. */
/* creation doesn't strdup so strings mustn't change later. */
func lookup(prefix [NPREF]string, create bool) *State {
	var sp *State
	var i int
	var p string
	var match bool = false

	h := hash(prefix)

	for sp = statetab[h]; sp != nil; sp = sp.next {
		match = true
		for i, p = range prefix {
			if p != sp.pref[i] {
				match = false
				break
			}
		}
		if match {
			return sp
		}
	}
	if create {
		sp = &State{}
		for i := range prefix {
			sp.pref[i] = prefix[i]
		}
		sp.next = statetab[h]
		statetab[h] = sp
	}
	return sp
}

/* addsuffix: add to state. */
func addsuffix(sp *State, suffix string) {
	suf := &Suffix{}
	suf.word = suffix
	suf.next = sp.suf
	sp.suf = suf
}

/* add: add word to suffix list, update prefix. */
func add(prefix *[NPREF]string, suffix string) {
	var sp *State
	sp = lookup(*prefix, true)
	addsuffix(sp, suffix)
	/* move the words down the prefix. */
	copy(prefix[:], prefix[1:])
	prefix[NPREF-1] = suffix
}

/* build: read input, build prefix table */
func build(prefix *[NPREF]string, reader io.Reader) {
	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		add(prefix, scanner.Text())
	}
}

const NONWORD = "@"

/* generate: produce output */
func generate(nwords int) {
	var sp *State
	var suf *Suffix
	var prefix [NPREF]string
	var w string

	for i := 0; i < NPREF; i++ {
		prefix[i] = NONWORD
	}
	for i := 0; i < nwords; i++ {
		sp = lookup(prefix, false)
		nmatch := 1
		for suf = sp.suf; suf != nil; suf = suf.next {
			if rand.Intn(255)%nmatch == 0 {
				w = suf.word
			}
			nmatch++
		}
		if w == NONWORD {
			break
		}

		fmt.Printf("%v ", w)
		copy(prefix[:], prefix[1:])
		prefix[NPREF-1] = w
	}
}

/* markov chain: markov-chain random text generation */
func main() {
	var prefix [NPREF]string
	input, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}

	for i := 0; i < NPREF; i++ {
		prefix[i] = NONWORD
	}
	build(&prefix, input)
	fmt.Println("build done")
	add(&prefix, NONWORD)
	generate(MAXGEN)
	return
}

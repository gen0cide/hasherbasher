package main

import (
	"math/rand"
	"unicode/utf8"
)

const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func stringWithCharset(length int, charset string, src *rand.Rand) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[src.Intn(len(charset))]
	}
	return string(b)
}

func randString(length int, src *rand.Rand) string {
	return stringWithCharset(length, charset, src)
}

func locateInjection(s []byte) (end int) {
	end = -1
	var r rune
	var rlen int
	i := 0
	lazy := false
	type jmp struct{ s, i int }
	var lazyArr [3]jmp
	lazyStack := lazyArr[:0]
	_, _, _ = r, rlen, i
	switch {
	case i == 0:
		goto s2
	}
	goto bt
s2:
	if lazy {
		lazy = false
		goto s3
	}
	lazyStack = append(lazyStack, jmp{s: 2, i: i})
	r, rlen = utf8.DecodeRune(s[i:])
	if rlen == 0 {
		goto bt
	}
	i += rlen
	switch {
	case r == 39:
		goto s5
	}
	goto bt
s3:
	r, rlen = utf8.DecodeRune(s[i:])
	if rlen == 0 {
		goto bt
	}
	i += rlen
	switch {
	case r <= 9 || r >= 11:
		goto s4
	}
	goto bt
s4:
	if lazy {
		lazy = false
		goto s3
	}
	lazyStack = append(lazyStack, jmp{s: 4, i: i})
	r, rlen = utf8.DecodeRune(s[i:])
	if rlen == 0 {
		goto bt
	}
	i += rlen
	switch {
	case r == 39:
		goto s5
	}
	goto bt
s5:
	r, rlen = utf8.DecodeRune(s[i:])
	if rlen == 0 {
		goto bt
	}
	i += rlen
	switch {
	case r == 111:
		goto s13
	case r == 124:
		goto s16
	case r == 79:
		goto s6
	}
	goto bt
s6:
	r, rlen = utf8.DecodeRune(s[i:])
	if rlen == 0 {
		goto bt
	}
	i += rlen
	switch {
	case r == 82 || r == 114:
		goto s7
	}
	goto bt
s7:
	r, rlen = utf8.DecodeRune(s[i:])
	if rlen == 0 {
		goto bt
	}
	i += rlen
	switch {
	case r == 39:
		goto s8
	}
	goto bt
s8:
	r, rlen = utf8.DecodeRune(s[i:])
	if rlen == 0 {
		goto bt
	}
	i += rlen
	switch {
	case r >= 49 && r <= 57:
		goto s9
	}
	goto bt
s9:
	if lazy {
		lazy = false
		goto s10
	}
	lazyStack = append(lazyStack, jmp{s: 9, i: i})
	switch {
	case i == len(s):
		end = i
		goto bt
	}
	r, rlen = utf8.DecodeRune(s[i:])
	if rlen == 0 {
		goto bt
	}
	i += rlen
	switch {
	case r <= 9 || r >= 11:
		goto s12
	}
	goto bt
s10:
	r, rlen = utf8.DecodeRune(s[i:])
	if rlen == 0 {
		goto bt
	}
	i += rlen
	switch {
	case r >= 49 && r <= 57:
		goto s9
	}
	goto bt
s12:
	switch {
	case i == len(s):
		end = i
		goto bt
	}
	r, rlen = utf8.DecodeRune(s[i:])
	if rlen == 0 {
		goto bt
	}
	i += rlen
	switch {
	case r <= 9 || r >= 11:
		goto s12
	}
	goto bt
s13:
	r, rlen = utf8.DecodeRune(s[i:])
	if rlen == 0 {
		goto bt
	}
	i += rlen
	switch {
	case r == 82:
		goto s14
	case r == 114:
		goto s15
	}
	goto bt
s14:
	r, rlen = utf8.DecodeRune(s[i:])
	if rlen == 0 {
		goto bt
	}
	i += rlen
	switch {
	case r == 39:
		goto s8
	}
	goto bt
s15:
	r, rlen = utf8.DecodeRune(s[i:])
	if rlen == 0 {
		goto bt
	}
	i += rlen
	switch {
	case r == 39:
		goto s8
	}
	goto bt
s16:
	r, rlen = utf8.DecodeRune(s[i:])
	if rlen == 0 {
		goto bt
	}
	i += rlen
	switch {
	case r == 124:
		goto s17
	}
	goto bt
s17:
	r, rlen = utf8.DecodeRune(s[i:])
	if rlen == 0 {
		goto bt
	}
	i += rlen
	switch {
	case r == 39:
		goto s8
	}
bt:
	if end >= 0 || len(lazyStack) == 0 {
		return
	}
	var to jmp
	to, lazyStack = lazyStack[len(lazyStack)-1], lazyStack[:len(lazyStack)-1]
	lazy = true
	i = to.i
	switch to.s {
	case 2:
		goto s2
	case 4:
		goto s4
	case 9:
		goto s9
	}
	return
}

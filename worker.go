package main

import (
	"crypto/md5"
	"math/rand"
	"sync/atomic"
	"time"
)

// Result is a result type
type Result struct {
	Input  string
	Output string
}

func worker(res chan Result, counter *int64, minLength, maxLength int, finished *bool) {
	var randomSource = rand.New(rand.NewSource(time.Now().UnixNano()))
	var matchedString string
	var inputString string
	md5hasher := md5.New()
	delta := maxLength - minLength
	var randFunc func() string

	if delta > 0 {
		randFunc = func() string {
			return randString(minLength+randomSource.Intn(maxLength-minLength), randomSource)
		}
	} else {
		randFunc = func() string {
			return randString(minLength, randomSource)
		}
	}

	for {
		if *finished {
			return
		}
		md5hasher.Reset()
		inputString = randFunc()
		md5hasher.Write([]byte(inputString))
		result := md5hasher.Sum(nil)
		if locateInjection(result) != -1 {
			matchedString = string(result)
			resultObj := Result{
				Input:  inputString,
				Output: matchedString,
			}
			*finished = true
			res <- resultObj
			return
		}
		atomic.AddInt64(counter, 1)
	}
}

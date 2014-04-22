// Copyright 2014 Ryan Rogers. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fuzzers

import (
	"fmt"
	"strconv"
)

const (
	incrementFuzzerName           = "Increment"
	incrementFuzzerRegisteredName = "increment"
)

type IncrementFuzzer struct {
	finished          bool
	start, stop, step int64
	defaults          struct {
		start, stop, step int64
	}
}

func init() {
	register(incrementFuzzerRegisteredName, NewIncrementFuzzer)
}

func NewIncrementFuzzer() Fuzzer {
	return &IncrementFuzzer{}
}

func (a *IncrementFuzzer) Type() string {
	return incrementFuzzerRegisteredName
}

func (a *IncrementFuzzer) Configure(config interface{}) (err error) {
	c, ok := config.(map[string]interface{})
	if !ok {
		return fmt.Errorf(errConfigureWrongType, incrementFuzzerName, "config")
	}

	var start, stop, step int64

	switch i := c["start"].(type) {
	case int:
		start = int64(i)
	case float64:
		start = int64(i)
	default:
		return fmt.Errorf(errConfigureWrongType, incrementFuzzerName, "config['start']")
	}

	switch i := c["stop"].(type) {
	case int:
		stop = int64(i)
	case float64:
		stop = int64(i)
	default:
		return fmt.Errorf(errConfigureWrongType, incrementFuzzerName, "config['stop']")
	}

	if start == stop {
		return fmt.Errorf("%s Configure: start and stop can not be the same.", incrementFuzzerName)
	}

	switch i := c["step"].(type) {
	case int:
		step = int64(i)
	case float64:
		step = int64(i)
	default:
		return fmt.Errorf(errConfigureWrongType, incrementFuzzerName, "config['step']")
	}
	if step == 0 {
		return fmt.Errorf("%s Configure: step can not be zero.", incrementFuzzerName)
	}

	a.start = start
	a.stop = stop + step
	a.step = step

	a.defaults.start = a.start
	a.defaults.stop = a.stop
	a.defaults.step = a.step

	return
}

func (a *IncrementFuzzer) Reset() {
	a.finished = false
	a.start = a.defaults.start
	a.stop = a.defaults.stop
	a.step = a.defaults.step
}

func (a *IncrementFuzzer) Next() (value string, finished bool) {
	// FIXME: Verify that Configure() was successfully called?
	finished = a.finished
	if !a.finished {
		value = strconv.FormatInt(a.start, 10)
		a.start += a.step
		if a.start == a.stop {
			a.finished = true
		}
	}
	return
}

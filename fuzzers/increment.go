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

func (f *IncrementFuzzer) Type() string {
	return incrementFuzzerRegisteredName
}

func (f *IncrementFuzzer) Configure(config interface{}) (err error) {
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

	f.start = start
	f.stop = stop + step
	f.step = step

	f.defaults.start = f.start
	f.defaults.stop = f.stop
	f.defaults.step = f.step

	return
}

func (f *IncrementFuzzer) Reset() {
	f.finished = false
	f.start = f.defaults.start
	f.stop = f.defaults.stop
	f.step = f.defaults.step
}

func (f *IncrementFuzzer) Next() (value string, finished bool) {
	// FIXME: Verify that Configure() was successfully called?
	finished = f.finished
	if !f.finished {
		value = strconv.FormatInt(f.start, 10)
		f.start += f.step
		if f.start == f.stop {
			f.finished = true
		}
	}
	return
}

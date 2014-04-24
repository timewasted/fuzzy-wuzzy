// Copyright 2014 Ryan Rogers. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fuzzers

import "fmt"

const (
	valueFuzzerName           = "ValueFuzzer"
	valueFuzzerRegisteredName = "value"
)

type ValueFuzzer struct {
	finished bool
	value    string
}

func init() {
	register(valueFuzzerRegisteredName, NewValueFuzzer)
}

func NewValueFuzzer() Fuzzer {
	return &ValueFuzzer{}
}

func (f *ValueFuzzer) Type() string {
	return valueFuzzerRegisteredName
}

func (f *ValueFuzzer) Configure(config interface{}) (err error) {
	value, ok := config.(string)
	if !ok {
		return fmt.Errorf(errConfigureWrongType, valueFuzzerName, "config")
	}
	f.value = value

	return
}

func (f *ValueFuzzer) Reset() {
	f.finished = false
}

func (f *ValueFuzzer) Next() (value string, finished bool) {
	// FIXME: Verify that Configure() was successfully called?
	finished = f.finished
	if !f.finished {
		f.finished = true
		value = f.value
	}
	return
}

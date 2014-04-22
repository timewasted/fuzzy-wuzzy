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

func (a *ValueFuzzer) Type() string {
	return valueFuzzerRegisteredName
}

func (a *ValueFuzzer) Configure(config interface{}) (err error) {
	value, ok := config.(string)
	if !ok {
		return fmt.Errorf(errConfigureWrongType, valueFuzzerName, "config")
	}
	a.value = value

	return
}

func (a *ValueFuzzer) Reset() {
	a.finished = false
}

func (a *ValueFuzzer) Next() (value string, finished bool) {
	// FIXME: Verify that Configure() was successfully called?
	finished = a.finished
	if !a.finished {
		a.finished = true
		value = a.value
	}
	return
}

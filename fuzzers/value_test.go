// Copyright 2014 Ryan Rogers. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fuzzers

import (
	"testing"
)

func TestValue_Register(t *testing.T) {
	if _, registered := registeredFuzzers[valueFuzzerRegisteredName]; !registered {
		t.Errorf("%s is not a registered fuzzer.", valueFuzzerRegisteredName)
	}
}

func TestValue_Configure(t *testing.T) {
	var err error
	f := NewValueFuzzer()

	if err = f.Configure(1); err == nil {
		t.Error("Expected Configure to fail due to invalid type.")
	}
	if err = f.Configure("valid"); err != nil {
		t.Errorf("Expected Configure to succeed, instead received error '%s'.", err)
	}
}

func TestValue_Reset(t *testing.T) {
	f := NewValueFuzzer()
	f.Configure("valid")
	if _, finished := f.Next(); finished {
		t.Fatal("Expected Next to not be finished.")
	}
	if _, finished := f.Next(); !finished {
		t.Fatal("Expected Next to be finished.")
	}
	f.Reset()
	if _, finished := f.Next(); finished {
		t.Fatal("Expected Next to not be finished.")
	}
}

func TestValue_Next(t *testing.T) {
	var input = "valid"
	var output = []string{}

	f := NewValueFuzzer()
	f.Configure(input)
	for {
		value, finished := f.Next()
		if finished {
			break
		}
		output = append(output, value)
	}
	if len(output) != 1 {
		t.Fatalf("Expected Next to output 1 line, instead received %d.", len(output))
	}
	if output[0] != input {
		t.Errorf("Expected Next to return '%s', instead received '%s'.", input, output[0])
	}
}

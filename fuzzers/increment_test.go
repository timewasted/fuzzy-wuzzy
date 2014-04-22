// Copyright 2014 Ryan Rogers. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fuzzers

import (
	"testing"
)

func TestIncrement_Register(t *testing.T) {
	if _, registered := registeredFuzzers[incrementFuzzerRegisteredName]; !registered {
		t.Errorf("%s is not a registered fuzzer.", incrementFuzzerRegisteredName)
	}
}

func TestIncrement_Configure(t *testing.T) {
	f := NewIncrementFuzzer()

	// Valid data.
	if err := f.Configure(map[string]interface{}{
		"start": 0,
		"stop":  10,
		"step":  1,
	}); err != nil {
		t.Errorf("Expected Configure to succeed, instead received error '%s'.", err)
	}

	// Invalid start type.
	if err := f.Configure(map[string]interface{}{
		"start": "0",
		"stop":  10,
		"step":  1,
	}); err == nil {
		t.Errorf("Expected Configure to fail due to invalid start type.")
	}

	// Invalid stop type.
	if err := f.Configure(map[string]interface{}{
		"start": 0,
		"stop":  "10",
		"step":  1,
	}); err == nil {
		t.Errorf("Expected Configure to fail due to invalid stop type.")
	}

	// Same start and stop values.
	if err := f.Configure(map[string]interface{}{
		"start": 0,
		"stop":  0,
		"step":  1,
	}); err == nil {
		t.Errorf("Expected Configure to fail due to same start and stop values.")
	}

	// Invalid step type.
	err := f.Configure(map[string]interface{}{
		"start": 0,
		"stop":  10,
		"step":  "1",
	})
	if err == nil {
		t.Errorf("Expected Configure to fail due to invalid step type.")
	}

	// Invalid step value.
	if err := f.Configure(map[string]interface{}{
		"start": 0,
		"stop":  10,
		"step":  0,
	}); err == nil {
		t.Errorf("Expected Configure to fail due to invalid step value.")
	}
}

func TestIncrement_Reset(t *testing.T) {
	var expected = []string{"0", "2", "4", "6", "8", "10"}
	var output = []string{}

	f := NewIncrementFuzzer()
	f.Configure(map[string]interface{}{
		"start": 0,
		"stop":  10,
		"step":  2,
	})

	for {
		value, finished := f.Next()
		if finished {
			break
		}
		output = append(output, value)
	}
	if len(output) != len(expected) {
		t.Fatalf("Expected Next to return %d elements, instead received %d.", len(expected), len(output))
	}
	for pos := range output {
		if output[pos] != expected[pos] {
			t.Errorf("Expected Next to return '%s' in position %d, instead received '%s'.", expected[pos], pos, output[pos])
		}
	}

	f.Reset()
	output = []string{}
	for {
		value, finished := f.Next()
		if finished {
			break
		}
		output = append(output, value)
	}
	if len(output) != len(expected) {
		t.Fatalf("Expected Next to return %d elements, instead received %d.", len(expected), len(output))
	}
	for pos := range output {
		if output[pos] != expected[pos] {
			t.Errorf("Expected Next to return '%s' in position %d, instead received '%s'.", expected[pos], pos, output[pos])
		}
	}
}

func TestIncrement_Next(t *testing.T) {
	var expected = []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "10"}
	var output = []string{}

	f := NewIncrementFuzzer()
	f.Configure(map[string]interface{}{
		"start": 0,
		"stop":  10,
		"step":  1,
	})
	for {
		value, finished := f.Next()
		if finished {
			break
		}
		output = append(output, value)
	}

	if len(output) != len(expected) {
		t.Fatalf("Expected Next to return %d elements, instead received %d.", len(expected), len(output))
	}
	for pos := range output {
		if output[pos] != expected[pos] {
			t.Errorf("Expected Next to return '%s' in position %d, instead received '%s'.", expected[pos], pos, output[pos])
		}
	}
}

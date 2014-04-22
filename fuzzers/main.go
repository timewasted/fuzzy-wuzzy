// Copyright 2014 Ryan Rogers. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fuzzers

const (
	errConfigureWrongType = "%s Configure: received invalid type for '%s'."
)

type Fuzzer interface {
	Type() string
	Configure(config interface{}) (err error)
	Reset()
	Next() (value string, finished bool)
}

type newFuzzerFunction func() Fuzzer
type registeredFuzzersList map[string]newFuzzerFunction

var registeredFuzzers = make(registeredFuzzersList)

func register(name string, newFn newFuzzerFunction) {
	registeredFuzzers[name] = newFn
}

func Registered() registeredFuzzersList {
	return registeredFuzzers
}

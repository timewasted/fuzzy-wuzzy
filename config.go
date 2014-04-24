// Copyright 2014 Ryan Rogers. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/timewasted/fuzzy-wuzzy/fuzzers"
)

const (
	errNoFuzzersSpecified = "config: there are no parameters to be fuzzed."
	errInvalidHost        = "config: host '%s' is invalid."
	errInvalidPort        = "config: port '%d' is invalid."
	errInvalidFuzzer      = "config: %s fuzzer '%s' for parameter '%s' is not recognized."
	errInvalidFuzzerCount = "config: %s fuzzer has invalid number of operations for parameter '%s'."
)

type fuzzerConfiguration []map[string]interface{}
type parameterList map[string]fuzzerConfiguration

type configuration struct {
	Host         string        `json:"host"`
	Port         uint          `json:"port"`
	TLS          bool          `json:"tls"`
	Path         string        `json:"path"`
	PathParams   parameterList `json:"path-params"`
	URLParams    parameterList `json:"url-params"`
	HeaderParams parameterList `json:"header-params"`
	CookieParams parameterList `json:"cookie-params"`
}

type paramFuzzer struct {
	paramType, param string
	fuzzers.Fuzzer
}

var paramFuzzers []paramFuzzer

func loadConfig(filePath string) (config *configuration, err error) {
	file, err := os.Open(filePath)
	if err != nil {
		return
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	if err = decoder.Decode(&config); err != nil {
		return
	}

	if err = validateConfig(config); err != nil {
		return
	}

	return
}

func validateConfig(config *configuration) (err error) {
	if len(config.Host) > 0 && config.Host[len(config.Host)-1:] == "/" {
		config.Host = config.Host[:len(config.Host)-1]
	}
	if len(config.Host) == 0 {
		return fmt.Errorf(errInvalidHost, config.Host)
	}

	if config.Port == 0 {
		if config.TLS {
			config.Port = 443
		} else {
			config.Port = 80
		}
	}
	if config.Port == 0 || config.Port > 65535 {
		return fmt.Errorf(errInvalidPort, config.Port)
	}

	if config.Path == "" {
		config.Path = "/"
	}

	if err = validateFuzzers("path", config.PathParams); err != nil {
		return
	}
	if err = validateFuzzers("url", config.URLParams); err != nil {
		return
	}
	if err = validateFuzzers("header", config.HeaderParams); err != nil {
		return
	}
	if err = validateFuzzers("cookie", config.CookieParams); err != nil {
		return
	}

	if len(paramFuzzers) == 0 {
		return fmt.Errorf(errNoFuzzersSpecified)
	}

	return
}

func validateFuzzers(paramType string, params parameterList) (err error) {
	validFuzzers := fuzzers.Registered()
	for param, fuzzerList := range params {
		for _, fuzzerConfig := range fuzzerList {
			if len(fuzzerConfig) > 0 && len(fuzzerConfig) != 1 {
				return fmt.Errorf(errInvalidFuzzerCount, paramType, param)
			}
			for name, config := range fuzzerConfig {
				name = strings.ToLower(name)
				newFuzzerFn, valid := validFuzzers[name]
				if !valid {
					return fmt.Errorf(errInvalidFuzzer, paramType, name, param)
				}

				fuzzer := newFuzzerFn()
				if err = fuzzer.Configure(config); err != nil {
					return
				}

				paramFuzzers = append(paramFuzzers, paramFuzzer{
					paramType: paramType,
					param:     string(param),
					Fuzzer:    fuzzer,
				})
			}
		}
	}

	return
}

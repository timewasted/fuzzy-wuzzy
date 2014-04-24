// Copyright 2014 Ryan Rogers. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"errors"
	"flag"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

// NOTE: 6 was chosen because in theory it would allow performing 2 requests
// per URL/Header/Cookie parameter at a time.
const defaultConcurrency = 6

type requestParameters map[string]map[string]string

type httpResponse struct {
	*http.Request
	*http.Response
	err error
}

var (
	configFile  = flag.String("config", "./config.json", "Path to the fuzzer configuration file.")
	concurrency = flag.Int("concurrency", defaultConcurrency, "Maximum number of concurrent requests to make.")
)

var httpClient = &http.Client{
	Transport: &http.Transport{},
}

var errExitingEarly = errors.New("Exiting.")

func main() {
	config, err := loadConfig(*configFile)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Configure the HTTP client.
	jar, err := cookiejar.New(nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	httpClient.Jar = jar

	if err := doFuzzing(config); err != nil {
		fmt.Println(err)
		return
	}
}

func doFuzzing(config *configuration) error {
	done := make(chan struct{})
	defer close(done)

	requests, errc := queueRequests(config, done)

	responses := make(chan httpResponse)
	var wg sync.WaitGroup
	wg.Add(*concurrency)
	for i := 0; i < *concurrency; i++ {
		go func() {
			processRequests(done, requests, responses)
			wg.Done()
		}()
	}
	go func() {
		wg.Wait()
		close(responses)
	}()
	go func() {
		// FIXME: Is there a cleaner way to do this?
		osSignals := make(chan os.Signal)
		signal.Notify(osSignals, os.Interrupt, syscall.SIGQUIT, syscall.SIGTERM)
		select {
		case <-osSignals:
			done <- struct{}{}
		}
	}()

	for response := range responses {
		if response.err != nil {
			return response.err
		}
		// FIXME: Actually do something with the results.
	}

	if err := <-errc; err != nil {
		return err
	}

	return nil
}

func queueRequests(config *configuration, done <-chan struct{}) (<-chan *http.Request, <-chan error) {
	requests := make(chan *http.Request)
	errc := make(chan error, 1)

	go func() {
		defer close(requests)

		// Build the basic URL.
		reqUrl := &url.URL{
			Path: config.Path,
		}
		if !config.TLS {
			reqUrl.Scheme = "http"
		} else {
			reqUrl.Scheme = "https"
		}
		if (!config.TLS && config.Port == 80) || (config.TLS && config.Port == 443) {
			reqUrl.Host = config.Host
		} else {
			reqUrl.Host = fmt.Sprintf("%s:%d", config.Host, config.Port)
		}

		var reqParams = requestParameters{
			"url":    map[string]string{},
			"header": map[string]string{},
			"cookie": map[string]string{},
		}

		// Build and send the initial request.
		var prevParam string
		for _, fuzzer := range paramFuzzers {
			if prevParam == fuzzer.param {
				continue
			}
			prevParam = fuzzer.param

			value, _ := fuzzer.Next()
			reqParams[fuzzer.paramType][fuzzer.param] = value
		}
		request, err := newRequest(reqUrl, reqParams)
		if err != nil {
			errc <- err
			return
		}
		select {
		case requests <- request:
		case <-done:
			errc <- errExitingEarly
			return
		}

		// Build and send the subsequent requests.
		var fuzzerIndex int
		for {
			if fuzzerIndex == len(paramFuzzers) {
				break
			}
			fuzzer := paramFuzzers[fuzzerIndex]

			value, finished := fuzzer.Next()
			if finished {
				fuzzer.Reset()
				value, _ := fuzzer.Next()
				reqParams[fuzzer.paramType][fuzzer.param] = value

				fuzzerIndex++
				continue
			}
			reqParams[fuzzer.paramType][fuzzer.param] = value
			request, err = newRequest(reqUrl, reqParams)
			if err != nil {
				errc <- err
				return
			}
			select {
			case requests <- request:
			case <-done:
				errc <- errExitingEarly
				return
			}
		}

		errc <- nil
	}()

	return requests, errc
}

func newRequest(reqUrl *url.URL, reqParams requestParameters) (req *http.Request, err error) {
	urlValues := url.Values{}
	for key, value := range reqParams["url"] {
		urlValues.Add(key, value)
	}
	reqUrl.RawQuery = urlValues.Encode()
	req, err = http.NewRequest("GET", reqUrl.String(), nil)
	if err != nil {
		return
	}

	for key, value := range reqParams["header"] {
		req.Header.Add(key, value)
	}

	for key, value := range reqParams["cookie"] {
		req.AddCookie(&http.Cookie{
			Name:   key,
			Value:  value,
			Path:   reqUrl.Path,
			Domain: reqUrl.Host,
		})
	}

	return
}

func processRequests(done <-chan struct{}, requests <-chan *http.Request, responses chan<- httpResponse) {
	for request := range requests {
		response, err := httpClient.Do(request)
		select {
		case responses <- httpResponse{request, response, err}:
		case <-done:
			return
		}
	}
}

/*
 * Copyright 2019 Rackspace US, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package restclient_test

import (
	"encoding/xml"
	"fmt"
	"github.com/racker/go-restclient"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
)

//noinspection GoUnhandledErrorResult
func Example_post() {
	// Setup a test HTTP server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bytes, _ := ioutil.ReadAll(r.Body)
		fmt.Printf("RECV BODY %s", string(bytes))
		fmt.Fprintf(w, `{"Msg":"greetings via %s"}`, r.URL.Path)
	}))
	defer ts.Close()

	// Real example starts here
	client := restclient.New()
	client.SetBaseUrl(ts.URL)

	type MsgHolder struct {
		Msg string
	}

	req := &MsgHolder{Msg: "hello"}
	var resp MsgHolder

	err := client.Exchange("POST", "/ping", nil,
		restclient.NewJsonEntity(req), restclient.NewJsonEntity(&resp))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(resp.Msg)
	// Output:
	// RECV BODY {"Msg":"hello"}
	// greetings via /ping
}

//noinspection GoUnhandledErrorResult
func Example_getWithQuery() {
	// Setup a test HTTP server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("RECV QUERY %s\n", r.URL.RawQuery)
		fmt.Fprint(w, `{"Rows":["green apple","mandarin orange"]}`)
	}))
	defer ts.Close()

	// Real example starts here
	client := restclient.New()
	client.SetBaseUrl(ts.URL)

	query := make(url.Values)
	query.Set("q", "select all")
	query.Add("filter", "apple")
	query.Add("filter", "orange")

	type Resp struct {
		Rows []string
	}
	var resp Resp

	err := client.Exchange("GET", "/db", query, nil,
		restclient.NewJsonEntity(&resp))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%+v\n", resp)
	// Output:
	// RECV QUERY filter=apple&filter=orange&q=select+all
	// {Rows:[green apple mandarin orange]}
}

//noinspection GoUnhandledErrorResult
func Example_postText() {
	// Setup a test HTTP server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bytes, _ := ioutil.ReadAll(r.Body)
		fmt.Printf("RECV BODY %s\n", string(bytes))
		fmt.Printf("RECV TYPE %s\n", r.Header.Get("Content-Type"))
	}))
	defer ts.Close()

	// Real example starts here
	client := restclient.New()
	client.SetBaseUrl(ts.URL)

	err := client.Exchange("POST", "/ingest", nil,
		restclient.NewTextEntity("some plain text"), nil)
	if err != nil {
		log.Fatal(err)
	}

	// Output:
	// RECV BODY some plain text
	// RECV TYPE text/plain
}

//noinspection GoUnhandledErrorResult
func Example_interceptorSetHeader() {
	// Setup a test HTTP server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("RECV HEADER %s\n", r.Header.Get("x-req-id"))
	}))
	defer ts.Close()

	// Real example starts here
	client := restclient.New()
	client.SetBaseUrl(ts.URL)
	client.AddInterceptor(func(req *http.Request, next restclient.NextCallback) (*http.Response, error) {
		req.Header.Set("x-req-id", "123")
		return next(req)
	})

	type Event struct {
		Msg string
	}
	event := &Event{Msg: "first"}

	err := client.Exchange("POST", "/events", nil,
		restclient.NewJsonEntity(&event), nil)
	if err != nil {
		log.Fatal(err)
	}

	// Output:
	// RECV HEADER 123
}

//noinspection GoUnhandledErrorResult
func Example_interceptorLogging() {
	// Setup a test HTTP server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// copy request body back to response body
		io.Copy(w, r.Body)
	}))
	defer ts.Close()

	// Real example starts here
	client := restclient.New()
	client.SetBaseUrl(ts.URL)
	client.AddInterceptor(func(req *http.Request, next restclient.NextCallback) (*http.Response, error) {
		fmt.Printf("OUT %s %s\n", req.Method, req.URL.Path)
		response, err := next(req)
		if err == nil {
			fmt.Printf("IN %s\n", response.Status)
		}
		return response, err
	})

	resp := restclient.NewTextEntity("")
	err := client.Exchange("POST", "/ping", nil,
		restclient.NewTextEntity("greetings"), resp)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(resp.Content)
	// Output:
	// OUT POST /ping
	// IN 200 OK
	// greetings
}

//noinspection GoUnhandledErrorResult
func Example_externalEncoding() {
	// Setup a test HTTP server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Content-Type = %s\n", r.Header.Get("Content-Type"))
		fmt.Printf("Accept = %s\n", r.Header.Get("Accept"))
		// copy request body back to response body
		io.Copy(w, r.Body)
	}))
	defer ts.Close()

	// Real example starts here
	client := restclient.New()
	client.SetBaseUrl(ts.URL)

	type MsgHolder struct {
		Msg string
	}

	reqBytes, err := xml.Marshal(&MsgHolder{Msg: "hello"})
	if err != nil {
		log.Fatal(err)
	}

	req := &restclient.Entity{
		ContentType: "application/xml",
		Content:     reqBytes,
	}
	resp := &restclient.Entity{
		ContentType: "application/xml",
		Content:     []byte{}, // placeholder to convey type
	}
	err = client.Exchange("POST", "/ping", nil, req, resp)
	if err != nil {
		log.Fatal(err)
	}

	var respMsg MsgHolder
	respContentBytes := resp.Content.([]byte)
	err = xml.Unmarshal(respContentBytes, &respMsg)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(respContentBytes))
	fmt.Println(respMsg.Msg)
	// Output:
	// Content-Type = application/xml
	// Accept = application/xml
	// <MsgHolder><Msg>hello</Msg></MsgHolder>
	// hello
}

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

/*
Package restclient provides a higher-order type that simplifies calling REST APIs.

It wraps the standard http and json packages to provide the following features:
JSON encoding of request entity,
JSON parsing of response entity,
conversion of non-2xx responses into a structured error type including response body extraction,
timeout management,
easier request URL building given a client-wide base URL and query values,
auto-closing of the response body,
interceptors to inject authentication tokens, etc

Quick example

	client := restclient.NewClient()
	client.SetBaseUrl("http://your.own.domain")

	type MsgHolder struct {
		Msg string
	}

	req := &MsgHolder{Msg:"hello"}
	var resp MsgHolder

	client.Exchange("POST", "/ping", nil,
		restclient.NewJsonEntity(req), restclient.NewJsonEntity(&resp))

  	fmt.Println(resp.Msg)
*/
package restclient

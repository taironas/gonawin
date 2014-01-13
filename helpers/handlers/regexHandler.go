/*
 * Copyright (c) 2013 Santiago Arias | Remy Jourde | Carlos Bernal
 *
 * Permission to use, copy, modify, and distribute this software for any
 * purpose with or without fee is hereby granted, provided that the above
 * copyright notice and this permission notice appear in all copies.
 *
 * THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES
 * WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF
 * MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR
 * ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES
 * WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN
 * ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF
 * OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.
 */

package handlers

import (
	"net/http"
	"regexp"
)

// inspired by the following sources with some small changes:
//http://stackoverflow.com/questions/6564558/wildcards-in-the-pattern-for-http-handlefunc
//https://github.com/raymi/quickerreference
type route struct {
    pattern *regexp.Regexp
    handler http.Handler
}

type RegexpHandler struct {
    routes []*route
}


// Handler that appends a new pattern, handler pair to the RegexpHandler routes.
func (h *RegexpHandler) Handler(pattern *regexp.Regexp, handler http.Handler) {
    h.routes = append(h.routes, &route{pattern, handler})
}

// main handler function used, it encapsulate string pattern start and end.
func (h *RegexpHandler) HandleFunc(strPattern string, handler func(http.ResponseWriter, *http.Request)) {
	// encapsulate string pattern with start and end constraints
	// so that HandleFunc would work as for Python GAE
	pattern := regexp.MustCompile("^"+strPattern+"$")
	h.routes = append(h.routes, &route{pattern, http.HandlerFunc(handler)})
}

// looks for a matching route among the regexpHandler routes returns 404 if no match is found
func (h *RegexpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, route := range h.routes {
			if route.pattern.MatchString(r.URL.Path) {
					route.handler.ServeHTTP(w, r)
					return
			}
	}
	// no pattern matched; send 404 response
	http.NotFound(w, r)
}

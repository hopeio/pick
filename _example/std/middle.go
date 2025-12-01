/*
 * Copyright 2024 hopeio. All rights reserved.
 * Licensed under the MIT License that can be found in the LICENSE file.
 * @Created by jyb
 */

package main

import (
	"log"
	"net/http"
)

func Log(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Log", r.Method, r.RequestURI)
		handler.ServeHTTP(w, r)
		log.Println("Log End", r.Method, r.RequestURI)
	})
}

func Log2(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Log", r.Method, r.RequestURI)
		handler.ServeHTTP(w, r)
		log.Println("Log End", r.Method, r.RequestURI)
	})
}

func Log3(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Log", r.Method, r.RequestURI)
		handler.ServeHTTP(w, r)
		log.Println("Log End", r.Method, r.RequestURI)
	})
}

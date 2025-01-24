/*
 * Copyright 2024 hopeio. All rights reserved.
 * Licensed under the MIT License that can be found in the LICENSE file.
 * @Created by jyb
 */

package pick

import (
	"net/http"
)

type Writer struct {
	http.ResponseWriter
}

func (w Writer) Status(code int) {
	w.ResponseWriter.WriteHeader(code)
}

func (w Writer) Set(k, v string) {
	w.ResponseWriter.Header().Set(k, v)
}

func (w Writer) Write(p []byte) (int, error) {
	return w.ResponseWriter.Write(p)
}

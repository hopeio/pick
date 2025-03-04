/*
 * Copyright 2024 hopeio. All rights reserved.
 * Licensed under the MIT License that can be found in the LICENSE file.
 * @Created by jyb
 */

package pick

import (
	httpi "github.com/hopeio/utils/net/http"
	"net/http"
)

type Writer struct {
	http.ResponseWriter
}

func (w Writer) Status(code int) {
	w.ResponseWriter.WriteHeader(code)
}

func (w Writer) Header() httpi.Header {
	return httpi.HttpHeader(w.ResponseWriter.Header())
}

func (w Writer) Write(p []byte) (int, error) {
	return w.ResponseWriter.Write(p)
}

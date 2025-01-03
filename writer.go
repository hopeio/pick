/*
 * Copyright 2024 hopeio. All rights reserved.
 * Licensed under the MIT License that can be found in the LICENSE file.
 * @Created by jyb
 */

package pick

import (
	"encoding/json"
	"github.com/hopeio/utils/errors/errcode"

	"github.com/hopeio/utils/log"
	httpi "github.com/hopeio/utils/net/http"
	http_fs "github.com/hopeio/utils/net/http/fs"
	"go.uber.org/zap"
	"io"
	"net/http"
	"reflect"
)

func ResWriteReflect(w http.ResponseWriter, traceId string, result []reflect.Value) {
	if !result[1].IsNil() {
		err := errcode.ErrHandle(result[1].Interface())
		log.Errorw(err.Error(), zap.String(log.FieldTraceId, traceId))
		json.NewEncoder(w).Encode(err)
		return
	}
	if info, ok := result[0].Interface().(*http_fs.File); ok {
		header := w.Header()
		header.Set(httpi.HeaderContentType, httpi.ContentTypeOctetStream)
		header.Set(httpi.HeaderContentDisposition, "attachment;filename="+info.Name)
		io.Copy(w, info.File)
		if flusher, canFlush := w.(http.Flusher); canFlush {
			flusher.Flush()
		}
		info.File.Close()
		return
	}
	json.NewEncoder(w).Encode(httpi.ResAnyData{
		Msg:  "OK",
		Data: result[0].Interface(),
	})
}

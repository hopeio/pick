/*
 * Copyright 2024 hopeio. All rights reserved.
 * Licensed under the MIT License that can be found in the LICENSE file.
 * @Created by jyb
 */

package pick

import (
	"encoding/json"
	"io"
	"reflect"

	"github.com/hopeio/gox/errors"
	"github.com/hopeio/gox/log"
	httpx "github.com/hopeio/gox/net/http"
	http_fs "github.com/hopeio/gox/net/http/fs"
	"go.uber.org/zap"
)

var (
	ErrRepType = reflect.TypeOf((*ErrRep)(nil))
)

type ErrRep errors.ErrRep

func Response(w httpx.ICommonResponseWriter, traceId string, result []reflect.Value) error {
	if !result[1].IsNil() {
		err := ErrRepFrom(result[1].Interface())
		log.Errorw(err.Error(), zap.String(log.FieldTraceId, traceId))
		w.Header().Set(httpx.HeaderContentType, httpx.ContentTypeJsonUtf8)
		return json.NewEncoder(w).Encode(err)
	}
	data := result[0].Interface()
	if info, ok := data.(*http_fs.File); ok {
		header := w.Header()
		header.Set(httpx.HeaderContentType, httpx.ContentTypeOctetStream)
		header.Set(httpx.HeaderContentDisposition, "attachment;filename="+info.Name)
		defer info.File.Close()
		_, err := io.Copy(w, info.File)
		return err
	}
	if info, ok := data.(httpx.ICommonResponseTo); ok {
		_, err := info.CommonResponse(w)
		return err
	}

	w.Header().Set(httpx.HeaderContentType, httpx.ContentTypeJsonUtf8)
	return json.NewEncoder(w).Encode(httpx.RespAnyData{
		Data: data,
	})
}

func ErrRepFrom(err any) *errors.ErrRep {
	if err == nil {
		return nil
	}
	switch e := err.(type) {
	case *ErrRep:
		return (*errors.ErrRep)(e)
	case *httpx.ErrRep:
		return (*errors.ErrRep)(e)
	case errors.IErrRep:
		return e.ErrRep()
	case *errors.ErrRep:
		return e
	case errors.ErrCode:
		return e.ErrRep()
	case error:
		return errors.ErrRepFrom(e)
	case string:
		return errors.NewErrRep(errors.Unknown, e)
	}
	return errors.Unknown.ErrRep()
}

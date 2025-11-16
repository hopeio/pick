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
	ErrRespType = reflect.TypeOf((*ErrResp)(nil))
)

type ErrResp errors.ErrResp

func Respond(w httpx.ICommonResponseWriter, traceId string, result []reflect.Value) error {
	if !result[1].IsNil() {
		err := ErrRespFrom(result[1].Interface())
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
	if info, ok := data.(httpx.ICommonRespond); ok {
		_, err := info.CommonRespond(w)
		return err
	}

	w.Header().Set(httpx.HeaderContentType, httpx.ContentTypeJsonUtf8)
	return json.NewEncoder(w).Encode(httpx.RespAnyData{
		Data: data,
	})
}

func ErrRespFrom(err any) *errors.ErrResp {
	if err == nil {
		return nil
	}
	switch e := err.(type) {
	case *ErrResp:
		return (*errors.ErrResp)(e)
	case *httpx.ErrResp:
		return (*errors.ErrResp)(e)
	case errors.IErrResp:
		return e.ErrResp()
	case *errors.ErrResp:
		return e
	case errors.ErrCode:
		return e.ErrResp()
	case error:
		return errors.ErrRespFrom(e)
	case string:
		return errors.NewErrResp(errors.Unknown, e)
	}
	return errors.Unknown.ErrResp()
}

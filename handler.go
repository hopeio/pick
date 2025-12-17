/*
 * Copyright 2024 hopeio. All rights reserved.
 * Licensed under the MIT License that can be found in the LICENSE file.
 * @Created by jyb
 */

package pick

import (
	"context"
	"io"
	"net/http"
	"reflect"
	"strconv"

	"github.com/hopeio/gox/errors"
	"github.com/hopeio/gox/log"
	httpx "github.com/hopeio/gox/net/http"
	"github.com/hopeio/gox/net/http/apidoc"
	http_fs "github.com/hopeio/gox/net/http/fs"
	"go.uber.org/zap"
)

var (
	ErrRespType = reflect.TypeOf((*ErrResp)(nil))
)

type ErrResp errors.ErrResp

func Respond(ctx context.Context, w httpx.CommonResponseWriter, traceId string, result []reflect.Value) error {
	if !result[1].IsNil() {
		return RespondError(ctx, w, result[1].Interface(), traceId)
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
	if info, ok := data.(httpx.CommonResponder); ok {
		_, err := info.CommonRespond(ctx, w)
		return err
	}
	w.Header().Set(httpx.HeaderContentType, DefaultMarshaler.ContentType(data))
	buf, err := DefaultMarshaler.Marshal(data)
	if err != nil {
		buf = []byte(err.Error())
	}
	_, err = w.Write(buf)
	return err
}

func RespondError(ctx context.Context, w httpx.CommonResponseWriter, err any, traceId string) error {
	errresp := ErrRespFrom(err)
	log.Errorw(errresp.Error(), zap.String(log.FieldTraceId, traceId))
	w.Header().Set(httpx.HeaderContentType, DefaultMarshaler.ContentType(errresp))
	w.Header().Set(httpx.HeaderErrorCode, strconv.Itoa(int(errresp.Code)))
	buf, err1 := DefaultMarshaler.Marshal(errresp)
	if err1 != nil {
		buf = []byte(err1.Error())
	}
	_, err1 = w.Write(buf)
	return err1
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

func OpenApi(addr string) {
	Log(http.MethodGet, addr+apidoc.UriPrefix, "apidoc list")
	Log(http.MethodGet, addr+apidoc.UriPrefix+"/openapi/*file", "openapi")
	go func() {
		err := http.ListenAndServe(addr, http.DefaultServeMux)
		if err != nil {
			log.Fatal(err)
		}
	}()
}

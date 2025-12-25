/*
 * Copyright 2024 hopeio. All rights reserved.
 * Licensed under the MIT License that can be found in the LICENSE file.
 * @Created by jyb
 */

package pick

import (
	"context"
	"net/http"
	"reflect"
	"strconv"

	"github.com/hopeio/gox/errors"
	"github.com/hopeio/gox/log"
	httpx "github.com/hopeio/gox/net/http"
	"github.com/hopeio/gox/net/http/apidoc"
	"go.uber.org/zap"
)

var (
	ErrRespType = reflect.TypeOf((*ErrResp)(nil))
)

type ErrResp errors.ErrResp

func Respond(ctx context.Context, w http.ResponseWriter, traceId string, result []reflect.Value) (int, error) {
	if !result[1].IsNil() {
		return RespondError(ctx, w, result[1].Interface(), traceId)
	}
	data := result[0].Interface()

	if info, ok := data.(httpx.Responder); ok {
		return info.Respond(ctx, w)
	}
	buf, contentType := DefaultMarshaler(ctx, data)
	if wx, ok := w.(httpx.ResponseWriter); ok {
		wx.HeaderX().Set(httpx.HeaderContentType, contentType)
	} else {
		w.Header().Set(httpx.HeaderContentType, contentType)
	}
	ow := w
	if uw, ok := w.(httpx.Unwrapper); ok {
		ow = uw.Unwrap()
	}
	if recorder, ok := ow.(httpx.RecordBody); ok {
		recorder.RecordBody(buf, data)
	}
	return w.Write(buf)
}

func RespondError(ctx context.Context, w http.ResponseWriter, err any, traceId string) (int, error) {
	errresp := ErrRespFrom(err)
	log.Errorw(errresp.Error(), zap.String(log.FieldTraceId, traceId))

	buf, contentType := DefaultMarshaler(ctx, errresp)
	if wx, ok := w.(httpx.ResponseWriter); ok {
		header := wx.HeaderX()
		header.Set(httpx.HeaderContentType, contentType)
		header.Set(httpx.HeaderErrorCode, strconv.Itoa(int(errresp.Code)))
	} else {
		header := w.Header()
		header.Set(httpx.HeaderContentType, contentType)
		header.Set(httpx.HeaderErrorCode, strconv.Itoa(int(errresp.Code)))
	}
	ow := w
	if uw, ok := w.(httpx.Unwrapper); ok {
		w = uw.Unwrap()
	}
	if recorder, ok := ow.(httpx.RecordBody); ok {
		recorder.RecordBody(buf, errresp)
	}
	return w.Write(buf)
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

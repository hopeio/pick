/*
 * Copyright 2024 hopeio. All rights reserved.
 * Licensed under the MIT License that can be found in the LICENSE file.
 * @Created by jyb
 */

package pick

import (
	"encoding/json"
	"github.com/hopeio/context/httpctx"
	"github.com/hopeio/utils/errors/errcode"
	"github.com/hopeio/utils/log"
	httpi "github.com/hopeio/utils/net/http"
	"github.com/hopeio/utils/net/http/binding"
	"github.com/hopeio/utils/net/http/consts"
	http_fs "github.com/hopeio/utils/net/http/fs"
	"go.uber.org/zap"
	"io"
	"net/http"
	"reflect"
)

var (
	HttpContextType = reflect.TypeOf((*httpctx.Context)(nil))
	ErrRepType      = reflect.TypeOf((*ErrRep)(nil))
)

type ErrRep errcode.ErrRep

func CommonHandler(w http.ResponseWriter, req *http.Request, handle *reflect.Value) {
	handleTyp := handle.Type()
	handleNumIn := handleTyp.NumIn()
	if handleNumIn != 0 {
		params := make([]reflect.Value, handleNumIn)
		ctxi := httpctx.FromRequest(httpctx.RequestCtx{
			Request:  req,
			Response: w,
		})
		defer ctxi.RootSpan().End()
		for i := 0; i < handleNumIn; i++ {
			if handleTyp.In(i).ConvertibleTo(HttpContextType) {
				params[i] = reflect.ValueOf(ctxi)
			} else {
				params[i] = reflect.New(handleTyp.In(i).Elem())
				binding.Bind(req, params[i].Interface())
			}
		}
		result := handle.Call(params)
		RespWriteReflect(Writer{w}, ctxi.TraceID(), result)
	}
}

func RespWriteReflect(w httpi.ICommonResponseWriter, traceId string, result []reflect.Value) error {
	if !result[1].IsNil() {
		err := ErrHandle(result[1].Interface())
		log.Errorw(err.Error(), zap.String(log.FieldTraceId, traceId))
		w.Header().Set(consts.HeaderContentType, consts.ContentTypeJsonUtf8)
		return json.NewEncoder(w).Encode(err)
	}
	data := result[0].Interface()
	if info, ok := data.(*http_fs.File); ok {
		header := w.Header()
		header.Set(consts.HeaderContentType, consts.ContentTypeOctetStream)
		header.Set(consts.HeaderContentDisposition, "attachment;filename="+info.Name)
		defer info.File.Close()
		_, err := io.Copy(w, info.File)
		return err
	}
	if info, ok := data.(httpi.IHttpResponse); ok {
		_, err := httpi.CommonResponseWrite(w, info)
		return err
	}
	if info, ok := data.(httpi.ICommonResponseTo); ok {
		_, err := info.Response(w)
		return err
	}

	w.Header().Set(consts.HeaderContentType, consts.ContentTypeJsonUtf8)
	return json.NewEncoder(w).Encode(httpi.RespAnyData{
		Data: data,
	})
}

func ErrHandle(err any) *errcode.ErrRep {
	if err == nil {
		return nil
	}
	switch e := err.(type) {
	case *ErrRep:
		return (*errcode.ErrRep)(e)
	case *httpi.ErrRep:
		return (*errcode.ErrRep)(e)
	case errcode.IErrRep:
		return e.ErrRep()
	case *errcode.ErrRep:
		return e
	case errcode.ErrCode:
		return e.ErrRep()
	case error:
		return errcode.Unknown.Msg(e.Error())
	}
	return errcode.Unknown.ErrRep()
}

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
	http_fs "github.com/hopeio/utils/net/http/fs"
	"github.com/hopeio/utils/reflect/mtos"
	"go.uber.org/zap"
	"io"
	"net/http"
	"reflect"
)

var (
	HttpContextType = reflect.TypeOf((*httpctx.Context)(nil))
)

// Param is a single URL parameter, consisting of a key and a value.
type Param struct {
	Key   string
	Value string
}

// Params is a Param-slice, as returned by the router.
// The slice is ordered, the first URL parameter is also the first slice value.
// It is therefore safe to read values by the index.
type Params []Param

func CommonHandler(w http.ResponseWriter, req *http.Request, handle *reflect.Value, ps *Params) {
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
				if ps != nil || req.URL.RawQuery != "" {
					src := req.URL.Query()
					if ps != nil {
						pathParam := *ps
						if len(pathParam) > 0 {
							for i := range pathParam {
								src.Set(pathParam[i].Key, pathParam[i].Value)
							}
						}
					}
					mtos.Decode(params[i], src)
				}
				if req.Method != http.MethodGet {
					json.NewDecoder(req.Body).Decode(params[i].Interface())
				}
			}
		}
		result := handle.Call(params)
		ResWriteReflect(Writer{w}, ctxi.TraceID(), result)
	}
}

func ResWriteReflect(w httpi.ICommonResponseWriter, traceId string, result []reflect.Value) error {
	if !result[1].IsNil() {
		err := errcode.ErrHandle(result[1].Interface())
		log.Errorw(err.Error(), zap.String(log.FieldTraceId, traceId))
		w.Set(httpi.HeaderContentType, httpi.ContentTypeJsonUtf8)
		return json.NewEncoder(w).Encode(err)
	}
	data := result[0].Interface()
	if info, ok := data.(*http_fs.File); ok {
		w.Set(httpi.HeaderContentType, httpi.ContentTypeOctetStream)
		w.Set(httpi.HeaderContentDisposition, "attachment;filename="+info.Name)
		defer info.File.Close()
		_, err := io.Copy(w, info.File)
		return err
	}
	if info, ok := data.(httpi.IHttpResponse); ok {
		_, err := httpi.CommonResponseWrite(w, info)
		return err
	}
	if info, ok := data.(httpi.IHttpResponseTo); ok {
		if rw, ok := w.(http.ResponseWriter); ok {
			_, err := info.Response(rw)
			return err
		}
	}

	w.Set(httpi.HeaderContentType, httpi.ContentTypeJsonUtf8)
	return json.NewEncoder(w).Encode(httpi.ResAnyData{
		Data: data,
	})
}

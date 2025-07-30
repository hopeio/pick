/*
 * Copyright 2024 hopeio. All rights reserved.
 * Licensed under the MIT License that can be found in the LICENSE file.
 * @Created by jyb
 */

package pick

import (
	"encoding/json"
	"github.com/hopeio/gox/errors/errcode"
	"github.com/hopeio/gox/log"
	httpi "github.com/hopeio/gox/net/http"
	"github.com/hopeio/gox/net/http/consts"
	http_fs "github.com/hopeio/gox/net/http/fs"
	"go.uber.org/zap"
	"io"
	"reflect"
)

var (
	ErrRepType = reflect.TypeOf((*ErrRep)(nil))
)

type ErrRep errcode.ErrRep

func Response(w httpi.ICommonResponseWriter, traceId string, result []reflect.Value) error {
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
	if info, ok := data.(httpi.ICommonResponseTo); ok {
		_, err := info.CommonResponse(w)
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

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
	"reflect"
)

func ResWriteReflect(w httpi.ICommonResponseWriter, traceId string, result []reflect.Value) error {
	if !result[1].IsNil() {
		err := errcode.ErrHandle(result[1].Interface())
		log.Errorw(err.Error(), zap.String(log.FieldTraceId, traceId))
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
	w.Set(httpi.HeaderContentType, httpi.ContentTypeJsonUtf8)
	return json.NewEncoder(w).Encode(httpi.ResAnyData{
		Data: data,
	})
}

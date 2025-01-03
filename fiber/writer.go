/*
 * Copyright 2024 hopeio. All rights reserved.
 * Licensed under the MIT License that can be found in the LICENSE file.
 * @Created by jyb
 */

package pickfiber

import (
	"encoding/json"
	"github.com/gofiber/fiber/v3"
	"github.com/hopeio/utils/errors/errcode"

	"github.com/hopeio/utils/log"
	httpi "github.com/hopeio/utils/net/http"
	http_fs "github.com/hopeio/utils/net/http/fs"
	"go.uber.org/zap"
	"io"
	"net/http"
	"reflect"
)

func ResWriterReflect(ctx fiber.Ctx, traceId string, result []reflect.Value) error {
	writer := ctx.Response().BodyWriter()
	if !result[1].IsNil() {
		err := errcode.ErrHandle(result[1].Interface())
		log.Errorw(err.Error(), zap.String(log.FieldTraceId, traceId))
		json.NewEncoder(ctx.Response().BodyWriter()).Encode(err)
	}
	if info, ok := result[0].Interface().(*http_fs.File); ok {
		header := ctx.Response().Header
		header.Set(httpi.HeaderContentType, httpi.ContentTypeOctetStream)
		header.Set(httpi.HeaderContentDisposition, "attachment;filename="+info.Name)
		io.Copy(writer, info.File)
		if flusher, canFlush := writer.(http.Flusher); canFlush {
			flusher.Flush()
		}
		return info.File.Close()
	}
	return ctx.JSON(httpi.ResAnyData{
		Code: 0,
		Msg:  "success",
		Data: result[0].Interface(),
	})
}

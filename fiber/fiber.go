package pickfiber

import (
	"encoding/json"
	"github.com/gofiber/fiber/v3"
	"github.com/hopeio/cherry/context/fiber_context"
	http_fs "github.com/hopeio/cherry/utils/net/http/fs"
	"github.com/hopeio/pick"
	"io"
	"net/http"
	"reflect"

	"github.com/hopeio/cherry/protobuf/errorcode"
	"github.com/hopeio/cherry/utils/log"
	httpi "github.com/hopeio/cherry/utils/net/http"
	"github.com/hopeio/cherry/utils/net/http/api/apidoc"
	fiber_build "github.com/hopeio/cherry/utils/net/http/fasthttp/fiber"
)

func fiberResHandler(ctx fiber.Ctx, result []reflect.Value) error {
	writer := ctx.Response().BodyWriter()
	if !result[1].IsNil() {
		return json.NewEncoder(writer).Encode(errorcode.ErrHandle(result[1].Interface()))
	}
	if info, ok := result[0].Interface().(*http_fs.File); ok {
		header := ctx.Response().Header
		header.Set(httpi.HeaderContentType, httpi.ContentBinaryHeaderValue)
		header.Set(httpi.HeaderContentDisposition, "attachment;filename="+info.Name)
		io.Copy(writer, info.File)
		if flusher, canFlush := writer.(http.Flusher); canFlush {
			flusher.Flush()
		}
		return info.File.Close()
	}
	return ctx.JSON(httpi.ResAnyData{
		Code:    0,
		Message: "success",
		Details: result[0].Interface(),
	})
}

// 复用pick service，不支持单个接口的中间件
func Start(engine *fiber.App, genApi bool, modName string) {

	for _, v := range Svcs {
		describe, preUrl, middleware := v.Service()
		value := reflect.ValueOf(v)
		if value.Kind() != reflect.Ptr {
			log.Fatal("必须传入指针")
		}
		var infos []*pick.ApiDocInfo
		group := engine.Group(preUrl, middleware...)
		for j := 0; j < value.NumMethod(); j++ {
			method := value.Type().Method(j)
			methodInfo := pick.GetMethodInfo(&method, preUrl, FiberContextType)
			if methodInfo == nil {
				continue
			}
			if err := methodInfo.Check(); err != nil {
				log.Fatal(err)
			}

			methodType := method.Type
			methodValue := method.Func
			in2Type := methodType.In(2)
			methodInfoExport := methodInfo.GetApiInfo()
			group.Add([]string{methodInfoExport.Method}, methodInfoExport.Path[len(preUrl):], func(ctx fiber.Ctx) error {
				ctxi, span := fiber_context.ContextFromRequest(ctx, true)
				if span != nil {
					defer span.End()
				}
				in1 := reflect.ValueOf(ctxi)
				in2 := reflect.New(in2Type.Elem())
				if err := fiber_build.Bind(ctx, in2.Interface()); err != nil {
					return ctx.Status(http.StatusBadRequest).JSON(errorcode.InvalidArgument.ErrRep())
				}
				result := methodValue.Call([]reflect.Value{value, in1, in2})
				return fiberResHandler(ctx, result)
			})
			infos = append(infos, &pick.ApiDocInfo{ApiInfo: methodInfo, Method: method.Type})
		}
		pick.GroupApiInfos = append(pick.GroupApiInfos, &pick.GroupApiInfo{Describe: describe, Infos: infos})
	}
	if genApi {
		filePath := apidoc.FilePath
		pick.Markdown(filePath, modName)
		pick.Swagger(filePath, modName)
		//gin_build.OpenApi(engine, filePath)
	}
	pick.Registered(Svcs)
}

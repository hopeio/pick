package fiber

import (
	"context"
	"encoding/json"
	"github.com/hopeio/lemon/context/fasthttp_context"
	http_fs "github.com/hopeio/lemon/utils/net/http/fs"
	"github.com/hopeio/pick"
	"io"
	"net/http"
	"reflect"

	"github.com/gofiber/fiber/v2"
	"github.com/hopeio/lemon/protobuf/errorcode"
	"github.com/hopeio/lemon/utils/log"
	httpi "github.com/hopeio/lemon/utils/net/http"
	"github.com/hopeio/lemon/utils/net/http/api/apidoc"
	fiber_build "github.com/hopeio/lemon/utils/net/http/fasthttp/fiber"
)

type Service interface {
	//返回描述，url的前缀，中间件
	FiberService() (describe, prefix string, middleware []fiber.Handler)
}

var fiberSvcs = make([]Service, 0)

func RegisterFiberService(svc ...Service) {
	fiberSvcs = append(fiberSvcs, svc...)
}

var faberIsRegistered = false

func faberRegistered() {
	faberIsRegistered = true
	fiberSvcs = nil
}

func Api(f func()) {
	if !faberIsRegistered {
		f()
	}
}

func fiberResHandler(ctx *fiber.Ctx, result []reflect.Value) error {
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
func RegisterWithCtx(engine *fiber.App, genApi bool, modName string) {

	for _, v := range fiberSvcs {
		describe, preUrl, middleware := v.FiberService()
		value := reflect.ValueOf(v)
		if value.Kind() != reflect.Ptr {
			log.Fatal("必须传入指针")
		}
		var infos []*pick.ApiDocInfo
		engine.Group(preUrl, middleware...)
		for j := 0; j < value.NumMethod(); j++ {
			method := value.Type().Method(j)
			methodInfo := pick.GetMethodInfo(&method, preUrl, pick.ClaimsType)
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
			engine.Add(methodInfoExport.Method, methodInfoExport.Path, func(ctx *fiber.Ctx) error {
				in1 := reflect.ValueOf(fasthttp_context.ContextWithRequest(context.Background(), ctx.Request()))
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
	faberRegistered()
}

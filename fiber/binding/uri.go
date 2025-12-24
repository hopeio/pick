/*
 * Copyright 2024 hopeio. All rights reserved.
 * Licensed under the MIT License that can be found in the LICENSE file.
 * @Created by jyb
 */

package binding

import (
	"reflect"

	"github.com/gofiber/fiber/v3"
	"github.com/hopeio/gox/kvstruct"
)

type uriSource struct {
	fiber.Ctx
}

// TrySet tries to set a value by request's form source (likes map[string][]string)
func (s uriSource) TrySet(value reflect.Value, field *reflect.StructField, key string, opt *kvstruct.Options) (isSet bool, err error) {
	return kvstruct.SetValueByKVs(value, field, s, key, opt)
}

func (s uriSource) GetVs(key string) ([]string, bool) {
	v := s.Params(key)
	return []string{v}, v != ""
}

func (s uriSource) Has(key string) bool {
	v := s.Params(key)
	return v != ""
}

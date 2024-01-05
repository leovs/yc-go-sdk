// Copyright 2023 ztlcloud.com
// leovs @2023.12.21

package runtime

import (
	"reflect"
)

type ReflectFunc struct {
	Type  reflect.Type
	Value reflect.Value
}

func (r *ReflectFunc) Of(f any) (params []any) {
	r.Type = reflect.TypeOf(f)
	r.Value = reflect.ValueOf(f)
	if r.Type.Kind() != reflect.Func {
		return
	}
	for i := 0; i < r.Type.NumIn(); i++ {
		params = append(params, reflect.New(r.Type.In(i)).Interface())
	}
	return
}

func (r *ReflectFunc) Call(params ...any) []reflect.Value {
	var args = make([]reflect.Value, 0)
	for i := 0; i < len(params); i++ {
		args = append(args, reflect.ValueOf(params[i]).Elem())
	}
	return r.Value.Call(args)
}

func (r *ReflectFunc) CallMethod(name string, params ...any) []reflect.Value {
	m := r.Value.MethodByName(name)
	var args = make([]reflect.Value, 0)
	for i := 0; i < len(params); i++ {
		args = append(args, reflect.ValueOf(params[i]))
	}
	return m.Call(args)
}

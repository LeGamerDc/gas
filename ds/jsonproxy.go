package ds

import (
	"reflect"

	jsoniter "github.com/json-iterator/go"
)

type (
	id interface {
		Id() int32
	}

	Proxy[T any] struct {
		ptr T
	}

	mm struct {
		m map[int32]interface{}
	}

	proxyManager struct {
		mms map[reflect.Type]*mm
	}
)

var _manager = proxyManager{mms: make(map[reflect.Type]*mm)}

func RegisterProxy[T id](p T) {
	var (
		zero T
		tp   = reflect.TypeOf(zero)
	)

	if m, ok := _manager.mms[tp]; ok {
		m.m[p.Id()] = p
		return
	}
	_manager.mms[tp] = &mm{m: map[int32]interface{}{
		p.Id(): p,
	}}
}

func lookupPtr[T any](id int32) (p T) {
	var tp = reflect.TypeOf(p)
	if m, ok := _manager.mms[tp]; ok {
		p = m.m[id].(T)
	}
	return
}

func (p *Proxy[T]) Get() T {
	return p.ptr
}

func (p *Proxy[T]) UnmarshalJSON(b []byte) (e error) {
	var x int32
	if e = jsoniter.Unmarshal(b, &x); e != nil {
		return
	}
	p.ptr = lookupPtr[T](x)
	return nil
}

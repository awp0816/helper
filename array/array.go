package array

import (
	"reflect"
	"sync"
)

type CheckType interface {
	GetValueType(in interface{}) string
	Len() int
	Check(in interface{}) bool
	Add(in ...interface{})
	Delete(in interface{})
}

type Type struct {
	locker sync.Mutex
	values map[interface{}]bool
}

func NewArrayType() *Type {
	return &Type{
		values: make(map[interface{}]bool),
	}
}

func (e *Type) GetValueType(in interface{}) string {
	defer e.locker.Unlock()
	e.locker.Lock()
	return reflect.TypeOf(in).String()
}

func (e *Type) _add(in interface{}) {
	defer e.locker.Unlock()
	e.locker.Lock()
	e.values[in] = true
}

func (e *Type) Len() int {
	defer e.locker.Unlock()
	e.locker.Lock()
	return len(e.values)
}

func (e *Type) Check(in interface{}) bool {
	defer e.locker.Unlock()
	e.locker.Lock()
	if v, ok := e.values[in]; ok && v {
		return true
	}
	return false
}

func (e *Type) Add(in ...interface{}) {
	for _, v := range in {
		e._add(v)
	}
}

func (e *Type) Delete(in interface{}) {
	defer e.locker.Unlock()
	e.locker.Lock()
	delete(e.values, in)
}

package main

import (
	"fmt"
	"io"
	"reflect"
	"testing"
)

func TestStructInnerPtr(t *testing.T) {
	type A struct {
		Name *string
	}
	a := A{}
	ref := reflect.ValueOf(&a).Elem()
	name := "赵云"
	ref.Field(0).Set(reflect.ValueOf(&name))
	fmt.Println(*a.Name)
	*a.Name = "关云长"
	fmt.Println(*a.Name)
	fmt.Println(name)
	//赵云
	//关云长
	//关云长
}

func TestNil1(t *testing.T) {
	var m io.Reader
	v := reflect.ValueOf(m)
	fmt.Println(v.IsNil())
	fmt.Println(v.IsValid())
	fmt.Println(v.CanSet())
}

func TestNil2(t *testing.T) {
	var m []byte
	v := reflect.ValueOf(m)
	fmt.Println(v.IsNil())
	fmt.Println(v.IsValid())
	fmt.Println(v.CanSet())
}

func TestNil3(t *testing.T) {
	var m map[string]int
	v := reflect.ValueOf(m)
	fmt.Println(v.IsNil())
	fmt.Println(v.IsValid())
	fmt.Println(v.CanSet())
}

func TestStr(t *testing.T) {
	type A struct {
	}
	s := make([]interface{}, 7, 10)
	s[0] = nil
	r := reflect.ValueOf(&s).Elem()
	i0 := r.Index(0)
	fmt.Println(i0.IsValid())
	fmt.Println(i0.IsNil())
	fmt.Println(i0.CanSet())
}

func TestNOvalue(t *testing.T) {
	type A struct {
	}
	a := (*A)(nil)
	fmt.Println(a == nil)
	v := reflect.ValueOf(&a).Elem()
	fmt.Println(v.IsValid())
	fmt.Println(v.IsNil())
	fmt.Println(v.CanSet())
}

func TestNIlINterfa(t *testing.T) {
	var a interface{}
	v := reflect.ValueOf(a)
	fmt.Println(v.IsValid()) //false
	fmt.Println(v.IsNil())   //panic
}

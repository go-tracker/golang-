package main

import (
	"fmt"
	"reflect"
	"unsafe"
)

func AssignPrivateVariable() {
	type BStruct struct {
		UnionId string
	}

	type AStruct struct {
		aBStruct BStruct
		aInt     int
		aString  string
	}
	typeA := reflect.TypeOf(AStruct{})

	fmt.Println("AStruct 内存大小：", typeA.Size())

	var obj = AStruct{}
	ptrStartOffset := uintptr(unsafe.Pointer(&obj))

	typeB := reflect.TypeOf(BStruct{})
	*((*string)(unsafe.Pointer(ptrStartOffset + typeA.Field(0).Offset + typeB.Field(0).Offset))) = "BStruct.UnionId" //aBStruct

	*((*int)(unsafe.Pointer(ptrStartOffset + typeA.Field(1).Offset))) = 666 //aInt

	*((*string)(unsafe.Pointer(ptrStartOffset + typeA.Field(2).Offset))) = "aString" //aString

	fmt.Printf("%#v\n", obj)

}

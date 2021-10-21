package main

import (
	"fmt"
	"reflect"
)

// 反射模板代码，通用，亦可学习
// 想要"搞好"反射，只有一个"根原理"：完整判断类型，直到递归到基本数据类型，站在上层封装类型搞反射，总会遇到意想不到的问题，除非"框架反射"有明确的"数据结构约定"。
// 只处理"约定"的数据类型，开发者保证传值正确，否则"框架panic不符合约定的数据类型"。没有万能反射，有仅有"约定反射"，要使用反射接口，就必须遵从反射约定。
func main() {
	//  要想设置值，必须可寻址，即反射对象代表本身内存值，而不是副本值（golang的值拷贝模式）
	//  简单理解可寻址，就是指针，指针不管拷贝多少次，任何拷贝的指针值，都指向同一个内存实体值，而不是"实体值"的副本。副本便是不可寻址
	//  偏门理解"可寻址"就是"羊毛出在羊身上"，更该值一定是更该的同一只羊身上的东西，而不是克隆版本的羊身上去了
	var obj ReflectValue
	ref := reflect.ValueOf(&obj)
	handleReflect(ref)
}

// ref 参数值拷贝，由于reflect.Value内部是指针(ptr unsafe.Pointer)，发生拷贝依然可以寻址。但是常规类型，就算内部是指针，如果初始反射值不是取指针反射，依然是panic（无法寻址）
func handleReflect(ref reflect.Value) {
	// 第一步，去除多余指针，得到一级指针（正常写代码只用一级指针）
	for ref.Kind() == reflect.Ptr && !ref.IsNil() {
		ref = ref.Elem()
	} // 循环结束后：
	// 1.如果ref还是指针，那么绝对是空指针（*ReReflectType）
	// 2.如果ref不是指针，那么绝对是反射类型值（ReflectValue），此时的值也可能是nil，但是这个"值"绝对是你函数传参的"本体值"
	//   var list []string
	//   callReflectFunc(&list)
	//   这个"本体值"就是指 list 变量自己，在callReflectFunc函数里面操作list和函数外面操作list是操作"同一个变量"，而不是"同一个指针"，
	//   那么操作代码就跟写在callReflectFunc函数外面是一样的。

	// 第二步，断言上面结论
	if ref.Kind() == reflect.Ptr && ref.Elem().IsValid() {
		panic("此处Ptr绝对是空指针")
	}
	if ref.Kind() != reflect.Ptr && ref.IsValid() {
		fmt.Println(fmt.Sprintf("非nil值%s\n", ref.Kind()))
	}
	if ref.Kind() != reflect.Ptr && !ref.IsValid() {
		fmt.Println(fmt.Sprintf("nil值%s\n", ref.Kind()))
	}

	// 第三步，由于存在nil和非nil的值，所以这里必须用type来判定类型，才好决定后续操作(对nil的情况是初始化还是return)
	typ := ref.Type()
	switch typ.Kind() {
	case reflect.Invalid:
		// nil值的type, 任何反射方法都会panic，除了String()方法
	//
	//————————————————基础类型，调用对应数据类型方法设置值————————————————————————————————————————————————————
	case reflect.Bool:
		ref.SetBool(true)
		b := ref.Bool()
		fmt.Println("读取布尔值：", b)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		ref.SetInt(int64(64)) //根据int长度设置
		it := ref.Int()
		fmt.Println("只有Int()返回64位，自己强转：", int32(it))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		ref.SetUint(uint64(64)) //根据uint长度设置
		ut := ref.Uint()
		fmt.Println("只有UInt()返回64位，自己强转：", uint32(ut))
	case reflect.Float32, reflect.Float64:
		ref.SetFloat(float64(64))
		fl := ref.Float()
		fmt.Println("只有Float()返回64位，自己强转：", fl)
	case reflect.String:
		ref.SetString("hello world")
		s := ref.String()
		fmt.Println("读取字符串内容：", s)
		strLen := ref.Len()
		fmt.Println("字符串长度：", strLen)
		v := ref.Index(0)
		if v.Kind() != reflect.Uint8 {
			panic("此处v绝对是个byte")
		}
		strHello := ref.Slice(0, 5)
		if strHello.Interface() != reflect.ValueOf("hello").Interface() {
			// 使用 Interface() 比较反射值 Value, 不可以使用 == 比较
			panic("此处必是 hello 字符串")
		}
	//————————————————————————————————————————————————————————————————————
	case reflect.Map:
		// 可能是nil未初始化
		if ref.IsValid() {
			fmt.Println("已经是初始化分配了内存的实体值")
		}
		if !ref.IsValid() {
			fmt.Println("未初始化的Map")
			// 初始化Map
			ref.Set(reflect.MakeMapWithSize(typ, 2)) //初始容量n根据需求自定义
		}
		if !ref.IsValid() {
			panic("此处Map绝对是已经初始化的Map")
		}
		{ // 查看map定义
			// 长度
			mpLen := ref.Len()
			fmt.Println("Map长度：", mpLen)
			//——————————————— map[string]interface{} —————————————————————————
			kType := typ.Key()
			fmt.Println("Map key类型：", kType)
			vType := typ.Elem()
			fmt.Println("Map value类型：", vType)
			//————————————————————————————————————————————————————————————————
			// 遍历map 方式1 ——————————————————————————
			ks := ref.MapKeys()
			for _, k := range ks {
				fmt.Println("Map k类型：", k.Type())
				v := ref.MapIndex(k)
				fmt.Println("Map v类型：", v.Type())
				if !v.IsValid() {
					fmt.Println("Map不存在这个key（或者key对应值就是nil？）")
				}
				if v.IsValid() {
					fmt.Println("Map存在这个key：")
				}
			} // 遍历map 方式1 ——————————————————————————
			// 遍历map 方式2 ——————————————————————————
			iter := ref.MapRange()
			for iter.Next() {
				k := iter.Key()
				v := iter.Value()
				fmt.Println("Map k类型：", k.Type())
				fmt.Println("Map v类型：", v.Type())
				if !v.IsValid() {
					fmt.Println("Map不存在这个key（或者key对应值就是nil？）")
				}
				if v.IsValid() {
					fmt.Println("Map存在这个key：")
				}
			} // 遍历map 方式2 ——————————————————————————
		} // 查看map定义
		// 设置map值，根据map k类型和v类型及对应业务设置对应值
		ref.SetMapIndex(reflect.New(typ.Key()).Elem(), reflect.New(typ.Elem()).Elem())
		// 比如：
		if typ.Key().Kind() == reflect.String && typ.Elem().Kind() == reflect.String {
			mp := ref.Interface().(map[string]string)
			mp["zhaoyun"] = "赵云"
			mp["guanyu"] = "关羽"
		}
		if typ.Key().Kind() == reflect.String && typ.Elem().Kind() == reflect.String {
			ref.SetMapIndex(reflect.ValueOf("mingren"), reflect.ValueOf("鸣人"))
			ref.SetMapIndex(reflect.ValueOf("zuozhu"), reflect.ValueOf("佐助"))
		}
		// 复杂点，可以递归设置map元素内部值，比如某个元素是结构体或者指定结构体名字
		iter := ref.MapRange()
		if iter.Next() {
			val := iter.Value()
			if val.IsValid() { // 非nil值
				handleReflect(val) //递归处理内部值
			}
		}
	case reflect.Ptr:
		// 绝对是nil未初始化的指针
		if !ref.IsNil() {
			panic("此处绝对是nil指针")
		}
		// 初始化
		switch typ.Elem().Kind() {
		// 所有反射类型处理，但实际开发中基本是"基础数据类型指针，结构体类型指针"，其他"非主流指针可以直接panic，降低开发难度，提升实用性"
		case reflect.Map:
			panic("非主流指针类型")
		}
		ref.Set(reflect.New(typ.Elem()))
		if !ref.IsNil() {
			panic("此处绝对是非nil指针")
		}
		fmt.Println("当前指针类型：", ref.Type())
		memoryValue := ref.Elem()
		// 此时可以单层判断值类型处理，也可以递归处理复杂业务。推荐递归
		handleReflect(memoryValue)
	case reflect.Slice:
		// 可能是nil未初始化
		if ref.IsValid() {
			fmt.Println("已初始化的切片")
		}
		if !ref.IsValid() {
			fmt.Println("未初始化的切片")
			//初始化切片
			ref.Set(reflect.MakeSlice(typ, 5, 8)) //根据业务设置长度和容量
		}
		if !ref.IsValid() {
			panic("此时切片绝对是初始化的")
		}
		// 设置长度
		ref.SetLen(5)
		// 设置容量
		ref.SetCap(8)
		{ // 查看Slice
			sliLen := ref.Len()
			fmt.Println("切片长度：", sliLen)
			sliCap := ref.Cap()
			fmt.Println("切片容量：", sliCap)
			subSlice := ref.Slice(0, 2)
			if subSlice.Len() != 2 {
				panic("长度必是2")
			}
			subSlice2 := ref.Slice3(0, 2, 3)
			if subSlice2.Cap() != 3 {
				panic("容量必是3")
			}
			// 只能遍历
			for i := 0; i < ref.Len(); i++ {
				v := ref.Index(i)
				fmt.Println("切片类型：", v.Type())
				fmt.Println("切片值：", v.Interface())
				// 如果切片元素是结构体类型，结构体指针类型，或者符合业务数据类型，都可以递归处理
				if v.Kind() == reflect.Struct || v.Kind() == reflect.Ptr && v.Type().Elem().Kind() == reflect.Struct {
					handleReflect(v)
				}
			}
		} // 查看Slice
		{ // 设置slice
			// 方式1：append追加元素
			ref.Set(reflect.Append(ref, reflect.New(typ.Elem()).Elem(), reflect.New(typ.Elem()).Elem())) // 追加一个或多个元素
			ref.Set(reflect.AppendSlice(ref, reflect.MakeSlice(reflect.SliceOf(typ.Elem()), 1, 2)))      // 追加一个切片
			// 方式2：设置元素
			v := ref.Index(0)
			// v.IsNil() 说明 chan, func, interface, map, pointer, or slice 是否是nil值，nil的interface会panic，其他类型不会。调用Nil()方法很复杂
			// v.IsNil() 会不会panic主要看v是不是底层代表nil值，比如一开始就用v=reflect.ValueOf(nil)来得到反射值，v直接代表nil值。但是Slice，map里面的nil值不一样，不能等同直接nil值（slice，map是有类型的nil值）
			// 主要是"直接nil值"有个特点"无类型&无值"，切片，Map的"间接nil值"有个特点"有类型&无值"。在GO里面"无类型无值的变量"才视作nil，有类型无值的变量不是nil
			// 空接口（interface{}，io.Reader）会panic可以这样想：函数实现接口是用来调用的，而不是用来发射设置值的，所以一开始就要指定函数，反射直接调用，而不是通过反射去设置函数变量的值。
			switch v.Kind() {
			case reflect.Chan, reflect.Func, reflect.Map, reflect.Ptr, reflect.Slice:
				if v.IsNil() { //调用Nil之前，必须是这些类型的判定
					fmt.Println("切片存在nil值")
					if !(v.IsValid() == true || v.IsNil() == true || v.CanSet() == true) {
						panic("切片元素任何类型为nil，都可以设置")
					}
				}
				// 递归设置
				handleReflect(v)
			case reflect.Interface:
			}
			if v.Kind() == reflect.Struct {
				v.Set(reflect.New(v.Type()).Elem())
			}
			if v.Kind() == reflect.Ptr && v.Type().Elem().Kind() == reflect.Struct {
				v.Set(reflect.New(v.Type().Elem()))
			}
		} // 设置slice
	case reflect.Struct:
		// 绝对不是nil，至少是零值结构体
		if !ref.IsValid() {
			panic("结构体绝对是有效非zero Value")
		}
		fmt.Println("字段数量：", ref.NumField())
		{ // 遍历读字段
			for i := 0; i < ref.NumField(); i++ {
				fv := ref.Field(i) //字段value
				ft := typ.Field(i) //字段信息StructField
				fmt.Println("结构体字段名字：", ft.Name)
				fmt.Println("结构体字段值：", fv.Interface())
				if ft.Anonymous {
					fmt.Println("匿名字段类型: ", fv.Type())
					// 递归处理或者根据业务接口实现或者数据类型判断决定实际处理
					if ft.Type.Kind() != reflect.Ptr {
						handleReflect(fv)
					}
				}
			}
		} // 遍历字段
		{ // 设置值
			if ref.Field(0).Kind() == reflect.Int {
				ref.Field(0).SetInt(996)
			}
			if ref.Field(0).Kind() == reflect.Struct {
				handleReflect(ref.Field(0)) // 递归设置结构体值
			}
			if ref.Field(0).Kind() == reflect.Ptr && ref.Field(0).IsNil() {
				switch t := ref.Field(0).Type().Elem(); t.Kind() {
				// 指针指向类型
				case reflect.Int:
				case reflect.Struct:
					ref.Field(0).Set(reflect.New(t).Elem())
				case reflect.Map:
				case reflect.Slice:
					// 如果是字节数组
					if t.Elem().Kind() == reflect.Uint8 {
						ref.Field(0).Set(reflect.ValueOf(make([]byte, 0, 8)))
					}
				}
			}
		} // 设置值
	case reflect.Array, reflect.Chan, reflect.Func, reflect.Interface, reflect.UnsafePointer, reflect.Uintptr:
		// 对于反射设置值的目的来说，这些类型用不上，不去关注
	case reflect.Complex64, reflect.Complex128:
		// 虚数业务系统实际开发也用不上，不去关注
		ref.SetComplex(complex(64, 64))
	}
}

// 用于文档书写，代表一个真实的反射对象值（切片，Map，结构体，数字，字符串，布尔）
type ReflectValue int
type ReflectType int

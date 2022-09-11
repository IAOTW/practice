package reflect

import (
	"fmt"
	"reflect"
	"testing"
)

func TestReflect(t *testing.T) {
	// 声明一个空结构体
	type cat struct {
		a string
	}
	// 创建cat的实例，它是一个指针
	ins := cat{a: "a"}
	// 对指针取反射对象
	typeOfCat := reflect.TypeOf(ins)
	// 显示反射类型对象的名称和种类，name:'' kind:'ptr'
	fmt.Printf("name:'%v' kind:'%v'\n", typeOfCat.Name(), typeOfCat.Kind())
	// 取类型的元素，typeOfCat是一个指针对反射对象，Elem()是会拿到该指针所指向对值对反射对象
	//typeOfCat = typeOfCat.Elem()
	//// 显示反射类型对象的名称和种类，element name: 'cat', element kind: 'struct'
	//fmt.Printf("element name: '%v', element kind: '%v'\n", typeOfCat.Name(), typeOfCat.Kind())

	valueOfCat := reflect.ValueOf(ins)
	// 显示反射类型对象的名称和种类，name:'*main.cat' kind:'ptr'
	fmt.Printf("name:'%v' kind:'%v'\n", valueOfCat.Type(), valueOfCat.Kind())
	// 取类型的元素，valueOfCat是一个指针对反射对象，Elem()是会拿到该指针所指向对值对反射对象
	//valueOfCat = valueOfCat.Elem()
	//// 显示反射类型对象的名称和种类，element name: 'main.cat', element kind: 'struct'
	//fmt.Printf("element name: '%v', element kind: '%v'\n", valueOfCat.Type(), valueOfCat.Kind())

	//
	//// 创建cat的实例，它是一个指针
	//ins := &cat{a: "a"}
	//// 对指针取反射对象
	//typeOfCat := reflect.TypeOf(ins)
	//// 显示反射类型对象的名称和种类，name:'' kind:'ptr'
	//fmt.Printf("name:'%v' kind:'%v'\n", typeOfCat.Name(), typeOfCat.Kind())
	//// 取类型的元素，typeOfCat是一个指针对反射对象，Elem()是会拿到该指针所指向对值对反射对象
	//typeOfCat = typeOfCat.Elem()
	//// 显示反射类型对象的名称和种类，element name: 'cat', element kind: 'struct'
	//fmt.Printf("element name: '%v', element kind: '%v'\n", typeOfCat.Name(), typeOfCat.Kind())
	//
	//valueOfCat := reflect.ValueOf(ins)
	//// 显示反射类型对象的名称和种类，name:'*main.cat' kind:'ptr'
	//fmt.Printf("name:'%v' kind:'%v'\n", valueOfCat.Type(), valueOfCat.Kind())
	//// 取类型的元素，valueOfCat是一个指针对反射对象，Elem()是会拿到该指针所指向对值对反射对象
	//valueOfCat = valueOfCat.Elem()
	//// 显示反射类型对象的名称和种类，element name: 'main.cat', element kind: 'struct'
	//fmt.Printf("element name: '%v', element kind: '%v'\n", valueOfCat.Type(), valueOfCat.Kind())
}

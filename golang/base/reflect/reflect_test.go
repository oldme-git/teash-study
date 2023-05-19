package reflect

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

// value有可以取出type
func ExampleValueOfType() {
	var i int8
	fmt.Println(reflect.ValueOf(i).Type() == reflect.TypeOf(i))

	// OutPut:
	// true
}

// value有可以取出type elem
func ExampleValueOfTypeElem() {
	var i int8
	// elem
	vOf := reflect.ValueOf(&i)
	tOf := reflect.TypeOf(&i)
	fmt.Println(vOf.Elem().Type() == tOf.Elem())
	fmt.Println(vOf.Type().Elem() == tOf.Elem())

	// OutPut:
	// true
	// true
}

func ExampleTypeElem() {
	var (
		i       int
		tOf     = reflect.TypeOf(i)
		tOfElem = reflect.TypeOf(&i).Elem()
	)
	fmt.Println(tOf == tOfElem)

	// OutPut:
	// true
}

// value和type的kind
func ExampleKindWithValueAndType() {
	var i int8
	fmt.Println(reflect.TypeOf(i).Kind() == reflect.ValueOf(i).Kind())

	// OutPut:
	// true
}

// name获取类型名称，kind获取类型种类
func ExampleKind() {
	type myInt int
	var (
		// 自定义类型
		myint myInt
		// int基础类型
		i      int
		tOfi   = reflect.TypeOf(i)
		tOfMyi = reflect.TypeOf(myint)
	)

	fmt.Printf("int name: %v，int kind: %v\n", tOfi.Name(), tOfi.Kind())
	fmt.Printf("int name: %v，int kind: %v", tOfMyi.Name(), tOfMyi.Kind())

	// OutPut:
	// int name: int，int kind: int
	// int name: myInt，int kind: int
}

// 从value中恢复值
func ExampleValueInt() {
	var i int
	i = 1
	fmt.Println(reflect.ValueOf(i).Int())

	// OutPut:
	// 1
}

// 使用Interface获取空接口，然后使用类型断言获取原始值
func ExampleInterface() {
	var (
		i      int = 1
		origin     = reflect.ValueOf(i).Interface().(int)
	)

	fmt.Println(origin)

	// OutPut:
	// 1
}

// canset
func TestCanSet(t *testing.T) {
	var i int = 1
	fmt.Println("int:", reflect.ValueOf(i).CanSet())

	var s string = "abc"
	fmt.Println("string:", reflect.ValueOf(s).CanSet())

	var b bool = true
	fmt.Println("bool:", reflect.ValueOf(b).CanSet())

	var f float64 = 1.1
	fmt.Println("float:", reflect.ValueOf(f).CanSet())

	// ptr
	fmt.Println("floatPtr:", reflect.ValueOf(&f).CanSet())
	// elem
	fmt.Println("floatElem:", reflect.ValueOf(&f).Elem().CanSet())
}

// valueElem修改值
func TestValueElem(t *testing.T) {
	var i int = 1
	vOf := reflect.ValueOf(&i)
	if b := vOf.Elem().CanSet(); b {
		vOf.Elem().SetInt(2)
	}
	fmt.Println(i)
}

// struct test
type myStr string

type Internal struct {
	OneFloat float64
	OneMap   map[string]interface{}
}

type External struct {
	OneStr string
	int
	MyStr    myStr `are_you_ok:"I'm Ok!"`
	Internal Internal
}

var myStruct = External{
	OneStr: "one str",
	int:    1,
	MyStr:  "my str",
	Internal: Internal{
		OneFloat: 1.1,
		OneMap: map[string]interface{}{
			"key1": 1,
			"key2": "2",
		},
	},
}

// Type与Struct，主要是StructField结构体
func TestTypeStruct(t *testing.T) {
	tOf := reflect.TypeOf(myStruct)

	// NumField 获取结构体成员数量
	// Field 根据索引获取StructField
	for i := 0; i < tOf.NumField(); i++ {
		sf := tOf.Field(i)

		fmt.Printf("%+v\n", sf)
		//fmt.Println(sf.Name)      //  string		字段名
		//fmt.Println(sf.PkgPath)   //  string		字段路径
		//fmt.Println(sf.Type)      //  Type		字段反射类型对象
		//fmt.Println(sf.Tag)       //  StructTag	字段的结构体标签
		//fmt.Println(sf.Offset)    //  uintptr		字段在结构体中的相对偏移
		//fmt.Println(sf.Index)     //  []int		Type.FieldByIndex中的返回的索引值
		//fmt.Println(sf.Anonymous) //  bool 		是否为匿名字段
	}

	fmt.Println("-----------")

	// FindByName 通过Name获取StructField，没有找到时b为false
	if myFindByName, b := tOf.FieldByName("MyStr"); b {
		fmt.Printf("%+v\n", myFindByName)

		fmt.Println("-----------")
		// 取tag
		myTag := myFindByName.Tag.Get("are_you")
		fmt.Printf("MyStr的tag: %s\n", myTag)
	}

	fmt.Println("-----------")

	// FieldByIndex 多层结构体访问
	myFieldByIndex := tOf.FieldByIndex([]int{3, 0})
	fmt.Printf("%+v\n", myFieldByIndex)

	fmt.Println("-----------")

	// FieldByNameFunc 根据匹配函数寻找，如果匹配到多个则视为没有匹配到
	if myFieldByNameFunc, b := tOf.FieldByNameFunc(func(s string) bool {
		return strings.HasPrefix(s, "O")
	}); b {
		fmt.Printf("%+v\n", myFieldByNameFunc)
	} else {
		fmt.Println("没有匹配到")
	}
}

// Value与Struct
func TestValueStruct(t *testing.T) {
	vOf := reflect.ValueOf(myStruct)
	// Value对于结构体的操作和Type基本一致，只不过返回的数据类型不是StructField，而是结构体字段的Value
	myStrValue := vOf.FieldByName("MyStr")
	fmt.Println(myStrValue.Interface())
}

// 调用函数
func add(x, y int) int {
	return x + y
}
func TestFunc(t *testing.T) {
	var (
		vOf   = reflect.ValueOf(add)
		param = []reflect.Value{reflect.ValueOf(1), reflect.ValueOf(2)}
		res   = vOf.Call(param)
	)
	fmt.Println(res[0].Int())
}
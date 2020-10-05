package configUtil

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

func loadConfig(data interface{}) {
	file, err := os.OpenFile("reflect/ini/conf.ini", os.O_RDONLY, 0644)
	if err != nil {
		fmt.Println("read err:", err)
		return
	}
	var b [512]byte
	readSize, err := file.Read(b[:])
	var structName string
	//判断传进来的参数是不是指针类型且是结构体
	t := reflect.ValueOf(data)
	if t.Kind() != reflect.Ptr || t.IsNil() {
		fmt.Println("type err:", reflect.TypeOf(data))
		return
	}

	if t.Elem().Kind() != reflect.Struct {
		fmt.Println("type err:", reflect.TypeOf(t.Elem()))
		return
	}

	//将字节类型的文件内容转换成字符串
	split := strings.Split(string(b[:readSize]), "\r\n")
	for idx, line := range split {
		//去除每一行的空格
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, ";") || strings.HasPrefix(line, "#") || line == "" {
			continue
		}
		//判断是否以[开头
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			//内容是否为空
			sectionName := strings.TrimSpace(line[1 : len(line)-1])
			if len(sectionName) == 0 {
				fmt.Printf("syntax error in line:%d,%#v\n", idx+1, line)
				return
			}
			//拿到名称
			v := reflect.TypeOf(data)
			for i := 0; i < v.Elem().NumField(); i++ {
				field := v.Elem().Field(i)
				if sectionName == field.Tag.Get("ini") {
					structName = field.Name
				}
			}
		} else {
			//判断是否为键值对
			keyVal := strings.Split(line, "=")
			if len(keyVal) < 2 {
				fmt.Printf("syntax error in line:%d,%#v\n", idx+1, line)
				return
			}
			//取出对象
			structObj := t.Elem().FieldByName(structName)
			//判断是否为结构体
			if structObj.Kind() != reflect.Struct {
				fmt.Printf("%v is not a struct \n", structName)
				return
			}
			var fieldName string
			Obj := structObj.Type()
			for i := 0; i < Obj.NumField(); i++ {
				field := Obj.Field(i)
				if field.Tag.Get("ini") == strings.TrimSpace(keyVal[0]) {
					fieldName = field.Name
					break
				}
			}

			//没有找到对应的字段
			if fieldName == "" {
				continue
			}

			//根据名称找到对应的字段
			name, _ := Obj.FieldByName(fieldName)

			//根据字段类型赋值
			switch name.Type.Kind() {
			case reflect.String:
				structObj.FieldByName(fieldName).SetString(keyVal[1])
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				parseInt, err := strconv.ParseInt(keyVal[1], 10, 64)
				typeErr(err, idx, line)
				structObj.FieldByName(fieldName).SetInt(parseInt)
			case reflect.Bool:
				parseBool, err := strconv.ParseBool(keyVal[1])
				typeErr(err, idx, line)
				structObj.FieldByName(fieldName).SetBool(parseBool)
			case reflect.Float32, reflect.Float64:
				parseFloat, err := strconv.ParseFloat(keyVal[1], 64)
				typeErr(err, idx, line)
				structObj.FieldByName(fieldName).SetFloat(parseFloat)
			default:
				break
			}
		}
	}

}

func typeErr(err error, line int, msg string) {
	if err != nil {
		fmt.Printf("type error in line:%d,%#v\n", line+1, msg)
		return
	}
}

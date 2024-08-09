package ini_analyzer

import (
	"bufio"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

func LoadIni(filepath string, data any) error {
	t := reflect.TypeOf(data)
	if t.Kind() != reflect.Ptr {
		err := fmt.Errorf("data需要是一个指针类型")
		return err
	}
	if t.Elem().Kind() != reflect.Struct {
		err := fmt.Errorf("data指针指向的类型必须是结构体")
		return err
	}

	//定于结构体名称，方便后续根据结构体名称进行修改值
	var structName string
	//解析ini文件并赋值到data指针指向的结构体中
	file, err := os.Open(filepath)
	defer file.Close()
	if err != nil {
		return err
	}
	scanner := bufio.NewScanner(file)
	var currentSection string

	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		//忽略空行和注释行
		// 忽略空行和注释行
		if line == "" || strings.HasPrefix(line, ";") || strings.HasPrefix(line, "#") {
			continue
		}
		// 检查是否为新的节
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			if strings.Count(line, "[") != 1 || strings.Count(line, "]") != 1 {
				return fmt.Errorf("error: invalid section format, line: %s", line)
			}
			section := strings.Trim(line, "[]")
			if section == "" {
				err = fmt.Errorf("error: syntax error:section canot be empty, line: %s", line)
				return err
			}

			currentSection = section
			// 根据当前节的名称，找到对应的字段，并设置值
			for i := 0; i < t.Elem().NumField(); i++ {
				field := t.Elem().Field(i)
				filedName := field.Tag.Get("ini")
				if currentSection == filedName {
					structName = field.Name // 记录结构体名称
				}
			}
			continue
		}
		// 解析键值对
		if strings.Index(line, "=") == -1 {
			return fmt.Errorf("error: invalid key:value format, line: %s", line)
		}
		keyValue := strings.Split(line, "=")
		key := strings.TrimSpace(keyValue[0])
		value := strings.TrimSpace(keyValue[1])
		if key == "" {
			return fmt.Errorf("error: syntax error: key or value canot be empty, line: %s", line)
		}
		//根据获取到的结构体名称，通过反射给结构体赋值
		v := reflect.ValueOf(data)
		sValue := v.Elem().FieldByName(structName)
		sType := sValue.Type()
		var fieldName string
		var structfield reflect.StructField
		if sType.Kind() != reflect.Struct {
			return fmt.Errorf("error: struct kind error, line: %s", line)
		}

		for i := 0; i < sType.NumField(); i++ {
			field := sType.Field(i)
			if field.Tag.Get("ini") == key {
				//找到了字段
				fieldName = field.Name
				structfield = field
			}
		}

		//给字段赋值
		//根据fieldName取出字段的值
		field := sValue.FieldByName(fieldName)
		//fmt.Println("找到了字段", fieldName, structfield.Type.Kind())
		switch structfield.Type.Kind() {
		case reflect.String:
			field.SetString(value)
		case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			valueInt, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				return fmt.Errorf("error: invalid int value, line: %s", line)
			}
			field.SetInt(valueInt)
		case reflect.Bool:
			valueBool, err := strconv.ParseBool(value)
			if err != nil {
				return fmt.Errorf("error: invalid boolean value, line: %s", line)
			}
			field.SetBool(valueBool)
		case reflect.Float32, reflect.Float64:
			valueFloat, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return fmt.Errorf("error: invalid float value, line: %s", line)
			}
			field.SetFloat(valueFloat)
		}
	}
	return nil
}

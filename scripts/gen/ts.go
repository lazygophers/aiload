package main

import (
	"bytes"
	"fmt"
	"github.com/emicklei/proto"
	"github.com/lazygophers/codegen/codegen"
	"github.com/lazygophers/log"
	"github.com/lazygophers/utils/anyx"
	"github.com/lazygophers/utils/candy"
	"github.com/lazygophers/utils/osx"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// 生成 TypeScript 代码的逻辑
func GenTs(pb *codegen.PbPackage, baseTsDir string) error {
	pkgName := pb.PackageName()
	if pkgName == "protobuf" {
		pkgName = strings.TrimSuffix(pb.ProtoFileName(), filepath.Ext(pb.ProtoFileName()))
	}

	filePath := filepath.Join(baseTsDir, pkgName+".d.ts")

	if !osx.IsDir(filepath.Dir(filePath)) {
		err := os.MkdirAll(filepath.Dir(filePath), 0777)
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}
	}

	b := log.GetBuffer()
	defer log.PutBuffer(b)

	b.WriteString("export namespace ")
	b.WriteString(pkgName)
	b.WriteString(" {\n\n")

	// 生成枚举类型
	{
		var source, enum, options bytes.Buffer

		for _, e := range pb.Enums() {
			var isListOption bool

			if e.FullName == "ErrCode" {
				continue
			}

			if strings.EqualFold(e.Name, "ListOption") {
				isListOption = true
			}

			source.Reset()
			enum.Reset()
			options.Reset()

			source.WriteString("\texport const enum ")
			source.WriteString(e.FullName)
			source.WriteString(" {\n")

			if !isListOption {
				enum.WriteString("\texport const ")
				enum.WriteString(e.FullName)
				enum.WriteString("Enum: {\n")
			}

			if !isListOption {
				options.WriteString("\texport const ")
				options.WriteString(e.FullName)
				options.WriteString("Options: [\n")
			}

			for _, field := range e.FieldList() {
				desc := strings.TrimSpace(field.Desc())

				source.WriteString("\t\t")
				source.WriteString(field.Name)
				source.WriteString(" ")
				source.WriteString("=")
				source.WriteString(" ")
				source.WriteString(anyx.ToString(field.Value))
				source.WriteString(",")

				if desc != "" {
					source.WriteString(" // ")
					source.WriteString(desc)
				}

				source.WriteString("\n")

				if !isListOption && desc == "" {
					if strings.EqualFold(field.Name, "Nil") ||
						field.Value == 0 {
						continue
					}

					return fmt.Errorf("%s.%s not found desc", e.FullName, field.FieldName())
				} else if isListOption && field.Value == 0 {
					continue
				}

				if !isListOption {
					enum.WriteString("\t\t")
					enum.WriteString("[")
					enum.WriteString(e.FullName)
					enum.WriteString(".")
					enum.WriteString(field.Name)
					enum.WriteString("]: '")
					enum.WriteString(desc)
					enum.WriteString("',\n")
				}

				if !isListOption {
					options.WriteString("\t\t")
					options.WriteString("{label: '")
					options.WriteString(desc)
					options.WriteString("', value: ")
					options.WriteString(e.FullName)
					options.WriteString(".")
					options.WriteString(field.Name)
					options.WriteString("},\n")
				}
			}
			source.WriteString("\t}\n\n")

			if !isListOption {
				enum.WriteString("\t};\n\n")
			}

			if !isListOption {
				options.WriteString("\t];\n\n")
			}

			if source.Len() > 0 {
				_, _ = source.WriteTo(b)
			}

			if enum.Len() > 0 {
				_, _ = enum.WriteTo(b)
			}

			if options.Len() > 0 {
				_, _ = options.WriteTo(b)
			}
		}
	}

	var importPackageList []string

	//getTsFieldType := func(typ string) string {
	//	switch typ {
	//	case "int32", "int64", "uint32", "uint64":
	//		return "bigint"
	//
	//	case "float", "double":
	//		return "number"
	//
	//	case "string", "bytes":
	//		return "string"
	//
	//	case "bool":
	//		return "boolean"
	//	default:
	//		return "any"
	//	}
	//}

	getTsFieldKeyType := func(typ string) string {
		switch typ {
		case "int32", "int64", "uint32", "uint64":
			fallthrough

		case "float", "double":
			return "number"

		case "string", "bytes":
			return "string"

		case "bool":
			//return "boolean"
			fallthrough
		default:
			log.Panicf("unsuppost typ %s", typ)
			return "any"
		}
	}

	getTsNormalFieldType := func(field *codegen.PbNormalField) string {
		switch field.Field().Type {
		case "int32", "int64", "uint32", "uint64":
			return "bigint"

		case "float", "double":
			return "number"

		case "string", "bytes":
			return "string"

		case "bool":
			return "boolean"
		default:
			if field.Field().Type == "google.protobuf.Any" {
				return "any"
			}

			// 尝试看看是不是枚举的数据
			if strings.Contains(field.FullType(), ".") {
				if pb.GetEnum(strings.ReplaceAll(field.FullType(), ",", "_")) != nil {
					return strings.ReplaceAll(field.FullType(), ",", "_")
				}
			} else {
				var names []string
				checkEnum := func() bool {
					return pb.GetEnum(strings.Join(candy.Reverse(names), "_")) != nil
				}

				var walk func(e proto.Visitee) (stop bool)
				walk = func(e proto.Visitee) (stop bool) {
					switch x := e.(type) {
					case *proto.Message:
						names = append(names, x.Name)
						if checkEnum() {
							return true
						}
						if walk(x.Parent) {
							return true
						}

					case *proto.Enum:
						names = append(names, x.Name)
						if checkEnum() {
							return true
						}
						if walk(x.Parent) {
							return true
						}

					case *proto.EnumField:
						if x.Parent != nil {
							if walk(x.Parent.(*proto.Enum).Parent) {
								return true
							}
						}

					case *proto.Proto:

					case *proto.NormalField:
						if walk(x.Parent) {
							return true
						}

					default:
						log.Panicf("unknown parent type:%T", x)
					}

					return false
				}
				names = append(names, field.Type())
				if walk(field.Field()) {
					return strings.Join(candy.Reverse(names), "_")
				}
			}

			switch strings.Count(field.Field().Type, ".") {
			case 0:
				return field.Field().Type
			case 1:
				before, after, _ := strings.Cut(field.Field().Type, ".")
				importPackageList = append(importPackageList, before)
				return fmt.Sprintf("%s.%s", before, after)
			default:
				index := strings.LastIndex(field.Field().Type, ".")
				after := field.Field().Type[index+1:]
				before := field.Field().Type[:index]

				index = strings.LastIndex(before, ".")
				before = before[index+1:]

				importPackageList = append(importPackageList, before)
				return fmt.Sprintf("%s.%s", before, after)
			}

			// lazygophers.lrpc.core.Paginate
			// ListSubscribeRsp
			return "any"
		}
	}

	getTsMapFieldType := func(field *codegen.PbMapField) string {
		switch field.Field().Type {
		case "int32", "int64", "uint32", "uint64":
			return "bigint"

		case "float", "double":
			return "number"

		case "string", "bytes":
			return "string"

		case "bool":
			return "boolean"
		default:
			if field.Field().Type == "google.protobuf.Any" {
				return "any"
			}

			// 尝试看看是不是枚举的数据
			if strings.Contains(field.FieldType(), ".") {
				if pb.GetEnum(strings.ReplaceAll(field.FieldType(), ",", "_")) != nil {
					return strings.ReplaceAll(field.FieldType(), ",", "_")
				}
			} else {
				var names []string
				checkEnum := func() bool {
					return pb.GetEnum(strings.Join(candy.Reverse(names), "_")) != nil || pb.GetMessage(strings.Join(candy.Reverse(names), "_")) != nil
				}

				var walk func(e proto.Visitee) (stop bool)
				walk = func(e proto.Visitee) (stop bool) {
					switch x := e.(type) {
					case *proto.Message:
						names = append(names, x.Name)
						if checkEnum() {
							return true
						}
						if walk(x.Parent) {
							return true
						}

					case *proto.Enum:
						names = append(names, x.Name)
						if checkEnum() {
							return true
						}
						if walk(x.Parent) {
							return true
						}

					case *proto.EnumField:
						if x.Parent != nil {
							if walk(x.Parent.(*proto.Enum).Parent) {
								return true
							}
						}

					case *proto.Proto:

					case *proto.MapField:
						if walk(x.Parent) {
							return true
						}

					default:
						log.Panicf("unknown parent type:%T", x)
					}

					return false
				}
				names = append(names, field.FieldType())
				if walk(field.Field()) {
					return strings.Join(candy.Reverse(names), "_")
				}
			}

			switch strings.Count(field.Field().Type, ".") {
			case 0:
				return field.Field().Type
			case 1:
				before, after, _ := strings.Cut(field.Field().Type, ".")
				importPackageList = append(importPackageList, before)
				return fmt.Sprintf("%s.%s", before, after)
			default:
				index := strings.LastIndex(field.Field().Type, ".")
				after := field.Field().Type[index+1:]
				before := field.Field().Type[:index]

				index = strings.LastIndex(before, ".")
				before = before[index+1:]

				importPackageList = append(importPackageList, before)
				return fmt.Sprintf("%s.%s", before, after)
			}

			// lazygophers.lrpc.core.Paginate
			// ListSubscribeRsp
			return "any"
		}
	}

	// 遍历所有类型对于每个类型，生成相应的 TypeScript 常量定义
	for _, message := range pb.Messages() {
		if message.Message().IsExtend {
			continue
		}

		b.WriteString("\t export class ")
		b.WriteString(message.Name)
		b.WriteString(" {\n")

		// 遍历字段
		for _, field := range message.NormalFields() {
			b.WriteString("\t\t")
			b.WriteString(field.FieldName())
			b.WriteString("?:")
			b.WriteString("\t")
			if field.IsSlice() {
				b.WriteString("[")
				b.WriteString(getTsNormalFieldType(field))
				b.WriteString("]")
			} else {
				b.WriteString(getTsNormalFieldType(field))
			}
			b.WriteString(";")

			b.WriteString("\n")
		}

		// readonly [P in keyof T]: T[P];
		for _, field := range message.MapFields() {
			b.WriteString("\t\t")
			b.WriteString(field.FieldName())
			b.WriteString("?:")
			b.WriteString("\t")

			b.WriteString("Record<")
			b.WriteString(getTsFieldKeyType(field.KeyType()))
			b.WriteString(",")
			b.WriteString(getTsMapFieldType(field))
			b.WriteString(">")

			b.WriteString(";")
			b.WriteString("\n")
		}

		//for _, field := range message.EnumFields() {
		//	log.Info(field.Name)
		//	log.Info(field.FieldType())
		//}

		b.WriteString("\t}\n\n")
	}

	b.WriteString("}")

	importPackageList = candy.Unique(importPackageList)
	importPackageList = candy.FilterNot(importPackageList, func(s string) bool {
		return s == ""
	})
	importPackageList = candy.Sort(importPackageList)

	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	defer file.Close()

	//file.WriteString("import { translate } from '../i18n/tran';\n")

	if len(importPackageList) > 0 {
		for _, pkg := range importPackageList {
			file.WriteString("import { ")
			file.WriteString(pkg)
			file.WriteString(" } from \"./")
			file.WriteString(pkg)
			file.WriteString("\";\n")
		}
		file.WriteString("\n")
	}

	_, err = io.Copy(file, b)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	return nil
}

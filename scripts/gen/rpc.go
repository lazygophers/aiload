package main

import (
	"bytes"
	"github.com/lazygophers/codegen/codegen"
	"github.com/lazygophers/log"
	"github.com/lazygophers/utils/candy"
	"github.com/lazygophers/utils/stringx"
	"os"
	"strings"
)

func GenRpc(pb *codegen.PbPackage) error {
	var b bytes.Buffer

	file, err := os.Open("scripts/gen/api.js")
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	_, err = file.WriteTo(&b)
	_ = file.Close()
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	b.WriteString("\n")

	splitType := func(typ string) (string, string) {
		var before, after string
		index := strings.LastIndex(typ, ".")
		if index > 0 {
			after = typ[index+1:]
			before = typ[:index]
		} else {
			after = typ
		}
		index = strings.LastIndex(before, ".")
		if index > 0 {
			before = before[index+1:]
		}

		if before == "" {
			before = pb.PackageName()
		}

		return before, after
	}

	var importPackageList []string
	for _, rpc := range pb.RPCs() {
		b.WriteString("export async function ")
		b.WriteString(stringx.ToSmallCamel(rpc.Name))

		// 处理请求的类型
		{
			b.WriteString("(body")

			b.WriteString(": ")
			before, after := splitType(rpc.RequestType())
			importPackageList = append(importPackageList, before)
			b.WriteString(before)
			b.WriteString(".")
			b.WriteString(after)
			b.WriteString(")")
		}

		// 处理返回类型
		{
			b.WriteString(": ")
			before, after := splitType(rpc.ReturnsType())
			importPackageList = append(importPackageList, before)
			b.WriteString(before)
			b.WriteString(".")
			b.WriteString(after)
			b.WriteString(" {\n")
		}

		b.WriteString("\treturn ")
		b.WriteString(strings.ToLower(rpc.Method()))
		b.WriteString("('")
		b.WriteString(rpc.Path())
		b.WriteString("', body);\n")
		b.WriteString("}\n")
		b.WriteString("\n")
	}

	importPackageList = candy.Unique(importPackageList)
	importPackageList = candy.FilterNot(importPackageList, func(s string) bool {
		return s == ""
	})
	importPackageList = candy.Sort(importPackageList)

	file, err = os.OpenFile("assets/src/utils/api.js", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	defer file.Close()

	for _, pkg := range importPackageList {
		file.WriteString("import {")
		file.WriteString(pkg)
		file.WriteString("} from \"@constant/")
		file.WriteString(pkg)
		file.WriteString("\";\n")
	}

	b.WriteTo(file)

	return nil
}

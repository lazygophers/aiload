package main

import (
	"fmt"
	"github.com/lazygophers/codegen/codegen"
	"github.com/lazygophers/log"
	"github.com/lazygophers/utils/anyx"
	"gopkg.in/yaml.v3"
	"os"
)

func GenErrorCode(pb *codegen.PbPackage) error {
	// 尝试获取一下 errcode

	enum := pb.GetEnum("ErrCode")
	if enum == nil {
		return fmt.Errorf("can't find enum 'ErrCode'")
	}

	buffer, err := os.ReadFile("internal/state/localize/zh.yaml")
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	var m map[string]any
	err = yaml.Unmarshal(buffer, &m)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	if len(m) == 0 {
		m = map[string]any{}
	}

	em := anyx.ToMapInt32String(m["error"])

	for _, field := range enum.FieldList() {
		desc := field.Desc()
		if desc != "" {
			em[field.Value] = desc
		}
	}

	m["error"] = em

	buffer, err = yaml.Marshal(m)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	err = os.WriteFile("internal/state/localize/zh.yaml", buffer, 0600)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	return nil
}

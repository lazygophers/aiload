package main

import (
	"github.com/lazygophers/codegen/codegen"
	"github.com/lazygophers/codegen/state"
	"github.com/lazygophers/log"
	"os"
)

func main() {
	err := state.Load()
	if err != nil {
		log.Errorf("err:%v", err)
		os.Exit(1)
		return
	}

	pb, err := codegen.ParseProto("aiload.proto")
	if err != nil {
		log.Errorf("err:%v", err)
		os.Exit(1)
		return
	}

	err = GenEnum(pb)
	if err != nil {
		log.Errorf("err:%v", err)
		os.Exit(1)
		return
	}

	err = GenJsEnum(pb)
	if err != nil {
		log.Errorf("err:%v", err)
		os.Exit(1)
		return
	}

	err = GenErrorCode(pb)
	if err != nil {
		log.Errorf("err:%v", err)
		os.Exit(1)
		return
	}

	err = GenRpc(pb)
	if err != nil {
		log.Errorf("err:%v", err)
		os.Exit(1)
		return
	}

	err = GenJsI18n()
	if err != nil {
		log.Errorf("err:%v", err)
		os.Exit(1)
		return
	}

	err = GenGoI18n()
	if err != nil {
		log.Errorf("err:%v", err)
		os.Exit(1)
		return
	}

	err = GenTs(pb, "assets/src/constant")
	if err != nil {
		log.Errorf("err:%v", err)
		os.Exit(1)
		return
	}
}

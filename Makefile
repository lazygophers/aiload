LANG := zh_CN.UTF-8
LC_ALL := zh_CN.UTF-8
TZ := UTC

PWD=$(shell pwd)
Organization=lazygophers
ProjectName=aiload
Version=$(shell git rev-parse --short=12 HEAD)
MihomoVersion=$(shell cd /Users/luoxin/persons/go/lazygophers/barbecue/pkg/mihomo && git rev-parse --short=12 HEAD)

BUILDTIME=$(shell date -u '+%Y-%m-%d %H:%M')
GOMODCACHE=$(shell go env GOMODCACHE)
GOCACHE=$(shell go env GOCACHE)
GONOPROXY=$(shell go env GONOPROXY)
GONOSUMDB=$(shell go env GONOSUMDB)
GOPRIVATE=$(shell go env GOPRIVATE)
GOOS=$(shell go env GOOS)
GOARCH=$(shell go env GOARCH)
GOARM=$(shell go env GOARM)
GOAMD64=$(shell go env GOAMD64)
GOMIPS=$(shell go env GOMIPS)
GOVERSION := $(shell go version | awk -F'go' '{print $$3}' | awk '{print $$1}')

.PHONY: gen
gen: ## 代码生成
	codegen g pb -i ./aiload.proto -d \
	--go-module-prefix="github.com/lazygophers/" \
	--add-proto-files="../proto"

	codegen g impl -i ./aiload.proto -d \
	--go-module-prefix="github.com/lazygophers/" \
	--template-impl-route="scripts/template/rpc_route.gtpl" \
	--template-impl-path="scripts/template/rpc_path.gtpl"

	codegen g table -i ./aiload.proto -d \
	--go-module-prefix="github.com/lazygophers/"

	#go run -v ./scripts/gen/

.PHONY: fmt
fmt: ## 格式化
	gofmt -w .

 .PHONY: run
 run: ## 运行
	go run -v ./cmd

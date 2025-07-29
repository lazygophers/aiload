LANG := zh_CN.UTF-8
LC_ALL := zh_CN.UTF-8
TZ := UTC

PWD=$(shell pwd)
Organization=lazygophers
ProjectName=aiload
Version=$(shell git rev-parse --short=12 HEAD)

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
	@$(MAKE) fmt

	#go run -v ./scripts/gen/

.PHONY: fmt
fmt: ## 格式化
	gofmt -w .

 .PHONY: run
 run: fmt ## 运行
	CGO_ENABLED=0 \
	GODEBUG=madvdontneed=1,asyncpreemptoff=1 \
	go build \
		-trimpath \
		--gcflags '-N -l' \
		--tags netgo,osusergo \
		--ldflags '-checklinkname=0 -s -w --extldflags "-fpic" -X "github.com/lazygophers/utils/app.Name=${ProjectName}" -X "github.com/lazygophers/utils/app.Version=${Version}" -X "github.com/lazygophers/utils/app.Organization=${Organization}" -X "github.com/lazygophers/utils/app.GoVersion=${GOVERSION}" -X "github.com/lazygophers/utils/app.GoOS=${GOOS}" -X "github.com/lazygophers/utils/app.Goarch=${GOARCH}"' \
		 ./cmd


.PHONY: help
help: ## 显示此帮助消息
	@awk 'BEGIN {FS = ":.*##"; printf "Usage: make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2 }' $(MAKEFILE_LIST)

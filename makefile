SVC = auth
BUILD_DIR = build
CGO_ENABLED ?= 0
GOOS ?= linux
GOARCH ?= amd64
VERSION ?= $(shell git describe --abbrev=0 --tags || echo "none")
COMMIT ?= $(shell git rev-parse HEAD)
TIME ?= $(shell date +%F_%T)


define compile_service
	CGO_ENABLED=$(CGO_ENABLED) GOOS=$(GOOS) GOARCH=$(GOARCH) GOARM=$(GOARM) \
	go build -mod=mod -ldflags "-s -w \
	-X 'github.com/ianmuchyri/jwt-Authentication.BuildTime=$(TIME)' \
	-X 'github.com/ianmuchyri/jwt-Authentication.Version=$(VERSION)' \
	-X 'github.com/ianmuchyri/jwt-Authentication.Commit=$(COMMIT)'" \
	-o $(BUILD_DIR)/$(SVC) main.go
endef

all: $(SVC)

.PHONY: $(SVC)

clean:
	rm -rf $(BUILD_DIR)

install: 
	cp $(BUILD_DIR)/$(SVC) $(GOBIN)

auth:
	$(call compile_service,$(@))

run:
	${BUILD_DIR}/$(SVC)

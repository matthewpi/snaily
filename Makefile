CGO    = 0
GOOS   = linux
GOARCH = amd64

OUTPUT_FILE = snaily

BUILD_VERSION := `git describe --abbrev=0  2> /dev/null || echo ""`
BUILD_BRANCH  := `git rev-parse --abbrev-ref HEAD`
BUILD_COMMIT  := `git rev-parse HEAD | head -c 8`
BUILD_DATE    := `date +%Y-%m-%d`

LDFLAGS = -ldflags "-s -w -X main.buildVersion=$(BUILD_VERSION) -X main.buildBranch=$(BUILD_BRANCH) -X main.buildCommit=$(BUILD_COMMIT) -X main.buildDate=$(BUILD_DATE)"

all: clean build

clean:
	@go clean
	@rm $(OUTPUT_FILE) -f

build:
	@CGO_ENABLED=$(CGO) GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o $(OUTPUT_FILE) -v $(LDFLAGS) main.go

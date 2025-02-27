MODULE=github.com/sean9999/pear
#REPO=$$(git remote -v | head -n 1 | cut -f 2 | cut -f 1 -d ' '| cut -c 5-)
SEMVER := $$(git tag --sort=-version:refname | head -n 1)
BRANCH := $$(git branch --show-current)
REF := $$(git describe --dirty --tags --always)
GOPROXY=proxy.golang.org

info:
	@printf "MODULE:\t%s\nSEMVER:\t%s\nBRANCH:\t%s\nREF:\t%s\n" $(MODULE) $(SEMVER) $(BRANCH) $(REF)

tidy:
	go mod tidy

clean:
	go clean
	go clean -modcache
	rm bin/*

pkgsite:
	if [ -z "$$(command -v pkgsite)" ]; then go install golang.org/x/pkgsite/cmd/pkgsite@latest; fi

docs: pkgsite
	pkgsite -open .

publish:
	GOPROXY=https://$(GOPROXY),direct go list -m $(MODULE)@$(SEMVER)

test:
	git restore testdata
	go test ./...
	git restore testdata

.PHONY: test

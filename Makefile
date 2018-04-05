ZIPLIST=\
examples/specifications/petstore \
examples/apikey_injection \
examples/guides \
examples/metadata \
examples/overlay \
assets 

UNIX_LIST=${ZIPLIST} run_example.sh
WIN_LIST=${ZIPLIST} run_example.bat

VERSION=$(shell git rev-parse --abbrev-ref HEAD 2> /dev/null | sed s/release.v//)

BZW=./buildzip $@ dapperdox.exe $+
BZU=./buildzip $@ dapperdox     $+

STEM=dist/dapperdox-${VERSION}

#LDFLAGS=-ldflags "-X main.VERSION=${VERSION}"

all:
	@echo "Build DapperDox..."; 
	go get && go build ${LDFLAGS}

release: distribution \
	${STEM}.linux-x86.tgz \
	${STEM}.linux-amd64.tgz \
	${STEM}.darwin-amd64.tgz \
	${STEM}.linux-arm.tgz \
	${STEM}.linux-arm64.tgz \
	${STEM}.windows-x86.zip \
	${STEM}.windows-amd64.zip \
	releaseTable

distribution:
	@mkdir -p dist; \
	rm -f dist/*

releaseTable: dist/release-table-${VERSION}.md

dist/release-table-${VERSION}.md:
	rm $@; \
	./createReleaseTable.sh ${VERSION} > $@

${STEM}.linux-arm.tgz: dapperdox_linux_arm.exe ${UNIX_LIST}
	@${BZU}

${STEM}.linux-arm64.tgz: dapperdox_linux_arm64.exe ${UNIX_LIST}
	@${BZU}

${STEM}.linux-amd64.tgz: dapperdox_linux_amd64.exe ${UNIX_LIST}
	@${BZU}

${STEM}.linux-x86.tgz: dapperdox_linux_x86.exe ${UNIX_LIST}
	@${BZU}

${STEM}.darwin-amd64.tgz: dapperdox_darwin_amd64.exe ${UNIX_LIST}
	@${BZU}

${STEM}.windows-x86.zip: dapperdox_win_x86.exe ${WIN_LIST}
	@${BZW}

${STEM}.windows-amd64.zip: dapperdox_win_amd64.exe ${WIN_LIST}
	@${BZW}
	
dapperdox_linux_x86.exe: main.go
	GOOS=linux GOARCH=386 go build ${LDFLAGS} -o $@

dapperdox_linux_amd64.exe: main.go
	GOOS=linux GOARCH=amd64 go build ${LDFLAGS} -o $@

dapperdox_linux_arm64.exe: main.go
	GOOS=linux GOARCH=arm64 go build ${LDFLAGS} -o $@

dapperdox_linux_arm.exe: main.go
	GOOS=linux GOARCH=arm go build ${LDFLAGS} -o $@

dapperdox_darwin_amd64.exe: main.go
	GOOS=darwin GOARCH=amd64 go build ${LDFLAGS} -o $@

dapperdox_win_x86.exe: main.go
	GOOS=windows GOARCH=386 go build ${LDFLAGS} -o $@

dapperdox_win_amd64.exe: main.go
	GOOS=windows GOARCH=amd64 go build ${LDFLAGS} -o $@

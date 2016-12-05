default: darwin linux windows

darwin:
	GOOS=darwin GOARCH=amd64 go build -o weatherpresenter_loadtest.darwin-amd64

linux:
	GOOS=linux GOARCH=amd64 go build -o weatherpresenter_loadtest.linux-amd64

windows:
	GOOS=windows GOARCH=amd64 go build -o weatherpresenter_loadtest.windows-amd64

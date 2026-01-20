$env:GOOS="linux"; $env:GOARCH="amd64"; $env:CGO_ENABLED="0"; go build -o rdr-linux-amd64
$env:GOOS="linux"; $env:GOARCH="arm64"; $env:CGO_ENABLED="0"; go build -o rdr-linux-arm64
$env:GOOS="darwin"; $env:GOARCH="amd64"; $env:CGO_ENABLED="0"; go build -o rdr-darwin-amd64
$env:GOOS="darwin"; $env:GOARCH="arm64"; $env:CGO_ENABLED="0"; go build -o rdr-darwin-arm64
$env:GOOS="windows"; $env:GOARCH="amd64"; $env:CGO_ENABLED="0"; go build -o rdr-windows-amd64
$env:GOOS="windows"; $env:GOARCH="arm64"; $env:CGO_ENABLED="0"; go build -o rdr-windows-arm64
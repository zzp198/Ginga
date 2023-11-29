SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=amd64

SET name=Crond

go build -ldflags="-s -w" -o %name%

SSH root@223.240.111.27 "pkill %name%"
SCP %name% root@223.240.111.27:
DEL %name%
SSH root@223.240.111.27 "chmod +x %name% && ./%name%"

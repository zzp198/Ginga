SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=amd64

SET name=Ginga
SET ip=223.240.111.27

go build -ldflags="-s -w" -o %name%

SSH root@%ip% "pkill %name%"
SCP %name% root@%ip%:/root/%name%/%name%
DEL %name%
SSH root@%ip% "cd %name% && chmod +x %name% && ./%name%"

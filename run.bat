SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=amd64

SET name=Ginga
SET ip=223.240.111.27

@REM go build -ldflags="-s -w" -o %name%
go build -o %name%/%name%

SSH root@%ip% "pkill %name%"
SCP -r %name% root@%ip%:/root
SSH root@%ip% "cd %name% && chmod +x %name% && ./%name%"

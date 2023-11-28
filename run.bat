SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=amd64

go build -ldflags="-s -w" -o Crond/Crond

@REM ssh root@223.240.111.27
scp.exe Crond/Crond root@223.240.111.27:/root/Crond
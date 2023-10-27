SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=amd64
go build -ldflags="-s -w" -o Quark\Quark
@REM upx.exe --best Quark\Quark
@REM upx.exe Quark\Quark

cd Quark
wsl pkill Quark
wsl ./Quark
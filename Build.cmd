@echo off

if not exist go.mod (
	echo Initializing go module...
	go mod init main 2> nul
)

if not exist go.sum (
	echo Tidying go module...
	go mod tidy 2> nul
)

echo Configuring file names...
for /f "tokens=*" %%0 in ('dir /b ^| findstr /e .flac') do (
	set FLAC=%%0
	set WAV=%%~n0.wav
	set app=%%~n0
)

echo Writing "embed.go"...
(
	echo package main
	echo import _ "embed"
	echo //go:embed "%FLAC%"
	echo var flacRaw []byte
	echo var wavName string = "%WAV%"
) > embed.go

timeout /t 3 /nobreak > nul

:Menu
echo.
echo Generate executable for which operating system and architecture?
echo 1.] Windows x86_64
echo 2.] Windows x86
echo 3.] Linux x86_64
echo 4.] Linux x86
echo 5.] Darwin [Mac] x86_64
echo.
echo X.] Exit
choice /c 12345x /n
goto %errorlevel%

:1
set GOARCH=amd64
set GOOS=windows
set file=%app%.exe
goto Build

:2
set GOARCH=386
set GOOS=windows
set file=%app%.exe
goto Build

:3
set GOARCH=amd64
set GOOS=linux
set file=%app%
goto Build

:4
set GOARCH=386
set GOOS=linux
set file=%app%
goto Build

:5
set GOARCH=amd64
set GOOS=darwin
set file=%app%.app
goto Build

:6
exit /b

:Build
echo Building "Release/%file%"...
call go build -ldflags="-s -w" -o "Release/%file%" flacsfx.go embed.go

if %errorlevel%==0 (echo Build successful!
) else echo Build unsuccessful!

goto Menu
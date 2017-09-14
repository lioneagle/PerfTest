@echo off
setlocal
if exist install.bat goto ok
echo install.bat must be run from its folder
goto end
: ok
call env.bat
gofmt -w src

if "%2" == "" (
	set output32=.\bin\%132.exe
	set output64=.\bin\%164.exe
) else (
	set output32=.\bin\%232.exe
	set output64=.\bin\%264.exe
)

if %GOARCH% == amd64 (
	echo installing 64-bit ......
	go install %1

	if exist .\bin\%1.exe (
		copy .\bin\%1.exe  %output64%
		del .\bin\%1.exe
		
		rem echo building 32-bit ......
		
		rem set GOARCH=386
		rem go build -o=%output32% %1
		rem set GOARCH=amd64
	)

) else (

	go install %1
	copy .\bin\%1.exe  %output32%
	del .\bin\%1.exe

)

:end
echo finished
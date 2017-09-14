@echo off
setlocal
if exist coverage.bat goto ok
echo coverage.bat must be run from its folder
goto end
: ok

call env.bat

if not exist test_temp mkdir test_temp

if exist .\test_temp\coverage.out  del .\test_temp\coverage.out

set param1=%1
set param2=%2
set param3=%3
set param4=%4

if "%1" == "-html" set param1=
if "%2" == "-html" set param2=
if "%3" == "-html" set param3=
if "%4" == "-html" set param4=

if "%1" == "-html" set html=1
if "%2" == "-html" set html=1
if "%3" == "-html" set html=1
if "%4" == "-html" set html=1

go test -coverprofile=./test_temp/coverage.out %param1% %param2% %param3% %param4% 
if not exist .\test_temp\coverage.out goto end

if "%html%" == "1" (
	go tool cover -html=./test_temp/coverage.out -o ./test_temp/coverage.html
	if exist .\test_temp\coverage.html  .\test_temp\coverage.html
)

:end
echo finished
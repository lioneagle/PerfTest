@echo off
setlocal
if exist coverage_all.bat goto ok
echo coverage_all.bat must be run from its folder
goto end
: ok

call env.bat

if not exist test_temp mkdir test_temp

if exist .\test_temp\coverage.out  del .\test_temp\coverage.out
if exist .\test_temp\package.html del .\test_temp\package.txt

go list ./... | findstr -v "main github vendor" >> .\test_temp\package.txt

for /f %%a in (.\test_temp\package.txt) do (
	if exist .\test_temp\coverage.out (
		go test -coverprofile=./test_temp/coverage1.out %%a
		if exist .\test_temp\coverage1.out (
			findstr -v "mode": .\test_temp\coverage1.out >> .\test_temp\coverage.out
			@echo off
			del .\test_temp\coverage1.out
		)
	) else (
		go test -coverprofile=./test_temp/coverage.out %%a
	)
)
if exist .\test_temp\package.txt del .\test_temp\package.txt

if exist .\test_temp\coverage.out (
	go tool cover -func=./test_temp/coverage.out -o ./test_temp/coverage.txt
	findstr "total" .\test_temp\coverage.txt >> .\test_temp\coverage2.txt
	del .\test_temp\coverage.txt
	
	for /f "tokens=1,2,3 delims=	" %%a in (.\test_temp\coverage2.txt) do (
	    echo %%a %%c of statements
	)
	del .\test_temp\coverage2.txt
	
	if "%1" == "-html" (
		go tool cover -html=./test_temp/coverage.out -o ./test_temp/coverage.html
		if exist .\test_temp\coverage.html .\test_temp\coverage.html
	)
)

:end
echo finished
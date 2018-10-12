@echo off
title Forecast_app batch file
:: See the tite at the top.
echo compiling go files...
go build

echo Installing dependencies....
go get github.com/tools/godep
godep save ./
for %%I in (.) do set CurrDirName=%%~nxI
start %CurrDirName%.exe

@echo off
SET pro=gateway
echo '��ʼ������Ŀ'
set goos=linux
echo '===10%==='
 if exist %pro% (
    del %pro%
 )
echo '===20%==='
go build -o %pro%
echo '===90%==='
set goos=windows
echo '===100%===='
echo '��Ŀ�������'
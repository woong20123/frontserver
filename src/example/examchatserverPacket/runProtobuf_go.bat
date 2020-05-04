:: 현재 경로 지정
set CURRENT_PATH=%~dp0

:: proto 파일 컴파일 
protoc -I%CURRENT_PATH% --go_out=%CURRENT_PATH% %CURRENT_PATH%LogicPacket.proto
pause()
set FILE_PATH=%~dp0
echo %FILE_PATH%

protoc -I%FILE_PATH% --go_out=%FILE_PATH% %FILE_PATH%LogicPacket.proto

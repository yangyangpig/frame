@echo off
protoc --gofast_out=. match.proto -I=. -I=../../../ -I=../../../github.com/gogo/protobuf/protobuf
echo "success"
pause
